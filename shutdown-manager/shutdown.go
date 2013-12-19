package main

import (
	"dbus/org/freedesktop/upower"
	"dbus/org/gnome/sessionmanager"
	"dlib/dbus"
	"os/exec"
)

type Manager struct{}

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

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.ShutdownManager",
		"/com/deepin/daemon/ShutdownManager",
		"com.deepin.daemon.ShutdownManager",
	}
}

func (m *Manager) RequestLogout() {
	dShut.Logout(1)
}

func (m *Manager) Logout() {
	ExecCommand(_LOGOUT_EXEC)
}

func (shudown *Manager) CanShutdown() bool {
	return dShut.CanShutdown()
}

func (m *Manager) Shutdown() {
	ExecCommand(_SHUTDOWN_EXEC)
}

func (m *Manager) RequestShutdown() {
	dShut.RequestShutdown()
}

func (m *Manager) Reboot() {
	ExecCommand(_REBOOT_EXEC)
}

func (m *Manager) RequestReboot() {
	dShut.RequestReboot()
}

func (m *Manager) CanSuspend() bool {
	return dPower.SuspendAllowed()
}

func (m *Manager) RequestSuspend() {
	dPower.Suspend()
}

func (m *Manager) CanHibernate() bool {
	return dPower.HibernateAllowed()
}

func (m *Manager) RequestHibernate() {
	dPower.Hibernate()
}

func (m *Manager) RequestLock() {
	ExecCommand(_LOCK_EXEC)
}

func (m *Manager) PowerOffChoose() {
	ExecCommand(_POWER_CHOOSE_EXEC)
}

func NewManager() *Manager {
	return &Manager{}
}

func ExecCommand(cmd string) {
	cmdExec := exec.Command(cmd)
	cmdExec.Run()
}

func main() {
	shutdown := NewManager()
	dbus.InstallOnSession(shutdown)
	select {}
}
