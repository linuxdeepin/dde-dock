package main

import "github.com/BurntSushi/xgb/randr"
import "encoding/json"
import "fmt"
import "os"
import "io/ioutil"
import "sync"
import "strings"
import "sort"

const (
	DPModeUnknow  = -100
	DPModeMirrors = -1
	DPModeNormal  = 0
	DPModeOnlyOne = 1
)

var _ConfigPath = os.Getenv("HOME") + "/.config/deepin_monitors.json"
var configLock sync.RWMutex

func (dpy *Display) QueryCurrentPlanName() string {
	names := make([]string, 0)
	for name, _ := range GetDisplayInfo().outputNames {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
	//return base64.NewEncoding("1").EncodeToString([]byte(strings.Join(names, ",")))
}

func (cfg *ConfigDisplay) attachCurrentMonitor(dpy *Display) {
	cfg.CurrentPlanName = dpy.QueryCurrentPlanName()
	if _, ok := cfg.Monitors[cfg.CurrentPlanName]; ok {
		cfg.ensureValid(dpy)
		return
	}
	Logger.Warning("attachCurrentMonitor: build info")

	//grab and build monitors information
	monitors := make(map[string]*ConfigMonitor)
	for _, op := range GetDisplayInfo().outputNames {
		mcfg, err := CreateConfigMonitor(dpy, op)
		if err != nil {
			Logger.Warning("skip invalid monitor", op)
			continue
		}
		monitors[mcfg.Name] = mcfg
	}

	//save it at CurrentPlanName slot
	cfg.Monitors[cfg.CurrentPlanName] = monitors

	cfg.Primary = guestPrimaryName()

	//query brightness information
	for name, op := range GetDisplayInfo().outputNames {
		var support bool
		if support, cfg.Brightness[name] = supportedBacklight(xcon, op); support {
			//Assume the brightness is 1.0 if there hasn't any saved information
			GetDisplayInfo().backlightLevel[name] = uint32(queryBacklightRange(xcon, op))
		} else {
			cfg.Brightness[name] = 1
		}
	}
	cfg.ensureValid(dpy)
}

func createConfigDisplay(dpy *Display) *ConfigDisplay {
	cfg := &ConfigDisplay{}
	cfg.Monitors = make(map[string]map[string]*ConfigMonitor)
	cfg.Brightness = make(map[string]float64)
	cfg.DisplayMode = DPModeNormal

	cfg.attachCurrentMonitor(dpy)
	return cfg
}

func (cfg *ConfigDisplay) updateMonitorPlan(dpy *Display) {
}

func (cfg *ConfigDisplay) ensureValid(dpy *Display) {
	var opend []*ConfigMonitor
	var any *ConfigMonitor

	for _, m := range cfg.Monitors[cfg.CurrentPlanName] {
		any = m
		if m.Enabled {
			opend = append(opend, m)
		}

		//1.1. ensure the output support the mode which be matched with the width/height
		valid := false
		for _, opName := range m.Outputs {
			op := GetDisplayInfo().outputNames[opName]
			oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
			if err != nil {
				Logger.Error("ensureValid failed:", err)
				continue
			}
			if len(oinfo.Modes) == 0 {
				Logger.Error("ensureValid failed:", opName, "hasn't any mode info")
				continue
			} else {
				m.bestMode = oinfo.Modes[0]
			}
			for _, id := range oinfo.Modes {
				minfo := GetDisplayInfo().modes[id]
				if minfo.Width == m.Width && minfo.Height == m.Height {
					m.currentMode = randr.Mode(minfo.ID)
					valid = true
					break
				}
			}
		}
		if !valid {
		}
	}
	if any == nil {
		Logger.Fatal("Can't find any ConfigMonitor at ", cfg.CurrentPlanName)
	}
	//1. ensure there has a opened monitor.
	if len(opend) == 0 {
		any.Enabled = true
		opend = append(opend, any)
	}

	//2. ensure primary is opened
	primaryOk := false
	for _, m := range opend {
		if cfg.Primary == m.Name {
			primaryOk = true
			break
		}
	}
	if !primaryOk {
		cfg.Primary = any.Name
	}

	//4. avoid monitor allocation overlay
	valid := true
	for _, m1 := range cfg.Monitors[cfg.CurrentPlanName] {
		for _, m2 := range cfg.Monitors[cfg.CurrentPlanName] {
			if m1 != m2 {
				if isOverlap(m1.X, m1.Y, m1.Width, m1.Height, m2.X, m2.Y, m2.Width, m2.Height) {
					Logger.Warningf("%s(%d,%d,%d,%d) is ovlerlap with %s(%d,%d,%d,%d)! **rearrange all monitor**\n",
						m1.Name, m1.X, m1.Y, m1.Width, m1.Height, m2.Name, m2.X, m2.Y, m2.Width, m2.Height)
					valid = false
					break
				}
			}
		}
	}
	if !valid {
		pm := cfg.Monitors[cfg.CurrentPlanName][cfg.Primary]
		cx, cy, pw, ph := int16(0), int16(0), pm.Width, pm.Height
		pm.X, pm.Y = 0, 0
		Logger.Infof("Rearrange %s to (%d,%d,%d,%d)\n", pm.Name, pm.X, pm.Y, pm.Width, pm.Height)
		for _, m := range cfg.Monitors[cfg.CurrentPlanName] {
			if m != pm {
				cx += int16(pw)
				cy += int16(ph)
				m.X, m.Y = cx, 0
				Logger.Infof("Rearrange %s to (%d,%d,%d,%d)\n", m.Name, m.X, m.Y, m.Width, m.Height)
			}
		}
	}
}

func LoadConfigDisplay(dpy *Display) (r *ConfigDisplay) {
	configLock.RLock()
	defer configLock.RUnlock()

	defer func() {
		if r == nil {
			r = createConfigDisplay(dpy)
		}
	}()

	if f, err := os.Open(_ConfigPath); err != nil {
		return nil
	} else {
		if data, err := ioutil.ReadAll(f); err != nil {
			return nil
		} else {
			cfg := &ConfigDisplay{
				Brightness: make(map[string]float64),
				Monitors:   make(map[string]map[string]*ConfigMonitor),
			}
			if err = json.Unmarshal(data, &cfg); err != nil {
				return nil
			}
			cfg.attachCurrentMonitor(dpy)
			return cfg
		}
	}
	return nil
}

type ConfigDisplay struct {
	DisplayMode     int16
	CurrentPlanName string
	Monitors        map[string]map[string]*ConfigMonitor

	Primary    string
	Brightness map[string]float64
}

func (c *ConfigDisplay) Compare(cfg *ConfigDisplay) bool {
	if c.CurrentPlanName != cfg.CurrentPlanName {
		Logger.Warning("Compare tow ConfigDisply which hasn't same CurrentPlaneName!")
		return false
	}

	if c.Primary != cfg.Primary {
		fmt.Println("Primary NootSame..")
		return false
	}

	for _, m1 := range c.Monitors[c.CurrentPlanName] {
		if m2, ok := cfg.Monitors[c.CurrentPlanName][m1.Name]; ok {
			return m1.Compare(m2)
		} else {
			return false
		}
	}

	return true
}
func (c *ConfigDisplay) Save() {
	configLock.Lock()
	defer configLock.Unlock()

	bytes, err := json.Marshal(c)
	if err != nil {
		Logger.Error("Can't save configure:", err)
		return
	}

	f, err := os.Create(_ConfigPath)
	if err != nil {
		Logger.Error("Cant create configure:", err)
	}
	defer f.Close()
	f.Write(bytes)
}

type ConfigMonitor struct {
	Name    string
	Outputs []string

	currentMode randr.Mode
	bestMode    randr.Mode

	Width, Height uint16
	RefreshRate   float64

	X, Y int16

	Enabled  bool
	Rotation uint16
	Reflect  uint16
}

func mergeConfigMonitor(dpy *Display, a *ConfigMonitor, b *ConfigMonitor) *ConfigMonitor {
	c := &ConfigMonitor{}
	c.Outputs = append(a.Outputs, b.Outputs...)
	c.Name = a.Name + joinSeparator + b.Name
	c.Reflect = 0
	c.Rotation = 1
	c.X, c.Y = 0, 0

	var ops []randr.Output
	for _, opName := range c.Outputs {
		if op, ok := GetDisplayInfo().outputNames[opName]; ok {
			ops = append(ops, op)
		}
	}
	c.Width, c.Height = getMatchedSize(ops)
	c.Enabled = true
	return c
}

func CreateConfigMonitor(dpy *Display, op randr.Output) (*ConfigMonitor, error) {
	cfg := &ConfigMonitor{}
	oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
	if err != nil {
		return nil, err
	}
	cfg.Name = string(oinfo.Name)
	cfg.Outputs = append(cfg.Outputs, cfg.Name)

	if oinfo.Crtc != 0 && oinfo.Connection == randr.ConnectionConnected {
		cinfo, err := randr.GetCrtcInfo(xcon, oinfo.Crtc, LastConfigTimeStamp).Reply()
		if err != nil {
			return nil, err
		}
		cfg.Width, cfg.Height = cinfo.Width, cinfo.Height

		cfg.Rotation, cfg.Reflect = parseRandR(cinfo.Rotation)
		cfg.currentMode = cinfo.Mode
		cfg.Enabled = true
	} else {
		if len(oinfo.Modes) == 0 {
			return nil, fmt.Errorf(string(oinfo.Name), "hasn't any mode info")
		}
		bestMode := oinfo.Modes[0]
		minfo := GetDisplayInfo().modes[bestMode]
		cfg.Width, cfg.Height = minfo.Width, minfo.Height
		cfg.Rotation, cfg.Reflect = 1, 0
		cfg.currentMode = bestMode
		cfg.Enabled = true

		randr.SetCrtcConfig(xcon, oinfo.Crtc, 0, LastConfigTimeStamp, cfg.X, cfg.Y, bestMode, 1, []randr.Output{op})
	}

	return cfg, nil
}

func (c *ConfigMonitor) Save() {
	cfg := LoadConfigDisplay(GetDisplay())
	configLock.Lock()
	defer configLock.Unlock()

	for i, m := range cfg.Monitors[cfg.CurrentPlanName] {
		if m.Name == c.Name {
			cfg.Monitors[cfg.CurrentPlanName][i] = c
			cfg.Save()
			return
		}
	}
	panic("not reached")
}

func (m1 *ConfigMonitor) Compare(m2 *ConfigMonitor) bool {
	if m1.Enabled != m2.Enabled {
		return false
	}
	if m1.Width != m2.Width || m1.Height != m2.Height {
		return false
	}
	if m1.X != m2.X || m1.Y != m2.Y {
		return false
	}
	if m1.Reflect != m2.Reflect {
		return false
	}
	if m1.Rotation != m2.Rotation {
		return false
	}
	return true
}

func (dpy *Display) saveBrightness(output string, v float64) {
	cfg := LoadConfigDisplay(dpy)
	cfg.Brightness[output] = v
	cfg.Save()
}
func (dpy *Display) savePrimary(output string) {
	cfg := LoadConfigDisplay(dpy)
	cfg.Primary = output
	cfg.Save()
}
