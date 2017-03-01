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
	"path/filepath"
	"pkg.deepin.io/lib/appinfo/desktopappinfo"
)

const desktopHashPrefix = "d:"

type AppInfo struct {
	*desktopappinfo.DesktopAppInfo
	innerId string
}

// dai != nil
func newAppInfo(dai *desktopappinfo.DesktopAppInfo) *AppInfo {
	ai := &AppInfo{DesktopAppInfo: dai}
	ai.genInnerId()
	return ai
}

func getDockedDesktopAppInfo(app string) *desktopappinfo.DesktopAppInfo {
	if app[0] != '/' || len(app) <= 3 {
		return desktopappinfo.NewDesktopAppInfo(app)
	}

	absPath := unzipDesktopPath(app)
	ai, err := desktopappinfo.NewDesktopAppInfoFromFile(absPath)
	if err != nil {
		logger.Warning(err)
		return nil
	}
	return ai
}

func NewDockedAppInfo(app string) *AppInfo {
	if app == "" {
		return nil
	}
	dai := getDockedDesktopAppInfo(app)
	if dai == nil {
		return nil
	}
	return newAppInfo(dai)
}

func NewAppInfo(id string) *AppInfo {
	if id == "" {
		return nil
	}
	dai := desktopappinfo.NewDesktopAppInfo(id)
	if dai == nil {
		// try scratch dir
		desktopFile := filepath.Join(scratchDir, addDesktopExt(id))
		logger.Debugf("scratch dir desktopFile : %q", desktopFile)
		var err error
		dai, err = desktopappinfo.NewDesktopAppInfoFromFile(desktopFile)
		if err != nil {
			logger.Debug("NewAppInfo scratchDir failed:", err)
			return nil
		}
	}
	return newAppInfo(dai)
}

func NewAppInfoFromFile(file string) *AppInfo {
	dai, _ := desktopappinfo.NewDesktopAppInfoFromFile(file)
	if dai == nil {
		return nil
	}

	if !dai.IsInstalled() {
		createdBy, _ := dai.GetString(desktopappinfo.MainSection, "X-Deepin-CreatedBy")
		if createdBy == launcherDest {
			appId, _ := dai.GetString(desktopappinfo.MainSection, "X-Deepin-AppID")
			dai1 := desktopappinfo.NewDesktopAppInfo(appId)
			if dai1 != nil {
				dai = dai1
			}
		}
	}
	return newAppInfo(dai)
}

func (ai *AppInfo) genInnerId() {
	cmdline := ai.GetCommandline()
	hasher := md5.New()
	hasher.Write([]byte(cmdline))
	ai.innerId = desktopHashPrefix + hex.EncodeToString(hasher.Sum(nil))
}

func (ai *AppInfo) String() string {
	if ai == nil {
		return "<nil>"
	}
	desktopFile := ai.GetFileName()
	icon := ai.GetIcon()
	id := ai.GetId()
	return fmt.Sprintf("<AppInfo id=%q hash=%q icon=%q desktop=%q>", id, ai.innerId, icon, desktopFile)
}
