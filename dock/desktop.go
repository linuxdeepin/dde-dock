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
	"gir/gio-2.0"
	"gir/glib-2.0"
	"regexp"
)

var (
	// gio support `Desktop Action`, not support `Shortcut Group`
	actionReg = regexp.MustCompile(`(?P<actionGroup>.*) Shortcut Group`)
)

type DesktopAppInfo struct {
	id string
	*gio.DesktopAppInfo
	*glib.KeyFile
	gioSupported bool
}

func (dai *DesktopAppInfo) init() *DesktopAppInfo {
	logger.Debugf("[init] %q", dai.id)
	if dai.DesktopAppInfo == nil {
		logger.Debugf("[init] %q failed", dai.id)
		return nil
	}

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

func NewDesktopAppInfo(desktopID string) *DesktopAppInfo {
	dai := &DesktopAppInfo{}
	dai.DesktopAppInfo = gio.NewDesktopAppInfo(desktopID)
	dai.id = desktopID
	return dai.init()
}

func NewDesktopAppInfoFromFilename(file string) *DesktopAppInfo {
	dai := &DesktopAppInfo{}
	dai.DesktopAppInfo = gio.NewDesktopAppInfoFromFilename(file)
	dai.id = file
	return dai.init()
}

func (dai *DesktopAppInfo) ListActions() []string {
	logger.Debugf("[ListActions] %q", dai.id)
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

func (dai *DesktopAppInfo) getGroupName(name string) string {
	if dai.gioSupported {
		return "Desktop Action " + name
	}
	return name + " Shortcut Group"
}

func (dai *DesktopAppInfo) GetActionName(actionGroup string) string {
	logger.Debugf("[GetActionName] %q", dai.id)
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

func (dai *DesktopAppInfo) LaunchAction(actionGroup string, ctx gio.AppLaunchContextLike) {
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

func (dai *DesktopAppInfo) Destroy() {
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
