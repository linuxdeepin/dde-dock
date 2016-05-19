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
	"regexp"
	"strings"

	"bytes"
	"gir/gio-2.0"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/xgraphics"
	"github.com/BurntSushi/xgbutil/xprop"
	"os"
	"path/filepath"
	"pkg.deepin.io/dde/daemon/appinfo"
)

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
	logger.Debugf("guess_desktop_id %q", appId)
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

func dataUriToFile(dataUri, path string) (string, error) {
	commaIndex := strings.Index(dataUri, ",")
	img, err := base64.StdEncoding.DecodeString(dataUri[commaIndex+1:])
	if err != nil {
		return path, err
	}

	return path, ioutil.WriteFile(path, img, 0744)
}

func getWmName(xu *xgbutil.XUtil, win xproto.Window) string {
	// get _NET_WM_NAME
	name, err := ewmh.WmNameGet(xu, win)
	if err != nil || name == "" {
		// get WM_NAME
		name, _ = icccm.WmNameGet(xu, win)
	}
	return name
}

func getWmPid(xu *xgbutil.XUtil, win xproto.Window) uint {
	pid, _ := ewmh.WmPidGet(xu, win)
	return pid
}

func getWmCommand(xu *xgbutil.XUtil, win xproto.Window) ([]string, error) {
	command, err := xprop.PropValStrs(xprop.GetProperty(xu, win, "WM_COMMAND"))
	return command, err
}

func getProcessCmdline(pid uint) ([]string, error) {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	bytes, err := ioutil.ReadFile(cmdlinePath)
	if err != nil {
		return nil, err
	}
	content := string(bytes)
	parts := strings.Split(content, "\x00")
	length := len(parts)
	if length >= 2 && parts[length-1] == "" {
		return parts[:length-1], nil
	}
	return parts, nil
}

func getProcessCwd(pid uint) (string, error) {
	cwdPath := fmt.Sprintf("/proc/%d/cwd", pid)
	cwd, err := os.Readlink(cwdPath)
	return cwd, err
}

func getProcessExe(pid uint) (string, error) {
	exePath := fmt.Sprintf("/proc/%d/exe", pid)
	exe, err := filepath.EvalSymlinks(exePath)
	return exe, err
}

func getIconFromWindow(xu *xgbutil.XUtil, win xproto.Window) string {
	icon, err := xgraphics.FindIcon(xu, win, 48, 48)
	// FIXME: gets empty icon for minecraft
	if err == nil {
		buf := bytes.NewBuffer(nil)
		icon.WritePng(buf)
		return "data:image/png;base64," + base64.StdEncoding.EncodeToString(buf.Bytes())
	}

	logger.Debug("get icon from X failed:", err)
	logger.Debug("get icon name from _NET_WM_ICON_NAME")
	name, _ := ewmh.WmIconNameGet(XU, win)
	return name
}
