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
	"errors"
	"io/ioutil"
	"strings"

	"pkg.deepin.io/dde/daemon/grub_common"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName = "com.deepin.daemon.Grub2"
	dbusPath        = "/com/deepin/daemon/Grub2"
	dbusInterface   = "com.deepin.daemon.Grub2"

	polikitActionIdCommon               = "com.deepin.daemon.Grub2"
	polikitActionIdPrepareGfxmodeDetect = "com.deepin.daemon.grub2.prepare-gfxmode-detect"

	timeoutMax = 10
)

func (*Grub2) GetInterfaceName() string {
	return dbusInterface
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

func (g *Grub2) GetAvailableGfxmodes(sender dbus.Sender) ([]string, *dbus.Error) {
	g.service.DelayAutoQuit()
	gfxmodes, err := g.getAvailableGfxmodes(sender)
	if err != nil {
		logger.Warning(err)
		return nil, dbusutil.ToError(err)
	}
	gfxmodes.SortDesc()
	result := make([]string, len(gfxmodes))
	for idx, m := range gfxmodes {
		result[idx] = m.String()
	}
	return result, nil
}

func (g *Grub2) SetDefaultEntry(sender dbus.Sender, entry string) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender, polikitActionIdCommon)
	if err != nil {
		return dbusutil.ToError(err)
	}

	idx := g.defaultEntryStr2Idx(entry)
	if idx == -1 {
		return dbusutil.ToError(errors.New("invalid entry"))
	}

	g.PropsMu.Lock()
	if g.setPropDefaultEntry(entry) {
		g.addModifyTask(getModifyTaskDefaultEntry(idx))
	}
	g.PropsMu.Unlock()
	return nil
}

var errInGfxmodeDetect = errors.New("in gfxmode detection mode")

func (g *Grub2) SetEnableTheme(sender dbus.Sender, enabled bool) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender, polikitActionIdCommon)
	if err != nil {
		return dbusutil.ToError(err)
	}

	lang, err := g.getSenderLang(sender)
	if err != nil {
		logger.Warning("failed to get sender lang:", err)
	}

	g.PropsMu.Lock()

	if g.setPropEnableTheme(enabled) {
		var themeFile string
		if enabled {
			if g.gfxmodeDetectState == gfxmodeDetectStateNone {
				themeFile = defaultGrubTheme
			} else {
				themeFile = fallbackGrubTheme
			}
		}
		g.setPropThemeFile(themeFile)
		task := getModifyTaskEnableTheme(enabled, lang, g.gfxmodeDetectState)
		logger.Debugf("SetEnableTheme task: %#v", task)
		g.addModifyTask(task)
	}

	g.PropsMu.Unlock()
	return nil
}

func (g *Grub2) SetGfxmode(sender dbus.Sender, gfxmode string) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender, polikitActionIdCommon)
	if err != nil {
		return dbusutil.ToError(err)
	}

	err = checkGfxmode(gfxmode)
	if err != nil {
		return dbusutil.ToError(err)
	}

	lang, err := g.getSenderLang(sender)
	if err != nil {
		logger.Warning("failed to get sender lang:", err)
	}

	g.PropsMu.Lock()
	if g.gfxmodeDetectState == gfxmodeDetectStateDetecting {
		g.PropsMu.Unlock()
		return dbusutil.ToError(errInGfxmodeDetect)
	}

	if g.setPropGfxmode(gfxmode) {
		g.addModifyTask(getModifyTaskGfxmode(gfxmode, lang))
	}
	g.PropsMu.Unlock()
	return nil
}

func (g *Grub2) SetTimeout(sender dbus.Sender, timeout uint32) *dbus.Error {
	g.service.DelayAutoQuit()

	/*
		err := g.checkAuth(sender, polikitActionIdCommon)
		if err != nil {
			return dbusutil.ToError(err)
		}
	*/

	if timeout > timeoutMax {
		return dbusutil.ToError(errors.New("exceeded the maximum value"))
	}

	g.PropsMu.Lock()
	if g.setPropTimeout(timeout) {
		g.addModifyTask(getModifyTaskTimeout(timeout))
	}
	g.PropsMu.Unlock()
	return nil
}

// Reset reset all configuration.
func (g *Grub2) Reset(sender dbus.Sender) *dbus.Error {
	g.service.DelayAutoQuit()

	const defaultEnableTheme = true

	err := g.checkAuth(sender, polikitActionIdCommon)
	if err != nil {
		return dbusutil.ToError(err)
	}

	lang, err := g.getSenderLang(sender)
	if err != nil {
		logger.Warning("failed to get sender lang:", err)
	}

	var modifyTasks []modifyTask

	g.PropsMu.Lock()
	if g.setPropTimeout(defaultGrubTimeoutInt) {
		modifyTasks = append(modifyTasks, getModifyTaskTimeout(defaultGrubTimeoutInt))
	}

	if g.setPropEnableTheme(defaultEnableTheme) {
		modifyTasks = append(modifyTasks,
			getModifyTaskEnableTheme(defaultEnableTheme, lang, g.gfxmodeDetectState))
	}

	g.setPropThemeFile(defaultGrubTheme)

	cfgDefaultEntry, _ := g.defaultEntryIdx2Str(defaultGrubDefaultInt)
	if g.setPropDefaultEntry(cfgDefaultEntry) {
		modifyTasks = append(modifyTasks, getModifyTaskDefaultEntry(defaultGrubDefaultInt))
	}
	g.PropsMu.Unlock()

	if len(modifyTasks) > 0 {
		compoundModifyFunc := func(params map[string]string) {
			for _, task := range modifyTasks {
				task.paramsModifyFunc(params)
			}
		}
		g.addModifyTask(modifyTask{
			paramsModifyFunc: compoundModifyFunc,
			adjustTheme:      true,
		})
	}

	return nil
}

func (g *Grub2) PrepareGfxmodeDetect(sender dbus.Sender) *dbus.Error {
	g.service.DelayAutoQuit()

	err := g.checkAuth(sender, polikitActionIdPrepareGfxmodeDetect)
	if err != nil {
		return dbusutil.ToError(err)
	}

	gfxmodes, err := g.getGfxmodesFromXRandr(sender)
	if err != nil {
		logger.Debug("failed to get gfxmodes from XRandr:", err)
	}

	gfxmodes.SortDesc()
	gfxmodesStr := joinGfxmodesForDetect(gfxmodes)

	g.PropsMu.RLock()
	gfxmodeDetectState := g.gfxmodeDetectState
	g.PropsMu.RUnlock()

	params, err := grub_common.LoadGrubParams()
	if err != nil {
		logger.Warning("failed to load grub params:", err)
		return dbusutil.ToError(err)
	}

	if gfxmodeDetectState == gfxmodeDetectStateDetecting {
		return dbusutil.ToError(errors.New("already in detection mode"))
	} else if gfxmodeDetectState == gfxmodeDetectStateFailed {
		cur, _, err := grub_common.GetBootArgDeepinGfxmode()
		if err == nil {
			if params[grubGfxmode] == gfxmodesStr ||
				(len(gfxmodes) > 0 && gfxmodes[0] == cur) {
				g.finishGfxmodeDetect(params)
				return nil
			}
		}
	}

	themeFile := getTheme(params)
	g.PropsMu.Lock()
	g.gfxmodeDetectState = gfxmodeDetectStateDetecting
	if themeFile != "" {
		g.setPropThemeFile(fallbackGrubTheme)
		g.setPropEnableTheme(true)
	} else {
		g.setPropThemeFile("")
		g.setPropEnableTheme(false)
	}
	g.setPropGfxmode(gfxmodesStr)
	g.PropsMu.Unlock()
	g.theme.emitSignalBackgroundChanged()

	g.addModifyTask(getModifyTaskPrepareGfxmodeDetect(gfxmodesStr))

	err = ioutil.WriteFile(grub_common.GfxmodeDetectReadyPath, nil, 0644)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}
