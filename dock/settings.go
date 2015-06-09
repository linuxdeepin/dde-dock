package dock

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"sync"
)

const (
	HideModeKey    string = "hide-mode"
	DisplayModeKey string = "display-mode"
	ClockTypeKey   string = "clock-type"
	DisplayDateKey string = "display-date"
	DisplayWeekKey string = "display-week"
)

type HideModeType int32

const (
	HideModeKeepShowing HideModeType = iota
	HideModeKeepHidden
	HideModeAutoHide
	HideModeSmartHide
)

func (t HideModeType) String() string {
	switch t {
	case HideModeKeepShowing:
		return "Keep showing mode"
	case HideModeKeepHidden:
		return "Keep hidden mode"
	case HideModeAutoHide:
		return "Auto hide mode"
	case HideModeSmartHide:
		return "Smart hide mode"
	default:
		return "Unknown mode"
	}
}

type DisplayModeType int32

const (
	DisplayModeModernMode DisplayModeType = iota
	DisplayModeEfficientMode
	DisplayModeClassicMode
)

func (t DisplayModeType) String() string {
	switch t {
	case DisplayModeModernMode:
		return "Fashion mode"
	case DisplayModeEfficientMode:
		return "Efficient mode"
	case DisplayModeClassicMode:
		return "Classic mode"
	default:
		return "Unknown mode"
	}
}

type ClockType int32

const (
	ClockTypeDigit ClockType = iota
	ClockTypeAnalog
)

func (self ClockType) String() string {
	switch self {
	case ClockTypeDigit:
		return "digit clock"
	case ClockTypeAnalog:
		return "analog clock"
	default:
		return "unknown type clock"
	}
}

type Setting struct {
	core *gio.Settings

	hideModeLock sync.RWMutex
	hideMode     HideModeType

	displayModeLock sync.RWMutex
	displayMode     DisplayModeType

	clockTypeLock sync.RWMutex
	clockType     ClockType

	displayDateLock sync.RWMutex
	displayDate     bool

	displayWeekLock sync.RWMutex
	displayWeek     bool

	HideModeChanged    func(mode int32)
	DisplayModeChanged func(mode int32)
	ClockTypeChanged   func(mode int32)
	DisplayDateChanged func(bool)
	DisplayWeekChanged func(bool)
}

func NewSetting() *Setting {
	s := &Setting{}
	if s.init() {
		return s
	}
	return nil
}

func (s *Setting) init() bool {
	s.core = gio.NewSettings(SchemaId)
	if s.core == nil {
		return false
	}

	s.listenSettingChange(HideModeKey, func(g *gio.Settings, key string) {
		s.hideModeLock.Lock()
		defer s.hideModeLock.Unlock()

		value := HideModeType(g.GetEnum(key))
		s.hideMode = value
		logger.Debug(key, "changed to", key)
		dbus.Emit(s, "HideModeChanged", int32(value))
	})

	s.listenSettingChange(DisplayModeKey, func(g *gio.Settings, key string) {
		s.displayModeLock.Lock()
		defer s.displayModeLock.Unlock()

		value := DisplayModeType(g.GetEnum(key))
		logger.Debug(key, "changed to", value)

		s.displayMode = value

		for _, rApp := range ENTRY_MANAGER.runtimeApps {
			rebuildXids := []xproto.Window{}
			for xid, _ := range rApp.xids {
				if _, err := xprop.PropValStr(
					xprop.GetProperty(
						XU,
						xid,
						"_DDE_DOCK_APP_ID",
					),
				); err != nil {
					continue
				}

				rebuildXids = append(rebuildXids, xid)
				rApp.detachXid(xid)
			}

			l := len(rebuildXids)
			if l == 0 {
				continue
			}

			if len(rApp.xids) == 0 {
				ENTRY_MANAGER.destroyRuntimeApp(rApp)
			}

			newApp := ENTRY_MANAGER.createRuntimeApp(rebuildXids[0])
			for i := 0; i < l; i++ {
				newApp.attachXid(rebuildXids[i])
			}

			activeXid, err := ewmh.ActiveWindowGet(XU)
			if err != nil {
				continue
			}

			for xid, _ := range newApp.xids {
				logger.Debugf("through new app xids")
				if activeXid == xid {
					logger.Debugf("0x%x(a), 0x%x(x)",
						activeXid, xid)
					newApp.setLeader(xid)
					newApp.updateState(xid)
					ewmh.ActiveWindowSet(XU, xid)
					break
				}
			}
		}

		dockProperty.updateDockHeight(value)
		dbus.Emit(s, "DisplayModeChanged", int32(value))
	})
	s.listenSettingChange(ClockTypeKey, func(*gio.Settings, string) {
		s.clockTypeLock.Lock()
		defer s.clockTypeLock.Unlock()
		s.clockType = ClockType(s.core.GetEnum(ClockTypeKey))
		dbus.Emit(s, "ClockTypeChanged", int32(s.clockType))
	})

	s.listenSettingChange(DisplayDateKey, func(*gio.Settings, string) {
		s.displayDateLock.Lock()
		defer s.displayDateLock.Unlock()
		s.displayDate = s.core.GetBoolean(DisplayDateKey)
		dbus.Emit(s, "DisplayDateChanged", s.displayDate)
	})
	s.listenSettingChange(DisplayWeekKey, func(*gio.Settings, string) {
		s.displayWeekLock.Lock()
		defer s.displayWeekLock.Unlock()
		s.displayWeek = s.core.GetBoolean(DisplayWeekKey)
		dbus.Emit(s, "DisplayWeekChanged", s.displayWeek)
	})

	// at least one read operation must be called after signal connected, otherwise,
	// the signal connection won't work from glib 2.43.
	// NB: https://github.com/GNOME/glib/commit/8ff5668a458344da22d30491e3ce726d861b3619
	s.displayMode = DisplayModeType(s.core.GetEnum(DisplayModeKey))
	s.hideMode = HideModeType(s.core.GetEnum(HideModeKey))
	if s.hideMode == HideModeAutoHide {
		s.hideMode = HideModeSmartHide
		s.core.SetEnum(HideModeKey, int32(HideModeSmartHide))
	}
	s.clockType = ClockType(s.core.GetEnum(ClockTypeKey))
	s.displayDate = s.core.GetBoolean(DisplayDateKey)
	s.displayWeek = s.core.GetBoolean(DisplayWeekKey)

	return true
}

func (s *Setting) listenSettingChange(key string, handler func(*gio.Settings, string)) {
	signalDetial := fmt.Sprintf("changed::%s", key)
	logger.Debugf("connect to %s signal", signalDetial)
	s.core.Connect(signalDetial, handler)
}

func (s *Setting) GetHideMode() int32 {
	return int32(s.hideMode)
}

func (s *Setting) SetHideMode(_mode int32) bool {
	mode := HideModeType(_mode)
	logger.Debug("[Setting.SetHideMode]:", mode)
	ok := s.core.SetEnum(HideModeKey, int32(mode))
	return ok
}

func (s *Setting) GetDisplayMode() int32 {
	return int32(s.displayMode)
}

func (s *Setting) SetDisplayMode(_mode int32) bool {
	mode := DisplayModeType(_mode)
	logger.Debug("[Setting.SetDisplayMode]:", mode)
	ok := s.core.SetEnum(DisplayModeKey, int32(mode))
	return ok
}

func (s *Setting) GetClockType() int32 {
	return int32(s.clockType)
}

func (s *Setting) SetClockType(_clockType int32) bool {
	clockType := ClockType(_clockType)
	logger.Debug("clock type changed to:", clockType)
	ok := s.core.SetEnum(ClockTypeKey, int32(clockType))
	return ok
}

func (s *Setting) GetDisplayDate() bool {
	return s.displayDate
}

func (s *Setting) SetDisplayDate(d bool) bool {
	return s.core.SetBoolean(DisplayDateKey, d)
}

func (s *Setting) GetDisplayWeek() bool {
	return s.displayWeek
}

func (s *Setting) SetDisplayWeek(d bool) bool {
	return s.core.SetBoolean(DisplayWeekKey, d)
}

func (s *Setting) destroy() {
	if s.core != nil {
		s.core.Unref()
		s.core = nil
	}
	dbus.UnInstallObject(s)
}
