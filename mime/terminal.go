/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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
	"gir/gio-2.0"
	"strings"
)

const (
	terminalSchema = "com.deepin.desktop.default-applications.terminal"
	gsKeyExec      = "exec"
	gsKeyExecArg   = "exec-arg"

	cateKeyTerminal  = "TerminalEmulator"
	execKeyXTerminal = "x-terminal-emulator"
)

// ignore quake-style terminal emulator
var termBlackList = []string{
	"guake.desktop",
	"tilda.desktop",
	"org.kde.yakuake.desktop",
	"qterminal_drop.desktop",
	"Terminal.desktop",
}

func resetTerminal() {
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	s.Reset(gsKeyExec)
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
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	for _, info := range getTerminalInfos() {
		if info.Id == id {
			exec := strings.Split(info.Exec, " ")[0]
			s.SetString(gsKeyExec, exec)
			s.SetString(gsKeyExecArg, getExecArg(exec))
			return nil
		}
	}
	return fmt.Errorf("Invalid terminal id '%s'", id)
}

func getDefaultTerminal() (*AppInfo, error) {
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	exec := s.GetString(gsKeyExec)
	for _, info := range getTerminalInfos() {
		if exec == strings.Split(info.Exec, " ")[0] {
			return info, nil
		}
	}

	return nil, fmt.Errorf("Not found app id for '%s'", exec)
}

func getTerminalInfos() AppInfos {
	infos := gio.AppInfoGetAll()
	defer unrefAppInfos(infos)

	var list AppInfos
	for _, info := range infos {
		if !isTerminalApp(info.GetId()) {
			continue
		}

		var tmp = &AppInfo{
			Id:          info.GetId(),
			Name:        info.GetName(),
			DisplayName: info.GetDisplayName(),
			Description: info.GetDescription(),
			Exec:        info.GetCommandline(),
		}
		iconObj := info.GetIcon()
		if iconObj != nil {
			tmp.Icon = iconObj.ToString()
			iconObj.Unref()
		}
		list = append(list, tmp)
	}
	return list
}

func isTerminalApp(id string) bool {
	if isStrInList(id, termBlackList) {
		return false
	}

	ginfo := gio.NewDesktopAppInfo(id)
	defer ginfo.Unref()
	cates := ginfo.GetCategories()
	if !strings.Contains(cates, cateKeyTerminal) {
		return false
	}

	exec := ginfo.GetCommandline()
	if strings.Contains(exec, execKeyXTerminal) {
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

func unrefAppInfos(infos []*gio.AppInfo) {
	for _, info := range infos {
		info.Unref()
	}
}
