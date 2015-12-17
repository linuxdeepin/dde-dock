package inputdevices

import (
	"bufio"
	"dbus/com/deepin/api/greeterhelper"
	"fmt"
	"os"
	"os/user"
	"path"
	"pkg.deepin.io/dde/api/dxinput"
	"pkg.deepin.io/lib/dbus/property"
	"gir/gio-2.0"
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
	greeter       *greeterhelper.GreeterHelper
	layoutDescMap map[string]string
}

var _kbd *Keyboard

func getKeyboard() *Keyboard {
	if _kbd == nil {
		_kbd = NewKeyboard()

		_kbd.init()
		_kbd.handleGSettings()
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

	kbd.greeter, err = greeterhelper.NewGreeterHelper(
		"com.deepin.api.GreeterHelper",
		"/com/deepin/api/GreeterHelper")
	if err != nil {
		logger.Warning("Create 'GreeterHelper' failed:", err)
		kbd.greeter = nil
	}

	return kbd
}

func (kbd *Keyboard) init() {
	group, _ := getUsername()
	value, _ := getGreeterLayout(kbdGreeterConfig, group)
	if len(value) != 0 && value != kbd.CurrentLayout.Get() {
		kbd.CurrentLayout.Set(value)
	}

	kbd.setLayout()
	kbd.setOptions()
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
}

func (kbd *Keyboard) setOptions() {
	// clear old options
	var cmd = fmt.Sprintf("%s -option \"\"", cmdSetKbd)
	doAction(cmd)

	for _, option := range kbd.UserOptionList.Get() {
		err := doSetOption(option)
		if err != nil {
			logger.Debugf("Set option '%s' failed: %v", option, err)
		}
	}
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
	if kbd.greeter == nil {
		return
	}

	name, _ := getUsername()
	if isInvalidUser(name) {
		return
	}

	err := kbd.greeter.SetLayout(name, kbd.CurrentLayout.Get())
	if err != nil {
		logger.Debugf("Set '%s' greeter layout failed: %v", name, err)
	}
}

func (kbd *Keyboard) setGreeterLayoutList() {
	if kbd.greeter == nil {
		return
	}

	name, _ := getUsername()
	if isInvalidUser(name) {
		return
	}

	err := kbd.greeter.SetLayoutList(name, kbd.UserLayoutList.Get())
	if err != nil {
		logger.Debugf("Set '%s' greeter layout list failed: %v",
			name, err)
	}
}

func (kbd *Keyboard) setRepeat() {
	err := dxinput.SetKeyboardRepeat(kbd.RepeatEnabled.Get(),
		kbd.RepeatDelay.Get(), kbd.RepeatInterval.Get())
	if err != nil {
		logger.Debug("Set kbd repeat failed:", err)
	}
}

func doSetLayout(value string) error {
	array := strings.Split(value, layoutDelim)
	if len(array) != 2 {
		return fmt.Errorf("Invalid layout: %s", value)
	}

	var cmd = fmt.Sprintf("%s -layout \"%s\" -variant \"%s\"",
		cmdSetKbd, array[0], array[1])
	return doAction(cmd)
}

func doSetOption(option string) error {
	var cmd = fmt.Sprintf("%s -option \"%s\"", cmdSetKbd, option)
	return doAction(cmd)
}

func setQtCursorBlink(rate int32, file string) error {
	ok := dutils.WriteKeyToKeyFile(file, "Qt", "cursorFlashTime", rate)
	if !ok {
		return fmt.Errorf("Write failed")
	}

	return nil
}

func getGreeterLayout(file, group string) (string, error) {
	if isInvalidUser(group) {
		return "", fmt.Errorf("Invalid group: %s", group)
	}

	kfile, err := dutils.NewKeyFileFromFile(file)
	if err != nil {
		return "", err
	}
	defer kfile.Free()

	value, err := kfile.GetString(group, "KeyboardLayout")
	if err != nil {
		return "", err
	}

	array := strings.Split(value, "|")
	if len(array) != 2 {
		return "", fmt.Errorf("Invalid kbd layout: %v", value)
	}

	return array[0] + layoutDelim + array[1], nil
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

		regLayout = regexp.MustCompile(`^XKBLAYOUT=`)
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

func getUsername() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}

	return u.Username, nil
}
