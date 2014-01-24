package main

import (
	"dlib/dbus"
	"errors"
	"fmt"
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
			logError("%v", err) // TODO
		}
	}()
	switch name {
	case "DefaultEntry":
		grub.setDefaultEntry(grub.DefaultEntry)
	case "Timeout":
		grub.setTimeout(grub.Timeout)
	case "Gfxmode":
		grub.setGfxmode(grub.Gfxmode)
	}
}

func (grub *Grub2) Save() error {
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

// Get all entry titles in level one.
func (grub *Grub2) GetSimpleEntryTitles() ([]string, error) {
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.parentSubMenu == nil && entry.entryType == MENUENTRY {
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
