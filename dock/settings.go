package dock

import (
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"fmt"
)

const (
	HideModeKey string = "hide-mode"

	HideModeKeepShowing = "keep-showing"
	HideModeKeepHidden  = "keep-hidden"
	HideModeAutoHide    = "auto-hide"
)

type Setting struct {
	core *gio.Settings

	HideModeChanged func(mode string)
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

	logger.Debug("connect to changed::", HideModeKey, "signal")
	signalDetial := fmt.Sprintf("changed::%s", HideModeKey)
	s.core.Connect(signalDetial, func(g *gio.Settings, key string) {
		value := g.GetString(key)
		logger.Info(key, "changed to", value)
		s.HideModeChanged(value)
	})
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
