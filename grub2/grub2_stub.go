package main

import (
	"dlib/dbus"
	"errors"
	"fmt"
)

func (grub *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Grub2",
		"/com/deepin/daemon/Grub2",
		"com.deepin.daemon.Grub2",
	}
}

func (grub *Grub2) OnPropertiesChanged(name string, oldv interface{}) {
	defer func() {
		if err := recover(); err != nil {
			logError("%v", err)
		}
	}()
	switch name {
	case "DefaultEntry":
		grub.setDefaultEntry(grub.DefaultEntry)
	case "Timeout":
		grub.setTimeout(grub.Timeout)
	}
}

// Get entry titles in level one.
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
