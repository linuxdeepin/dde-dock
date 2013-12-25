package main

import (
	"dbus/org/freedesktop/consolekit"
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
	dShut, _    = sessionmanager.NewSessionManager("/org/gnome/SessionManager")
	dConsole, _ = consolekit.NewManager("/org/freedesktop/ConsoleKit/Manager")
	dPower, _   = upower.NewUpower("/org/freedesktop/UPower")
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.ShutdownManager",
		"/com/deepin/daemon/ShutdownManager",
		"com.deepin.daemon.ShutdownManager",
	}
}

func (m *Manager) CanLogout() bool {
	if IsInhibited(1) {
		return false
	}

	return true
}

func (m *Manager) Logout() {
	ExecCommand(_LOGOUT_EXEC)
}

func (m *Manager) RequestLogout() {
	dShut.Logout(1)
}

func (m *Manager) ForceLogout() {
	dShut.Logout(2)
}

func (shudown *Manager) CanShutdown() bool {
	if IsInhibited(1) {
		return false
	}

	return true
}

func (m *Manager) Shutdown() {
	ExecCommand(_SHUTDOWN_EXEC)
}

func (m *Manager) RequestShutdown() {
	dShut.RequestShutdown()
}

func (m *Manager) ForceShutdown() {
	dConsole.Stop()
}

func (shudown *Manager) CanReboot() bool {
	if IsInhibited(1) {
		return false
	}

	return true
}

func (m *Manager) Reboot() {
	ExecCommand(_REBOOT_EXEC)
}

func (m *Manager) RequestReboot() {
	dShut.RequestReboot()
}

func (m *Manager) ForceReboot() {
	dConsole.Restart()
}

func (m *Manager) CanSuspend() bool {
	if IsInhibited(4) {
		return false
	}

	return true
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

func IsInhibited(action uint32) bool {
	return dShut.IsInhibited(action)
}

func main() {
	shutdown := NewManager()
	dbus.InstallOnSession(shutdown)
	select {}
}
