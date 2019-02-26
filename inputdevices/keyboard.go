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

package inputdevices

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"sync"

	"pkg.deepin.io/gir/gio-2.0"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.accounts"
	"pkg.deepin.io/dde/api/dxinput"
	ddbus "pkg.deepin.io/dde/daemon/dbus"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/gsprop"
	"pkg.deepin.io/lib/dbusutil/proxy"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	kbdSchema = "com.deepin.dde.keyboard"

	kbdKeyRepeatEnable   = "repeat-enabled"
	kbdKeyRepeatInterval = "repeat-interval"
	kbdKeyRepeatDelay    = "delay"
	kbdKeyLayoutOptions  = "layout-options"
	kbdKeyCursorBlink    = "cursor-blink-time"
	kbdKeyCapslockToggle = "capslock-toggle"

	layoutDelim      = ";"
	kbdDefaultLayout = "us" + layoutDelim

	kbdSystemConfig = "/etc/default/keyboard"
	qtDefaultConfig = ".config/Trolltech.conf"

	cmdSetKbd = "/usr/bin/setxkbmap"
)

type Keyboard struct {
	service       *dbusutil.Service
	sysSigLoop    *dbusutil.SignalLoop
	PropsMu       sync.RWMutex
	CurrentLayout string `prop:"access:rw"`
	// dbusutil-gen: equal=nil
	UserLayoutList []string

	// dbusutil-gen: ignore-below
	RepeatEnabled  gsprop.Bool `prop:"access:rw"`
	CapslockToggle gsprop.Bool `prop:"access:rw"`

	CursorBlink gsprop.Int `prop:"access:rw"`

	RepeatInterval gsprop.Uint `prop:"access:rw"`
	RepeatDelay    gsprop.Uint `prop:"access:rw"`

	UserOptionList gsprop.Strv

	setting       *gio.Settings
	user          *accounts.User
	layoutDescMap map[string]string

	devNumber int

	methods *struct {
		AddLayoutOption    func() `in:"option"`
		DeleteLayoutOption func() `in:"option"`

		AddUserLayout    func() `in:"layout"`
		DeleteUserLayout func() `in:"layout"`

		GetLayoutDesc func() `in:"layout" out:"description"`
		LayoutList    func() `out:"layout_list"`
	}
}

func newKeyboard(service *dbusutil.Service) *Keyboard {
	var kbd = new(Keyboard)

	kbd.service = service
	kbd.setting = gio.NewSettings(kbdSchema)
	kbd.RepeatEnabled.Bind(kbd.setting, kbdKeyRepeatEnable)
	kbd.RepeatInterval.Bind(kbd.setting, kbdKeyRepeatInterval)
	kbd.RepeatDelay.Bind(kbd.setting, kbdKeyRepeatDelay)
	kbd.CursorBlink.Bind(kbd.setting, kbdKeyCursorBlink)
	kbd.CapslockToggle.Bind(kbd.setting, kbdKeyCapslockToggle)
	kbd.UserOptionList.Bind(kbd.setting, kbdKeyLayoutOptions)

	var err error

	kbd.layoutDescMap, err = getLayoutListByFile(kbdLayoutsXml)
	if err != nil {
		logger.Error("failed to get layouts description:", err)
		return nil
	}

	sysConn, err := dbus.SystemBus()
	if err != nil {
		logger.Warning(err)
		return nil
	}
	kbd.sysSigLoop = dbusutil.NewSignalLoop(sysConn, 10)
	kbd.sysSigLoop.Start()
	kbd.initUser()

	if kbd.user != nil {
		// set current layout
		layout, err := kbd.user.Layout().Get(0)
		if err != nil {
			logger.Warning(err)
		} else {
			kbd.PropsMu.Lock()
			kbd.CurrentLayout = fixLayout(layout)
			kbd.PropsMu.Unlock()
		}

		// set layout list
		layoutList, err := kbd.user.HistoryLayout().Get(0)
		if err != nil {
			logger.Warning(err)
		} else {
			kbd.PropsMu.Lock()
			kbd.UserLayoutList = fixLayoutList(layoutList)
			kbd.PropsMu.Unlock()
		}
	}

	kbd.devNumber = getKeyboardNumber()
	return kbd
}

func fixLayout(layout string) string {
	if !strings.Contains(layout, layoutDelim) {
		return layout + layoutDelim
	}
	return layout
}

func fixLayoutList(layouts []string) []string {
	result := make([]string, len(layouts))
	for idx, layout := range layouts {
		result[idx] = fixLayout(layout)
	}
	return result
}

func (kbd *Keyboard) initUser() {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		logger.Warning(err)
		return
	}

	cur, err := user.Current()
	if err != nil {
		logger.Warning("failed to get current user:", err)
		return
	}

	kbd.user, err = ddbus.NewUserByUid(systemConn, cur.Uid)
	if err != nil {
		logger.Warningf("failed to new user by uid %s: %v", cur.Uid, err)
		return
	}
	kbd.user.InitSignalExt(kbd.sysSigLoop, true)
	kbd.user.Layout().ConnectChanged(func(hasValue bool, value string) {
		if !hasValue {
			return
		}

		kbd.PropsMu.Lock()
		kbd.setPropCurrentLayout(fixLayout(value))
		kbd.PropsMu.Unlock()

		kbd.applyLayout()
	})
	kbd.user.HistoryLayout().ConnectChanged(func(hasValue bool, value []string) {
		if !hasValue {
			return
		}

		kbd.PropsMu.Lock()
		kbd.setPropUserLayoutList(fixLayoutList(value))
		kbd.PropsMu.Unlock()
	})
}

func (kbd *Keyboard) destroy() {
	if kbd.user != nil {
		kbd.user.RemoveHandler(proxy.RemoveAllHandlers)
		kbd.user = nil
	}

	kbd.sysSigLoop.Stop()
}

func (kbd *Keyboard) init() {
	if kbd.user != nil {
		kbd.correctLayout()
	}
	kbd.applySettings()
}

func (kbd *Keyboard) applySettings() {
	kbd.applyLayout()
	kbd.applyOptions()
	kbd.applyRepeat()
}

func (kbd *Keyboard) correctLayout() {
	kbd.PropsMu.RLock()
	currentLayout := kbd.CurrentLayout
	kbd.PropsMu.RUnlock()

	if currentLayout == "" {
		layoutFromSysCfg, err := getSystemLayout(kbdSystemConfig)
		if err != nil {
			logger.Warning(err)
		}
		if layoutFromSysCfg == "" {
			layoutFromSysCfg = kbdDefaultLayout
		}

		kbd.setLayoutForAccountsUser(layoutFromSysCfg)
	}

	kbd.PropsMu.RLock()
	layoutList := kbd.UserLayoutList
	kbd.PropsMu.RUnlock()

	if len(layoutList) == 0 {
		kbd.setLayoutListForAccountsUser([]string{currentLayout})
	}
}

func (kbd *Keyboard) handleDeviceChanged() {
	num := getKeyboardNumber()
	logger.Debug("Keyboard changed:", num, kbd.devNumber)
	if num > kbd.devNumber {
		kbd.applySettings()
	}
	kbd.devNumber = num
}

func (kbd *Keyboard) applyLayout() {
	kbd.PropsMu.RLock()
	currentLayout := kbd.CurrentLayout
	kbd.PropsMu.RUnlock()

	err := applyLayout(currentLayout)
	if err != nil {
		logger.Warningf("failed to set layout to %q: %v", currentLayout, err)
		return
	}

	err = applyXmodmapConfig()
	if err != nil {
		logger.Warning("failed to apply xmodmap:", err)
	}
}

func (kbd *Keyboard) applyOptions() {
	options := kbd.UserOptionList.Get()
	if len(options) == 0 {
		return
	}

	// the old value wouldn't be cleared, so we will force clear it.
	err := doAction(cmdSetKbd + " -option")
	if err != nil {
		logger.Warning("failed to clear keymap option:", err)
		return
	}

	cmd := cmdSetKbd
	for _, opt := range options {
		cmd += fmt.Sprintf(" -option %q", opt)
	}
	err = doAction(cmd)
	if err != nil {
		logger.Warning("failed to set keymap options:", err)
	}
}

var errInvalidLayout = errors.New("invalid layout")

func (kbd *Keyboard) checkLayout(layout string) error {
	if layout == "" {
		return dbusutil.ToError(errInvalidLayout)
	}

	_, ok := kbd.layoutDescMap[layout]
	if !ok {
		return dbusutil.ToError(errInvalidLayout)
	}
	return nil
}

func (kbd *Keyboard) setCurrentLayout(write *dbusutil.PropertyWrite) *dbus.Error {
	layout := write.Value.(string)
	logger.Debugf("setCurrentLayout %q", layout)

	if kbd.user == nil {
		return dbusutil.ToError(errors.New("kbd.user is nil"))
	}

	err := kbd.checkLayout(layout)
	if err != nil {
		return dbusutil.ToError(err)
	}

	kbd.setLayoutForAccountsUser(layout)
	kbd.addUserLayout(layout)
	return nil
}

func (kbd *Keyboard) addUserLayout(layout string) {
	kbd.PropsMu.Lock()
	newLayoutList, added := addItemToList(layout, kbd.UserLayoutList)
	kbd.PropsMu.Unlock()

	if !added {
		return
	}

	kbd.setLayoutListForAccountsUser(newLayoutList)
}

func (kbd *Keyboard) delUserLayout(layout string) {
	if layout == "" {
		return
	}

	kbd.PropsMu.Lock()
	newLayoutList, deleted := delItemFromList(layout, kbd.UserLayoutList)
	kbd.PropsMu.Unlock()
	if !deleted {
		return
	}
	kbd.setLayoutListForAccountsUser(newLayoutList)
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

func (kbd *Keyboard) applyCursorBlink() {
	value := kbd.CursorBlink.Get()
	xsSetInt32(xsPropBlinkTimeut, value)

	err := setQtCursorBlink(value, path.Join(os.Getenv("HOME"),
		qtDefaultConfig))
	if err != nil {
		logger.Debugf("failed to set qt cursor blink to '%v': %v",
			value, err)
	}
}

func (kbd *Keyboard) setLayoutForAccountsUser(layout string) {
	if kbd.user == nil {
		return
	}

	err := kbd.user.SetLayout(0, layout)
	if err != nil {
		logger.Debug("failed to set layout for accounts user:", err)
	}
}

func (kbd *Keyboard) setLayoutListForAccountsUser(layoutList []string) {
	if kbd.user == nil {
		return
	}

	layoutList = filterSpaceStr(layoutList)

	err := kbd.user.SetHistoryLayout(0, layoutList)
	if err != nil {
		logger.Debug("failed to set layout list for accounts user:", err)
	}
}

func (kbd *Keyboard) applyRepeat() {
	var (
		repeat   = kbd.RepeatEnabled.Get()
		delay    = kbd.RepeatDelay.Get()
		interval = kbd.RepeatInterval.Get()
	)
	err := dxinput.SetKeyboardRepeat(repeat, delay, interval)
	if err != nil {
		logger.Debug("failed to set repeat:", err, repeat, delay, interval)
	}
	setWMKeyboardRepeat(repeat, delay, interval)
}

func applyLayout(value string) error {
	array := strings.Split(value, layoutDelim)
	if len(array) != 2 {
		return fmt.Errorf("invalid layout: %s", value)
	}

	layout, variant := array[0], array[1]
	if layout != "us" {
		layout += ",us"

		if variant != "" {
			variant += ","
		}
	}

	var cmd = fmt.Sprintf("%s -layout \"%s\" -variant \"%s\"",
		cmdSetKbd, layout, variant)
	return doAction(cmd)
}

func setQtCursorBlink(rate int32, file string) error {
	ok := dutils.WriteKeyToKeyFile(file, "Qt", "cursorFlashTime", rate)
	if !ok {
		return fmt.Errorf("write failed")
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
		return "", fmt.Errorf("not found default layout")
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

func applyXmodmapConfig() error {
	config := os.Getenv("HOME") + "/.Xmodmap"
	if !dutils.IsFileExist(config) {
		return nil
	}
	return doAction("xmodmap " + config)
}
