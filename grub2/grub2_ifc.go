/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	"fmt"
	"strings"

	"pkg.deepin.io/lib/dbus"
)

const (
	DBusDest      = "com.deepin.daemon.Grub2"
	DBusObjPath   = "/com/deepin/daemon/Grub2"
	DBusInterface = "com.deepin.daemon.Grub2"

	timeoutMax = 10
)

// GetDBusInfo implements interface of dbus.DBusObject.
func (_ *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DBusDest,
		ObjectPath: DBusObjPath,
		Interface:  DBusInterface,
	}
}

// GetSimpleEntryTitles return entry titles only in level one and will
// filter out some useless entries such as sub-menus and "memtest86+".
func (grub *Grub2) GetSimpleEntryTitles() ([]string, error) {
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
		err := fmt.Errorf("there is no menu entry in %s", grubScriptFile)
		return entryTitles, err
	}
	return entryTitles, nil
}

func (grub *Grub2) GetAvailableResolutions() (modesJSON string, err error) {
	// TODO
	type mode struct {
		Text, Value string
	}
	var modes []mode
	modes = append(modes, mode{Text: "Auto", Value: "auto"})
	modes = append(modes, mode{Text: "1024x768", Value: "1024x768"})
	modes = append(modes, mode{Text: "800x600", Value: "800x600"})
	data, err := json.Marshal(modes)
	modesJSON = string(data)
	return
}

func (g *Grub2) SetDefaultEntry(dbusMsg dbus.DMessage, v string) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	g.setPropMu.Lock()
	defer g.setPropMu.Unlock()

	if g.DefaultEntry == v {
		return
	}

	idx := g.defaultEntryStr2Idx(v)
	if idx == -1 {
		return errors.New("invalid entry")
	}

	g.DefaultEntry = v
	dbus.NotifyChange(g, propNameDefaultEntry)
	g.modifyFuncChan <- getModifyFuncDefaultEntry(idx)
	return
}

func (g *Grub2) SetEnableTheme(dbusMsg dbus.DMessage, v bool) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	g.setPropMu.Lock()
	defer g.setPropMu.Unlock()

	if g.EnableTheme == v {
		return
	}
	g.EnableTheme = v
	dbus.NotifyChange(g, propNameEnableTheme)
	g.modifyFuncChan <- getModifyFuncEnableTheme(v)
	return
}

func (g *Grub2) SetResolution(dbusMsg dbus.DMessage, v string) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	g.setPropMu.Lock()
	defer g.setPropMu.Unlock()

	err = checkResolution(v)
	if err != nil {
		return
	}

	if g.Resolution == v {
		return
	}
	g.Resolution = v
	dbus.NotifyChange(g, propNameResolution)
	g.modifyFuncChan <- getModifyFuncResolution(v)
	return
}

func (g *Grub2) SetTimeout(dbusMsg dbus.DMessage, v uint32) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	g.setPropMu.Lock()
	defer g.setPropMu.Unlock()

	if v > timeoutMax {
		return errors.New("exceeded the maximum value 10")
	}

	if g.Timeout == v {
		return
	}
	g.Timeout = v
	dbus.NotifyChange(g, propNameTimeout)
	g.modifyFuncChan <- getModifyFuncTimeout(v)
	return
}

// Reset reset all configuration.
func (g *Grub2) Reset(dbusMsg dbus.DMessage) (err error) {
	const defaultEnableTheme = true

	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	g.setPropMu.Lock()
	defer g.setPropMu.Unlock()

	g.theme.reset()

	var modifyFuncs []ModifyFunc
	if g.Timeout != defaultGrubTimeoutInt {
		g.Timeout = defaultGrubTimeoutInt
		dbus.NotifyChange(g, propNameTimeout)

		modifyFuncs = append(modifyFuncs, getModifyFuncTimeout(g.Timeout))
	}
	if g.EnableTheme != defaultEnableTheme {
		g.EnableTheme = defaultEnableTheme
		dbus.NotifyChange(g, propNameEnableTheme)

		modifyFuncs = append(modifyFuncs, getModifyFuncEnableTheme(g.EnableTheme))
	}

	cfgDefaultEntry, _ := g.defaultEntryIdx2Str(defaultGrubDefaultInt)
	if g.DefaultEntry != cfgDefaultEntry {
		g.DefaultEntry = cfgDefaultEntry
		dbus.NotifyChange(g, propNameDefaultEntry)

		modifyFuncs = append(modifyFuncs, getModifyFuncDefaultEntry(defaultGrubDefaultInt))
	}

	if len(modifyFuncs) > 0 {
		compoundModifyFunc := func(params map[string]string) {
			for _, fn := range modifyFuncs {
				fn(params)
			}
		}
		g.modifyFuncChan <- compoundModifyFunc
	}

	return
}
