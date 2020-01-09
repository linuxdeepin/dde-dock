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

package mime

import (
	"fmt"
	"strings"

	"pkg.deepin.io/gir/gio-2.0"
	"pkg.deepin.io/lib/appinfo/desktopappinfo"
	"pkg.deepin.io/lib/strv"
)

const (
	gsSchemaDefaultTerminal = "com.deepin.desktop.default-applications.terminal"
	gsKeyExec               = "exec"
	gsKeyExecArg            = "exec-arg"
	gsKeyAppId              = "app-id"

	categoryTerminalEmulator = "TerminalEmulator"
	execXTerminalEmulator    = "x-terminal-emulator"
	desktopExt               = ".desktop"
)

// ignore quake-style terminal emulator
var termBlackList = strv.Strv{
	"guake",
	"tilda",
	"org.kde.yakuake",
	"qterminal_drop",
	"Terminal",
}

func resetTerminal() {
	settings := gio.NewSettings(gsSchemaDefaultTerminal)

	settings.Reset(gsKeyExec)
	settings.Reset(gsKeyExecArg)
	settings.Reset(gsKeyAppId)

	settings.Unref()
}

// readonly
var execArgMap = map[string]string{
	"gnome-terminal": "-x",
	"mate-terminal":  "-x",
	"terminator":     "-x",
	"xfce4-terminal": "-x",

	//"deepin-terminal": "-e",
	//"xterm":  "-e",
	//"pterm":  "-e",
	//"uxterm": "-e",
	//"rxvt": "-e",
	//"urxvt": "-e",
	//"rxvt-unicode": "-e",
	//"konsole": "-e",
	//"roxterm": "-e",
	//"lxterminal": "-e",
	//"terminology": "-e",
	//"sakura": "-e",
	//"evilvte": "-e",
	//"qterminal": "-e",
	//"termit": "-e",
	//"vala-terminal": "-e",
}

func getExecArg(exec string) string {
	execArg := execArgMap[exec]
	if execArg != "" {
		return execArg
	}
	return "-e"
}

func setDefaultTerminal(id string) error {
	settings := gio.NewSettings(gsSchemaDefaultTerminal)
	defer settings.Unref()

	for _, info := range getTerminalInfos() {
		if info.Id == id {
			exec := strings.Split(info.Exec, " ")[0]
			settings.SetString(gsKeyExec, exec)
			settings.SetString(gsKeyExecArg, getExecArg(exec))

			id = strings.TrimSuffix(id, desktopExt)
			settings.SetString(gsKeyAppId, id)
			return nil
		}
	}
	return fmt.Errorf("invalid terminal id '%s'", id)
}

func getDefaultTerminal() (*AppInfo, error) {
	settings := gio.NewSettings(gsSchemaDefaultTerminal)
	appId := settings.GetString(gsKeyAppId)
	// add suffix .desktop
	if !strings.HasSuffix(appId, desktopExt) {
		appId = appId + desktopExt
	}
	settings.Unref()
	for _, info := range getTerminalInfos() {
		if info.Id == appId {
			return info, nil
		}
	}

	return nil, fmt.Errorf("not found app id for %q", appId)
}

func getTerminalInfos() AppInfos {
	appInfoList := desktopappinfo.GetAll(nil)

	var list AppInfos
	for _, appInfo := range appInfoList {
		if !isTerminalApp(appInfo) {
			continue
		}

		name := getAppName(appInfo)
		var tmp = &AppInfo{
			Id:          appInfo.GetId() + desktopExt,
			Name:        name,
			DisplayName: name,
			Description: appInfo.GetComment(),
			Exec:        appInfo.GetCommandline(),
			Icon:        appInfo.GetIcon(),
			fileName:    appInfo.GetFileName(),
		}
		list = append(list, tmp)
	}
	return list
}

func isTerminalApp(appInfo *desktopappinfo.DesktopAppInfo) bool {
	if termBlackList.Contains(appInfo.GetId()) {
		return false
	}

	categories := appInfo.GetCategories()
	if !strv.Strv(categories).Contains(categoryTerminalEmulator) {
		return false
	}

	exec := appInfo.GetCommandline()
	if strings.Contains(exec, execXTerminalEmulator) {
		return false
	}
	return true
}

func isStrInList(s string, list []string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}

	return false
}
