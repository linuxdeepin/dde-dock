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
