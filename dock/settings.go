package dock

import (
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xprop"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"sync"
)

const (
	HideModeKey    string = "hide-mode"
	DisplayModeKey string = "display-mode"
)

type HideModeType int32

const (
	HideModeKeepShowing HideModeType = iota
	HideModeKeepHidden
	HideModeAutoHide
)

func (t HideModeType) String() string {
	switch t {
	case HideModeKeepShowing:
		return "Keep showing mode"
	case HideModeKeepHidden:
		return "Keep hidden mode"
	case HideModeAutoHide:
		return "Auto hide mode"
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

type Setting struct {
	core *gio.Settings

	hideModeLock sync.RWMutex
	hideMode     HideModeType

	displayModeLock sync.RWMutex
	displayMode     DisplayModeType

	HideModeChanged    func(mode int32)
	DisplayModeChanged func(mode int32)
}

func NewSetting() *Setting {
	s := &Setting{}
	s.init()
	return s
}

func (s *Setting) init() {
	s.core = gio.NewSettings(SchemaId)
	if s.core == nil {
		return
	}

	s.displayMode = DisplayModeType(s.core.GetEnum(DisplayModeKey))
	s.hideMode = HideModeType(s.core.GetEnum(HideModeKey))

	s.listenSettingChange(HideModeKey, func(g *gio.Settings, key string) {
		s.hideModeLock.Lock()
		defer s.hideModeLock.Unlock()

		value := HideModeType(g.GetEnum(key))
		s.hideMode = value
		logger.Info(key, "changed to", key)
		s.HideModeChanged(int32(value))
	})

	s.listenSettingChange(DisplayModeKey, func(g *gio.Settings, key string) {
		s.displayModeLock.Lock()
		defer s.displayModeLock.Unlock()

		value := DisplayModeType(g.GetEnum(key))
		logger.Info(key, "changed to", value)

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
				logger.Warningf("through new app xids")
				if activeXid == xid {
					logger.Warningf("0x%x(a), 0x%x(x)",
						activeXid, xid)
					newApp.setLeader(xid)
					newApp.updateState(xid)
					ewmh.ActiveWindowSet(XU, xid)
					break
				}
			}
		}
		s.DisplayModeChanged(int32(value))
	})
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
	ok := s.core.SetEnum(HideModeKey, int(mode))
	if ok {
		s.HideModeChanged(_mode)
	}
	return ok
}

func (s *Setting) GetDisplayMode() int32 {
	return int32(s.displayMode)
}

func (s *Setting) SetDisplayMode(_mode int32) bool {
	mode := DisplayModeType(_mode)
	logger.Info("[Setting.SetDisplayMode]:", mode)
	ok := s.core.SetEnum(DisplayModeKey, int(mode))
	if ok {
		s.DisplayModeChanged(_mode)
	}
	return ok
}
