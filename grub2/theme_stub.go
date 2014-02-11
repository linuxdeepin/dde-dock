package main

import (
	"dbus/com/deepin/dde/api/image"
	"dlib/dbus"
)

var dimg *image.Image

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

// TODO Set the background source file, then generate the background
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

	// backup background source file and convert it to png format
	dimg.ConvertToPNG(imageFile, theme.bgSrcFile)
	theme.generateBackground()

	return true
}

// TODO remove
func (theme *Theme) Reset() error {
	tplJsonData, err := theme.getThemeTplJsonData()
	if err != nil {
		return err
	}

	// theme.relBgFile = tplJsonData.DefaultTplValue.Background // TODO
	theme.makeBackground()
	theme.ItemColor = tplJsonData.ItemColor
	theme.SelectedItemColor = tplJsonData.SelectedItemColor

	// theme.customTheme()			// TODO
	return nil
}

// TODO move Generate background to fit the monitor resolution.
func (theme *Theme) generateBackground() {
	screenWidth, screenHeight := getScreenResolution()
	imgWidth, imgHeight, err := dimg.GetImageSize(theme.bgSrcFile)
	if err != nil {
		panic(err)
	}
	w, h := getImgClipSizeByResolution(screenWidth, screenHeight, imgWidth, imgHeight)
	err = dimg.ClipPNG(theme.bgSrcFile, theme.bgFile, 0, 0, w, h)
	if err != nil {
		panic(err)
	}
}
