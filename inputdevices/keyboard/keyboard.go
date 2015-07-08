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

package keyboard

import (
	"bufio"
	"dbus/com/deepin/api/greeterutils"
	"dbus/com/deepin/sessionmanager"
	"fmt"
	"os"
	"os/exec"
	"path"
	"pkg.deepin.io/dde/daemon/inputdevices/wrapper"
	"pkg.deepin.io/lib/dbus/property"
	"pkg.deepin.io/lib/gio-2.0"
	"pkg.deepin.io/lib/log"
	dutils "pkg.deepin.io/lib/utils"
	"regexp"
	"strings"
)

const (
	kbdKeyRepeatEnable   = "repeat-enabled"
	kbdKeyRepeatInterval = "repeat-interval"
	kbdKeyRepeatDelay    = "delay"
	kbdKeyLayout         = "layout"
	KBD_KEY_LAYOUT_MODEL = "layout-model"
	kbdKeyLayoutOptions  = "layout-options"
	kbdKeyUserLayoutList = "user-layout-list"
	kbdKeyCursorBlink    = "cursor-blink-time"
	kbdKeyCapslockToggle = "capslock-toggle"
	kbdDefaultConfig     = "/etc/default/keyboard"

	kbdSetCommand = "/usr/bin/setxkbmap"
)

type Keyboard struct {
	RepeatEnabled  *property.GSettingsBoolProperty `access:"readwrite"`
	CapslockToggle *property.GSettingsBoolProperty `access:"readwrite"`

	CursorBlink *property.GSettingsIntProperty `access:"readwrite"`

	RepeatInterval *property.GSettingsUintProperty `access:"readwrite"`
	RepeatDelay    *property.GSettingsUintProperty `access:"readwrite"`

	CurrentLayout *property.GSettingsStringProperty `access:"readwrite"`

	UserLayoutList *property.GSettingsStrvProperty
	UserOptionList *property.GSettingsStrvProperty

	logger        *log.Logger
	settings      *gio.Settings
	greeter       *greeterutils.GreeterUtils
	xsettings     *sessionmanager.XSettings
	layoutDescMap map[string]string
}

func NewKeyboard(l *log.Logger) *Keyboard {
	kbd := &Keyboard{}

	kbd.settings = gio.NewSettings("com.deepin.dde.keyboard")
	kbd.CurrentLayout = property.NewGSettingsStringProperty(
		kbd, "CurrentLayout",
		kbd.settings, kbdKeyLayout)
	kbd.RepeatEnabled = property.NewGSettingsBoolProperty(
		kbd, "RepeatEnabled",
		kbd.settings, kbdKeyRepeatEnable)
	kbd.RepeatInterval = property.NewGSettingsUintProperty(
		kbd, "RepeatInterval",
		kbd.settings, kbdKeyRepeatInterval)
	kbd.RepeatDelay = property.NewGSettingsUintProperty(
		kbd, "RepeatDelay",
		kbd.settings, kbdKeyRepeatDelay)
	kbd.CursorBlink = property.NewGSettingsIntProperty(
		kbd, "CursorBlink",
		kbd.settings, kbdKeyCursorBlink)
	kbd.CapslockToggle = property.NewGSettingsBoolProperty(
		kbd, "CapslockToggle",
		kbd.settings, kbdKeyCapslockToggle)
	kbd.UserLayoutList = property.NewGSettingsStrvProperty(
		kbd, "UserLayoutList",
		kbd.settings, kbdKeyUserLayoutList)
	kbd.UserOptionList = property.NewGSettingsStrvProperty(
		kbd, "UserOptionList",
		kbd.settings, kbdKeyLayoutOptions)

	kbd.logger = l
	var err error
	kbd.layoutDescMap, err = getLayoutListByFile(kbdKeyLayoutXml)
	if err != nil {
		//TODO: handle error
		kbd.errorInfo("Get Layout Desc List Failed: %v", err)
		return nil
	}

	kbd.greeter, err = greeterutils.NewGreeterUtils(
		"com.deepin.api.GreeterUtils",
		"/com/deepin/api/GreeterUtils")
	if err != nil {
		kbd.warningInfo("Create GreeterUtils Failed: %v", err)
		kbd.greeter = nil
	}

	kbd.xsettings, err = sessionmanager.NewXSettings(
		"com.deepin.SessionManager",
		"/com/deepin/XSettings",
	)
	if err != nil {
		kbd.warningInfo("Create XSettings Failed: %v", err)
		kbd.xsettings = nil
	}

	kbd.handleGSettings()
	kbd.init()

	return kbd
}

func (kbd *Keyboard) setLayout() {
	if len(kbd.CurrentLayout.Get()) == 0 {
		kbd.settings.SetString(kbdKeyLayout,
			getLayoutFromFile(kbdDefaultConfig))
		return
	}

	err := setUserLayout(kbd.CurrentLayout.Get())
	if err != nil {
		kbd.debugInfo("Set Layout '%s' Failed: %v",
			kbd.CurrentLayout.Get(), err)
		kbd.settings.SetString(kbdKeyLayout,
			getLayoutFromFile(kbdDefaultConfig))
		return
	}
	//kbd.setLayoutOptions()

	value := kbd.CurrentLayout.Get()
	kbd.addUserLayoutToList(value)
	kbd.setGreeterLayoutList(kbd.settings.GetStrv(kbdKeyUserLayoutList))
}

func (kbd *Keyboard) addUserLayoutToList(value string) {
	if len(value) == 0 {
		return
	}

	list := kbd.settings.GetStrv(kbdKeyUserLayoutList)
	if isStrInList(value, list) {
		return
	}

	list = append(list, value)
	kbd.settings.SetStrv(kbdKeyUserLayoutList, list)
}

func (kbd *Keyboard) deleteUserLayoutFromList(value string) {
	list := kbd.settings.GetStrv(kbdKeyUserLayoutList)
	var tmpList []string
	for _, v := range list {
		if v == value {
			continue
		}
		tmpList = append(tmpList, v)
	}
	if len(list) == len(tmpList) {
		return
	}

	kbd.settings.SetStrv(kbdKeyUserLayoutList, tmpList)
}

func (kbd *Keyboard) setLayoutOptions() {
	options := kbd.settings.GetStrv(kbdKeyLayoutOptions)

	// clear layout option settings
	setOptionList([]string{""})
	setOptionList(options)
}

func (kbd *Keyboard) setCursorBlink(value uint32) {
	if kbd.xsettings != nil {
		kbd.xsettings.SetInteger("Net/CursorBlinkTime", value)
	}
	configDir := dutils.GetConfigDir()
	qtPath := path.Join(configDir, "Trolltech.conf")
	setQtCursorBlink(value, qtPath)
}

func (kbd *Keyboard) setGreeterLayoutList(list []string) {
	if kbd.greeter == nil {
		return
	}

	username := dutils.GetUserName()
	homeDir := dutils.GetHomeDir()
	if homeDir == path.Join("/tmp", username) && len(username) == 0 {
		return
	}
	kbd.greeter.SetKbdLayoutList(username, list)
}

func (kbd *Keyboard) init() {
	kbd.setLayout()
	kbd.setLayoutOptions()
	kbd.setCursorBlink(uint32(kbd.CursorBlink.Get()))
	setKbdRepeat(kbd.RepeatEnabled.Get(),
		kbd.RepeatDelay.Get(),
		kbd.RepeatInterval.Get())
}

func setKbdRepeat(enable bool, delay, interval uint32) {
	wrapper.SetKeyboardRepeat(enable, delay, interval)
}

func getLayoutFromFile(config string) string {
	layout := "us"
	variant := ""
	fp, err := os.Open(config)
	if err != nil {
		return layout + kbdKeyLayoutDelim + variant
	}

	scanner := bufio.NewScanner(fp)
	var found int
	for scanner.Scan() {
		if found == 2 {
			break
		}

		line := scanner.Text()
		ok, _ := regexp.MatchString(`^XKBLAYOUT=`, line)
		if ok {
			found += 1
			strs := strings.Split(line, "=")
			layout = strs[1]
			continue
		}

		ok, _ = regexp.MatchString(`^XKBVARIANT=`, line)
		if ok {
			found += 1
			strs := strings.Split(line, "=")
			variant = strs[1]
		}
	}
	fp.Close()

	layout = strings.Trim(layout, "\"")
	variant = strings.Trim(variant, "\"")
	if len(layout) == 0 {
		layout = "us"
		variant = ""
	}

	return layout + kbdKeyLayoutDelim + variant
}

func setUserLayout(value string) error {
	var (
		layout  string
		variant string
	)

	strs := strings.Split(value, kbdKeyLayoutDelim)
	if len(strs[0]) == 0 {
		layout = "us"
		variant = ""
	} else {
		layout = strs[0]
		variant = strs[1]
	}

	err := exec.Command("/bin/sh", "-c",
		kbdSetCommand+" -layout \""+layout+
			"\" -variant \""+variant+"\"").Run()
	return err
}

func setOptionList(list []string) {
	if len(list) == 0 {
		exec.Command("/bin/sh", "-c",
			kbdSetCommand+" -option \"\"").Run()
		return
	}

	for _, option := range list {
		exec.Command("/bin/sh", "-c",
			kbdSetCommand+" -option \""+option+"\"").Run()
		//TODO:handle error
	}
}

func setQtCursorBlink(rate uint32, filename string) error {
	ok := dutils.WriteKeyToKeyFile(filename, "Qt",
		"cursorFlashTime", rate)
	if !ok {
		return fmt.Errorf("Set Qt CursorBlink Failed")
	}

	return nil
}

func isStrInList(str string, list []string) bool {
	for _, l := range list {
		if str == l {
			return true
		}
	}

	return false
}
