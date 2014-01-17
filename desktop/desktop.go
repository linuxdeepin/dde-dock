package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
	"dlib/logger"
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
	_COMPIZ_SCALE_PATH        = "/org/compiz/profiles/deepin/plugins/_scale/"

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
			logger.Println("recover a error:", err)
		}
	}()

	initCompizGSettings()
	desk := NewManager()
	err := dbus.InstallOnSession(desk)
	if err != nil {
		logger.Println("Install Session DBus Failed:", err)
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	desk.printManager()
	dlib.StartLoop()
}

func (m *Manager) printManager() {
	logger.Printf("Top Action: %d\n", m.TopLeft)
	logger.Printf("Bottom Action: %d\n", m.BottomRight)
}
