package dock

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/gio-2.0"
)

const (
	HideModeKey         string = "hide-mode"
	HideModeKeepShowing        = "keep-showing"
	HideModeKeepHidden         = "keep-hidden"
	HideModeAutoHide           = "auto-hide"

	DisplayModeKey           string = "display-mode"
	DisplayModeModernModeStr        = "legacy"
	DisplayModeLegacyModeStr        = "modern"
)

type Setting struct {
	core *gio.Settings

	HideModeChanged    func(mode string)
	DisplayModeChanged func(mode string)
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
		value := g.GetString(key)
		logger.Info(key, "changed to", value)
		s.HideModeChanged(value)
	})

	s.listenGSettingChange(DisplayModeKey, func(g *gio.Settings, key string) {
		value := g.GetString(key)
		logger.Info(key, "changed to", value)
		s.DisplayModeChanged(value)
	})
}

func (s *Setting) listenGSettingChange(key string, handler func(*gio.Settings, string)) {
	signalDetial := fmt.Sprintf("changed::%s", HideModeKey)
	logger.Debugf("connect to %s signal", signalDetial)
	s.core.Connect(signalDetial, handler)
}

func (s *Setting) GetHideMode() string {
	return s.core.GetString(HideModeKey)
}

func (s *Setting) SetHideMode(mode string) bool {
	ok := s.core.SetString(HideModeKey, mode)
	if ok {
		s.HideModeChanged(mode)
	}
	return ok
}

func (s *Setting) GetDisplayMode() string {
	return s.core.GetString(DisplayModeKey)
}

func (s *Setting) SetDisplayMode(mode string) bool {
	ok := s.core.SetString(DisplayModeKey, mode)
	if ok {
		s.DisplayModeChanged(mode)
	}
	return ok
}
