package main

import (
	"dlib/dbus"
	"dlib/logger"
	"fmt"
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
	"os"
)

var (
	_              = fmt.Println
	DPY            *Display
	X, _           = xgb.NewConn()
	Root           = xproto.Setup(X).DefaultScreen(X).Root
	ScreenWidthMm  = xproto.Setup(X).DefaultScreen(X).WidthInMillimeters
	ScreenHeightMm = xproto.Setup(X).DefaultScreen(X).HeightInMillimeters

	LastConfigTimeStamp = xproto.Timestamp(0)

	MinWidth, MinHeight, MaxWidth, MaxHeight uint16

	Logger = logger.NewLogger("com.deepin.daemon.Display")
)

type Display struct {
	modes map[randr.Mode]Mode

	Monitors []*Monitor

	ScreenWidth  uint16
	ScreenHeight uint16

	//used by deepin-dock/launcher/desktop
	Primary        string
	PrimaryRect    xproto.Rectangle
	PrimaryChanged func(xproto.Rectangle)

	DisplayMode   int16
	BuiltinOutput *Monitor

	HasChanged bool
}

func (dpy *Display) JoinMonitor(a string, b string) error {
	var monitorA, monitorB *Monitor
	newMonitors := make([]*Monitor, 0)
	for _, m := range dpy.Monitors {
		if m.Name == a {
			monitorA = m
		} else if m.Name == b {
			monitorB = m
		} else {
			newMonitors = append(newMonitors, m)
		}
	}
	if monitorA == nil {
		return fmt.Errorf("cann't find %s", a)
	}
	if monitorB == nil {
		return fmt.Errorf("cann't find %s", b)
	}

	ops := monitorA.outputs
	ops = append(ops, monitorB.outputs...)
	monitor := NewMonitor(ops)
	monitor.SetMode(monitor.BestMode.ID)
	monitor.SetPos(0, 0)

	if monitor != nil {
		dpy.setPropMonitors(append(newMonitors, monitor))
		dpy.setPropHasChanged(true)
		return nil
	} else {
		return fmt.Errorf("can't create composted monitor")
	}
	return nil
}
func (dpy *Display) SplitMonitor(a string) error {
	newMonitors := make([]*Monitor, 0)
	var monitor *Monitor
	for _, m := range dpy.Monitors {
		if m.Name == a {
			monitor = m
		} else {
			newMonitors = append(newMonitors, m)
		}
	}
	if monitor == nil {
		return fmt.Errorf("can't find composited monitor: %s", a)
	}

	haveStaticPosMonitor := false
	for _, op := range monitor.outputs {
		m := NewMonitor([]randr.Output{op})
		if m == nil {
			return fmt.Errorf("can't create monitor: %d", op)
		}
		if cfg, ok := __CFG__.Monitors[m.Name]; ok {
			m.restore(cfg)
		} else {
			if haveStaticPosMonitor {
				m.SetRelativePos(_CurrentRight, "right-of")
			} else {
				m.SetPos(0, 0)
				haveStaticPosMonitor = true
			}
			_CurrentRight = m.Name
			_RightX += int16(m.CurrentMode.Width)
		}
		newMonitors = append(newMonitors, m)
	}
	dpy.setPropMonitors(newMonitors)
	dpy.setPropHasChanged(true)
	return nil
}

func (dpy *Display) SetPrimary(name string) error {
	fmt.Println("TrySetPrimary:", name)
	var validPrimary *Monitor
	for _, m := range dpy.Monitors {
		if name == m.Name && m.Opened {
			validPrimary = m
		}
	}
	if validPrimary == nil && len(dpy.Monitors) >= 1 {
		builtIn := guestBuiltIn(dpy.Monitors)
		if builtIn.Opened {
			validPrimary = builtIn
		} else {
			for _, m := range dpy.Monitors {
				if m.Opened {
					validPrimary = m
				}
			}
		}
	}
	if validPrimary != nil {
		dpy.setPropPrimary(validPrimary.Name)
		dpy.setPropPrimaryRect(xproto.Rectangle{validPrimary.X, validPrimary.Y, validPrimary.Width, validPrimary.Height})
		dpy.detectChanged()
		fmt.Println("SetTo:", name)
		return nil
	} else {
		dpy.setPropPrimaryRect(xproto.Rectangle{0, 0, dpy.ScreenWidth, dpy.ScreenHeight})
		return fmt.Errorf("Can't find this monitor: %s", name)
	}
}

func initDisplay() *Display {
	dpy := &Display{}
	DPY = dpy
	screen := xproto.Setup(X).DefaultScreen(X)
	dpy.setPropScreenWidth(screen.WidthInPixels)
	dpy.setPropScreenHeight(screen.HeightInPixels)
	dbus.InstallOnSession(dpy)

	__CFG__ = readCfg()
	dpy.Primary = __CFG__.Primary
	dpy.SetPrimary(dpy.Primary)

	dpy.updateInfo()
	dpy.updateMonitorList()
	dpy.SetPrimary(dpy.Primary)

	initBrightness(dpy.Monitors)

	randr.SelectInput(X, Root, randr.NotifyMaskOutputChange|randr.NotifyMaskOutputProperty|randr.NotifyMaskCrtcChange|randr.NotifyMaskScreenChange)
	go dpy.listener()

	dpy.workaroundBacklight()

	dpy.setPropHasChanged(false)
	__LastCode__ = dpy.generateShell()

	for _, m := range dpy.Monitors {
		if m.Name == dpy.Primary {
			dpy.setPropPrimaryRect(xproto.Rectangle{m.X, m.Y, m.Width, m.Height})
		}
	}

	return dpy
}

func (dpy *Display) updateInfo() {
	// update output list
	resources, err := randr.GetScreenResources(X, Root).Reply()

	if err != nil {
		panic("GetScreenResources failed:" + err.Error())
	}

	dpy.modes = make(map[randr.Mode]Mode)
	for _, minfo := range resources.Modes {
		dpy.modes[randr.Mode(minfo.Id)] = buildMode(minfo)
	}
	for _, m := range dpy.Monitors {
		m.updateInfo()
	}
}

func (dpy *Display) listener() {
	for {
		e, err := X.WaitForEvent()
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
					randr.SetCrtcConfig(X, info.Crtc, xproto.TimeCurrentTime, LastConfigTimeStamp, 0, 0, 0, randr.RotationRotate0, nil)
				}
			case randr.NotifyOutputProperty:
			}
		case randr.ScreenChangeNotifyEvent:
			dpy.setPropScreenWidth(ee.Width)
			dpy.setPropScreenHeight(ee.Height)

			dpy.updateInfo()
			if LastConfigTimeStamp < ee.ConfigTimestamp {
				fmt.Println("AAAAAAAAAAAAA")
				LastConfigTimeStamp = ee.ConfigTimestamp
				dpy.updateMonitorList()
			}
			dpy.SetPrimary(dpy.Primary)
		}
	}
}

func main() {
	randr.Init(X)
	ver, err := randr.QueryVersion(X, 1, 3).Reply()
	fmt.Println("VER:", ver)
	if err != nil {
		panic(fmt.Sprintln("randr.QueryVersion error:", err))
	}
	if ver.MajorVersion != 1 || ver.MinorVersion != 3 {
		panic(fmt.Sprintln("randr version is too low:", ver.MajorVersion, ver.MinorVersion, "this program require at least randr 1.3"))
	}

	if err != nil {
		panic(fmt.Sprintln("randr.GetSceenSizeRange failed :", err))
	}

	initDisplay()

	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		Logger.Error("lost dbus session:", err)
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
