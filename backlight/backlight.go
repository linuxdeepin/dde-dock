package backlight

import (
	"fmt"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

const (
	BacklightUnknow   string = "unknow"
	BacklightRaw             = "raw"
	BacklightPlatform        = "platform"
	BacklightFirmware        = "firmware"
)

const (
	lcdBacklightSysDir = "/sys/class/backlight"
	kbdBacklightSysDir = "/sys/class/leds"
)

type SyspathInfo struct {
	Path          string
	Type          string
	MaxBrightness int32
}
type SyspathInfos []*SyspathInfo

func ListLCDBacklight() SyspathInfos {
	var infos SyspathInfos
	paths := doListSyspath(lcdBacklightSysDir)
	for _, p := range paths {
		info, err := NewSyspathInfo(p)
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}

	if len(infos) > 1 {
		infos = infos.sortLCD()
	}
	return infos
}

func ListKbdBacklight() SyspathInfos {
	var infos SyspathInfos
	paths := doListSyspath(kbdBacklightSysDir)
	for _, p := range paths {
		if !strings.Contains(p, "kbd_backlight") {
			continue
		}
		info, err := NewSyspathInfo(p)
		if err != nil {
			continue
		}
		infos = append(infos, info)
	}
	return infos
}

func NewSyspathInfo(syspath string) (*SyspathInfo, error) {
	max, err := getMaxBrightness(syspath)
	if err != nil {
		return nil, err
	}

	ty, _ := getType(syspath)
	return &SyspathInfo{
		Path:          syspath,
		Type:          ty,
		MaxBrightness: max,
	}, nil
}

func (info *SyspathInfo) GetBrightness() (int32, error) {
	var file = path.Join(info.Path, "brightness")
	return doGetBrightness(file)
}

func (info *SyspathInfo) SetBrightness(value int32) error {
	v, err := info.GetBrightness()
	if err != nil {
		return err
	}

	if value == v {
		return nil
	}
	return doSetBrightness(info.Path, value)
}

func (infos SyspathInfos) Get(syspath string) (*SyspathInfo, error) {
	for _, info := range infos {
		if info.Path == syspath {
			return info, nil
		}
	}
	return nil, fmt.Errorf("Invalid syspath: %s", syspath)
}

// sort Sort by type, raw > firmware > platform
func (infos SyspathInfos) sortLCD() SyspathInfos {
	if len(infos) < 2 {
		return infos
	}

	var set = make(map[string]SyspathInfos)
	for _, info := range infos {
		set[info.Type] = append(set[info.Type], info)
	}

	var ret SyspathInfos
	ret = append(ret, set[BacklightRaw]...)
	ret = append(ret, set[BacklightPlatform]...)
	ret = append(ret, set[BacklightFirmware]...)
	return ret
}

func getType(syspath string) (string, error) {
	var file = path.Join(syspath, "type")
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return strings.TrimRight(string(content), "\n"), nil
}

func getMaxBrightness(syspath string) (int32, error) {
	var file = path.Join(syspath, "max_brightness")
	return doGetBrightness(file)
}

func doGetBrightness(file string) (int32, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}

	tmp := strings.TrimRight(string(content), "\n")
	v, err := strconv.ParseInt(tmp, 10, 64)
	return int32(v), err
}

func doSetBrightness(file string, value int32) error {
	return ioutil.WriteFile(path.Join(file, "brightness"),
		[]byte(fmt.Sprintf("%d\n", value)), 0644)
}

func doListSyspath(dir string) []string {
	finfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	var paths []string
	for _, finfo := range finfos {
		paths = append(paths, path.Join(dir, finfo.Name()))
	}
	return paths
}
