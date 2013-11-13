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
)

/*var defaults DesktopManager*/

type DesktopManager struct {
	ShowComputerIcon *property.GSettingsProperty
	ShowHomeIcon     *property.GSettingsProperty
	ShowTrashIcon    *property.GSettingsProperty
	ShowDSCIcon      *property.GSettingsProperty
	LeftEdgeAction   *property.GSettingsProperty
	RightEdgeAction  *property.GSettingsProperty
	DockMode         *property.GSettingsProperty
}

func (desk *DesktopManager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{_DESKTOP_DEST, _DESKTOP_PATH, _DESKTOP_IFC}
}

/*func (desk *DesktopManager) reset(propName string) {*/
/*}*/

/*
func InitDefaults() {
	defaults.ShowComputerIcon = true
	defaults.ShowHomeIcon = true
	defaults.ShowTrashIcon = true
	defaults.ShowDSCIcon = true
	defaults.DockMode = 0
	defaults.LeftEdgeAction = 10
	defaults.RightEdgeAction = 10
}*/

func GetDesktopSettings() DesktopManager {
	desk := DesktopManager{}
	busType, _ := dbus.SessionBus()

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

	gridSettings := dlib.NewSettingsWithPath("org.compiz.grid",
		"/org/compiz/profiles/deepin/plugins/grid/")
	desk.LeftEdgeAction = property.NewGSettingsPropertyFull(gridSettings,
		"left-edge-action", int32(0), busType, _DESKTOP_PATH,
		_DESKTOP_IFC, "LeftEdgeAction")
	desk.RightEdgeAction = property.NewGSettingsPropertyFull(gridSettings,
		"right-edge-action", int32(0), busType, _DESKTOP_PATH,
		_DESKTOP_IFC, "RightEdgeAction")

	return desk
}

func main() {
	desk := GetDesktopSettings()
	dbus.InstallOnSession(&desk)
	select {}
}
