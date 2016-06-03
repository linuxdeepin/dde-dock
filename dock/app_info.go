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

func NewAppInfo(id string) *AppInfo {
	if id == "" {
		return nil
	}
	dai := &AppInfo{}
	desktopId := id + ".desktop"
	desktopAppInfo := gio.NewDesktopAppInfo(desktopId)
	if desktopAppInfo == nil {
		// try guess_desktop_id
		logger.Debug("NewAppInfo add .desktop failed")
		desktopId = guess_desktop_id(id)
		desktopAppInfo = gio.NewDesktopAppInfo(desktopId)
		if desktopAppInfo == nil {
			logger.Debug("NewAppInfo guess_desktop_id failed")
			// try scratch dir
			desktopFile := filepath.Join(scratchDir, id+".desktop")
			logger.Debugf("scratch dir desktopFile : %q", desktopFile)
			desktopAppInfo = gio.NewDesktopAppInfoFromFilename(desktopFile)
			if desktopAppInfo == nil {
				logger.Debug("NewAppInfo scratchDir failed")
				logger.Warningf("NewAppInfo failed: id %q", id)
				return nil
			} else {
				desktopId = desktopAppInfo.GetId()
			}
		}
	}

	dai.DesktopAppInfo = desktopAppInfo
	dai.desktopId = desktopId
	dai.id = trimDesktop(normalizeAppID(desktopId))
	return dai.init()
}

func NewAppInfoFromFile(file string) *AppInfo {
	ai := &AppInfo{}
	ai.DesktopAppInfo = gio.NewDesktopAppInfoFromFilename(file)
	id := filepath.Base(file)

	// in kde4 dir
	dir := filepath.Dir(file)
	basedir := filepath.Base(dir)
	if basedir == "kde4" {
		id = "kde4-" + id
	}
	ai.desktopId = id
	ai.id = trimDesktop(normalizeAppID(id))
	return ai.init()
}

func (ai *AppInfo) genInnerId() {
	cmdline := ai.DesktopAppInfo.GetCommandline()
	hasher := md5.New()
	hasher.Write([]byte(cmdline))
	ai.innerId = desktopHashPrefix + hex.EncodeToString(hasher.Sum(nil))
}

func (ai *AppInfo) String() string {
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

	logger.Debug("ext:", ext)
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
