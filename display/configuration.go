package main

import "encoding/json"
import "strings"
import "github.com/BurntSushi/xgb/randr"
import "fmt"
import "os"
import "io/ioutil"

var __CFG__ = make(map[string]_MonitorConfiguration)

var _ConfigPath = os.Getenv("HOME") + "/.config/deepin_monitors.json"

type _MonitorConfiguration struct {
	Name string

	Width, Height uint16
	RefreshRate   float64

	X, Y     int16
	Primary  bool
	Enabled  bool
	Rotation uint16
	Reflect  uint16

	Brightness float64
}

func saveConfiguration() {
	bytes, err := json.Marshal(__CFG__)
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

func (dpy *Display) Apply() {
	code := "xrandr "
	hasPrimary := false
	for _, m := range dpy.Monitors {
		code += m.generateShell()
		if m.IsPrimary {
			hasPrimary = true
		}
	}
	if !hasPrimary {
		code += " --noprimary"
	}
	runCode(code)
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
				if m.IsPrimary {
					dpy.SetPrimary(m.Name)
				} else {
					dpy.SetPrimary("")
				}
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

func (dpy *Display) ResetChanged() {
	// dond't set the monitors which hasn't cfg information
	for _, cfg := range __CFG__ {
		for _, m := range dpy.Monitors {
			if m.Name == cfg.Name {
				m.SetPos(cfg.X, cfg.Y)
				fmt.Println("SetRotation:", m.Name, cfg.Rotation)
				m.ensureSize(cfg.Width, cfg.Height)
				m.SwitchOn(cfg.Enabled)
				m.setPrimary(cfg.Primary)
				m.setPropRotation(cfg.Rotation)
				m.setPropReflect(cfg.Reflect)
				m.setPropBrightness(cfg.Brightness)
			}
		}
	}
	dpy.Apply()
}
func (dpy *Display) SaveChanged() {
	__CFG__ = make(map[string]_MonitorConfiguration)
	for _, m := range dpy.Monitors {
		__CFG__[m.Name] = _MonitorConfiguration{
			Name:        m.Name,
			Width:       m.Width,
			Height:      m.Height,
			RefreshRate: m.Rate,
			X:           m.X,
			Y:           m.Y,
			Primary:     m.IsPrimary,
			Enabled:     m.Opened,
			Rotation:    m.Rotation,
			Reflect:     m.Reflect,
			Brightness:  m.Brightness,
		}
	}
	saveConfiguration()
}

func loadConfiguration(dpy *Display) {
	f, err := os.Open(_ConfigPath)
	if err != nil {
		fmt.Println("OpenFailed", err)
		return
	}
	data, err := ioutil.ReadAll(f)
	err = json.Unmarshal(data, &__CFG__)
	if err != nil {
		fmt.Println("Failed load displayConfiguration:", err, "Data:", string(data))
		return
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
		monitors = append(monitors, NewMonitor([]randr.Output{op}))
	}
	dpy.resetMonitors(monitors)
	for _, m := range __CFG__ {
		dpy.tryJoin(m.Name)
	}
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
