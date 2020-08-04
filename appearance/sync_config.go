package appearance

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
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
		err = m.doSetFontSize(v.FontSize)
		if err != nil {
			logger.Warning("failed to set font size:", err)
		} else {
			m.FontSize.Set(v.FontSize)
		}
	}

	if m.GtkTheme.Get() != v.GTK {
		err = m.doSetGtkTheme(v.GTK)
		if err != nil {
			logger.Warning("failed to set gtk theme:", err)
		} else {
			m.GtkTheme.Set(v.GTK)
		}
	}

	if m.IconTheme.Get() != v.Icon {
		err = m.doSetIconTheme(v.Icon)
		if err != nil {
			logger.Warning("failed to set icon theme:", err)
		} else {
			m.IconTheme.Set(v.Icon)
		}
	}

	if m.CursorTheme.Get() != v.Cursor {
		err = m.doSetCursorTheme(v.Cursor)
		if err != nil {
			logger.Warning("failed to set cursor theme:", err)
		} else {
			m.CursorTheme.Set(v.Cursor)
		}
	}

	if m.StandardFont.Get() != v.FontStandard {
		err = m.doSetStandardFont(v.FontStandard)
		if err != nil {
			logger.Warning("failed to set standard font:", err)
		} else {
			m.StandardFont.Set(v.FontStandard)
		}
	}

	if m.MonospaceFont.Get() != v.FontMonospace {
		err = m.doSetMonospaceFont(v.FontMonospace)
		if err != nil {
			logger.Warning("failed to set monospace font:", err)
		} else {
			m.MonospaceFont.Set(v.FontMonospace)
		}
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
	v.Version = "2.0"
	v.GreeterBackground = sc.m.greeterBg
	slideShow := sc.m.WallpaperSlideShow.Get()

	cfgSlideshow, _ := doUnmarshalWallpaperSlideshow(slideShow) // slideShow是一个map 格式为： "HDMI-0&&1":"600" 分别是屏幕名称&&工作区编号和自动切换壁纸配置
	uploadSlideShow := make(mapMonitorWorkspaceWSPolicy)
	for k, value := range cfgSlideshow { // 将具体的屏幕名称(例如"HDMI-0"或"VGA-0")转换为　主屏幕或副屏幕(Primary或Subsidiary0/Subsidiary1等)
		keySlice := strings.Split(k, "&&")
		if len(keySlice) < 2 {
			continue
		}
		index, err := strconv.Atoi(keySlice[1])
		if err != nil {
			logger.Warning(err)
			return nil, err
		}
		if int32(index) < 1 {
			return nil, errors.New("invalid workspace index")
		}
		monitorName := sc.m.monitorMap[keySlice[0]]
		key := monitorName + "&&" + keySlice[1]
		uploadSlideShow[key] = value
	}

	cfgWallpaperURIs, err := doUnmarshalMonitorWorkspaceWallpaperURIs(sc.m.WallpaperURIs.Get())
	if err != nil {
		logger.Warning(err)
		return nil, err
	}
	uploadWallpaperURIs := make(mapMonitorWorkspaceWallpaperURIs)
	for key, value := range cfgWallpaperURIs { // 对需要上传的壁纸信息进行过滤,只有当对应的自动切换的配置为空时(即：未配置自动切换相关内容),将该壁纸信息上传
		if uploadSlideShow[key] == "" {
			uploadWallpaperURIs[key] = value
		}
	}
	v.WallpaperURIs = uploadWallpaperURIs
	v.SlideShowConfig = uploadSlideShow

	return &v, nil
}

func (sc *backgroundSyncConfig) Set(data []byte) error {
	var v backgroundSyncData
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}
	m := sc.m
	if m.greeterBg != v.GreeterBackground {
		err = m.doSetGreeterBackground(v.GreeterBackground)
		if err != nil {
			logger.Warning(err)
		}
	}

	reverseMonitorMap := m.reverseMonitorMap() // 主副屏幕对应的具体屏幕名称的map 格式为{"Primary": "HDMI-0"}
	slideShow := make(mapMonitorWorkspaceWSPolicy)
	for k, value := range v.SlideShowConfig {
		keySlice := strings.Split(k, "&&")
		if len(keySlice) < 2 {
			continue
		}
		monitorName := reverseMonitorMap[keySlice[0]] //将主屏幕或副屏幕(Primary或Subsidiary0/Subsidiary1等)转换为具体的屏幕名称(例如"HDMI-0"或"VGA-0")
		index, err := strconv.Atoi(keySlice[1])
		if err != nil {
			logger.Warning(err)
			return err
		}
		if int32(index) < 1 {
			return errors.New("invalid workspace index")
		}
		key := monitorName + "&&" + keySlice[1]
		slideShow[key] = value
	}

	wallpaperSlideshow, err := doMarshalWallpaperSlideshow(slideShow)
	if err != nil {
		logger.Warning(err)
		return err
	}
	m.WallpaperSlideShow.Set(wallpaperSlideshow)

	uris, err := doMarshalMonitorWorkspaceWallpaperURIs(v.WallpaperURIs)
	if err != nil {
		logger.Warning(err)
		return err
	}
	m.WallpaperURIs.Set(uris)

	workspaceCount, _ := m.wm.WorkspaceCount(0) // 当前工作区数量
	for key, value := range v.WallpaperURIs {
		keySlice := strings.Split(key, "&&")
		if len(keySlice) < 2 {
			continue
		}
		monitorName := reverseMonitorMap[keySlice[0]] // 将主屏幕或副屏幕(Primary或Subsidiary0/Subsidiary1等)转换为具体的屏幕名称(例如"HDMI-0"或"VGA-0")
		index, err := strconv.Atoi(keySlice[1])
		if err != nil {
			logger.Warning(err)
			return err
		}
		if monitorName == "" {
			continue
		}
		if int32(index) < 1 { // index由1开始，代表工作区的编号，小于1代表编号错误
			return errors.New("invalid workspace index")
		}
		if int32(index) > workspaceCount {
			continue
		}
		err = m.wm.SetWorkspaceBackgroundForMonitor(0, int32(index), monitorName, value)
		if err != nil {
			logger.Warning("failed to set WorkspaceBackgroundForMonitor:", err)
		}
	}
	return nil
}

// version: 2.0
type backgroundSyncData struct {
	Version           string                           `json:"version"`
	GreeterBackground string                           `json:"greeter_background"`
	SlideShowConfig   mapMonitorWorkspaceWSPolicy      `json:"slide_show_config"`
	WallpaperURIs     mapMonitorWorkspaceWallpaperURIs `json:"wallpaper_uris"`
}
