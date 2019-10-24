package grub2

import (
	"syscall"

	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"
	"pkg.deepin.io/lib/dbus1"
)

func (m *Grub2) enableShutdown() {
	if m.inhibitFd != -1 {
		err := syscall.Close(int(m.inhibitFd))
		if err != nil {
			logger.Infof("enable shutdown: fd:%d, err:%s\n", m.inhibitFd, err)
		} else {
			logger.Infof("enable shutdown")
		}
		m.inhibitFd = -1
	}
}

func (m *Grub2) preventShutdown() {
	if m.inhibitFd == -1 {
		fd, err := inhibit("shutdown", dbusServiceName,
			"Updating the system, please shut down or reboot later.")
		logger.Infof("prevent shutdown: fd:%v\n", fd)
		if err != nil {
			logger.Infof("prevent shutdown failed: fd:%v, err:%v\n", fd, err)
			return
		}
		m.inhibitFd = fd
	}
}

func inhibit(what, who, why string) (dbus.UnixFD, error) {
	systemConn, err := dbus.SystemBus()
	if err != nil {
		return 0, err
	}
	m := login1.NewManager(systemConn)
	return m.Inhibit(0, what, who, why, "block")
}
