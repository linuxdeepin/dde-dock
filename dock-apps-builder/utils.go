package main

import (
	"dlib/gio-2.0"
	"path/filepath"
	"strings"
)

func guess_desktop_id(oldId string) string {
	allApp := gio.AppInfoGetAll()
	for _, app := range allApp {
		baseName := filepath.Base(gio.ToDesktopAppInfo(app).GetFilename())
		if oldId == strings.ToLower(baseName) {
			return baseName
		}
	}

	return ""
}

func getAppIcon(core *gio.DesktopAppInfo) string {
	gioIcon := core.GetIcon()
	if gioIcon == nil {
		return ""
	}

	LOGGER.Debug("GetIcon:", gioIcon.ToString())
	icon := get_theme_icon(gioIcon.ToString(), 48)
	if icon == "" {
		return ""
	}

	LOGGER.Debug("get_theme_icon:", icon)
	// the filepath.Ext return ".xxx"
	ext := filepath.Ext(icon)[1:]
	LOGGER.Debug("ext:", ext)
	if strings.EqualFold(ext, "xpm") {
		LOGGER.Debug("change xpm to data uri")
		return xpm_to_dataurl(icon)
	}

	return icon
}
