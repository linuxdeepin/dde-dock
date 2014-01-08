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

const (
	GRUB_MENU   = "/boot/grub/grub.cfg"
	GRUB_CONFIG = "/etc/default/grub"
)

type Grub2 struct {
	grubMenuFile   string
	grubConfigFile string
	entries        []string
	settings       map[string]string
}

func NewGrub2() *Grub2 {
	// TODO
	grub := &Grub2{}
	return grub
}

func (grub *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRUB2_DEST,
		_GRUB2_PATH,
		_GRUB2_IFC,
	}
}

func (grub *Grub2) Init() {
	// TODO
	grub.grubMenuFile = GRUB_MENU
	grub.grubConfigFile = GRUB_CONFIG
}

func (grub *Grub2) Load() {
	// TODO
}

func (grub *Grub2) Save() {
	// TODO
}

func (grub *Grub2) GetEntries() []string {
	// TODO
	// return []string{"a", "b", "c"}
	return grub.entries
}

func (grub *Grub2) SetDefaultEntry(index uint32) {
	indexStr := strconv.FormatInt(int64(index), 10)
	grub.settings["GRUB_DEFAULT"] = indexStr
}

func (grub *Grub2) GetDefaultEntry() uint32 {
	index, err := strconv.ParseInt(grub.settings["GRUB_DEFAULT"], 10, 32)
	if err != nil {
		logError(fmt.Sprintf(`valid value, settings["GRUB_DEFAULT"]=%s`, grub.settings["GRUB_DEFAULT"])) // TODO
		return 0
	}
	return uint32(index)
}

func (grub *Grub2) SetTimeout(timeout int32) {
	timeoutStr := strconv.FormatInt(int64(timeout), 10)
	grub.settings["GRUB_TIMEOUT"] = timeoutStr
}

func (grub *Grub2) GetTimeout() int32 {
	timeout, err := strconv.ParseInt(grub.settings["GRUB_TIMEOUT"], 10, 32)
	if err != nil {
		logError(fmt.Sprintf(`valid value, settings["GRUB_TIMEOUT"]=%s`, grub.settings["GRUB_TIMEOUT"])) // TODO
		return 5
	}
	return int32(timeout)
}

func (grub *Grub2) SetGfxmode(gfxmode string) {
	grub.settings["GRUB_GFXMODE"] = gfxmode
}

func (grub *Grub2) GetGfxmode() string {
	return grub.settings["GRUB_GFXMODE"]
}

func (grub *Grub2) SetBackground(imageFile string) {
	grub.settings["GRUB_BACKGROUND"] = imageFile
}

func (grub *Grub2) GetBackground() string {
	return grub.settings["GRUB_BACKGROUND"]
}

func (grub *Grub2) SetTheme(themeFile string) {
	grub.settings["GRUB_THEME"] = themeFile
}

func (grub *Grub2) GetTheme() string {
	return grub.settings["GRUB_THEME"]
}

func main() {
	grub := NewGrub2()
	err := dbus.InstallOnSession(grub)
	if err != nil {
		panic(err) // TODO
	}
	select {}
}
