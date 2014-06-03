package dock

import (
	"dlib/gio-2.0"
	"fmt"
)

const (
	HideModeKey string = "hide-mode"
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
	if s.core != nil {
		logger.Info("connect to changed signal")
		signalDetial := fmt.Sprint("changed::%s", HideModeKey)
		s.core.Connect(signalDetial, func(g *gio.Settings, key string) {
			value := g.GetString(key)
			logger.Info("HideMode changed to", value)
			s.HideModeChanged(value)
		})
	}
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
