package main

import (
	"dbus/com/deepin/daemon/dock"
)

func loadAll() []string {
	manager, err := dock.NewDockedAppManager(
		"com.deepin.daemon.Dock",
		"/dde/dock/DockedAppManager",
	)
	if err != nil {
		LOGGER.Warning("get DockedAppManager failed", err)
		return make([]string, 0)
	}

	l, _ := manager.DockedAppList()
	return l

	return []string{"firefox"}
	// return []string{
	// 	"libreoffice-impress",
	// 	"gnome-system-monitor",
	// 	"google-chrome",
	// 	"d-feet",
	// 	"nautilus",
	// 	"deepin-music-player",
	// 	"deepin-terminal",
	// }
}
