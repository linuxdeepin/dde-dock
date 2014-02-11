package main

import (
	"dlib/dbus"
)

func (theme *Theme) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Grub2",
		"/com/deepin/daemon/Grub2/Theme",
		"com.deepin.daemon.Grub2.Theme",
	}
}

func (theme *Theme) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logError("%v", err)
		}
	}()
	switch name {
	case "ItemColor":
		theme.setItemColor(theme.ItemColor)
	case "SelectedItemColor":
		theme.setSelectedItemColor(theme.SelectedItemColor)
	}
}

// Set the background source file, then generate the background
// to fit the screen resolution, support png and jpeg image format
func (theme *Theme) SetBackgroundSourceFile(imageFile string) bool {
	// check image size
	w, h, err := dimg.GetImageSize(imageFile)
	if err != nil {
		return false
	}
	if w < 800 || h < 600 {
		logError("image size too small") // TODO
		return false
	}

	// backup background source file
	_, err = copyFile(imageFile, theme.bgSrcFile)
	if err != nil {
		return false
	}

	theme.generateBackground()

	return true
}
