package main

import (
        "dlib/dbus/property"
        "dlib/gio-2.0"
        libutils "dlib/utils"
)

type ExtDevManager struct {
        DevInfoList []ExtDeviceInfo
}

type ExtDeviceInfo struct {
        Path    string
        Type    string
}

type MouseEntry struct {
        UseHabit       *property.GSettingsBoolProperty  `access:"readwrite"`
        MoveSpeed      *property.GSettingsFloatProperty `access:"readwrite"`
        MoveAccuracy   *property.GSettingsFloatProperty `access:"readwrite"`
        ClickFrequency *property.GSettingsIntProperty   `access:"readwrite"`
        DeviceID       string
}

type TPadEntry struct {
        TPadEnable     *property.GSettingsBoolProperty   `access:"readwrite"`
        UseHabit       *property.GSettingsStringProperty `access:"readwrite"`
        MoveSpeed      *property.GSettingsFloatProperty  `access:"readwrite"`
        MoveAccuracy   *property.GSettingsFloatProperty  `access:"readwrite"`
        ClickFrequency *property.GSettingsIntProperty    `access:"readwrite"`
        DragDelay      *property.GSettingsIntProperty    `access:"readwrite"`
        DeviceID       string
}

type KeyboardEntry struct {
        RepeatDelay    *property.GSettingsUintProperty `access:"readwrite"`
        RepeatSpeed    *property.GSettingsUintProperty `access:"readwrite"`
        CursorBlink    *property.GSettingsIntProperty  `access:"readwrite"`
        DisableTPad    *property.GSettingsBoolProperty `access:"readwrite"`
        CurrentLayout  string                          `access:"readwrite"`
        UserLayoutList []string
        DeviceID       string
}

const (
        _EXT_DEV_NAME = "com.deepin.daemon.InputDevices"
        _EXT_DEV_PATH = "/com/deepin/daemon/InputDevices"
        _EXT_DEV_IFC  = "com.deepin.daemon.InputDevices"

        _EXT_ENTRY_PATH = "/com/deepin/daemon/InputDevice/"
        _EXT_ENTRY_IFC  = "com.deepin.daemon.InputDevice."

        _KEYBOARD_REPEAT_SCHEMA = "com.deepin.dde.peripherals.keyboard"
        _LAYOUT_SCHEMA          = "org.gnome.libgnomekbd.keyboard"
        _DESKTOP_INFACE_SCHEMA  = "org.gnome.desktop.interface"
        _MOUSE_SCHEMA           = "com.deepin.dde.peripherals.mouse"
        _TPAD_SCHEMA            = "com.deepin.dde.peripherals.touchpad"
)

var (
        _mouseGSettings     *gio.Settings
        _tpadGSettings      *gio.Settings
        _infaceGSettings    *gio.Settings
        _layoutGSettings    *gio.Settings
        _keyRepeatGSettings *gio.Settings
        objUtils            = libutils.NewUtils()
)

func InitGSettings() bool {
        _mouseGSettings = gio.NewSettings(_MOUSE_SCHEMA)
        _tpadGSettings = gio.NewSettings(_TPAD_SCHEMA)
        _infaceGSettings = gio.NewSettings(_DESKTOP_INFACE_SCHEMA)
        _layoutGSettings = gio.NewSettings(_LAYOUT_SCHEMA)
        _keyRepeatGSettings = gio.NewSettings(_KEYBOARD_REPEAT_SCHEMA)
        return true
}

func NewKeyboardEntry() *KeyboardEntry {
        keyboard := &KeyboardEntry{}

        keyboard.DeviceID = "Keyboard"
        keyboard.RepeatDelay = property.NewGSettingsUintProperty(keyboard,
                "RepeatDelay", _keyRepeatGSettings, "delay")
        keyboard.RepeatSpeed = property.NewGSettingsUintProperty(keyboard,
                "RepeatSpeed", _keyRepeatGSettings, "repeat-interval")
        keyboard.DisableTPad = property.NewGSettingsBoolProperty(keyboard,
                "DisableTPad", _tpadGSettings, "disable-while-typing")
        keyboard.CursorBlink = property.NewGSettingsIntProperty(keyboard,
                "CursorBlink", _infaceGSettings, "cursor-blink-time")
        keyboard.setPropName("CurrentLayout")
        keyboard.setPropName("UserLayoutList")
        keyboard.appendUserLayout(keyboard.CurrentLayout)
        keyboard.listenLayoutChanged()

        keyboard.applyPropValue("CursorBlink", keyboard.CursorBlink.GetValue())
        keyboard.applyPropValue("CurrentLayout", keyboard.CurrentLayout)

        return keyboard
}

func (keyboard *KeyboardEntry) LayoutList() map[string]string {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Info("recover error in get layout list:",
                                err)
                }
        }()
        datas := ParseXML(_LAYOUT_XML_PATH)
        layouts := GetLayoutList(datas)

        return layouts
}

func (keyboard *KeyboardEntry) GetDescByLayout(layout string) string {
        listMap := keyboard.LayoutList()
        desc, ok := listMap[layout]
        if !ok {
                logObject.Warningf("'%s' invalid layout", layout)
                return ""
        }

        return desc
}

func (keyboard *KeyboardEntry) AddUserLayout(strs string) {
        keyboard.appendUserLayout(strs)
}

func (keyboard *KeyboardEntry) DeleteUserLayout(strs string) {
        if len(strs) <= 0 {
                return
        }
        keyboard.deleteUserLayout(strs)
}

func NewMouseEntry() *MouseEntry {
        mouse := &MouseEntry{}

        mouse.DeviceID = "Mouse"
        mouse.UseHabit = property.NewGSettingsBoolProperty(mouse,
                "UseHabit", _mouseGSettings, "left-handed")
        mouse.MoveSpeed = property.NewGSettingsFloatProperty(mouse,
                "MoveSpeed", _mouseGSettings, "motion-acceleration")
        mouse.MoveAccuracy = property.NewGSettingsFloatProperty(mouse,
                "MoveAccuracy", _mouseGSettings, "motion-threshold")
        mouse.ClickFrequency = property.NewGSettingsIntProperty(mouse,
                "ClickFrequency", _mouseGSettings, "double-click")

        return mouse
}

func NewTPadEntry() *TPadEntry {
        tpad := &TPadEntry{}

        tpad.DeviceID = "TouchPad"
        tpad.TPadEnable = property.NewGSettingsBoolProperty(tpad,
                "TPadEnable", _tpadGSettings, "touchpad-enabled")
        tpad.UseHabit = property.NewGSettingsStringProperty(tpad,
                "UseHabit", _tpadGSettings, "left-handed")
        tpad.MoveSpeed = property.NewGSettingsFloatProperty(tpad,
                "MoveSpeed", _tpadGSettings, "motion-acceleration")
        tpad.MoveAccuracy = property.NewGSettingsFloatProperty(tpad,
                "MoveAccuracy", _tpadGSettings, "motion-threshold")
        tpad.DragDelay = property.NewGSettingsIntProperty(tpad,
                "DragDelay", _mouseGSettings, "drag-threshold")
        tpad.ClickFrequency = property.NewGSettingsIntProperty(tpad,
                "ClickFrequency", _mouseGSettings, "double-click")

        return tpad
}

func NewExtDevManager() *ExtDevManager {
        return &ExtDevManager{}
}
