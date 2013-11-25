package main

import (
	"dbus-gen/gsession"
	"dbus-gen/upower"
	"dlib/dbus"
)

type DShutdown struct{}

var (
	dShut  = gsession.GetSessionManager("/org/gnome/SessionManager")
	dPower = upower.GetUpower("/org/freedesktop/UPower")
)

func (shutdown *DShutdown) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.DShutdown",
		"/com/deepin/daemon/DShutdown",
		"com.deepin.daemon.DShutdown",
	}
}

func (shutdown *DShutdown) Logout() {
	dShut.Logout(1)
}

func (shutdown *DShutdown) Shutdown() bool {
	if !dShut.CanShutdown() {
		return false
	}

	dShut.RequestShutdown()
	return true
}

func (shutdown *DShutdown) Reboot() {
	dShut.RequestReboot()
}

func (shutdown *DShutdown) Suspend() bool {
	if !dPower.SuspendAllowed() {
		return false
	}
	dPower.Suspend()
	return true
}

func (shutdown *DShutdown) Hibernate() bool {
	if !dPower.HibernateAllowed() {
		return false
	}
	dPower.Hibernate()
	return true
}

func NewShutdown() *DShutdown {
	return &DShutdown{}
}

func main() {
	shutdown := NewShutdown()
	dbus.InstallOnSession(shutdown)
	select {}
}
