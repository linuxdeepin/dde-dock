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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"gir/gio-2.0"
	"gir/glib-2.0"
	"path/filepath"
	"pkg.deepin.io/dde/daemon/appinfo"
	"regexp"
	"strings"
)

const desktopHashPrefix = "d:"

var (
	// gio support `Desktop Action`, not support `Shortcut Group`
	actionReg = regexp.MustCompile(`(?P<actionGroup>.*) Shortcut Group`)
)

type AppInfo struct {
	id        string
	desktopId string
	innerId   string
	*gio.DesktopAppInfo
	*glib.KeyFile
	gioSupported bool
}

func (dai *AppInfo) init() *AppInfo {
	logger.Debugf("[init] %q", dai.id)
	if dai.DesktopAppInfo == nil {
		logger.Debugf("[init] %q failed", dai.id)
		return nil
	}

	file := dai.DesktopAppInfo.GetFilename()
	dai.desktopId = getDesktopIdByFilePath(file)
	if dai.desktopId != "" {
		dai.id = dai.desktopId
	} else {
		dai.id = filepath.Base(file)
	}
	dai.id = appinfo.NormalizeAppID(trimDesktopExt(dai.id))

	dai.genInnerId()

	if len(dai.DesktopAppInfo.ListActions()) != 0 {
		dai.gioSupported = true
	}

	dai.KeyFile = glib.NewKeyFile()
	if ok, err := dai.LoadFromFile(dai.GetFilename(), glib.KeyFileFlagsNone); !ok {
		logger.Warning(err)
		dai.Destroy()
		return nil
	}
	return dai
}

func getDockedDesktopAppInfo(app string) *gio.DesktopAppInfo {
	if app[0] != '/' || len(app) <= 3 {
		return gio.NewDesktopAppInfo(addDesktopExt(app))
	}

	absPath := unzipDesktopPath(app)
	return gio.NewDesktopAppInfoFromFilename(absPath)
}

func NewDockedAppInfo(app string) *AppInfo {
	if app == "" {
		return nil
	}
	dai := &AppInfo{}
	desktopAppInfo := getDockedDesktopAppInfo(app)
	if desktopAppInfo == nil {
		return nil
	}
	dai.DesktopAppInfo = desktopAppInfo
	return dai.init()
}

func NewAppInfo(id string) *AppInfo {
	if id == "" {
		return nil
	}
	dai := &AppInfo{}
	desktopAppInfo := gio.NewDesktopAppInfo(addDesktopExt(id))
	if desktopAppInfo == nil {
		logger.Debug("NewAppInfo add .desktop failed")
		// try scratch dir
		desktopFile := filepath.Join(scratchDir, addDesktopExt(id))
		logger.Debugf("scratch dir desktopFile : %q", desktopFile)
		desktopAppInfo = gio.NewDesktopAppInfoFromFilename(desktopFile)
		if desktopAppInfo == nil {
			logger.Debug("NewAppInfo scratchDir failed")
			logger.Warningf("NewAppInfo failed: id %q", id)
			return nil
		}
	}

	dai.DesktopAppInfo = desktopAppInfo
	return dai.init()
}

func NewAppInfoFromFile(file string) *AppInfo {
	ai := &AppInfo{}
	ai.DesktopAppInfo = gio.NewDesktopAppInfoFromFilename(file)
	return ai.init()
}

func (ai *AppInfo) genInnerId() {
	cmdline := ai.DesktopAppInfo.GetCommandline()
	hasher := md5.New()
	hasher.Write([]byte(cmdline))
	ai.innerId = desktopHashPrefix + hex.EncodeToString(hasher.Sum(nil))
}

func (ai *AppInfo) String() string {
	if ai == nil {
		return "<nil>"
	}
	desktopFile := ai.GetFilePath()
	gioIcon := ai.DesktopAppInfo.GetIcon()
	var icon string
	if gioIcon == nil {
		logger.Warning("get icon from appinfo failed")
		icon = ""
	} else {
		icon = gioIcon.ToString()
	}
	return fmt.Sprintf("<AppInfo id=%q hash=%q icon=%q desktop=%q>", ai.id, ai.innerId, icon, desktopFile)
}

func (ai *AppInfo) GetId() string {
	return ai.id
}

func (ai *AppInfo) GetDesktopId() string {
	return ai.desktopId
}

func (ai *AppInfo) GetFilePath() string {
	if ai.DesktopAppInfo == nil {
		return ""
	}
	return ai.DesktopAppInfo.GetFilename()
}

func (dai *AppInfo) ListActions() []string {
	logger.Debugf("[ListActions] %q", dai.id)
	if dai.DesktopAppInfo == nil {
		return nil
	}

	if dai.gioSupported {
		logger.Debug("ListActions gio support")
		return dai.DesktopAppInfo.ListActions()
	}

	logger.Debug("ListActions gio not support")
	actions := make([]string, 0)
	_, groups := dai.GetGroups()
	for _, groupName := range groups {
		if tmp := actionReg.FindStringSubmatch(groupName); len(tmp) > 0 {
			actions = append(actions, tmp[1])
		}
	}

	return actions
}

func (dai *AppInfo) getGroupName(name string) string {
	if dai.gioSupported {
		return "Desktop Action " + name
	}
	return name + " Shortcut Group"
}

func (dai *AppInfo) GetActionName(actionGroup string) string {
	logger.Debugf("[GetActionName] %q", dai.id)
	if dai.DesktopAppInfo == nil {
		return ""
	}

	if dai.gioSupported {
		logger.Debug("GetActionName gio support")
		return dai.DesktopAppInfo.GetActionName(actionGroup)
	}

	logger.Debug("GetActionName gio not support")
	langs := GetLanguageNames()
	var str string
	groupName := dai.getGroupName(actionGroup)
	for _, lang := range langs {
		str, _ = dai.KeyFile.GetLocaleString(groupName, glib.KeyFileDesktopKeyName, lang)
		if str != "" {
			return str
		}
	}

	str, _ = dai.KeyFile.GetString(groupName, glib.KeyFileDesktopKeyName)
	return str
}

func (dai *AppInfo) GetIcon() string {
	if dai.DesktopAppInfo == nil {
		return ""
	}

	gioIcon := dai.DesktopAppInfo.GetIcon()
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

	iconPath := getIconFilePath(icon)
	logger.Debugf("icon file path %q", iconPath)
	if iconPath == "" {
		logger.Warningf("get icon %q file path failed", icon)
		// return a empty string might be a better idea here.
		// However, gtk will get theme icon failed sometimes for unknown reason.
		// frontend must make a validity check for icon.
		iconPath = icon
	}

	ext := filepath.Ext(iconPath)
	logger.Debug("ext:", ext)
	if ext == "" {
		logger.Info("get app icon:", icon)
		return icon
	}

	if strings.EqualFold(ext, ".xpm") {
		logger.Info("transform xpm to data uri")
		return xpm_to_dataurl(iconPath)
	}

	logger.Debug("get app icon:", icon)
	return icon
}

func (ai *AppInfo) GetDisplayName() string {
	if ai.DesktopAppInfo == nil {
		return ai.id
	}
	return ai.DesktopAppInfo.GetDisplayName()
}

func (dai *AppInfo) LaunchAction(actionGroup string, ctx gio.AppLaunchContextLike) {
	if dai.DesktopAppInfo == nil {
		return
	}

	logger.Debugf("[LaunchAction] %q action: %q", dai.id, actionGroup)
	if dai.gioSupported {
		logger.Info("LaunchAction gio support")
		dai.DesktopAppInfo.LaunchAction(actionGroup, ctx)
		return
	}

	logger.Debug("LaunchAction gio not support")
	exec, _ := dai.KeyFile.GetString(dai.getGroupName(actionGroup), glib.KeyFileDesktopKeyExec)
	logger.Infof("exec: %q", exec)
	cmdAppInfo, err := gio.AppInfoCreateFromCommandline(
		exec,
		"",
		gio.AppInfoCreateFlagsNone,
	)
	if err != nil {
		logger.Warning("Launch App Falied: ", err)
		return
	}

	defer cmdAppInfo.Unref()
	_, err = cmdAppInfo.Launch(nil, ctx)
	if err != nil {
		logger.Warning("Launch App Failed: ", err)
	}
}

func (dai *AppInfo) Destroy() {
	logger.Debugf("[Destroy] %q", dai.id)
	if dai.DesktopAppInfo != nil {
		dai.DesktopAppInfo.Unref()
		dai.DesktopAppInfo = nil
	}
	if dai.KeyFile != nil {
		dai.KeyFile.Free()
		dai.KeyFile = nil
	}
}
