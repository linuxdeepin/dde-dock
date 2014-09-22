/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package inputdevices

import (
	"dbus/com/deepin/api/greeterutils"
	"io/ioutil"
	"os/exec"
	"path"
	"pkg.linuxdeepin.com/lib/dbus/property"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"regexp"
	"strings"
)

const (
	KBD_KEY_REPEAT_ENABLE    = "repeat-enabled"
	KBD_KEY_REPEAT_INTERVAL  = "repeat-interval"
	KBD_KEY_DELAY            = "delay"
	KBD_KEY_LAYOUT           = "layout"
	KBD_KEY_LAYOUT_MODEL     = "layout-model"
	KBD_KEY_LAYOUT_OPTIONS   = "layout-options"
	KBD_KEY_USER_LAYOUT_LIST = "user-layout-list"
	KBD_CURSOR_BLINK_TIME    = "cursor-blink-time"
	KBD_KEY_CAPSLOCK_TOGGLE  = "capslock-toggle"
	KBD_DEFAULT_FILE         = "/etc/default/keyboard"

	CMD_SETXKB = "/usr/bin/setxkbmap"
)

var kbdSettings = gio.NewSettings("com.deepin.dde.keyboard")

type KeyboardManager struct {
	RepeatEnabled  *property.GSettingsBoolProperty `access:"readwrite"`
	CapslockToggle *property.GSettingsBoolProperty `access:"readwrite"`

	CursorBlink *property.GSettingsIntProperty `access:"readwrite"`

	RepeatInterval *property.GSettingsUintProperty `access:"readwrite"`
	RepeatDelay    *property.GSettingsUintProperty `access:"readwrite"`

	CurrentLayout *property.GSettingsStringProperty `access:"readwrite"`

	UserLayoutList *property.GSettingsStrvProperty
	UserOptionList *property.GSettingsStrvProperty

	layoutDescMap map[string]string
	greeterObj    *greeterutils.GreeterUtils

	listenFlag bool
}

var _kbdManager *KeyboardManager

func GetKeyboardManager() *KeyboardManager {
	if _kbdManager == nil {
		_kbdManager = newKeyboardManager()
	}

	return _kbdManager
}

func newKeyboardManager() *KeyboardManager {
	kbdManager := &KeyboardManager{}

	kbdManager.CurrentLayout = property.NewGSettingsStringProperty(
		kbdManager, "CurrentLayout",
		kbdSettings, KBD_KEY_LAYOUT)
	kbdManager.RepeatEnabled = property.NewGSettingsBoolProperty(
		kbdManager, "RepeatEnabled",
		kbdSettings, KBD_KEY_REPEAT_ENABLE)
	kbdManager.RepeatInterval = property.NewGSettingsUintProperty(
		kbdManager, "RepeatInterval",
		kbdSettings, KBD_KEY_REPEAT_INTERVAL)
	kbdManager.RepeatDelay = property.NewGSettingsUintProperty(
		kbdManager, "RepeatDelay",
		kbdSettings, KBD_KEY_DELAY)
	kbdManager.CursorBlink = property.NewGSettingsIntProperty(
		kbdManager, "CursorBlink",
		kbdSettings, KBD_CURSOR_BLINK_TIME)
	kbdManager.CapslockToggle = property.NewGSettingsBoolProperty(
		kbdManager, "CapslockToggle",
		kbdSettings, KBD_KEY_CAPSLOCK_TOGGLE)
	kbdManager.UserLayoutList = property.NewGSettingsStrvProperty(
		kbdManager, "UserLayoutList",
		kbdSettings, KBD_KEY_USER_LAYOUT_LIST)
	kbdManager.UserOptionList = property.NewGSettingsStrvProperty(
		kbdManager, "UserOptionList",
		kbdSettings, KBD_KEY_LAYOUT_OPTIONS)

	kbdManager.layoutDescMap = make(map[string]string)
	kbdManager.layoutDescMap = kbdManager.getLayoutList()

	var err error
	if kbdManager.greeterObj, err = greeterutils.NewGreeterUtils(
		"com.deepin.api.GreeterUtils",
		"/com/deepin/api/GreeterUtils"); err != nil {
		logger.Warning("New GreeterUtils failed:", err)
		kbdManager.greeterObj = nil
	}

	kbdManager.listenFlag = false

	kbdManager.init()

	return kbdManager
}

func (kbdManager *KeyboardManager) correctCurrentLayout() {
	curLayout := kbdManager.CurrentLayout.Get()
	if len(curLayout) < 1 {
		v := kbdManager.getDefaultLayout()
		if _, ok := kbdManager.layoutDescMap[v]; ok {
			curLayout = v
		} else {
			strs := strings.Split(v, LAYOUT_DELIM)
			curLayout = strs[0] + LAYOUT_DELIM
			if len(strs[1]) > 0 {
				list := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)
				list = append(list, strs[1])
				kbdSettings.SetStrv(KBD_KEY_LAYOUT_OPTIONS, list)
			}
		}
	}

	kbdManager.CurrentLayout.Set(curLayout)
	kbdSettings.SetString(KBD_KEY_LAYOUT, curLayout)
}

func (kbdManager *KeyboardManager) getDefaultLayout() string {
	layout := "us"
	option := ""
	contents, err := ioutil.ReadFile(KBD_DEFAULT_FILE)
	if err != nil {
		logger.Debug("ReadFile Failed:", err)
		return layout + LAYOUT_DELIM + option
	}

	lines := strings.Split(string(contents), "\n")
	for _, line := range lines {
		if ok, _ := regexp.MatchString(`^XKBLAYOUT`, line); ok {
			layout = strings.Split(line, "=")[1]
		} else if ok, _ := regexp.MatchString(`^XKBOPTIONS`, line); ok {
			option = strings.Split(line, "=")[1]
		}
	}

	layout = strings.Trim(layout, "\"")
	option = strings.Trim(option, "\"")
	if len(layout) < 1 {
		layout = "us"
		option = ""
	}

	return layout + LAYOUT_DELIM + option
}

func (kbdManager *KeyboardManager) isStrInList(str string, list []string) bool {
	for _, l := range list {
		if str == l {
			return true
		}
	}

	return false
}

func (kbdManager *KeyboardManager) setLayoutOptions() {
	options := kbdSettings.GetStrv(KBD_KEY_LAYOUT_OPTIONS)

	kbdManager.setLayoutOption("")
	for _, v := range options {
		kbdManager.setLayoutOption(v)
	}
}

func (kbdManager *KeyboardManager) setLayoutOption(option string) bool {
	if err := exec.Command(CMD_SETXKB,
		"-option", option).Run(); err != nil {
		logger.Warningf("Set option '%s' failed: %v", option, err)
		return false
	}

	return true
}

func (kbdManager *KeyboardManager) setLayout(value string) {
	layout := ""
	option := ""

	if len(value) < 1 || !strings.Contains(value, LAYOUT_DELIM) {
		layout = "us"
		option = ""
	} else {
		strs := strings.Split(value, LAYOUT_DELIM)
		if len(strs[0]) < 1 {
			layout = "us"
			option = ""
		} else {
			layout = strs[0]
			option = strs[1]
		}
	}

	kbdManager.setLayoutOptions()

	if err := exec.Command(CMD_SETXKB,
		"-layout", layout,
		"-option", option).Run(); err != nil {
		logger.Warningf("Set layout '%s ; %s' failed: %v",
			layout, option, err)
		return
	}

	value = layout + LAYOUT_DELIM + option
	if len(value) > 0 {
		list := kbdSettings.GetStrv(KBD_KEY_USER_LAYOUT_LIST)
		if !kbdManager.isStrInList(value, list) {
			list = append(list, value)
			kbdSettings.SetStrv(KBD_KEY_USER_LAYOUT_LIST, list)
		}

		kbdManager.setGreeterLayout(value)
	}
}

func (kbdManager *KeyboardManager) setQtCursorBlink(rate uint32) {
	if configDir := dutils.GetConfigDir(); len(configDir) > 0 {
		qtPath := path.Join(configDir, "Trolltech.conf")
		dutils.WriteKeyToKeyFile(qtPath, "Qt",
			"cursorFlashTime", rate)
	}
}

func (kbdManager *KeyboardManager) setCursorBlink(value uint32) {
	if xsObj != nil {
		xsObj.SetInterger("Net/CursorBlinkTime", value)
	}
	kbdManager.setQtCursorBlink(value)
}

func (kbdManager *KeyboardManager) setGreeterLayout(layout string) {
	if kbdManager.greeterObj == nil {
		return
	}

	username := dutils.GetUserName()
	homeDir := dutils.GetHomeDir()
	if homeDir != path.Join("/tmp", username) && len(username) > 0 {
		kbdManager.greeterObj.SetKbdLayout(username, layout)
	}
}

func (kbdManager *KeyboardManager) setGreeterLayoutList(list []string) {
	if kbdManager.greeterObj == nil {
		return
	}

	username := dutils.GetUserName()
	homeDir := dutils.GetHomeDir()
	if homeDir != path.Join("/tmp", username) && len(username) > 0 {
		kbdManager.greeterObj.SetKbdLayoutList(username, list)
	}
}

func (kbdManager *KeyboardManager) init() {
	if !GetManager().versionRight {
		kbdSettings.Reset(KBD_KEY_DELAY)
		kbdSettings.Reset(KBD_KEY_REPEAT_INTERVAL)
	}

	setKeyboardRepeat(kbdManager.RepeatEnabled.Get(),
		kbdManager.RepeatDelay.Get(),
		kbdManager.RepeatInterval.Get())
	kbdManager.correctCurrentLayout()
	kbdManager.setLayout(kbdManager.CurrentLayout.Get())
	kbdManager.setLayoutOptions()
	kbdManager.setCursorBlink(uint32(kbdManager.CursorBlink.Get()))

	if !kbdManager.listenFlag {
		kbdManager.listenGSettings()
	}
}
