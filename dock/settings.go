package dock

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/gio-2.0"
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

	s.listenGSettingChange(HideModeKey, func(g *gio.Settings, key string) {
		value := HideModeType(g.GetEnum(key))
		logger.Info(key, "changed to", key)
		s.HideModeChanged(int32(value))
	})

	s.listenGSettingChange(DisplayModeKey, func(g *gio.Settings, key string) {
		value := DisplayModeType(g.GetEnum(key))
		logger.Info(key, "changed to", value)
		s.DisplayModeChanged(int32(value))
	})
}

func (s *Setting) listenGSettingChange(key string, handler func(*gio.Settings, string)) {
	signalDetial := fmt.Sprintf("changed::%s", key)
	logger.Debugf("connect to %s signal", signalDetial)
	s.core.Connect(signalDetial, handler)
}

func (s *Setting) GetHideMode() int32 {
	return int32(s.core.GetEnum(HideModeKey))
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
	return int32(s.core.GetEnum(DisplayModeKey))
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
