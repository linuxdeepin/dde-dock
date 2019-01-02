package appearance

import (
	"encoding/json"

	"pkg.deepin.io/lib/strv"
)

type syncConfig struct {
	m *Manager
}

func (sc *syncConfig) Get() (interface{}, error) {
	var v syncData
	v.Version = "1.0"
	v.FontSize = sc.m.FontSize.Get()
	v.GTK = sc.m.GtkTheme.Get()
	v.Icon = sc.m.IconTheme.Get()
	v.Cursor = sc.m.CursorTheme.Get()
	v.FontStandard = sc.m.StandardFont.Get()
	v.FontMonospace = sc.m.MonospaceFont.Get()
	v.BackgroundURIs = sc.m.getBackgroundURIs()
	return v, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	m := sc.m
	if m.FontSize.Get() != v.FontSize {
		err = m.doSetFontSize(v.FontSize)
		if err != nil {
			logger.Warning(err)
		}
	}

	if m.GtkTheme.Get() != v.GTK {
		err = m.doSetGtkTheme(v.GTK)
		if err != nil {
			logger.Warning(err)
		}
	}

	if m.IconTheme.Get() != v.Icon {
		err = m.doSetIconTheme(v.Icon)
		if err != nil {
			logger.Warning(err)
		}
	}

	if m.CursorTheme.Get() != v.Cursor {
		err = m.doSetCursorTheme(v.Cursor)
		if err != nil {
			logger.Warning(err)
		}
	}

	if m.StandardFont.Get() != v.FontStandard {
		err = m.doSetStandardFont(v.FontStandard)
		if err != nil {
			logger.Warning(err)
		}
	}

	if m.MonospaceFont.Get() != v.FontMonospace {
		err = m.doSetMonospaceFont(v.FontMonospace)
		if err != nil {
			logger.Warning(err)
		}
	}

	bgs := m.getBackgroundURIs()
	if !strv.Strv(bgs).Equal(v.BackgroundURIs) {
		m.setting.SetStrv(gsKeyBackgroundURIs, v.BackgroundURIs)
	}
	return nil
}

// version: 1.0
type syncData struct {
	Version        string   `json:"version"`
	FontSize       float64  `json:"font_size"`
	GTK            string   `json:"gtk"`
	Icon           string   `json:"icon"`
	Cursor         string   `json:"cursor"`
	FontStandard   string   `json:"font_standard"`
	FontMonospace  string   `json:"font_monospace"`
	BackgroundURIs []string `json:"background_uris"`
}
