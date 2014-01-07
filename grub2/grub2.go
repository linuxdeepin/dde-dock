package main

import (
	"dlib/dbus"
)

const (
	_GRUB2_SERV = "com.deepin.daemon.Grub2"
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
	grub := &Grub2{}
	grub.grubMenuFile = GRUB_MENU
	grub.grubConfigFile = GRUB_CONFIG
	return grub
}

func (grub *Grub2) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		_GRUB2_SERV,
		_GRUB2_PATH,
		_GRUB2_IFC,
	}
}

func (grub *Grub2) Init() {
	// TODO
}

func (grub *Grub2) Load() {
	// TODO
}

func (grub *Grub2) Save() {
	// TODO
}

func (grub *Grub2) GetEntries() []string {
	// TODO
	return []string{"a", "b", "c"}
}

func (grub *Grub2) SetDefaultEntry(entry string) {
	// TODO
}

func (grub *Grub2) GetDefaultEntry() string {
	// TODO
	return ""
}

func (grub *Grub2) SetTimeout(timeout int32) {
	// TODO
}

func (grub *Grub2) GetTimeout() int32 {
	// TODO
	return 0
}

func (grub *Grub2) SetBackground(imageFile string) {
	// TODO
}

func (grub *Grub2) GetBackground() string {
	// TODO
	return ""
}

func (grub *Grub2) SetTheme(themeFile string) {
	// TODO
}

func (grub *Grub2) GetTheme() string {
	// TODO
	return ""
}

func main() {
	grub := NewGrub2()
	err := dbus.InstallOnSession(grub)
	if err != nil {
		panic(err) // TODO
	}
	select {}
}
