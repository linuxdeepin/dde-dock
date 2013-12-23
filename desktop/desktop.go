package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
	"dlib/gio-2.0"
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

type DesktopManager struct {
	ShowComputerIcon *property.GSettingsBoolProperty   `access:"readwrite"`
	ShowHomeIcon     *property.GSettingsBoolProperty   `access:"readwrite"`
	ShowTrashIcon    *property.GSettingsBoolProperty   `access:"readwrite"`
	ShowDSCIcon      *property.GSettingsBoolProperty   `access:"readwrite"`
	DockMode         *property.GSettingsStringProperty `access:"readwrite"`
	TopLeft          int32
	BottomRight      int32
}

func (desk *DesktopManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DESKTOP_DEST, _DESKTOP_PATH, _DESKTOP_IFC}
}

func (desk *DesktopManager) SetTopLeftAction(index int32) {
	if index == ACTION_NONE {
		rightTmp := desk.BottomRight
		_compizIntegrated.SetString("command-11", "")
		_compizCommand.SetString("run-command10-edge", "")
		_compizScale.SetString("initiate-edge", "")

		if rightTmp == ACTION_OPENED_WINDOWS {
			_compizScale.SetString("initiate-edge", "BottomRight")
		}
	} else if index == ACTION_OPENED_WINDOWS {
		if desk.BottomRight == ACTION_OPENED_WINDOWS {
			desk.BottomRight = ACTION_LAUNCHER
			_compizIntegrated.SetString("command-12", _LAUNCHER_CMD)
			_compizCommand.SetString("run-command11-edge", "BottomRight")
		}

		_compizIntegrated.SetString("command-11", "")
		_compizCommand.SetString("run-command10-edge", "")
		_compizScale.SetString("initiate-edge", "TopLeft")
	} else if index == ACTION_LAUNCHER {
		if desk.BottomRight == ACTION_LAUNCHER {
			desk.BottomRight = ACTION_OPENED_WINDOWS
			_compizIntegrated.SetString("command-12", "")
			_compizCommand.SetString("run-command11-edge", "")
			_compizScale.SetString("initiate-edge", "BottomRight")
		}

		_compizIntegrated.SetString("command-11", _LAUNCHER_CMD)
		_compizCommand.SetString("run-command10-edge", "TopLeft")
	}
}

func (desk *DesktopManager) SetBottomRightAction(index int32) {
	if index == ACTION_NONE {
		leftTmp := desk.TopLeft
		_compizIntegrated.SetString("command-12", "")
		_compizCommand.SetString("run-command11-edge", "")
		_compizScale.SetString("initiate-edge", "")

		if leftTmp == ACTION_OPENED_WINDOWS {
			_compizScale.SetString("initiate-edge", "TopLeft")
		}
	} else if index == ACTION_OPENED_WINDOWS {
		if desk.TopLeft == ACTION_OPENED_WINDOWS {
			desk.TopLeft = ACTION_LAUNCHER
			_compizIntegrated.SetString("command-11", _LAUNCHER_CMD)
			_compizCommand.SetString("run-command10-edge", "TopLeft")
		}

		_compizIntegrated.SetString("command-12", "")
		_compizCommand.SetString("run-command11-edge", "")
		_compizScale.SetString("initiate-edge", "BottomRight")
	} else if index == ACTION_LAUNCHER {
		if desk.TopLeft == ACTION_LAUNCHER {
			desk.TopLeft = ACTION_OPENED_WINDOWS
			_compizIntegrated.SetString("command-11", "")
			_compizCommand.SetString("run-command10-edge", "")
			_compizScale.SetString("initiate-edge", "TopLeft")
		}

		_compizIntegrated.SetString("command-12", _LAUNCHER_CMD)
		_compizCommand.SetString("run-command11-edge", "BottomRight")
	}
}

func NewDesktopManager() (*DesktopManager, error) {
	desk := &DesktopManager{}

	deskSettings := gio.NewSettings(_DESKTOP_SCHEMA)
	desk.ShowComputerIcon = property.NewGSettingsBoolProperty(desk, "ShowComputerIcon", deskSettings, "show-computer-icon")
	desk.ShowHomeIcon = property.NewGSettingsBoolProperty(desk, "ShowHomeIcon", deskSettings, "show-home-icon")
	desk.ShowTrashIcon = property.NewGSettingsBoolProperty(desk, "ShowTrashIcon", deskSettings, "show-trash-icon")
	desk.ShowDSCIcon = property.NewGSettingsBoolProperty(desk, "ShowDSCIcon", deskSettings, "show-dsc-icon")
	desk.DockMode = property.NewGSettingsStringProperty(desk, "DockMode", gio.NewSettings(_DOCK_SCHEMA), "hide-mode")

	InitCompizGSettings()
	ListenCompizGSettings(desk)
	GetEdgeAction(desk)

	return desk, nil
}

func InitCompizGSettings() {
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

func ListenCompizGSettings(desk *DesktopManager) {
	_compizIntegrated.Connect("changed::command-11", func(s *gio.Settings, name string) {
		_runCommand11 = s.GetString("command-11")
		GetEdgeAction(desk)
	})
	_compizIntegrated.Connect("changed::command-12", func(s *gio.Settings, name string) {
		_runCommand12 = s.GetString("command-12")
		GetEdgeAction(desk)
	})
	_compizCommand.Connect("changed::run-command10-edge", func(s *gio.Settings, name string) {
		_runCommandEdge10 = s.GetString("run-command10-edge")
		GetEdgeAction(desk)
	})
	_compizCommand.Connect("changed::run-command11-edge", func(s *gio.Settings, name string) {
		_runCommandEdge11 = s.GetString("run-command11-edge")
		GetEdgeAction(desk)
	})
	_compizScale.Connect("changed::initiate-edge", func(s *gio.Settings, name string) {
		_scale = s.GetString("initiate-edge")
		GetEdgeAction(desk)
	})
}

func GetEdgeAction(desk *DesktopManager) {
	if _runCommand11 == "" && _runCommandEdge10 == "" && _scale == "" {
		desk.TopLeft = ACTION_NONE
	} else if _scale == "TopLeft" && _runCommandEdge10 == "" {
		desk.TopLeft = ACTION_OPENED_WINDOWS
	} else if _runCommand11 == "launcher" && _runCommandEdge10 == "TopLeft" {
		desk.TopLeft = ACTION_LAUNCHER
	}

	if _runCommand12 == "" && _runCommandEdge11 == "" && _scale == "" {
		desk.BottomRight = ACTION_NONE
	} else if _scale == "BottomRight" && _runCommand12 == "" {
		desk.BottomRight = ACTION_OPENED_WINDOWS
	} else if _runCommand12 == "launcher" && _runCommandEdge11 == "BottomRight" {
		desk.BottomRight = ACTION_LAUNCHER
	}

	dbus.NotifyChange(desk, "TopLeft")
	dbus.NotifyChange(desk, "BottomRight")
}

func main() {
	desk, err := NewDesktopManager()
	if err != nil {
		return
	}
	dbus.InstallOnSession(desk)
	dlib.StartLoop()
}
