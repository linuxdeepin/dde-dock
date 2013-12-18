package main

import (
	"dbus/org/freedesktop/upower"
	"dbus/org/gnome/sessionmanager"
	"dlib/dbus"
	"os/exec"
)

type DShutdown struct{}

const (
	_LOCK_EXEC         = "/usr/bin/dlock"
	_REBOOT_EXEC       = "/usr/lib/deepin-daemon/dreboot"
	_LOGOUT_EXEC       = "/usr/lib/deepin-daemon/dlogout"
	_SHUTDOWN_EXEC     = "/usr/lib/deepin-daemon/dshutdown"
	_POWER_CHOOSE_EXEC = "/usr/lib/deepin-daemon/dpowerchoose"
)

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
	ExecCommand(_LOGOUT_EXEC)
}

func (shudown *DShutdown) CanShutdown() bool {
	return dShut.CanShutdown()
}

func (shutdown *DShutdown) Shutdown() {
	ExecCommand(_SHUTDOWN_EXEC)
}

func (shutdown *DShutdown) RequestShutdown() {
	dShut.RequestShutdown()
}

func (shutdown *DShutdown) Reboot() {
	ExecCommand(_REBOOT_EXEC)
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

func (shutdown *DShutdown) RequestLock() {
	ExecCommand(_LOCK_EXEC)
}

func (shutdown *DShutdown) PowerOffChoose() {
	ExecCommand(_POWER_CHOOSE_EXEC)
}

func NewShutdown() *DShutdown {
	return &DShutdown{}
}

func ExecCommand(cmd string) {
	cmdExec := exec.Command(cmd)
	cmdExec.Run()
}

func main() {
	shutdown := NewShutdown()
	dbus.InstallOnSession(shutdown)
	select {}
}
