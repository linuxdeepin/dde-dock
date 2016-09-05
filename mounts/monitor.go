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
		logger.Debug("[Event] volume added:", info.Name, info.Type, info.Id)
		if volume.ShouldAutomount() && m.isAutoMount() {
			m.mountVolume(info.Id, volume)
			volume.Unref()
			return
		}
		volume.Unref()
		m.refreshDiskList()
		dbus.Emit(m, "Changed", EventTypeVolumeAdded, info.Id)
	})

	m.monitor.Connect("volume-removed", func(monitor *gio.VolumeMonitor,
		volume *gio.Volume) {
		logger.Debug("[Event] volume removed:", getVolumeId(volume))
		volume.Unref()
		soundutils.PlaySystemSound(soundutils.EventDeviceUnplug,
			"", false)
		oldInfos := m.DiskList.duplicate()
		m.refreshDiskList()
		removed := findChangedId(oldInfos, m.DiskList, false)
		m.emitChanged(removed, EventTypeVolumeRemoved)
		oldInfos = nil
	})

	m.monitor.Connect("volume-changed", func(monitor *gio.VolumeMonitor,
		volume *gio.Volume) {
		id := getVolumeId(volume)
		logger.Debug("[Event] volume changed:", id)
		volume.Unref()
		oldInfos := m.DiskList.duplicate()
		m.refreshDiskList()
		added, removed := compareDiskList(oldInfos, m.DiskList)
		logger.Debug("[Event] after compare:", added, removed)
		m.emitChanged(added, EventTypeVolumeAdded)
		m.emitChanged(removed, EventTypeVolumeRemoved)
		if len(added) == 0 && len(removed) == 0 {
			dbus.Emit(m, "Changed", EventTypeVolumeChanged, id)
		}
	})

	m.monitor.Connect("mount-added", func(monitor *gio.VolumeMonitor,
		mount *gio.Mount) {
		info := newDiskInfoFromMount(mount)
		if info == nil {
			mount.Unref()
			return
		}
		logger.Debug("[Event] mount added:", info.Name, info.Id, info.CanEject)

		volume := mount.GetVolume()
		mount.Unref()
		var autoOpen bool = false
		if volume != nil && volume.Object.C != nil {
			if volume.ShouldAutomount() && m.isAutoOpen() {
				autoOpen = true
			}
			volume.Unref()
		}

		m.refreshDiskList()
		dbus.Emit(m, "Changed", EventTypeMountAdded, info.Id)

		if autoOpen {
			go exec.Command("/bin/sh", "-c",
				fmt.Sprintf("gvfs-open %s",
					info.MountPoint)).Run()
		}
	})

	m.monitor.Connect("mount-removed", func(monitor *gio.VolumeMonitor,
		mount *gio.Mount) {
		logger.Debug("[Event] mount removed:", getMountId(mount))
		if mount == nil || mount.Object.C == nil {
			logger.Warning("Invalid GMount Object")
			return
		}

		root := mount.GetRoot()
		point := root.GetUri()
		info, err := m.DiskList.getByMountPoint(point)
		root.Unref()
		if err != nil {
			// fixed phone device
			m.refreshDiskList()
			dbus.Emit(m, "Changed", EventTypeMountRemoved, getMountId(mount))
			mount.Unref()
			logger.Warning(err)
			return
		}
		mount.Unref()
		oldLen := len(m.DiskList)
		m.refreshDiskList()
		if oldLen != len(m.DiskList) {
			logger.Debug("Mount removed && volume removed")
			// fixed for smb
			dbus.Emit(m, "Changed", EventTypeMountRemoved, info.Id)
			return
		}

		logger.Debug("Only mount removed:", info.Id)
		dbus.Emit(m, "Changed", EventTypeMountRemoved, info.Id)
	})
}

func (m *Manager) emitChanged(ids []string, event int32) {
	for _, id := range ids {
		dbus.Emit(m, "Changed", event, id)
	}
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

func findChangedId(oldInfos, newInfos DiskInfos, added bool) []string {
	var ret []string
	if added {
		for _, info := range newInfos {
			_, err := oldInfos.get(info.Id)
			if err != nil {
				ret = append(ret, info.Id)
			}
		}
	} else {
		for _, info := range oldInfos {
			_, err := newInfos.get(info.Id)
			if err != nil {
				ret = append(ret, info.Id)
			}
		}
	}
	return ret
}

func compareDiskList(oldInfos, newInfos DiskInfos) ([]string, []string) {
	return findChangedId(oldInfos, newInfos, true),
		findChangedId(oldInfos, newInfos, false)
}
