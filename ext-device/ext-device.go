package main

import (
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
)

type ExtDevManager struct {
	DevInfoList []ExtDeviceInfo `access:"read"`
}

type ExtDeviceInfo struct {
	DevicePath string
	DeviceType string
}

type MouseEntry struct {
	UseHabit       dbus.Property
	MoveSpeed      dbus.Property
	MoveAccuracy   dbus.Property
	ClickFrequency dbus.Property
	DeviceID       string `access:"read"`
}

type TPadEntry struct {
	UseHabit       dbus.Property
	MoveSpeed      dbus.Property
	MoveAccuracy   dbus.Property
	ClickFrequency dbus.Property
	DragDelay      dbus.Property
	DeviceID       string `access:"read"`
}

type KeyboardEntry struct {
	RepeatDelay    dbus.Property
	RepeatSpeed    dbus.Property
	CursorBlink    dbus.Property
	DisableTPad    dbus.Property
	KeyboardLayout dbus.Property
	DeviceID       string `access:"read"`
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
	busConn            *dbus.Conn
	mouseGSettings     *gio.Settings
	tpadGSettings      *gio.Settings
	infaceGSettings    *gio.Settings
	layoutGSettings    *gio.Settings
	keyRepeatGSettings *gio.Settings
)

func InitGSettings() bool {
	var dbusError error
	busConn, dbusError = dbus.SessionBus()
	if dbusError != nil {
		return false
	}
	mouseGSettings = gio.NewSettings(_MOUSE_SCHEMA)
	tpadGSettings = gio.NewSettings(_TPAD_SCHEMA)
	infaceGSettings = gio.NewSettings(_DESKTOP_INFACE_SCHEMA)
	layoutGSettings = gio.NewSettings(_LAYOUT_SCHEMA)
	keyRepeatGSettings = gio.NewSettings(_KEYBOARD_REPEAT_SCHEMA)
	return true
}

func NewKeyboardEntry() *KeyboardEntry {
	keyboard := KeyboardEntry{}

	keyboard.DeviceID = "Keyboard"
	keyboard.RepeatDelay = property.NewGSettingsPropertyFull(
		keyRepeatGSettings, "delay", uint32(0), busConn,
		_EXT_ENTRY_PATH+keyboard.DeviceID,
		_EXT_ENTRY_IFC+keyboard.DeviceID, "RepeatDelay")
	keyboard.RepeatSpeed = property.NewGSettingsPropertyFull(
		keyRepeatGSettings, "repeat-interval", uint32(0), busConn,
		_EXT_ENTRY_PATH+keyboard.DeviceID,
		_EXT_ENTRY_IFC+keyboard.DeviceID, "RepeatSpeed")
	keyboard.DisableTPad = property.NewGSettingsPropertyFull(
		tpadGSettings, "disable-while-typing", true, busConn,
		_EXT_ENTRY_PATH+keyboard.DeviceID,
		_EXT_ENTRY_IFC+keyboard.DeviceID, "DisableTPad")
	keyboard.CursorBlink = property.NewGSettingsPropertyFull(
		infaceGSettings, "cursor-blink-time", int32(0), busConn,
		_EXT_ENTRY_PATH+keyboard.DeviceID,
		_EXT_ENTRY_IFC+keyboard.DeviceID, "CursorBlink")
	keyboard.KeyboardLayout = property.NewGSettingsPropertyFull(
		layoutGSettings, "layouts", []string{}, busConn,
		_EXT_ENTRY_PATH+keyboard.DeviceID,
		_EXT_ENTRY_IFC+keyboard.DeviceID, "KeyboardLayout")
	return &keyboard
}

func (keyboard *KeyboardEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + keyboard.DeviceID,
		_EXT_ENTRY_IFC + keyboard.DeviceID,
	}
}

func NewMouseEntry() *MouseEntry {
	mouse := MouseEntry{}

	mouse.DeviceID = "Mouse"
	mouse.UseHabit = property.NewGSettingsPropertyFull(mouseGSettings,
		"left-handed", false, busConn, 
		_EXT_ENTRY_PATH+mouse.DeviceID,
		_EXT_ENTRY_IFC+mouse.DeviceID, "UseHabit")
	mouse.MoveSpeed = property.NewGSettingsPropertyFull(mouseGSettings,
		"motion-acceleration", float64(0), busConn,
		_EXT_ENTRY_PATH+mouse.DeviceID,
		_EXT_ENTRY_IFC+mouse.DeviceID, "MoveSpeed")
	mouse.MoveAccuracy = property.NewGSettingsPropertyFull(mouseGSettings,
		"motion-threshold", int64(0), busConn,
		_EXT_ENTRY_PATH+mouse.DeviceID,
		_EXT_ENTRY_IFC+mouse.DeviceID, "MoveAccuracy")
	mouse.ClickFrequency = property.NewGSettingsPropertyFull(mouseGSettings,
		"double-click", int64(0), busConn,
		_EXT_ENTRY_PATH+mouse.DeviceID,
		_EXT_ENTRY_IFC+mouse.DeviceID,
		"ClickFrequency")

	return &mouse
}

func (mouse *MouseEntry) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_EXT_DEV_NAME,
		_EXT_ENTRY_PATH + mouse.DeviceID,
		_EXT_ENTRY_IFC + mouse.DeviceID,
	}
}

func NewTPadEntry() *TPadEntry {
	tpad := TPadEntry{}

	tpad.DeviceID = "TouchPad"
	tpad.UseHabit = property.NewGSettingsPropertyFull(tpadGSettings,
		"left-handed", "", busConn,
		_EXT_ENTRY_PATH+tpad.DeviceID,
		_EXT_ENTRY_IFC+tpad.DeviceID, "UseHabit")
	tpad.MoveSpeed = property.NewGSettingsPropertyFull(tpadGSettings,
		"motion-acceleration", float64(0), busConn,
		_EXT_ENTRY_PATH+tpad.DeviceID,
		_EXT_ENTRY_IFC+tpad.DeviceID, "MoveSpeed")
	tpad.MoveAccuracy = property.NewGSettingsPropertyFull(tpadGSettings,
		"motion-threshold", int64(0), busConn,
		_EXT_ENTRY_PATH+tpad.DeviceID,
		_EXT_ENTRY_IFC+tpad.DeviceID, "MoveAccuracy")
	tpad.DragDelay = property.NewGSettingsPropertyFull(mouseGSettings,
		"drag-threshold", int64(0), busConn,
		_EXT_ENTRY_PATH+tpad.DeviceID,
		_EXT_ENTRY_IFC+tpad.DeviceID, "DragDelay")
	tpad.ClickFrequency = property.NewGSettingsPropertyFull(mouseGSettings,
		"double-click", int64(0), busConn,
		_EXT_ENTRY_PATH+tpad.DeviceID,
		_EXT_ENTRY_IFC+tpad.DeviceID, "ClickFrequency")

	return &tpad
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
