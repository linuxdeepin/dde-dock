package main

import (
        "dlib"
        "dlib/dbus"
        "dlib/dbus/property"
        "dlib/gio-2.0"
        "dlib/logger"
        "os"
)

const (
        _DESKTOP_DEST = "com.deepin.daemon.Desktop"
        _DESKTOP_PATH = "/com/deepin/daemon/Desktop"
        _DESKTOP_IFC  = "com.deepin.daemon.Desktop"

        _DESKTOP_SCHEMA = "com.deepin.dde.desktop"
        _DOCK_SCHEMA    = "com.deepin.dde.dock"

        _COMPIZ_INTEGRATED_SCHEMA = "org.compiz.integrated"
        _COMPIZ_COMMANDS_SCHEMA   = "org.compiz.commands"
        _COMPIZ_SCALE_SCHEMA      = "org.compiz.scale"
        _COMPIZ_COMMAND_PATH      = "/org/compiz/profiles/deepin/plugins/commands/"
        _COMPIZ_SCALE_PATH        = "/org/compiz/profiles/deepin/plugins/scale/"

        _LAUNCHER_CMD = "launcher"
)

const (
        ACTION_NONE           = int32(0)
        ACTION_OPENED_WINDOWS = int32(1)
        ACTION_LAUNCHER       = int32(2)
)

var (
        _compizIntegrated *gio.Settings
        _compizCommand    *gio.Settings
        _compizScale      *gio.Settings

        _runCommand11     string
        _runCommand12     string
        _runCommandEdge10 string
        _runCommandEdge11 string
        _scale            string

        logObject = logger.NewLogger("daemon/desktop")
)

type Manager struct {
        ShowComputerIcon *property.GSettingsBoolProperty   `access:"readwrite"`
        ShowHomeIcon     *property.GSettingsBoolProperty   `access:"readwrite"`
        ShowTrashIcon    *property.GSettingsBoolProperty   `access:"readwrite"`
        ShowDSCIcon      *property.GSettingsBoolProperty   `access:"readwrite"`
        DockMode         *property.GSettingsStringProperty `access:"readwrite"`
        TopLeft          int32                             `access:"readwrite"`
        BottomRight      int32                             `access:"readwrite"`
}

func NewManager() *Manager {
        desk := &Manager{}

        deskSettings := gio.NewSettings(_DESKTOP_SCHEMA)
        desk.ShowComputerIcon = property.NewGSettingsBoolProperty(desk, "ShowComputerIcon", deskSettings, "show-computer-icon")
        desk.ShowHomeIcon = property.NewGSettingsBoolProperty(desk, "ShowHomeIcon", deskSettings, "show-home-icon")
        desk.ShowTrashIcon = property.NewGSettingsBoolProperty(desk, "ShowTrashIcon", deskSettings, "show-trash-icon")
        desk.ShowDSCIcon = property.NewGSettingsBoolProperty(desk, "ShowDSCIcon", deskSettings, "show-dsc-icon")
        desk.DockMode = property.NewGSettingsStringProperty(desk, "DockMode", gio.NewSettings(_DOCK_SCHEMA), "hide-mode")

        desk.getEdgeAction()
        desk.listenCompizGSettings()

        return desk
}

func initCompizGSettings() {
        _compizIntegrated = gio.NewSettings(_COMPIZ_INTEGRATED_SCHEMA)
        _compizCommand = gio.NewSettingsWithPath(_COMPIZ_COMMANDS_SCHEMA,
                _COMPIZ_COMMAND_PATH)
        _compizScale = gio.NewSettingsWithPath(_COMPIZ_SCALE_SCHEMA,
                _COMPIZ_SCALE_PATH)

        _runCommand11 = _compizIntegrated.GetString("command-11")
        _runCommand12 = _compizIntegrated.GetString("command-12")
        _runCommandEdge10 = _compizCommand.GetString("run-command10-edge")
        _runCommandEdge11 = _compizCommand.GetString("run-command11-edge")
        _scale = _compizScale.GetString("initiate-edge")
}

func main() {
        defer func() {
                if err := recover(); err != nil {
                        logObject.Fatal("recover a error:", err)
                }
        }()

        logObject.SetRestartCommand("/usr/lib/deepin-daemon/desktop")
        initCompizGSettings()
        desk := NewManager()
        err := dbus.InstallOnSession(desk)
        if err != nil {
                logObject.Info("Install Session DBus Failed:", err)
                panic(err)
        }
        dbus.DealWithUnhandledMessage()

        //desk.printManager()
        go dlib.StartLoop()
        if err = dbus.Wait(); err != nil {
                logObject.Info("lost dbus session:", err)
                os.Exit(1)
        } else {
                os.Exit(0)
        }
}

func (m *Manager) printManager() {
        logObject.Infof("Top Action: %d", m.TopLeft)
        logObject.Infof("Bottom Action: %d", m.BottomRight)
}
