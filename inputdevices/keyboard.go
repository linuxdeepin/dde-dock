/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package inputdevices

import (
	"bufio"
	"dbus/com/deepin/daemon/accounts"
	"fmt"
	"gir/gio-2.0"
	"os"
	"os/user"
	"path"
	"pkg.deepin.io/dde/api/dxinput"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/lib/dbus/property"
	dutils "pkg.deepin.io/lib/utils"
	"regexp"
	"strings"
)

const (
	kbdSchema = "com.deepin.dde.keyboard"

	kbdKeyRepeatEnable   = "repeat-enabled"
	kbdKeyRepeatInterval = "repeat-interval"
	kbdKeyRepeatDelay    = "delay"
	kbdKeyLayout         = "layout"
	kbdKeyLayoutModel    = "layout-model"
	kbdKeyLayoutOptions  = "layout-options"
	kbdKeyUserLayoutList = "user-layout-list"
	kbdKeyCursorBlink    = "cursor-blink-time"
	kbdKeyCapslockToggle = "capslock-toggle"

	layoutDelim      = ";"
	kbdDefaultLayout = "us" + layoutDelim

	kbdSystemConfig  = "/etc/default/keyboard"
	kbdGreeterConfig = "/var/lib/greeter/users.ini"
	qtDefaultConfig  = ".config/Trolltech.conf"

	cmdSetKbd = "/usr/bin/setxkbmap"
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

	setting       *gio.Settings
	userObj       *accounts.User
	layoutDescMap map[string]string

	devNumber int
}

var _kbd *Keyboard

func getKeyboard() *Keyboard {
	if _kbd == nil {
		_kbd = NewKeyboard()
	}

	return _kbd
}

func NewKeyboard() *Keyboard {
	var kbd = new(Keyboard)

	kbd.setting = gio.NewSettings(kbdSchema)
	kbd.CurrentLayout = property.NewGSettingsStringProperty(
		kbd, "CurrentLayout",
		kbd.setting, kbdKeyLayout)
	kbd.RepeatEnabled = property.NewGSettingsBoolProperty(
		kbd, "RepeatEnabled",
		kbd.setting, kbdKeyRepeatEnable)
	kbd.RepeatInterval = property.NewGSettingsUintProperty(
		kbd, "RepeatInterval",
		kbd.setting, kbdKeyRepeatInterval)
	kbd.RepeatDelay = property.NewGSettingsUintProperty(
		kbd, "RepeatDelay",
		kbd.setting, kbdKeyRepeatDelay)
	kbd.CursorBlink = property.NewGSettingsIntProperty(
		kbd, "CursorBlink",
		kbd.setting, kbdKeyCursorBlink)
	kbd.CapslockToggle = property.NewGSettingsBoolProperty(
		kbd, "CapslockToggle",
		kbd.setting, kbdKeyCapslockToggle)
	kbd.UserLayoutList = property.NewGSettingsStrvProperty(
		kbd, "UserLayoutList",
		kbd.setting, kbdKeyUserLayoutList)
	kbd.UserOptionList = property.NewGSettingsStrvProperty(
		kbd, "UserOptionList",
		kbd.setting, kbdKeyLayoutOptions)

	var err error
	kbd.layoutDescMap, err = getLayoutListByFile(kbdLayoutsXml)
	if err != nil {
		logger.Error("Get layout desc list failed:", err)
		return nil
	}

	cur, err := user.Current()
	if err != nil {
		logger.Warning("Get current user info failed:", err)
	} else {
		kbd.userObj, err = ddbus.NewUserByUid(cur.Uid)
		if err != nil {
			logger.Warning("New user object failed:", cur.Name, err)
			kbd.userObj = nil
		}
	}

	kbd.devNumber = getKeyboardNumber()
	return kbd
}

func (kbd *Keyboard) init() {
	if kbd.userObj != nil {
		value := kbd.userObj.Layout.Get()
		if len(value) != 0 && value != kbd.CurrentLayout.Get() {
			kbd.CurrentLayout.Set(value)
		}
	}

	kbd.setLayout()
	err := kbd.setOptions()
	if err != nil {
		logger.Debugf("Init keymap options failed: %v", err)
	}
	kbd.setRepeat()
}

func (kbd *Keyboard) handleDeviceChanged() {
	num := getKeyboardNumber()
	logger.Debug("Keyboard changed:", num, kbd.devNumber)
	if num > kbd.devNumber {
		kbd.init()
	}
	kbd.devNumber = num
}

func (kbd *Keyboard) correctLayout() {
	current := kbd.CurrentLayout.Get()
	if len(current) != 0 {
		return
	}

	system, _ := getSystemLayout(kbdSystemConfig)
	if len(system) == 0 {
		kbd.CurrentLayout.Set(kbdDefaultLayout)
	} else {
		kbd.CurrentLayout.Set(system)
	}
}

func (kbd *Keyboard) setLayout() {
	kbd.correctLayout()
	err := doSetLayout(kbd.CurrentLayout.Get())
	if err != nil {
		logger.Debugf("Set layout to '%s' failed: %v",
			kbd.CurrentLayout.Get(), err)
		return
	}

	kbd.setGreeterLayout()
	kbd.addUserLayout(kbd.CurrentLayout.Get())

	err = applyXmodmapConfig()
	if err != nil {
		logger.Warning("Failed to apply xmodmap:", err)
	}
}

func (kbd *Keyboard) setOptions() error {
	options := kbd.UserOptionList.Get()

	if len(options) == 0 {
		return nil
	}

	// the old value wouldn't be cleared, so we will force clear it.
	doAction(cmdSetKbd + " -option")

	cmd := cmdSetKbd
	for _, opt := range options {
		cmd += fmt.Sprintf(" -option %q", opt)
	}
	return doAction(cmd)
}

func (kbd *Keyboard) addUserLayout(layout string) {
	if len(layout) == 0 {
		return
	}

	_, ok := kbd.layoutDescMap[layout]
	if !ok {
		return
	}

	ret, added := addItemToList(layout, kbd.UserLayoutList.Get())
	if !added {
		return
	}
	kbd.UserLayoutList.Set(filterSpaceStr(ret))
}

func (kbd *Keyboard) delUserLayout(layout string) {
	if len(layout) == 0 {
		return
	}

	ret, deleted := delItemFromList(layout, kbd.UserLayoutList.Get())
	if !deleted {
		return
	}
	kbd.UserLayoutList.Set(filterSpaceStr(ret))
}

func (kbd *Keyboard) addUserOption(option string) {
	if len(option) == 0 {
		return
	}

	// TODO: check option validity

	ret, added := addItemToList(option, kbd.UserOptionList.Get())
	if !added {
		return
	}
	kbd.UserOptionList.Set(ret)
}

func (kbd *Keyboard) delUserOption(option string) {
	if len(option) == 0 {
		return
	}

	ret, deleted := delItemFromList(option, kbd.UserOptionList.Get())
	if !deleted {
		return
	}
	kbd.UserOptionList.Set(ret)
}

func (kbd *Keyboard) setCursorBlink() {
	value := kbd.CursorBlink.Get()
	xsSetInt32(xsPropBlinkTimeut, value)

	err := setQtCursorBlink(value, path.Join(os.Getenv("HOME"),
		qtDefaultConfig))
	if err != nil {
		logger.Debugf("Set qt cursor blink to '%v' failed: %v",
			value, err)
	}
}

func (kbd *Keyboard) setGreeterLayout() {
	if kbd.userObj == nil {
		return
	}

	name := kbd.userObj.UserName.Get()
	if isInvalidUser(name) {
		return
	}

	err := kbd.userObj.SetLayout(kbd.CurrentLayout.Get())
	if err != nil {
		logger.Debugf("Set '%s' greeter layout failed: %v", name, err)
	}
}

func (kbd *Keyboard) setGreeterLayoutList() {
	if kbd.userObj == nil {
		return
	}

	name := kbd.userObj.UserName.Get()
	if isInvalidUser(name) {
		return
	}

	err := kbd.userObj.SetHistoryLayout(kbd.UserLayoutList.Get())
	if err != nil {
		logger.Debugf("Set '%s' greeter layout list failed: %v",
			name, err)
	}
}

func (kbd *Keyboard) setRepeat() {
	var (
		repeat   = kbd.RepeatEnabled.Get()
		delay    = kbd.RepeatDelay.Get()
		interval = kbd.RepeatInterval.Get()
	)
	err := dxinput.SetKeyboardRepeat(repeat, delay, interval)
	if err != nil {
		logger.Debug("Set kbd repeat failed:", err, repeat, delay, interval)
	}
	setWMKeyboardRepeat(repeat, delay, interval)
}

func doSetLayout(value string) error {
	array := strings.Split(value, layoutDelim)
	if len(array) != 2 {
		return fmt.Errorf("Invalid layout: %s", value)
	}

	layout, variant := array[0], array[1]
	if layout != "us" {
		layout += ",us"
	}

	var cmd = fmt.Sprintf("%s -layout \"%s\" -variant \"%s\"",
		cmdSetKbd, layout, variant)
	return doAction(cmd)
}

func setQtCursorBlink(rate int32, file string) error {
	ok := dutils.WriteKeyToKeyFile(file, "Qt", "cursorFlashTime", rate)
	if !ok {
		return fmt.Errorf("Write failed")
	}

	return nil
}

func getSystemLayout(file string) (string, error) {
	fr, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	var (
		found   int
		layout  string
		variant string

		regLayout  = regexp.MustCompile(`^XKBLAYOUT=`)
		regVariant = regexp.MustCompile(`^XKBVARIANT=`)

		scanner = bufio.NewScanner(fr)
	)
	for scanner.Scan() {
		if found == 2 {
			break
		}

		var line = scanner.Text()
		if regLayout.MatchString(line) {
			layout = strings.Trim(getValueFromLine(line, "="), "\"")
			found += 1
			continue
		}

		if regVariant.MatchString(line) {
			variant = strings.Trim(getValueFromLine(line, "="), "\"")
			found += 1
		}
	}

	if len(layout) == 0 {
		return "", fmt.Errorf("Not found default layout")
	}

	return layout + layoutDelim + variant, nil
}

func getValueFromLine(line, delim string) string {
	array := strings.Split(line, delim)
	if len(array) != 2 {
		return ""
	}

	return strings.TrimSpace(array[1])
}

func isInvalidUser(name string) bool {
	if len(name) == 0 {
		return true
	}

	if os.Getenv("HOME") == path.Join("/tmp", name) {
		return true
	}

	return false
}

func applyXmodmapConfig() error {
	config := os.Getenv("HOME") + "/.Xmodmap"
	if !dutils.IsFileExist(config) {
		return nil
	}
	return doAction("xmodmap " + config)
}
