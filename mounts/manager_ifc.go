/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mounts

import (
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"gir/gio-2.0"
	"gir/gobject-2.0"
)

// Eject disk.
//
// uuid: get from DiskList
func (m *Manager) DeviceEject(uuid string) (bool, error) {
	value := m.getDiskCache(uuid)
	if value == nil {
		var reason = fmt.Sprintf("Eject failed: invalid id '%s'", uuid)
		dbus.Emit(m, "Error", uuid, reason)
		return false, fmt.Errorf(reason)
	}

	switch value.Type {
	case diskTypeVolume:
		volume := value.Obj.(*gio.Volume)
		m.ejectVolume(uuid, volume)
	case diskTypeMount:
		mount := value.Obj.(*gio.Mount)
		m.ejectMount(uuid, mount)
	}

	return true, nil
}

// Mount disk.
func (m *Manager) DeviceMount(uuid string) (bool, error) {
	value := m.getDiskCache(uuid)
	if value == nil {
		var reason = fmt.Sprintf("Mount failed: invalid id '%s'", uuid)
		dbus.Emit(m, "Error", uuid, reason)
		return false, fmt.Errorf(reason)
	}

	switch value.Type {
	case diskTypeVolume:
		volume := value.Obj.(*gio.Volume)
		m.mountVolume(uuid, volume)
	case diskTypeMount:
		mount := value.Obj.(*gio.Mount)
		m.remountMount(uuid, mount)
	}

	return true, nil
}

// Unmount disk.
func (m *Manager) DeviceUnmount(uuid string) (bool, error) {
	value := m.getDiskCache(uuid)
	if value == nil {
		var reason = fmt.Sprintf("Unmount failed: invalid id '%s'", uuid)
		dbus.Emit(m, "Error", uuid, reason)
		return false, fmt.Errorf(reason)
	}

	switch value.Type {
	case diskTypeMount:
		mount := value.Obj.(*gio.Mount)
		m.unmountMount(uuid, mount)
	}

	return true, nil
}

func (m *Manager) ejectVolume(uuid string, volume *gio.Volume) {
	volume.Eject(gio.MountUnmountFlagsNone, nil,
		gio.AsyncReadyCallback(
			func(o *gobject.Object, res *gio.AsyncResult) {
				if volume == nil || volume.Object.C == nil {
					return
				}
				_, err := volume.EjectFinish(res)
				if err != nil {
					dbus.Emit(m, "Error", uuid, err.Error())
				}
			}))
}

func (m *Manager) ejectMount(uuid string, mount *gio.Mount) {
	mount.Eject(gio.MountUnmountFlagsNone, nil,
		gio.AsyncReadyCallback(
			func(o *gobject.Object, res *gio.AsyncResult) {
				if mount == nil || mount.Object.C == nil {
					return
				}
				_, err := mount.EjectFinish(res)
				if err != nil {
					dbus.Emit(m, "Error", uuid, err.Error())
				}
			}))
}

func (m *Manager) mountVolume(uuid string, volume *gio.Volume) {
	volume.Mount(gio.MountMountFlagsNone, nil, nil,
		gio.AsyncReadyCallback(
			func(o *gobject.Object, res *gio.AsyncResult) {
				if volume == nil || volume.Object.C == nil {
					return
				}
				_, err := volume.MountFinish(res)
				if err != nil {
					dbus.Emit(m, "Error", uuid, err.Error())
				}
			}))
}

func (m *Manager) remountMount(uuid string, mount *gio.Mount) {
	mount.Remount(gio.MountMountFlagsNone, nil, nil,
		gio.AsyncReadyCallback(
			func(o *gobject.Object, res *gio.AsyncResult) {
				if mount == nil || mount.Object.C == nil {
					return
				}
				_, err := mount.RemountFinish(res)
				if err != nil {
					dbus.Emit(m, "Error", uuid, err.Error())
				}
			}))
}

func (m *Manager) unmountMount(uuid string, mount *gio.Mount) {
	mount.Unmount(gio.MountUnmountFlagsNone, nil,
		gio.AsyncReadyCallback(
			func(o *gobject.Object, res *gio.AsyncResult) {
				if mount == nil || mount.Object.C == nil {
					return
				}
				_, err := mount.UnmountFinish(res)
				if err != nil {
					dbus.Emit(m, "Error", uuid, err.Error())
				}
			}))
}
