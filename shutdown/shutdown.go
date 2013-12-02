package main

import (
	"dbus/org/freedesktop/upower"
	"dbus/org/gnome/sessionmanager"
	"dlib/dbus"
)

type DShutdown struct{}

var (
	dShut  = sessionmanager.GetSessionManager("/org/gnome/SessionManager")
	dPower = upower.GetUpower("/org/freedesktop/UPower")
)

func (shutdown *DShutdown) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.DShutdown",
		"/com/deepin/daemon/DShutdown",
		"com.deepin.daemon.DShutdown",
	}
}

func (shutdown *DShutdown) RequestLogout() {
	dShut.Logout(1)
}

func (shutdown *DShutdown) Logout() {
	dShut.Logout(0)
}

func (shudown *DShutdown) CanShutdown() bool {
	return dShut.CanShutdown()
}

func (shutdown *DShutdown) Shutdown() {
	dShut.Shutdown()
}

func (shutdown *DShutdown) RequestShutdown() {
	dShut.RequestShutdown()
}

func (shutdown *DShutdown) Reboot() {
	dShut.Reboot()
}

func (shutdown *DShutdown) RequestReboot() {
	dShut.RequestReboot()
}

func (shutdown *DShutdown) CanSuspend() bool {
	return dPower.SuspendAllowed()
}

func (shutdown *DShutdown) RequestSuspend() {
	dPower.Suspend()
}

func (shutdown *DShutdown) CanHibernate() bool {
	return dPower.HibernateAllowed()
}

func (shutdown *DShutdown) RequestHibernate() {
	dPower.Hibernate()
}

func NewShutdown() *DShutdown {
	return &DShutdown{}
}

func main() {
	shutdown := NewShutdown()
	dbus.InstallOnSession(shutdown)
	select {}
}
