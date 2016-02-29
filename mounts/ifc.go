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
	"gir/gobject-2.0"
)

func (m *Manager) Eject(id string) error {
	info, err := m.DiskList.get(id)
	if err != nil {
		return err
	}

	switch info.object.Type {
	case diskObjVolume:
		m.ejectVolume(id, info.object)
	case diskObjMount:
		m.ejectMount(id, info.object)
	}
	return nil
}

func (m *Manager) Mount(id string) error {
	info, err := m.DiskList.get(id)
	if err != nil {
		return err
	}

	if info.object.Type != diskObjVolume {
		return fmt.Errorf("Not unmounted volume: %v", id)
	}

	m.mountVolume(id, info.object)
	return nil
}

func (m *Manager) Unmount(id string) error {
	info, err := m.DiskList.get(id)
	if err != nil {
		return err
	}

	if info.object.Type != diskObjMount {
		return fmt.Errorf("Not mounted volume: %v", id)
	}

	m.unmountMount(id, info.object)
	return nil
}

func (m *Manager) ejectVolume(id string, obj *diskObject) {
	volume := obj.getVolume()
	if volume == nil {
		return
	}

	if !volume.CanEject() {
		return
	}

	volume.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			if volume == nil || volume.Object.C == nil {
				return
			}

			_, err := volume.EjectFinish(ret)
			if err != nil {
				m.emitError(id, err.Error())
			}
		}))
}

func (m *Manager) ejectMount(id string, obj *diskObject) {
	mount := obj.getMount()
	if mount == nil {
		return
	}

	if !mount.CanEject() {
		return
	}

	mount.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			if mount == nil || mount.Object.C == nil {
				return
			}
			_, err := mount.EjectFinish(ret)
			if err != nil {
				m.emitError(id, err.Error())
			}
		}))
}

func (m *Manager) mountVolume(id string, obj *diskObject) {
	volume := obj.getVolume()
	if volume == nil {
		return
	}

	if !volume.CanMount() {
		return
	}

	volume.Mount(gio.MountMountFlagsNone, nil, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			if volume == nil || volume.Object.C == nil {
				return
			}
			_, err := volume.MountFinish(ret)
			if err != nil {
				m.emitError(id, err.Error())
			}
		}))
}

func (m *Manager) unmountMount(id string, obj *diskObject) {
	mount := obj.getMount()
	if mount == nil {
		return
	}

	if !mount.CanUnmount() {
		return
	}

	mount.Unmount(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			if mount == nil || mount.Object.C == nil {
				return
			}
			_, err := mount.UnmountFinish(ret)
			if err != nil {
				m.emitError(id, err.Error())
			}
		}))
}
