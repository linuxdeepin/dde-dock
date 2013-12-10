package main

import (
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
)

type ExtDevManager struct {
	DevInfoList []ExtDeviceInfo
}

type ExtDeviceInfo struct {
	DevicePath string
	DeviceType string
}

type MouseEntry struct {
	UseHabit       *property.GSettingsBoolProperty  `access:"readwrite"`
	MoveSpeed      *property.GSettingsFloatProperty `access:"readwrite"`
	MoveAccuracy   *property.GSettingsFloatProperty   `access:"readwrite"`
	ClickFrequency *property.GSettingsIntProperty   `access:"readwrite"`
	DeviceID       string
}

type TPadEntry struct {
	TPadEnable     *property.GSettingsBoolProperty   `access:"readwrite"`
	UseHabit       *property.GSettingsStringProperty `access:"readwrite"`
	MoveSpeed      *property.GSettingsFloatProperty  `access:"readwrite"`
	MoveAccuracy   *property.GSettingsFloatProperty    `access:"readwrite"`
	ClickFrequency *property.GSettingsIntProperty    `access:"readwrite"`
	DragDelay      *property.GSettingsIntProperty    `access:"readwrite"`
	DeviceID       string
}

type KeyboardEntry struct {
	RepeatDelay    *property.GSettingsUintProperty `access:"readwrite"`
	RepeatSpeed    *property.GSettingsUintProperty `access:"readwrite"`
	CursorBlink    *property.GSettingsIntProperty  `access:"readwrite"`
	DisableTPad    *property.GSettingsBoolProperty `access:"readwrite"`
	KeyboardLayout *property.GSettingsStrvProperty `access:"readwrite"`
	DeviceID       string
}

const (
	_EXT_DEV_NAME = "com.deepin.daemon.ExtDevManager"
	_EXT_DEV_PATH = "/com/deepin/daemon/ExtDevManager"
	_EXT_DEV_IFC  = "com.deepin.daemon.ExtDevManager"

	_EXT_ENTRY_PATH = "/com/deepin/daemon/ExtDevManager/"
	_EXT_ENTRY_IFC  = "com.deepin.daemon.ExtDevManager."

	_KEYBOARD_REPEAT_SCHEMA = "org.gnome.settings-daemon.peripherals.keyboard"
	_LAYOUT_SCHEMA          = "org.gnome.libgnomekbd.keyboard"
	_DESKTOP_INFACE_SCHEMA  = "org.gnome.desktop.interface"
	_MOUSE_SCHEMA           = "org.gnome.settings-daemon.peripherals.mouse"
	_TPAD_SCHEMA            = "org.gnome.settings-daemon.peripherals.touchpad"
)

var (
	mouseGSettings     *gio.Settings
	tpadGSettings      *gio.Settings
	infaceGSettings    *gio.Settings
	layoutGSettings    *gio.Settings
	keyRepeatGSettings *gio.Settings
)

func InitGSettings() bool {
	mouseGSettings = gio.NewSettings(_MOUSE_SCHEMA)
	tpadGSettings = gio.NewSettings(_TPAD_SCHEMA)
	infaceGSettings = gio.NewSettings(_DESKTOP_INFACE_SCHEMA)
	layoutGSettings = gio.NewSettings(_LAYOUT_SCHEMA)
	keyRepeatGSettings = gio.NewSettings(_KEYBOARD_REPEAT_SCHEMA)
	return true
}

func NewKeyboardEntry() *KeyboardEntry {
	keyboard := &KeyboardEntry{}

	keyboard.DeviceID = "Keyboard"
	keyboard.RepeatDelay = property.NewGSettingsUintProperty(keyboard,
		"RepeatDelay", keyRepeatGSettings, "delay")
	keyboard.RepeatSpeed = property.NewGSettingsUintProperty(keyboard,
		"RepeatSpeed", keyRepeatGSettings, "repeat-interval")
	keyboard.DisableTPad = property.NewGSettingsBoolProperty(keyboard,
		"DisableTPad", tpadGSettings, "disable-while-typing")
	keyboard.CursorBlink = property.NewGSettingsIntProperty(keyboard,
		"CursorBlink", infaceGSettings, "cursor-blink-time")
	keyboard.KeyboardLayout = property.NewGSettingsStrvProperty(keyboard,
		"KeyboardLayout", layoutGSettings, "layouts")
	return keyboard
}

func (keyboard *KeyboardEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + keyboard.DeviceID,
		_EXT_ENTRY_IFC + keyboard.DeviceID,
	}
}

func NewMouseEntry() *MouseEntry {
	mouse := &MouseEntry{}

	mouse.DeviceID = "Mouse"
	mouse.UseHabit = property.NewGSettingsBoolProperty(mouse,
		"UseHabit", mouseGSettings, "left-handed")
	mouse.MoveSpeed = property.NewGSettingsFloatProperty(mouse,
		"MoveSpeed", mouseGSettings, "motion-acceleration")
	mouse.MoveAccuracy = property.NewGSettingsFloatProperty(mouse,
		"MoveAccuracy", mouseGSettings, "motion-threshold")
	mouse.ClickFrequency = property.NewGSettingsIntProperty(mouse,
		"ClickFrequency", mouseGSettings, "double-click")

	return mouse
}

func (mouse *MouseEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + mouse.DeviceID,
		_EXT_ENTRY_IFC + mouse.DeviceID,
	}
}

func NewTPadEntry() *TPadEntry {
	tpad := &TPadEntry{}

	tpad.DeviceID = "TouchPad"
	tpad.TPadEnable = property.NewGSettingsBoolProperty(tpad,
		"TPadEnable", tpadGSettings, "touchpad-enabled")
	tpad.UseHabit = property.NewGSettingsStringProperty(tpad,
		"UseHabit", tpadGSettings, "left-handed")
	tpad.MoveSpeed = property.NewGSettingsFloatProperty(tpad,
		"MoveSpeed", tpadGSettings, "motion-acceleration")
	tpad.MoveAccuracy = property.NewGSettingsFloatProperty(tpad,
		"MoveAccuracy", tpadGSettings, "motion-threshold")
	tpad.DragDelay = property.NewGSettingsIntProperty(tpad,
		"DragDelay", mouseGSettings, "drag-threshold")
	tpad.ClickFrequency = property.NewGSettingsIntProperty(tpad,
		"ClickFrequency", mouseGSettings, "double-click")

	return tpad
}

func (tpad *TPadEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + tpad.DeviceID,
		_EXT_ENTRY_IFC + tpad.DeviceID,
	}
}

func NewExtDevManager() *ExtDevManager {
	return &ExtDevManager{}
}

func (dev *ExtDevManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_EXT_DEV_NAME, _EXT_DEV_PATH, _EXT_DEV_IFC}
}
