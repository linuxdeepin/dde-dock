package appearance

import (
	"encoding/json"

	"pkg.deepin.io/lib/strv"
)

const (
	backgroundDBusPath = dbusPath + "/Background"
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
	return &v, nil
}

func (sc *syncConfig) Set(data []byte) error {
	var v syncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	m := sc.m
	if m.FontSize.Get() != v.FontSize {
		m.FontSize.Set(v.FontSize)
	}

	if m.GtkTheme.Get() != v.GTK {
		m.GtkTheme.Set(v.GTK)
	}

	if m.IconTheme.Get() != v.Icon {
		m.IconTheme.Set(v.Icon)
	}

	if m.CursorTheme.Get() != v.Cursor {
		m.CursorTheme.Set(v.Cursor)
	}

	if m.StandardFont.Get() != v.FontStandard {
		m.StandardFont.Set(v.FontStandard)
	}

	if m.MonospaceFont.Get() != v.FontMonospace {
		m.MonospaceFont.Set(v.FontMonospace)
	}

	return nil
}

// version: 1.0
type syncData struct {
	Version       string  `json:"version"`
	FontSize      float64 `json:"font_size"`
	GTK           string  `json:"gtk"`
	Icon          string  `json:"icon"`
	Cursor        string  `json:"cursor"`
	FontStandard  string  `json:"font_standard"`
	FontMonospace string  `json:"font_monospace"`
}

type backgroundSyncConfig struct {
	m *Manager
}

func (sc *backgroundSyncConfig) Get() (interface{}, error) {
	var v backgroundSyncData
	v.Version = "1.0"
	v.GreeterBackground = sc.m.greeterBg
	v.SlideShow = sc.m.WallpaperSlideShow.Get()
	if v.SlideShow == "" {
		v.BackgroundURIs = sc.m.getBackgroundURIs()
	}
	return &v, nil
}

func (sc *backgroundSyncConfig) Set(data []byte) error {
	var v backgroundSyncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	m := sc.m
	m.WallpaperSlideShow.Set(v.SlideShow)

	if m.greeterBg != v.GreeterBackground {
		err = m.doSetGreeterBackground(v.GreeterBackground)
		if err != nil {
			logger.Warning(err)
		}
	}

	if v.SlideShow != "" {
		return nil
	}

	bgs := m.getBackgroundURIs()
	if strv.Strv(bgs).Equal(v.BackgroundURIs) {
		return nil
	}
	for i, uri := range v.BackgroundURIs {
		err := m.wm.SetWorkspaceBackground(0, int32(i+1), uri)
		if err != nil {
			logger.Warning("Failed to set workspace background:", i+1, uri)
		}
	}

	return nil
}

// version: 1.0
type backgroundSyncData struct {
	Version           string   `json:"version"`
	BackgroundURIs    []string `json:"background_uris"`
	GreeterBackground string   `json:"greeter_background"`
	SlideShow         string   `json:"slide_show"`
}
