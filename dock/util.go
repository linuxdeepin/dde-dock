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
	"io"
	"os"
	"path/filepath"
	"pkg.deepin.io/dde/daemon/appinfo"
	dutils "pkg.deepin.io/lib/utils"
	"text/template"
)

func iconifyWindow(win xproto.Window) {
	logger.Debug("iconifyWindow", win)
	ewmh.ClientEvent(XU, win, "WM_CHANGE_STATE", icccm.StateIconic)
}

func contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
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
	return strings.Replace(appinfo.NormalizeAppID(candidateID), " ", "-", -1)
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

func getWindowGtkApplicationId(xu *xgbutil.XUtil, win xproto.Window) string {
	gtkAppId, _ := xprop.PropValStr(xprop.GetProperty(xu, win, "_GTK_APPLICATION_ID"))
	return gtkAppId
}

func getWmWindowRole(xu *xgbutil.XUtil, win xproto.Window) string {
	role, _ := xprop.PropValStr(xprop.GetProperty(xu, win, "WM_WINDOW_ROLE"))
	return role
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

func getProcessEnvVars(pid uint) (map[string]string, error) {
	envPath := fmt.Sprintf("/proc/%d/environ", pid)
	bytes, err := ioutil.ReadFile(envPath)
	if err != nil {
		return nil, err
	}
	content := string(bytes)
	lines := strings.Split(content, "\x00")
	vars := make(map[string]string, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars, nil
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

func createScratchDesktopFile(id, title, icon, cmd string) error {
	logger.Debugf("create scratch file for %q", id)
	err := os.MkdirAll(scratchDir, 0775)
	if err != nil {
		logger.Warning("create scratch directory failed:", err)
		return err
	}
	f, err := os.OpenFile(filepath.Join(scratchDir, id+".desktop"),
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0744)
	if err != nil {
		logger.Warning("Open file for write failed:", err)
		return err
	}

	defer f.Close()
	temp := template.Must(template.New("docked_item_temp").Parse(dockedItemTemplate))
	dockedItem := dockedItemInfo{title, icon, cmd}
	logger.Debugf("dockedItem: %#v", dockedItem)
	err = temp.Execute(f, dockedItem)
	if err != nil {
		return err
	}
	return nil
}

func removeScratchFiles(id string) {
	extList := []string{"desktop", "sh", "png"}
	for _, ext := range extList {
		file := filepath.Join(scratchDir, id+"."+ext)
		if dutils.IsFileExist(file) {
			logger.Debugf("remove scratch file %q", file)
			err := os.Remove(file)
			if err != nil {
				logger.Warning("remove scratch file %q failed:", file, err)
			}
		}
	}
}

func strSliceEqual(sa, sb []string) bool {
	if len(sa) != len(sb) {
		return false
	}
	for i, va := range sa {
		vb := sb[i]
		if va != vb {
			return false
		}
	}
	return true
}

func uniqStrSlice(slice []string) []string {
	newSlice := make([]string, 0)
	for _, e := range slice {
		if !isStrInSlice(e, newSlice) {
			newSlice = append(newSlice, e)
		}
	}
	return newSlice
}

func isStrInSlice(v string, slice []string) bool {
	for _, e := range slice {
		if e == v {
			return true
		}
	}
	return false
}

func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
