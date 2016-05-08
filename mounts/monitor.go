/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mounts

import (
	"fmt"
	"gir/gio-2.0"
	"os/exec"
	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/dbus"
)

func (m *Manager) handleEvent() {
	m.monitor.Connect("volume-added", func(monitor *gio.VolumeMonitor,
		volume *gio.Volume) {
		soundutils.PlaySystemSound(soundutils.EventDevicePlug,
			"", false)
		info := newDiskInfoFromVolume(volume)
		logger.Debug("[Event] volume added:", info.Name, info.Type)
		if info.Type == DiskTypeRemovable && m.isAutoMount() {
			m.mountVolume(info.Id, info.object)
			return
		}
		logger.Debug("[Event] emit list changed signal")
		m.setPropDiskList(m.DiskList.add(info))
		dbus.Emit(m, "Changed", EventTypeVolumeAdded, info.Id)
	})

	m.monitor.Connect("volume-removed", func(monitor *gio.VolumeMonitor,
		volume *gio.Volume) {
		logger.Debug("[Event] volume removed")
		soundutils.PlaySystemSound(soundutils.EventDeviceUnplug,
			"", false)
		logger.Debug("[Event] emit list changed signal")
		oldInfos := m.DiskList.duplicate()
		m.refrashDiskList()
		id := findRemovedId(oldInfos, m.DiskList)
		if len(id) > 0 {
			dbus.Emit(m, "Changed", EventTypeVolumeRemoved, id)
		}
		oldInfos = nil
	})

	m.monitor.Connect("mount-added", func(monitor *gio.VolumeMonitor,
		mount *gio.Mount) {
		info := newDiskInfoFromMount(mount)
		logger.Debug("[Event] mount added:", info.Name, info.CanEject)
		if info.CanEject && m.isAutoOpen() {
			go exec.Command("/bin/sh", "-c",
				fmt.Sprintf("gvfs-open %s",
					info.MountPoint)).Run()
		}
		logger.Debug("[Event] emit list changed signal")
		m.setPropDiskList(m.DiskList.add(info))
		dbus.Emit(m, "Changed", EventTypeMountAdded, info.Id)
	})

	m.monitor.Connect("mount-removed", func(monitor *gio.VolumeMonitor,
		mount *gio.Mount) {
		logger.Debug("[Event] mount removed")
		logger.Debug("[Event] emit list changed signal")
		oldInfos := m.DiskList.duplicate()
		m.refrashDiskList()
		id := findRemovedId(oldInfos, m.DiskList)
		if len(id) > 0 {
			dbus.Emit(m, "Changed", EventTypeMountRemoved, id)
		}
		oldInfos = nil
	})
}

func (m *Manager) isAutoMount() bool {
	if m.setting == nil {
		return false
	}
	return m.setting.GetBoolean("automount")
}

func (m *Manager) isAutoOpen() bool {
	if m.setting == nil {
		return false
	}
	return m.setting.GetBoolean("automount-open")
}

func findRemovedId(oldInfos, newInfos DiskInfos) string {
	for _, info := range oldInfos {
		_, err := newInfos.get(info.Id)
		if err != nil {
			return info.Id
		}
	}
	return ""
}
