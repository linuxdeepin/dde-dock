package main

import "encoding/json"
import "strings"
import "github.com/BurntSushi/xgb/randr"
import "fmt"
import "os"
import "io/ioutil"

var __CFG__ _Configuration

var _ConfigPath = os.Getenv("HOME") + "/.config/deepin_monitors.json"

var _CurrentRight, _RightX = "", int16(0)

type _Configuration struct {
	Primary     string
	DisplayMode int16
	Monitors    map[string]*_MonitorConfiguration
}

type _MonitorConfiguration struct {
	Name string

	Width, Height uint16
	RefreshRate   float64

	X, Y         int16
	RelativeInfo [2]string

	Enabled  bool
	Rotation uint16
	Reflect  uint16

	Brightness map[string]float64
}

func readCfg() _Configuration {
	cfg := _Configuration{Monitors: make(map[string]*_MonitorConfiguration)}
	f, err := os.Open(_ConfigPath)
	if err != nil {
		Logger.Warning("Failed load displayConfiguration:", err)
		return cfg
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	err = json.Unmarshal(data, &cfg)
	return cfg
}

func writeCfg(cfg _Configuration) {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		panic("marshal display configuration failed:" + err.Error())
	}
	f, err := os.Create(_ConfigPath)
	if err != nil {
		fmt.Println("Couldn't save display configuration:", err)
		return
	}
	defer f.Close()
	f.Write(bytes)
}
func (dpy *Display) SaveChanged() {
	__CFG__.Monitors = make(map[string]*_MonitorConfiguration)
	for _, m := range dpy.Monitors {
		__CFG__.Monitors[m.Name] = m.saveStatus()
	}
	__CFG__.Primary = dpy.Primary
	writeCfg(__CFG__)
	dpy.SwitchMode(DisplayModeCustom)
	dpy.detectChanged()
}

func saveBrightness(name string, v float64) {
	cfg := readCfg()
	for _, mcfg := range cfg.Monitors {
		for _name, _ := range mcfg.Brightness {
			if _name == name {
				mcfg.Brightness[name] = v
			}
		}
	}
	writeCfg(cfg)
}
func initBrightness(monitors []*Monitor) {
	for _, m := range monitors {
		for _, mcfg := range __CFG__.Monitors {
			if mcfg.Name == m.Name {
				for k, v := range mcfg.Brightness {
					m.ChangeBrightness(k, v)
				}
			}
		}
	}
}

var __LastCode__ = ""

func (dpy *Display) Apply() {
	if dpy.HasChanged {
		__LastCode__ = dpy.generateShell()
		runCode(__LastCode__)
	}
	for _, m := range dpy.Monitors {
		m.updateInfo()
	}
	dpy.detectChanged()
}

func (dpy *Display) detectChanged() {
	if __LastCode__ != dpy.generateShell() {
		dpy.setPropHasChanged(true)
	} else {
		dpy.setPropHasChanged(false)
	}
}
func (dpy *Display) generateShell() string {
	code := "xrandr "
	for _, m := range dpy.Monitors {
		code += m.generateShell()
		if dpy.Primary == m.Name {
			code += " --primary"
		}
	}
	return code
}

const (
	DisplayModeUnknow  = -100
	DisplayModeMirrors = -1
	DisplayModeCustom  = 0
	DisplayModeOnlyOne = 1
)

func (dpy *Display) SwitchMode(mode int16) {
	if dpy.DisplayMode == mode {
		return
	}

	dpy.setPropDisplayMode(mode)

	if mode == DisplayModeMirrors {
		w, h := getMirrorSize(dpy.Monitors)
		for _, m := range dpy.Monitors {
			m.SwitchOn(true)
			m.SetPos(0, 0)
			m.ensureSize(w, h)
		}
	} else if mode == DisplayModeCustom {
		dpy.ResetChanged()
	} else if mode >= DisplayModeOnlyOne && int(mode) <= len(dpy.Monitors) {
		for i, m := range dpy.Monitors {
			if i+1 == int(mode) {
				m.SetPos(0, 0)
				m.SetMode(m.BestMode.ID)
				m.SwitchOn(true)
				dpy.SetPrimary(m.Name)
				fmt.Println("SetSwitch..", m.Name, m.Opened)
			} else {
				m.SwitchOn(false)
				fmt.Println("SetSwitch..", m.Name, m.Opened)
			}
		}
	} else {
		return
	}
	dpy.Apply()
}

func (m *Monitor) restore(cfg *_MonitorConfiguration) {
	m.setPropXY(cfg.X, cfg.Y)
	m.SetRelativePos(cfg.RelativeInfo[0], cfg.RelativeInfo[1])
	m.ensureSize(cfg.Width, cfg.Height)
	m.SwitchOn(cfg.Enabled)
	m.setPropRotation(cfg.Rotation)
	m.setPropReflect(cfg.Reflect)
	for k, v := range cfg.Brightness {
		if v != 0 {
			m.setPropBrightness(k, v)
		}
	}
}
func (dpy *Display) ResetChanged() {
	// dond't set the monitors which hasn't cfg information
	__CFG__ = readCfg()
	dpy.SetPrimary(__CFG__.Primary)
	for _, cfg := range __CFG__.Monitors {
		for _, m := range dpy.Monitors {
			if m.Name == cfg.Name {
				m.restore(cfg)
			}
		}
	}
	dpy.HasChanged = true //force apply saved configuration
	dpy.Apply()

	for _, m := range dpy.Monitors {
		m.updateInfo()
	}
	dpy.detectChanged()
}

func (m *Monitor) saveStatus() *_MonitorConfiguration {
	return &_MonitorConfiguration{
		Name:         m.Name,
		Width:        m.Width,
		Height:       m.Height,
		RefreshRate:  m.CurrentMode.Rate,
		X:            m.X,
		Y:            m.Y,
		RelativeInfo: m.relativePosInfo,
		Enabled:      m.Opened,
		Rotation:     m.Rotation,
		Reflect:      m.Reflect,
		Brightness:   m.Brightness,
	}
}

func (dpy *Display) updateMonitorList() {
	resources, err := randr.GetScreenResources(X, Root).Reply()
	if err != nil {
		return
	}
	monitors := make([]*Monitor, 0)
	for _, op := range resources.Outputs {
		oinfo, err := randr.GetOutputInfo(X, op, LastConfigTimeStamp).Reply()
		if err != nil || oinfo.Connection != randr.ConnectionConnected {
			continue
		}
		m := NewMonitor([]randr.Output{op})
		if mcfg, ok := __CFG__.Monitors[m.Name]; ok {
			m.relativePosInfo = mcfg.RelativeInfo
		}
		monitors = append(monitors, m)
	}
	setAutoFlag := len(dpy.Monitors) > len(monitors)
	dpy.setPropMonitors(monitors)
	for _, m := range __CFG__.Monitors {
		dpy.tryJoin(m.Name)
	}

	_CurrentRight, _RightX = "", 0
	for _, m := range dpy.Monitors {
		if _CurrentRight == "" {
			_CurrentRight = m.Name
		} else if m.X > _RightX {
			_CurrentRight = m.Name
			_RightX = m.X
		}

		if cfg, ok := __CFG__.Monitors[m.Name]; ok {
			if dpy.DisplayMode == DisplayModeCustom {
				m.restore(cfg)
			}
		} else {
			m.SwitchOn(true)

			m.SetRelativePos(_CurrentRight, "right-of")
			_CurrentRight = m.Name
			_RightX += int16(m.CurrentMode.Width)

			__CFG__.Monitors[m.Name] = m.saveStatus()
		}
	}
	if setAutoFlag {
		runCode("xrandr --auto")
	}
	dpy.Apply()
}

func (dpy *Display) tryJoin(name string) {
	names := strings.Split(name, joinSeparator)
	joined := names[0]
	for i := 1; i < len(names); i++ {
		dpy.JoinMonitor(joined, names[i])
		fmt.Println("TryJoin:", joined, names[i])
		joined += joinSeparator + names[i]
	}
}
