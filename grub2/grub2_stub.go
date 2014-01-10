package main

import (
	"dlib/dbus"
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

func (grub *Grub2) Load() {
	// TODO
	grub.readEntries()
	grub.readSettings()
}

func (grub *Grub2) Save() {
	// TODO
	grub.writeSettings()
	grub.generateGrubConfig()
}

func (grub *Grub2) GetEntryTitles() []string {
	// TODO
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.entryType == MENUENTRY {
			entryTitles = append(entryTitles, entry.getFullTitle())
		}
	}
	if len(entryTitles) == 0 {
		logError("there is no menu entry in %s", _GRUB_MENU)
	}
	return entryTitles
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
