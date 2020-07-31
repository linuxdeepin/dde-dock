/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package dock

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"pkg.deepin.io/lib/appinfo/desktopappinfo"
)

const desktopHashPrefix = "d:"

type AppInfo struct {
	*desktopappinfo.DesktopAppInfo
	identifyMethod string
	innerId        string
	name           string
}

func newAppInfo(dai *desktopappinfo.DesktopAppInfo) *AppInfo {
	if dai == nil {
		return nil
	}
	ai := &AppInfo{DesktopAppInfo: dai}
	xDeepinVendor, _ := dai.GetString(desktopappinfo.MainSection, "X-Deepin-Vendor")
	if xDeepinVendor == "deepin" {
		ai.name = dai.GetGenericName()
		if ai.name == "" {
			ai.name = dai.GetName()
		}
	} else {
		ai.name = dai.GetName()
	}
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
	return newAppInfo(getDockedDesktopAppInfo(app))
}

func NewAppInfo(id string) *AppInfo {
	if id == "" {
		return nil
	}
	return newAppInfo(desktopappinfo.NewDesktopAppInfo(id))
}

func NewAppInfoFromFile(file string) *AppInfo {
	if file == "" {
		return nil
	}
	dai, _ := desktopappinfo.NewDesktopAppInfoFromFile(file)
	if dai == nil {
		return nil
	}

	if !dai.IsInstalled() {
		appId, _ := dai.GetString(desktopappinfo.MainSection, "X-Deepin-AppID")
		if appId != "" {
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
	_, err := hasher.Write([]byte(cmdline))
	if err != nil {
		logger.Warning("Write error:",err)
	}
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
