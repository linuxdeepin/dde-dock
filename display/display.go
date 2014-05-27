package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/logger"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"os"
	"strings"
	"sync"
)

var (
	xcon, _ = xgb.NewConn()
	_       = initX11()

	Root           xproto.Window
	ScreenWidthMm  uint16
	ScreenHeightMm uint16

	LastConfigTimeStamp xproto.Timestamp

	MinWidth, MinHeight, MaxWidth, MaxHeight uint16

	Logger = logger.NewLogger("com.deepin.daemon.Display")
)

func initX11() bool {
	randr.Init(xcon)
	sinfo := xproto.Setup(xcon).DefaultScreen(xcon)
	Root = sinfo.Root
	ScreenWidthMm = sinfo.WidthInMillimeters
	ScreenHeightMm = sinfo.HeightInMillimeters
	LastConfigTimeStamp = xproto.Timestamp(0)

	ver, err := randr.QueryVersion(xcon, 1, 3).Reply()
	if err != nil {
		panic(fmt.Sprintln("randr.QueryVersion error:", err))
	}
	if ver.MajorVersion != 1 || ver.MinorVersion != 3 {
		panic(fmt.Sprintln("randr version is too low:", ver.MajorVersion, ver.MinorVersion, "this program require at least randr 1.3"))
	}
	if err != nil {
		panic(fmt.Sprintln("randr.GetSceenSizeRange failed :", err))
	}
	return true
}

var GetDisplay = func() func() *Display {
	dpy := &Display{}

	sinfo := xproto.Setup(xcon).DefaultScreen(xcon)
	dpy.setPropScreenWidth(sinfo.WidthInPixels)
	dpy.setPropScreenHeight(sinfo.HeightInPixels)
	GetDisplayInfo().update()
	dpy.setPropHasChanged(false)

	randr.SelectInputChecked(xcon, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|randr.NotifyMaskCrtcChange|randr.NotifyMaskScreenChange)

	return func() *Display {
		return dpy
	}
}()

type DisplayInfo struct {
	locker         sync.Mutex
	modes          map[randr.Mode]Mode
	outputNames    map[string]randr.Output
	backlightLevel map[string]uint32
}

var GetDisplayInfo = func() func() *DisplayInfo {
	info := &DisplayInfo{
		modes:          make(map[randr.Mode]Mode),
		outputNames:    make(map[string]randr.Output),
		backlightLevel: make(map[string]uint32),
	}
	info.update()
	return func() *DisplayInfo {
		return info
	}
}()

func (info *DisplayInfo) QueryModes(id randr.Mode) Mode {
	return Mode{}
}
func (info *DisplayInfo) QueryOutputs(name string) randr.Output {
	return 0
}
func (info *DisplayInfo) QueryBacklightLevel(name string) uint32 {
	return 0
}

func (info *DisplayInfo) update() {
	info.locker.Lock()
	defer info.locker.Unlock()

	resource, err := randr.GetScreenResources(xcon, Root).Reply()
	if err != nil {
		Logger.Error("GetScreenResouces failed", err)
		return
	}
	info.outputNames = make(map[string]randr.Output)
	info.backlightLevel = make(map[string]uint32)
	for _, op := range resource.Outputs {
		oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
		if err != nil {
			Logger.Warning("DisplayInfo.update filter:", err)
			continue
		}
		if oinfo.Connection != randr.ConnectionConnected {
			continue
		}

		info.outputNames[string(oinfo.Name)] = op
		info.backlightLevel[string(oinfo.Name)] = uint32(queryBacklightRange(xcon, op))
	}
	//if len(info.outputNames) != 2 {
	//Logger.Warning("XX", info.outputNames, resource.Outputs)
	//for _, op := range resource.Outputs {
	//oinfo, err := randr.GetOutputInfo(xcon, op, LastConfigTimeStamp).Reply()
	//if err != nil {
	//fmt.Println("XX E:", err)
	//}
	//fmt.Println("XX:", string(oinfo.Name), oinfo.Connection)
	//}
	//}

	info.modes = make(map[randr.Mode]Mode)
	for _, minfo := range resource.Modes {
		info.modes[randr.Mode(minfo.Id)] = buildMode(minfo)
	}
}

type Display struct {
	Monitors    []*Monitor
	monitorLock sync.RWMutex

	ScreenWidth  uint16
	ScreenHeight uint16

	//used by deepin-dock/launcher/desktop
	Primary        string
	PrimaryRect    xproto.Rectangle
	PrimaryChanged func(xproto.Rectangle)

	DisplayMode   int16
	BuiltinOutput *Monitor

	HasChanged bool

	Brightness map[string]float64
	cfg        *ConfigDisplay
}

func (dpy *Display) lockMonitors() {
	dpy.monitorLock.Lock()
}
func (dpy *Display) unlockMonitors() {
	dpy.monitorLock.Unlock()
}
func (dpy *Display) rLockMonitors() {
	dpy.monitorLock.RLock()
}
func (dpy *Display) rUnlockMonitors() {
	dpy.monitorLock.RUnlock()
}

func (dpy *Display) listener() {
	for {
		e, err := xcon.WaitForEvent()
		if err != nil {
			continue
		}
		switch ee := e.(type) {
		case randr.NotifyEvent:
			switch ee.SubCode {
			case randr.NotifyCrtcChange:
			case randr.NotifyOutputChange:
				info := ee.U.Oc
				if info.Connection != randr.ConnectionConnected && info.Mode != 0 {
					randr.SetCrtcConfig(xcon, info.Crtc, xproto.TimeCurrentTime, LastConfigTimeStamp, 0, 0, 0, randr.RotationRotate0, nil)
				}
			case randr.NotifyOutputProperty:
			}
		case randr.ScreenChangeNotifyEvent:
			dpy.setPropScreenWidth(ee.Width)
			dpy.setPropScreenHeight(ee.Height)

			GetDisplayInfo().update()

			if LastConfigTimeStamp < ee.ConfigTimestamp {
				LastConfigTimeStamp = ee.ConfigTimestamp
				if dpy.QueryCurrentPlanName() != dpy.cfg.CurrentPlanName {
					Logger.Info("Detect New ConfigTimestmap, try reset changes")
					dpy.ResetChanges()
				}
			}

			//sync Monitor's state
			for _, m := range dpy.Monitors {
				m.updateInfo()
			}

			//SetPrimary will try set the valid primary
			dpy.SetPrimary(dpy.Primary)
		}
	}
}

func (dpy *Display) ChangeBrightness(output string, v float64) {
	if v >= 0 && v <= 1 {
		if op, ok := GetDisplayInfo().outputNames[output]; ok {
			if max, ok := GetDisplayInfo().backlightLevel[output]; ok && max != 0 {
				setOutputBacklight(op, uint32(float64(max)*v))
			} else {
				setBrightness(xcon, op, v)
			}
			dpy.setPropBrightness(output, v)
		}
	} else {
		Logger.Warningf("Try change the brightness of %s to an invalid value(%v)", output, v)
	}

}
func (dpy *Display) ResetBrightness(output string) {
	if v, ok := LoadConfigDisplay(dpy).Brightness[output]; ok {
		dpy.SetBrightness(output, v)

	}
}
func (dpy *Display) SetBrightness(output string, v float64) {
	if v >= 0 && v <= 1 {
		dpy.ChangeBrightness(output, v)
		dpy.saveBrightness(output, v)
	} else {
		Logger.Warningf("Try set the brightness of %s to an invalid value(%v)", output, v)
	}
}

func (dpy *Display) JoinMonitor(a string, b string) error {
	dpy.lockMonitors()
	defer dpy.unlockMonitors()

	ms := dpy.cfg.Monitors[dpy.cfg.CurrentPlanName]
	if ma, ok := ms[a]; ok {
		if mb, ok := ms[b]; ok {
			mc := mergeConfigMonitor(dpy, ma, mb)
			delete(dpy.cfg.Monitors[dpy.cfg.CurrentPlanName], a)
			delete(dpy.cfg.Monitors[dpy.cfg.CurrentPlanName], b)
			dpy.cfg.Monitors[dpy.cfg.CurrentPlanName][mc.Name] = mc

			var newMonitors []*Monitor
			for _, m := range dpy.Monitors {
				if m.Name != a && m.Name != b {
					newMonitors = append(newMonitors, m)
				}
			}
			newMonitors = append(newMonitors, NewMonitor(dpy, mc))
			dpy.setPropMonitors(newMonitors)
		} else {
			return fmt.Errorf("Can't find Monitor %s\n", b)
		}
	} else {
		return fmt.Errorf("Can't find Monitor %s\n", a)
	}
	return nil
}
func (dpy *Display) SplitMonitor(a string) error {
	dpy.lockMonitors()
	defer dpy.unlockMonitors()

	var monitors []*Monitor
	found := false
	for _, m := range dpy.Monitors {
		if m.Name == a {
			submonitors := m.split(dpy)
			if submonitors == nil {
				return fmt.Errorf("Can't find composited monitor: %s", a)
			}
			found = true
			monitors = append(monitors, submonitors...)
		} else {
			monitors = append(monitors, m)
		}
	}
	if found {
		dpy.setPropMonitors(monitors)
		return nil
	} else {
		return fmt.Errorf("Can't find composited monitor: %s", a)
	}
}
func (m *Monitor) split(dpy *Display) (r []*Monitor) {
	if !strings.Contains(m.Name, joinSeparator) {
		return
	}

	delete(dpy.cfg.Monitors[dpy.QueryCurrentPlanName()], m.Name)
	dpyinfo := GetDisplayInfo()
	for _, name := range strings.Split(m.Name, joinSeparator) {
		if op, ok := dpyinfo.outputNames[name]; ok {
			mcfg, err := CreateConfigMonitor(dpy, op)
			if err != nil {
				Logger.Error("Failed createconfigmonitor at split", err, name, mcfg)
				continue
			}
			dpy.cfg.Monitors[dpy.QueryCurrentPlanName()][name] = mcfg

			minfo := dpyinfo.modes[mcfg.bestMode]
			mcfg.Width = minfo.Width
			mcfg.Height = minfo.Height
			mcfg.currentMode = mcfg.bestMode

			m := NewMonitor(dpy, mcfg)
			m.SetMode(m.BestMode.ID)
			r = append(r, m)
		}
	}
	return
}

func (dpy *Display) detectChanged() {
	dpy.setPropHasChanged(!dpy.cfg.Compare(LoadConfigDisplay(dpy)))
}

func (dpy *Display) SetPrimary(name string) error {
	if m, ok := dpy.cfg.Monitors[dpy.cfg.CurrentPlanName][name]; ok {
		if m.Enabled {
			dpy.setPropPrimary(name)
			dpy.cfg.Primary = name
			dpy.savePrimary(dpy.cfg.Primary)
			dpy.setPropPrimaryRect(xproto.Rectangle{m.X, m.Y, m.Width, m.Height})
			return nil
		}
	}

	if name != dpy.Primary {
		dpy.SetPrimary(dpy.Primary)
	}

	for _, m := range dpy.cfg.Monitors[dpy.cfg.CurrentPlanName] {
		if m.Name != name && m.Enabled {
			dpy.setPropPrimary(name)
			dpy.cfg.Primary = name
			dpy.savePrimary(dpy.cfg.Primary)
			dpy.setPropPrimaryRect(xproto.Rectangle{m.X, m.Y, m.Width, m.Height})
			return nil
		}
	}

	err := fmt.Errorf("Can't set primary to ", name)
	Logger.Fatal(err.Error())
	return err
}

func (dpy *Display) Apply() {
	dpy.apply(false)
}
func (dpy *Display) apply(auto bool) {
	dpy.rLockMonitors()
	defer dpy.rUnlockMonitors()

	code := "xrandr "
	for _, m := range dpy.Monitors {
		code += m.generateShell()
		if auto {
			code += " --auto"
		}

		if dpy.cfg.Primary == m.Name {
			code += " --primary"
		}
	}
	runCode(code)
}

func (dpy *Display) ResetChanges() {
	dpy.cfg = LoadConfigDisplay(dpy)

	//must be invoked after LoadConfigDisplay(dpy)
	var monitors []*Monitor
	for _, mcfg := range dpy.cfg.Monitors[dpy.cfg.CurrentPlanName] {
		m := NewMonitor(dpy, mcfg)
		monitors = append(monitors, m)
	}
	dpy.setPropMonitors(monitors)

	dpy.SetPrimary(dpy.cfg.Primary)

	//apply the saved configurations.
	dpy.apply(false)

	dpy.Brightness = make(map[string]float64)
	for name, v := range dpy.cfg.Brightness {
		dpy.Brightness[name] = v
		dpy.ChangeBrightness(name, v)

		//set brightness to 1, if the output support backlight feature
		if op, ok := GetDisplayInfo().outputNames[name]; ok {
			if max, ok := GetDisplayInfo().backlightLevel[name]; ok && max != 0 {
				setBrightness(xcon, op, 1)
			}
		}
	}
	dpy.detectChanged()
}

func (dpy *Display) SaveChanges() {
	dpy.cfg.Save()
}

func (dpy *Display) Reset() {
	dpy.rLockMonitors()
	defer dpy.rUnlockMonitors()

	for _, m := range dpy.Monitors {
		dpy.SetBrightness(m.Name, 1)
		m.SetReflect(0)
		m.SetRotation(1)
		m.SetMode(m.BestMode.ID)
	}
	dpy.apply(true)
}

func main() {
	defer Logger.EndTracing()

	if !dlib.UniqueOnSession("com.deepin.daemon.Display") {
		Logger.Warning("Another com.deepin.daemon.Display is running")
		return
	}

	dpy := GetDisplay()
	dpy.ResetChanges()
	go dpy.listener()

	for _, m := range dpy.Monitors {
		m.updateInfo()
	}

	err := dbus.InstallOnSession(dpy)
	if err != nil {
		Logger.Error("Can't install dbus display service on session:", err)
		return
	}
	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		Logger.Error("lost dbus session:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}

func (dpy *Display) QueryOutputFeature(name string) int32 {
	if max, ok := GetDisplayInfo().backlightLevel[name]; ok && max != 0 {
		return 1
	}
	return 0
}
