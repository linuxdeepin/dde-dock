package main

import (
	"dlib/dbus"
	"fmt"
)

const (
	_GRUB2_THEME_DEST = "com.deepin.daemon.Grub2"
	_GRUB2_THEME_PATH = "/com/deepin/daemon/Grub2/Theme"
	_GRUB2_THEME_IFC  = "com.deepin.daemon.Grub2.Theme"
)

func (theme *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRUB2_THEME_DEST,
		fmt.Sprintf("%s/%s", _GRUB2_THEME_PATH, theme.Name),
		_GRUB2_THEME_IFC,
	}
}

func (theme *Theme) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logError("%v", err) // TODO
		}
	}()
	switch name {
	case "Background":
		theme.setBackground(theme.Background)
	case "ItemColor":
		theme.setItemColor(theme.ItemColor)
	case "SelectedItemColor":
		theme.setSelectedItemColor(theme.SelectedItemColor)
	}
}

// TODO
func (t *Theme) Reset() error {
	return nil
}
