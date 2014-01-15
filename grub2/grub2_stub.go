package main

import (
	"dlib/dbus"
	"errors"
	"fmt"
	"strconv"
)

const (
	_GRUB2_DEST = "com.deepin.daemon.Grub2"
	_GRUB2_PATH = "/com/deepin/daemon/Grub2"
	_GRUB2_IFC  = "com.deepin.daemon.Grub2"
)

func (grub *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRUB2_DEST,
		_GRUB2_PATH,
		_GRUB2_IFC,
	}
}

func (grub *Grub2) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	switch name {
	case "DefaultEntry":
		grub.setDefaultEntry(grub.DefaultEntry)
	case "Timeout":
		grub.setTimeout(grub.Timeout)
	case "Gfxmode":
		grub.setGfxmode(grub.Gfxmode)
	case "Background":
		grub.setBackground(grub.Background)
	case "Theme":
		grub.setTheme(grub.Theme)
	}
}

func (grub *Grub2) Load() error {
	err := grub.readEntries()
	if err != nil {
		return err
	}
	err = grub.readSettings()
	if err != nil {
		return err
	}
	return nil
}

func (grub *Grub2) Save() error {
	// TODO
	err := grub.writeSettings()
	if err != nil {
		return err
	}
	grub.generateGrubConfig()
	if err != nil {
		return err
	}
	return nil
}

func (grub *Grub2) GetEntryTitles() ([]string, error) {
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.entryType == MENUENTRY {
			entryTitles = append(entryTitles, entry.getFullTitle())
		}
	}
	if len(entryTitles) == 0 {
		s := fmt.Sprintf("there is no menu entry in %s", _GRUB_MENU)
		logError(s)
		return entryTitles, errors.New(s)
	}
	return entryTitles, nil
}

func (grub *Grub2) setDefaultEntry(title string) {
	grub.settings["GRUB_DEFAULT"] = title
}

func (grub *Grub2) setTimeout(timeout int32) {
	timeoutStr := strconv.FormatInt(int64(timeout), 10)
	grub.settings["GRUB_TIMEOUT"] = timeoutStr
}

func (grub *Grub2) setGfxmode(gfxmode string) {
	grub.settings["GRUB_GFXMODE"] = gfxmode
}

func (grub *Grub2) setBackground(imageFile string) {
	grub.settings["GRUB_BACKGROUND"] = imageFile
}

func (grub *Grub2) setTheme(themeFile string) {
	grub.settings["GRUB_THEME"] = themeFile
}
