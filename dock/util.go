/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"

	"gir/gio-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xprop"
	"pkg.deepin.io/dde/daemon/appinfo"
)

func isEntryNameValid(name string) bool {
	if !strings.HasPrefix(name, entryDestPrefix) {
		return false
	}
	return true
}

func getEntryId(name string) (string, bool) {
	a := strings.SplitN(name, entryDestPrefix, 2)
	if len(a) >= 1 {
		return a[len(a)-1], true
	}
	return "", false
}

func trimDesktop(desktopID string) string {
	desktopIDLen := len(desktopID)
	if desktopIDLen == 0 {
		return ""
	}

	if desktopIDLen > 8 {
		return strings.TrimSuffix(desktopID, ".desktop")
	}

	panic(fmt.Sprintf("%q is not a desktop id", desktopID))
}

func normalizeAppID(candidateID string) string {
	return appinfo.NormalizeAppID(candidateID)
}

var _DesktopAppIdReg = regexp.MustCompile(`(?:[^.]+\.)*(?P<desktopID>[^.]+)\.desktop`)

func getAppIDFromDesktopID(candidateID string) string {
	desktopID := guess_desktop_id(candidateID)
	logger.Debug(fmt.Sprintf("get desktop id: %q", desktopID))
	if desktopID == "" {
		return ""
	}

	appID := normalizeAppID(trimDesktop(desktopID))
	return appID
}

// the key is appID
// the value is desktopID
var _appIDCache map[string]string = make(map[string]string)

func guess_desktop_id(appId string) string {
	logger.Debug(fmt.Sprintf("guess_desktop_id for %q", appId))
	if desktopID, ok := _appIDCache[appId]; ok {
		logger.Debug(appId, "is in cache")
		return desktopID
	}

	desktopID := appId + ".desktop"
	allApp := gio.AppInfoGetAll()

	defer func() {
		for _, app := range allApp {
			app.Unref()
		}
	}()

	for _, app := range allApp {
		_appInfo := gio.ToDesktopAppInfo(app)

		if _appInfo == nil {
			continue
		}

		_desktopID := _appInfo.GetId()
		normalizedDesktopID := normalizeAppID(_desktopID)
		if strings.HasSuffix(normalizedDesktopID, desktopID) {
			_appIDCache[appId] = _desktopID
			return _desktopID
		}

		// TODO: this is not a silver bullet, fix it later.
		appIDs := _DesktopAppIdReg.FindStringSubmatch(normalizedDesktopID)
		if len(appIDs) == 2 && appIDs[1] == appId {
			_appIDCache[appId] = _desktopID
			return _desktopID
		}
	}

	return ""
}

func getAppIcon(core *gio.DesktopAppInfo) string {
	gioIcon := core.GetIcon()
	if gioIcon == nil {
		logger.Warning("get icon from appinfo failed")
		return ""
	}

	icon := gioIcon.ToString()
	logger.Debug("GetIcon:", icon)
	if icon == "" {
		logger.Warning("gioIcon to string failed")
		return ""
	}

	iconPath := get_theme_icon(icon, 48)
	if iconPath == "" {
		logger.Warning("get icon from theme failed")
		// return a empty string might be a better idea here.
		// However, gtk will get theme icon failed sometimes for unknown reason.
		// frontend must make a validity check for icon.
		iconPath = icon
	}

	logger.Debug("get_theme_icon:", icon)
	ext := filepath.Ext(iconPath)
	if ext == "" {
		logger.Info("get app icon:", icon)
		return icon
	}

	// strip the '.' before extension name,
	// filepath.Ext function will return ".xxx"
	ext = ext[1:]
	logger.Debug("ext:", ext)
	if strings.EqualFold(ext, "xpm") {
		logger.Info("transform xpm to data uri")
		return xpm_to_dataurl(iconPath)
	}

	logger.Debug("get app icon:", icon)
	return icon
}

func dataUriToFile(dataUri, path string) (string, error) {
	commaIndex := strings.Index(dataUri, ",")
	img, err := base64.StdEncoding.DecodeString(dataUri[commaIndex+1:])
	if err != nil {
		return path, err
	}

	return path, ioutil.WriteFile(path, img, 0744)
}

func getWmName(xu *xgbutil.XUtil, win xproto.Window) (string, error) {
	wmname, err := xprop.PropValStr(xprop.GetProperty(xu, win, "_NET_WM_NAME"))
	if err != nil {
		wmname, err = xprop.PropValStr(xprop.GetProperty(xu, win, "WM_NAME"))
		if err != nil {
			return "", err
		}
	}
	return wmname, nil
}
