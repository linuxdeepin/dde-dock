package main

import (
	"dlib"
	"dlib/dbus"
	"dlib/dbus/property"
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
	compizIntegrated *dlib.Settings
	compizCommand    *dlib.Settings
	compizScale      *dlib.Settings

	runCommand11     string
	runCommand12     string
	runCommandEdge10 string
	runCommandEdge11 string
	scale            string
)

type DesktopManager struct {
	ShowComputerIcon dbus.Property
	ShowHomeIcon     dbus.Property
	ShowTrashIcon    dbus.Property
	ShowDSCIcon      dbus.Property
	DockMode         dbus.Property
	TopLeft          int32
	BottomRight      int32
}

func (desk *DesktopManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DESKTOP_DEST, _DESKTOP_PATH, _DESKTOP_IFC}
}

func (desk *DesktopManager) SetTopLeftAction(index int32) {
	if index == ACTION_NONE {
		compizIntegrated.SetString("command-11", "")
		compizCommand.SetString("run-command10-edge", "")
		compizScale.SetString("initiate-edge", "")

		if desk.BottomRight == ACTION_OPENED_WINDOWS {
			compizScale.SetString("initiate-edge", "BottomRight")
		}
	} else if index == ACTION_OPENED_WINDOWS {
		if desk.BottomRight == ACTION_OPENED_WINDOWS {
			desk.BottomRight = ACTION_LAUNCHER
			compizIntegrated.SetString("command-12", _LAUNCHER_CMD)
			compizCommand.SetString("run-command11-edge", "BottomRight")
		}

		compizIntegrated.SetString("command-11", "")
		compizCommand.SetString("run-command10-edge", "")
		compizScale.SetString("initiate-edge", "TopLeft")
	} else if index == ACTION_LAUNCHER {
		if desk.BottomRight == ACTION_LAUNCHER {
			desk.BottomRight = ACTION_OPENED_WINDOWS
			compizIntegrated.SetString("command-12", "")
			compizCommand.SetString("run-command11-edge", "")
			compizScale.SetString("initiate-edge", "BottomRight")
		}

		compizIntegrated.SetString("command-11", _LAUNCHER_CMD)
		compizCommand.SetString("run-command10-edge", "TopLeft")
	}
}

func (desk *DesktopManager) SetBottomRightAction(index int32) {
	if index == ACTION_NONE {
		compizIntegrated.SetString("command-12", "")
		compizCommand.SetString("run-command11-edge", "")
		compizScale.SetString("initiate-edge", "")

		if desk.TopLeft == ACTION_OPENED_WINDOWS {
			compizScale.SetString("initiate-edge", "TopLeft")
		}
	} else if index == ACTION_OPENED_WINDOWS {
		if desk.TopLeft == ACTION_OPENED_WINDOWS {
			desk.TopLeft = ACTION_LAUNCHER
			compizIntegrated.SetString("command-11", _LAUNCHER_CMD)
			compizCommand.SetString("run-command10-edge", "TopLeft")
		}

		compizIntegrated.SetString("command-12", "")
		compizCommand.SetString("run-command11-edge", "")
		compizScale.SetString("initiate-edge", "BottomRight")
	} else if index == ACTION_LAUNCHER {
		if desk.TopLeft == ACTION_LAUNCHER {
			desk.TopLeft = ACTION_OPENED_WINDOWS
			compizIntegrated.SetString("command-11", "")
			compizCommand.SetString("run-command10-edge", "")
			compizScale.SetString("initiate-edge", "TopLeft")
		}

		compizIntegrated.SetString("command-12", _LAUNCHER_CMD)
		compizCommand.SetString("run-command11-edge", "BottomRight")
	}
}

func NewDesktopManager() (*DesktopManager, error) {
	desk := DesktopManager{}
	busType, err := dbus.SessionBus()
	if err != nil {
		return nil, err
	}

	deskSettings := dlib.NewSettings(_DESKTOP_SCHEMA)
	desk.ShowComputerIcon = property.NewGSettingsPropertyFull(
		deskSettings, "show-computer-icon", true, busType,
		_DESKTOP_PATH, _DESKTOP_IFC, "ShowComputerIcon")
	desk.ShowHomeIcon = property.NewGSettingsPropertyFull(
		deskSettings, "show-home-icon", true, busType, _DESKTOP_PATH,
		_DESKTOP_IFC, "ShowHomeIcon")
	desk.ShowTrashIcon = property.NewGSettingsPropertyFull(
		deskSettings, "show-trash-icon", true, busType, _DESKTOP_PATH,
		_DESKTOP_IFC, "ShowTrashIcon")
	desk.ShowDSCIcon = property.NewGSettingsPropertyFull(
		deskSettings, "show-dsc-icon", true, busType, _DESKTOP_PATH,
		_DESKTOP_IFC, "ShowDSCIcon")

	dockSettings := dlib.NewSettings(_DOCK_SCHEMA)
	desk.DockMode = property.NewGSettingsPropertyFull(dockSettings,
		"hide-mode", "", busType, _DESKTOP_PATH, _DESKTOP_IFC, "DockMode")
	InitCompizGSettings()
	ListenCompizGSettings(&desk)
	desk.TopLeft, desk.BottomRight = GetEdgeAction()

	return &desk, nil
}

func InitCompizGSettings() {
	compizIntegrated = dlib.NewSettings(_COMPIZ_INTEGRATED_SCHEMA)
	compizCommand = dlib.NewSettingsWithPath(_COMPIZ_COMMANDS_SCHEMA,
		_COMPIZ_COMMAND_PATH)
	compizScale = dlib.NewSettingsWithPath(_COMPIZ_SCALE_SCHEMA,
		_COMPIZ_SCALE_PATH)

	runCommand11 = compizIntegrated.GetString("command-11")
	runCommand12 = compizIntegrated.GetString("command-12")
	runCommandEdge10 = compizCommand.GetString("run-command10-edge")
	runCommandEdge11 = compizCommand.GetString("run-command11-edge")
	scale = compizScale.GetString("initiate-edge")
}

func ListenCompizGSettings(desk *DesktopManager) {
	compizIntegrated.Connect("changed::command-11", func(s *dlib.Settings, name string) {
		runCommand11 = s.GetString("command-11")
		desk.TopLeft, desk.BottomRight = GetEdgeAction()
	})
	compizIntegrated.Connect("changed::command-12", func(s *dlib.Settings, name string) {
		runCommand12 = s.GetString("command-12")
		desk.TopLeft, desk.BottomRight = GetEdgeAction()
	})
	compizCommand.Connect("changed::run-command10-edge", func(s *dlib.Settings, name string) {
		runCommandEdge10 = s.GetString("run-command10-edge")
		desk.TopLeft, desk.BottomRight = GetEdgeAction()
	})
	compizCommand.Connect("changed::run-command11-edge", func(s *dlib.Settings, name string) {
		runCommandEdge11 = s.GetString("run-command11-edge")
		desk.TopLeft, desk.BottomRight = GetEdgeAction()
	})
	compizScale.Connect("changed::initiate-edge", func(s *dlib.Settings, name string) {
		scale = s.GetString("initiate-edge")
		desk.TopLeft, desk.BottomRight = GetEdgeAction()
	})
}

func GetEdgeAction() (topLeft, bottomRight int32) {
	if runCommand11 == "" && runCommandEdge10 == "" && scale == "" {
		topLeft = ACTION_NONE
	} else if scale == "TopLeft" && runCommandEdge10 == "" {
		topLeft = ACTION_OPENED_WINDOWS
	} else if runCommand11 == "launcher" && runCommandEdge10 == "TopLeft" {
		topLeft = ACTION_LAUNCHER
	}

	if runCommand12 == "" && runCommandEdge11 == "" && scale == "" {
		bottomRight = ACTION_NONE
	} else if scale == "BottomRight" && runCommand12 == "" {
		bottomRight = ACTION_OPENED_WINDOWS
	} else if runCommand12 == "launcher" && runCommandEdge11 == "BottomRight" {
		bottomRight = ACTION_LAUNCHER
	}

	return topLeft, bottomRight
}

func main() {
	go dlib.StartLoop()
	desk, err := NewDesktopManager()
	if err != nil {
		return
	}
	dbus.InstallOnSession(desk)
	select {}
}
