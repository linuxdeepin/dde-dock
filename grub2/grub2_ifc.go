/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package grub2

import (
	"encoding/json"
	"errors"
	"strings"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	DBusServiceName = "com.deepin.daemon.Grub2"
	DBusObjPath     = "/com/deepin/daemon/Grub2"
	DBusInterface   = "com.deepin.daemon.Grub2"

	timeoutMax = 10
)

func (*Grub2) GetDBusExportInfo() dbusutil.ExportInfo {
	return dbusutil.ExportInfo{
		Path:      DBusObjPath,
		Interface: DBusInterface,
	}
}

// GetSimpleEntryTitles return entry titles only in level one and will
// filter out some useless entries such as sub-menus and "memtest86+".
func (grub *Grub2) GetSimpleEntryTitles() ([]string, *dbus.Error) {
	grub.service.DelayAutoQuit()

	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.parentSubMenu == nil && entry.entryType == MENUENTRY {
			title := entry.getFullTitle()
			if !strings.Contains(title, "memtest86+") {
				entryTitles = append(entryTitles, title)
			}
		}
	}
	if len(entryTitles) == 0 {
		logger.Warningf("there is no menu entry in %q", grubScriptFile)
	}
	return entryTitles, nil
}

func (grub *Grub2) GetAvailableResolutions() (modeJSON string, err *dbus.Error) {
	grub.service.DelayAutoQuit()
	// TODO
	type mode struct {
		Text, Value string
	}
	var modes []mode
	modes = append(modes, mode{Text: "Auto", Value: "auto"})
	modes = append(modes, mode{Text: "1024x768", Value: "1024x768"})
	modes = append(modes, mode{Text: "800x600", Value: "800x600"})
	data, _ := json.Marshal(modes)
	return string(data), nil
}

func (g *Grub2) SetDefaultEntry(sender dbus.Sender, entry string) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	idx := g.defaultEntryStr2Idx(entry)
	if idx == -1 {
		return dbusutil.ToError(errors.New("invalid entry"))
	}

	g.PropsMu.Lock()
	if g.setPropDefaultEntry(entry) {
		g.modifyFuncChan <- getModifyFuncDefaultEntry(idx)
	}
	g.PropsMu.Unlock()
	return nil
}

func (g *Grub2) SetEnableTheme(sender dbus.Sender, enabled bool) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	g.PropsMu.Lock()
	if g.setPropEnableTheme(enabled) {
		g.modifyFuncChan <- getModifyFuncEnableTheme(enabled)
	}
	g.PropsMu.Unlock()
	return nil
}

func (g *Grub2) SetResolution(sender dbus.Sender, resolution string) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = checkResolution(resolution)
	if err != nil {
		return dbusutil.ToError(err)
	}

	g.PropsMu.Lock()
	if g.setPropResolution(resolution) {
		g.modifyFuncChan <- getModifyFuncResolution(resolution)
	}
	g.PropsMu.Unlock()
	return nil
}

func (g *Grub2) SetTimeout(sender dbus.Sender, timeout uint32) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	if timeout > timeoutMax {
		return dbusutil.ToError(errors.New("exceeded the maximum value"))
	}

	g.PropsMu.Lock()
	if g.setPropTimeout(timeout) {
		g.modifyFuncChan <- getModifyFuncTimeout(timeout)
	}
	g.PropsMu.Unlock()
	return nil
}

// Reset reset all configuration.
func (g *Grub2) Reset(sender dbus.Sender) *dbus.Error {
	g.service.DelayAutoQuit()

	const defaultEnableTheme = true

	err := g.checkAuth(sender)
	if err != nil {
		return dbusutil.ToError(err)
	}

	var modifyFuncs []ModifyFunc

	g.PropsMu.Lock()
	if g.setPropTimeout(defaultGrubTimeoutInt) {
		modifyFuncs = append(modifyFuncs, getModifyFuncTimeout(defaultGrubTimeoutInt))
	}

	if g.setPropEnableTheme(defaultEnableTheme) {
		modifyFuncs = append(modifyFuncs, getModifyFuncEnableTheme(defaultEnableTheme))
	}

	cfgDefaultEntry, _ := g.defaultEntryIdx2Str(defaultGrubDefaultInt)
	if g.setPropDefaultEntry(cfgDefaultEntry) {
		modifyFuncs = append(modifyFuncs, getModifyFuncDefaultEntry(defaultGrubDefaultInt))
	}
	g.PropsMu.Unlock()

	if len(modifyFuncs) > 0 {
		compoundModifyFunc := func(params map[string]string) {
			for _, fn := range modifyFuncs {
				fn(params)
			}
		}
		g.modifyFuncChan <- compoundModifyFunc
	}

	err = g.theme.reset()
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}
