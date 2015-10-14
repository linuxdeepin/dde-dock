package fonts

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"

	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/glib-2.0"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	fallbackStandard  = "Droid Sans"
	fallbackMonospace = "Droid Sans Mono"
	defaultDPI        = 96

	xsettingsSchema = "com.deepin.xsettings"
	gsKeyFontName   = "gtk-font-name"
)

var locker sync.Mutex

type Family struct {
	Id   string
	Name string

	Styles []string
	//Files  []string
}
type Families []*Family

func ListStandardFamily() Families {
	return ListFont().ListStandard().convertToFamilies()
}

func ListMonospaceFamily() Families {
	return ListFont().ListMonospace().convertToFamilies()
}

func ListAllFamily() Families {
	return ListFont().convertToFamilies()
}

func IsFontFamily(value string) bool {
	info := ListAllFamily().Get(value)
	if info != nil {
		return true
	}
	return false
}

func IsFontSizeValid(size int32) bool {
	if size >= 7 && size <= 22 {
		return true
	}
	return false
}

func SetFamily(standard, monospace string) error {
	standInfo := ListStandardFamily().Get(standard)
	if standInfo == nil {
		return fmt.Errorf("Invalid standard id '%s'", standard)
	}
	monoInfo := ListMonospaceFamily().Get(monospace)
	if monoInfo == nil {
		return fmt.Errorf("Invalid monospace id '%s'", monospace)
	}

	curStand := fcFontMatch("sans-serif")
	curMono := fcFontMatch("monospace")
	if (standInfo.Id == curStand || standInfo.Name == curStand) &&
		(monoInfo.Id == curMono || monoInfo.Name == curMono) {
		return nil
	}

	setFontByXSettings(standard, -1)
	return writeFontConfig(configContent(standard, monospace),
		path.Join(glib.GetUserConfigDir(), "fontconfig", "fonts.conf"))
}

func SetSize(size int32) error {
	return setFontByXSettings(fcFontMatch("sans-serif"), size)
}

func GetFontSize() int32 {
	setting, _ := dutils.CheckAndNewGSettings(xsettingsSchema)
	defer setting.Unref()

	return getFontSize(setting)
}

func (infos Families) GetIds() []string {
	var ids []string
	for _, info := range infos {
		ids = append(ids, info.Id)
	}
	return ids
}

func (infos Families) Get(id string) *Family {
	for _, info := range infos {
		if info.Id == id {
			return info
		}
	}
	return nil
}

func (infos Families) add(info *Family) Families {
	v := infos.Get(info.Id)
	if v == nil {
		infos = append(infos, info)
		return infos
	}

	v.Styles = compositeList(v.Styles, info.Styles)
	//v.Files = compositeList(v.Files, info.Files)
	return infos
}

func setFontByXSettings(name string, size int32) error {
	setting, err := dutils.CheckAndNewGSettings(xsettingsSchema)
	if err != nil {
		return err
	}
	defer setting.Unref()

	if size == -1 {
		size = getFontSize(setting)
	}
	v := fmt.Sprintf("%s %v", name, size)
	if v == setting.GetString(gsKeyFontName) {
		return nil
	}

	setting.SetString(gsKeyFontName, v)
	return nil
}

func getFontSize(setting *gio.Settings) int32 {
	value := setting.GetString(gsKeyFontName)
	if len(value) == 0 {
		return 0
	}
	array := strings.Split(value, " ")
	size, _ := strconv.ParseInt(array[len(array)-1], 10, 64)
	return int32(size)
}

func compositeList(l1, l2 []string) []string {
	for _, v := range l2 {
		if isItemInList(v, l1) {
			continue
		}
		l1 = append(l1, v)
	}
	return l1
}

func isItemInList(item string, list []string) bool {
	for _, v := range list {
		if item == v {
			return true
		}
	}
	return false
}

func writeFontConfig(content, file string) error {
	locker.Lock()
	defer locker.Unlock()
	err := os.MkdirAll(path.Dir(file), 0755)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(file, []byte(content), 0644)
}

// If set pixelsize, wps-office-wps will not show some text.
//
//func configContent(standard, mono string, pixel float64) string {
func configContent(standard, mono string) string {
	return fmt.Sprintf(`<?xml version="2.0"?>
<!DOCTYPE fontconfig SYSTEM "fonts.dtd">
<fontconfig>
    <match target="pattern">
        <test qual="any" name="family">
            <string>serif</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
	    <string>%s</string>
	    <string>%s</string>
	</edit>
    </match>
 
    <match target="pattern">
        <test qual="any" name="family">
            <string>sans-serif</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
	    <string>%s</string>
	    <string>%s</string>
	</edit>
    </match>

    <match target="pattern">
        <test qual="any" name="family">
            <string>monospace</string>
        </test>
        <edit name="family" mode="assign" binding="strong">
	    <string>%s</string>
	    <string>%s</string>
	</edit>
    </match>

    <match target="font">
	<edit name="antialias" mode="assign">
	    <bool>true</bool>
	</edit>
	<edit name="hinting" mode="assign">
	    <bool>true</bool>
	</edit>
	<edit name="hintstyle" mode="assign">
	    <const>hintfull</const>
        </edit>
	<edit name="rgba" mode="assign">
	    <const>rgb</const>
	</edit>
    </match>

</fontconfig>`, standard, fallbackStandard,
		standard, fallbackStandard,
		mono, fallbackMonospace)
}
