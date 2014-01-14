package main

import (
	"dbus/org/freedesktop/consolekit"
	"dbus/org/freedesktop/upower"
	"dbus/org/gnome/sessionmanager"
	"dlib/dbus"
	"fmt"
	"os/exec"
)

type Manager struct{}

const (
	_LOCK_EXEC         = "/usr/bin/dlock"
	_REBOOT_EXEC       = "/usr/lib/deepin-daemon/dpowerchoose --reboot"
	_LOGOUT_EXEC       = "/usr/lib/deepin-daemon/dpowerchoose --logout"
	_SHUTDOWN_EXEC     = "/usr/lib/deepin-daemon/dpowerchoose --shutdown"
	_POWER_CHOOSE_EXEC = "/usr/lib/deepin-daemon/dpowerchoose --choice"
)

var (
	dShut    *sessionmanager.SessionManager
	dConsole *consolekit.Manager
	dPower   *upower.Upower
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
	ok, _ := dPower.HibernateAllowed()
	return ok
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
	ok, err := dShut.IsInhibited(action)
	if err != nil {
		fmt.Println("IsInhibited Failed:", err)
		return true
	}

	return ok
}

func Init() {
	var err error

	dShut, err = sessionmanager.NewSessionManager("/org/gnome/SessionManager")
	if err != nil {
		fmt.Println("session: New SessionManager Failed:", err)
		return
	}

	dConsole, err = consolekit.NewManager("/org/freedesktop/ConsoleKit/Manager")
	if err != nil {
		fmt.Println("consolekit: New Manager Failed:", err)
		return
	}

	dPower, err = upower.NewUpower("/org/freedesktop/UPower")
	if err != nil {
		fmt.Println("upower: New Upower Failed:", err)
		return
	}
}

func main() {
	Init()
	shutdown := NewManager()
	dbus.InstallOnSession(shutdown)
	dbus.DealWithUnhandledMessage()
	select {}
}
