package screensaver

import (
	"github.com/BurntSushi/xgb/dpms"
	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/log"
	"sync"
)

var logger = log.NewLogger("org.freedesktop.ScreenSaver")

type inhibitor struct {
	cookie uint32
	name   string
	reason string
}
type ScreenSaver struct {
	xu *xgbutil.XUtil

	IdleOn      func()
	CycleActive func()
	IdleOff     func()

	blank        byte
	idleTime     uint32
	idleInterval uint32

	inhibitors  map[uint32]inhibitor
	counter     uint32
	counterLock sync.Mutex

	//Inhibit state, we need save the SetTimeout value,
	//so we can recover the correct state when enter UnInhibit state.
	lastVals *timeoutVals
}

type timeoutVals struct {
	seconds, interval uint32
	blank             bool
}

func (ss *ScreenSaver) Inhibit(name, reason string) uint32 {
	ss.counterLock.Lock()
	defer ss.counterLock.Unlock()

	ss.counter++

	ss.inhibitors[ss.counter] = inhibitor{ss.counter, name, reason}

	if len(ss.inhibitors) == 1 {
		ss.setTimeout(0, 0, false)
	}
	logger.Infof("\"%s\" want system enter inhibit, because: \"%s\"", name, reason)

	return ss.counter
}

func (ss *ScreenSaver) SimulateUserActivity() {
	xproto.ForceScreenSaver(ss.xu.Conn(), 0)
}

func (ss *ScreenSaver) UnInhibit(cookie uint32) {
	ss.counterLock.Lock()
	defer ss.counterLock.Unlock()

	inhibitor, ok := ss.inhibitors[cookie]
	if !ok {
		logger.Warning("no valid inhibit cookie", cookie)
		return
	}

	logger.Infof("\"%s\" no need inhibit.", inhibitor.name)

	delete(ss.inhibitors, cookie)
	if len(ss.inhibitors) == 0 {
		logger.Info("Enter uninhibit state")
		if ss.lastVals != nil {
			logger.Info("recover from ", ss.lastVals)
			ss.setTimeout(ss.lastVals.seconds, ss.lastVals.interval, ss.lastVals.blank)
			ss.lastVals = nil
		} else {
			ss.setTimeout(ss.idleTime, ss.idleInterval, ss.blank == 1)
		}
	}
}

func (ss *ScreenSaver) SetTimeout(seconds, interval uint32, blank bool) {
	if len(ss.inhibitors) > 0 {
		ss.lastVals = &timeoutVals{seconds, interval, blank}
		logger.Info("Current is inhibit state, the value", ss.lastVals, "will apply when in unhibit state")
	} else {
		ss.setTimeout(seconds, interval, blank)
	}
}

func (ss *ScreenSaver) setTimeout(seconds, interval uint32, blank bool) {
	if blank {
		ss.blank = 1
	} else {
		ss.blank = 0
	}
	xproto.SetScreenSaver(ss.xu.Conn(), int16(seconds), int16(interval), ss.blank, 0)
	dpms.SetTimeouts(ss.xu.Conn(), 0, 0, 0)
	logger.Info("SetTimeout to ", seconds, interval, blank)
}

func (*ScreenSaver) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"org.freedesktop.ScreenSaver",
		"/org/freedesktop/ScreenSaver",
		"org.freedesktop.ScreenSaver",
	}
}

func NewScreenSaver() *ScreenSaver {
	s := &ScreenSaver{inhibitors: make(map[uint32]inhibitor)}
	s.xu, _ = xgbutil.NewConn()
	screensaver.Init(s.xu.Conn())
	screensaver.QueryVersion(s.xu.Conn(), 1, 0)
	screensaver.SelectInput(s.xu.Conn(), xproto.Drawable(s.xu.RootWin()), screensaver.EventNotifyMask|screensaver.EventCycleMask)
	dpms.Init(s.xu.Conn())

	go s.loop()
	return s
}
func Start() {
	ssaver := NewScreenSaver()

	if err := dbus.InstallOnSession(ssaver); err != nil {
		return
	}
}
func Stop() {
}

func (ss *ScreenSaver) loop() {
	for {
		e, err := ss.xu.Conn().WaitForEvent()
		if err != nil {
			continue
		}
		switch ee := e.(type) {
		case screensaver.NotifyEvent:
			switch ee.State {
			case screensaver.StateCycle:
				if ss.CycleActive != nil {
					ss.CycleActive()
				}
			case screensaver.StateOn:
				if ss.IdleOn != nil {
					ss.IdleOn()
				}
			case screensaver.StateOff:
				if ss.IdleOff != nil {
					ss.IdleOff()
				}
			}
		}
	}
}
