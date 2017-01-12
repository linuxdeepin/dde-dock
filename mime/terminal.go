/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mime

import (
	"fmt"
	"gir/gio-2.0"
	"strings"
)

const (
	terminalSchema = "com.deepin.desktop.default-applications.terminal"
	gsKeyExec      = "exec"

	cateKeyTerminal  = "TerminalEmulator"
	execKeyXTerminal = "x-terminal-emulator"
)

var termBlackList = []string{
	"guake.desktop",
}

func resetTerminal() {
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	s.Reset(gsKeyExec)
}

func setDefaultTerminal(id string) error {
	s := gio.NewSettings(terminalSchema)
	defer s.Unref()

	for _, info := range getTerminalInfos() {
		if info.Id == id {
			s.SetString(gsKeyExec, strings.Split(info.Exec, " ")[0])
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

	return (isStrInList(id, termBlackList) == false)
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
