package dock

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

const (
	HideModeKey    string = "hide-mode"
	DisplayModeKey string = "display-mode"
)

const (
	HideModeKeepShowing int32 = iota
	HideModeKeepHidden
	HideModeAutoHide

	HideModeKeepShowingStr = "keep-showing"
	HideModeKeepHiddenStr  = "keep-hidden"
	HideModeAutoHideStr    = "auto-hide"
)

const (
	DisplayModeModernMode int32 = iota
	DisplayModeLegacyMode

	DisplayModeModernModeStr = "legacy"
	DisplayModeLegacyModeStr = "modern"
)

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
		value := int32(g.GetEnum(key))
		logger.Info(key, "changed to", g.GetString(key))
		s.HideModeChanged(value)
	})

	s.listenGSettingChange(DisplayModeKey, func(g *gio.Settings, key string) {
		value := int32(g.GetEnum(key))
		logger.Info(key, "changed to", g.GetString(key))
		s.DisplayModeChanged(value)
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

func (s *Setting) SetHideMode(mode int32) bool {
	ok := s.core.SetEnum(HideModeKey, int(mode))
	if ok {
		s.HideModeChanged(mode)
	}
	return ok
}

func (s *Setting) GetDisplayMode() int32 {
	return int32(s.core.GetEnum(DisplayModeKey))
}

func (s *Setting) SetDisplayMode(mode int32) bool {
	ok := s.core.SetEnum(DisplayModeKey, int(mode))
	if ok {
		s.DisplayModeChanged(mode)
	}
	return ok
}
