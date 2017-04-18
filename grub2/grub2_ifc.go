package grub2

import (
	"encoding/json"
	"errors"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"strings"
)

const (
	DBusDest      = "com.deepin.daemon.Grub2"
	DBusObjPath   = "/com/deepin/daemon/Grub2"
	DBusInterface = "com.deepin.daemon.Grub2"

	timeoutMax = 10
)

// GetDBusInfo implements interface of dbus.DBusObject.
func (_ *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       DBusDest,
		ObjectPath: DBusObjPath,
		Interface:  DBusInterface,
	}
}

// GetSimpleEntryTitles return entry titles only in level one and will
// filter out some useless entries such as sub-menus and "memtest86+".
func (grub *Grub2) GetSimpleEntryTitles() ([]string, error) {
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.parentSubMenu == nil && entry.entryType == MENUENTRY {
			title := entry.getFullTitle()
			if !strings.Contains(title, "memtest86+") {
				entryTitles = append(entryTitles, title)
			}
		}
	}
	if len(entryTitles) == 0 {
		err := fmt.Errorf("there is no menu entry in %s", grubScriptFile)
		return entryTitles, err
	}
	return entryTitles, nil
}

func (grub *Grub2) GetAvailableResolutions() (modesJSON string, err error) {
	// TODO
	type mode struct {
		Text, Value string
	}
	var modes []mode
	modes = append(modes, mode{Text: "Auto", Value: "auto"})
	modes = append(modes, mode{Text: "1024x768", Value: "1024x768"})
	modes = append(modes, mode{Text: "800x600", Value: "800x600"})
	data, err := json.Marshal(modes)
	modesJSON = string(data)
	return
}

func (g *Grub2) SetDefaultEntry(dbusMsg dbus.DMessage, v string) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	var titles []string
	titles, err = g.GetSimpleEntryTitles()
	if err != nil {
		return
	}

	if !isStringInArray(v, titles) {
		return fmt.Errorf("invalid entry %q", v)
	}

	if g.DefaultEntry == v {
		return
	}
	g.DefaultEntry = v
	dbus.NotifyChange(g, "DefaultEntry")
	g.saveConfig()
	return
}

func (g *Grub2) SetEnableTheme(dbusMsg dbus.DMessage, v bool) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	if g.EnableTheme == v {
		return
	}
	g.EnableTheme = v
	dbus.NotifyChange(g, "EnableTheme")
	g.saveConfig()
	return
}

func (g *Grub2) SetResolution(dbusMsg dbus.DMessage, v string) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	err = checkResolution(v)
	if err != nil {
		return
	}

	if g.Resolution == v {
		return
	}
	g.Resolution = v
	dbus.NotifyChange(g, "Resolution")
	g.saveConfig()
	return
}

func (g *Grub2) SetTimeout(dbusMsg dbus.DMessage, v uint32) (err error) {
	err = checkAuth(dbusMsg)
	if err != nil {
		return
	}

	if v > timeoutMax {
		return errors.New("exceeded the maximum value 10")
	}

	if g.Timeout == v {
		return
	}
	g.Timeout = v
	dbus.NotifyChange(g, "Timeout")
	g.saveConfig()
	return
}

// Reset reset all configuretion.
func (g *Grub2) Reset() {
	g.theme.reset()

	config := NewConfig()
	config.UseDefault()

	var changed bool
	if g.Timeout != config.Timeout {
		changed = true
		g.Timeout = config.Timeout
		dbus.NotifyChange(g, "Timeout")
	}

	if g.EnableTheme != config.EnableTheme {
		changed = true
		g.EnableTheme = config.EnableTheme
		dbus.NotifyChange(g, "EnableTheme")
	}

	if g.Resolution != config.Resolution {
		changed = true
		g.Resolution = config.Resolution
		dbus.NotifyChange(g, "Resolution")
	}

	cfgDefaultEntry, _ := g.defaultEntryIdx2Str(config.DefaultEntry)
	if g.DefaultEntry != cfgDefaultEntry {
		changed = true
		g.DefaultEntry = cfgDefaultEntry
		dbus.NotifyChange(g, "DefaultEntry")
	}

	if changed {
		g.saveConfig()
	}

}
