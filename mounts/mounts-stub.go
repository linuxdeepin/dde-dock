/**
 * Copyright (c) 2011 ~ 2013 Deepin, Inc.
 *               2011 ~ 2013 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"dlib/dbus"
	"dlib/gio-2.0"
	"dlib/logger"
)

const (
	DISK_INFO_DEST = "com.deepin.daemon.DiskMount"
	DISK_INFO_PATH = "/com/deepin/daemon/DiskMount"
	DISK_INFO_IFC  = "com.deepin.daemon.DiskMount"
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DISK_INFO_DEST,
		DISK_INFO_PATH,
		DISK_INFO_IFC,
	}
}

func (m *Manager) setPropName(name string) {
	switch name {
	case "DiskList":
		m.DiskList = getDiskInfoList()
		dbus.NotifyChange(m, name)
	default:
		logger.Printf("'%s': invalid mount property\n")
	}
}

func (m *Manager) listenSignalChanged() {
	monitor.Connect("mount-added", func(volumeMonitor *gio.VolumeMonitor, mount *gio.Mount) {
		// Judge whether the property 'mount_and_open' set true
		// if true, open the device use exec.Command("xdg-open", "device").Run()
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})
	monitor.Connect("mount-removed", func(volumeMonitor *gio.VolumeMonitor, mount *gio.Mount) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})
	monitor.Connect("mount-changed", func(volumeMonitor *gio.VolumeMonitor, mount *gio.Mount) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})

	monitor.Connect("volume-added", func(volumeMonitor *gio.VolumeMonitor, volume *gio.Volume) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})
	monitor.Connect("volume-removed", func(volumeMonitor *gio.VolumeMonitor, volume *gio.Volume) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})
	monitor.Connect("volume-changed", func(volumeMonitor *gio.VolumeMonitor, volume *gio.Volume) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})

	monitor.Connect("drive-disconnected", func(volumeMonitor *gio.VolumeMonitor, drive *gio.Drive) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})
	monitor.Connect("drive-connected", func(volumeMonitor *gio.VolumeMonitor, drive *gio.Drive) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})
	monitor.Connect("drive-changed", func(volumeMonitor *gio.VolumeMonitor, drive *gio.Drive) {
		m.setPropName("DiskList")
		//printDiskInfo(m.DiskList)
	})
}
