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
	. "pkg.deepin.io/lib/gettext"
)

func (m *Manager) ListDisk() DiskInfos {
	m.refreshLocker.Lock()
	defer m.refreshLocker.Unlock()
	return m.DiskList
}

func (m *Manager) QueryDisk(id string) (*DiskInfo, error) {
	m.refreshLocker.Lock()
	defer m.refreshLocker.Unlock()
	return m.DiskList.get(id)
}

func (m *Manager) Eject(id string) error {
	m.refreshLocker.Lock()
	defer m.refreshLocker.Unlock()

	mount := m.getMountById(id)
	if mount != nil {
		m.ejectMount(id, mount)
		return nil
	}

	volume := m.getVolumeById(id)
	if volume != nil {
		m.ejectVolume(id, volume)
		return nil
	}

	err := fmt.Errorf("Invalid disk id: %v", id)
	m.emitError(id, err.Error())
	return err
}

func (m *Manager) Mount(id string) error {
	m.refreshLocker.Lock()
	defer m.refreshLocker.Unlock()

	volume := m.getVolumeById(id)
	if volume != nil {
		m.mountVolume(id, volume)
		return nil
	}

	err := fmt.Errorf("Not found GVolume by '%s'", id)
	m.emitError(id, err.Error())
	return err
}

func (m *Manager) Unmount(id string) error {
	m.refreshLocker.Lock()
	defer m.refreshLocker.Unlock()

	mount := m.getMountById(id)
	if mount != nil {
		m.unmountMount(id, mount)
		mount.Unref()
		return nil
	}

	err := fmt.Errorf("Not found GMount by '%s'", id)
	m.emitError(id, err.Error())
	return err
}

func (m *Manager) ejectVolume(id string, volume *gio.Volume) {
	volume.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			if volume == nil || volume.Object.C == nil {
				return
			}

			_, err := volume.EjectFinish(ret)
			volume.Unref()
			if err != nil {
				m.emitError(id, err.Error())
			}
		}))
}

func (m *Manager) ejectMount(id string, mount *gio.Mount) {
	mount.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			if mount == nil || mount.Object.C == nil {
				return
			}
			_, err := mount.EjectFinish(ret)
			mount.Unref()
			if err != nil {
				m.emitError(id, err.Error())
			}
		}))
}

func (m *Manager) mountVolume(id string, volume *gio.Volume) {
	volume.Mount(gio.MountMountFlagsNone, nil, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			if volume == nil || volume.Object.C == nil {
				return
			}
			_, err := volume.MountFinish(ret)
			volume.Unref()
			if err != nil {
				m.emitError(id, err.Error())
			}
		}))
}

func (m *Manager) unmountMount(id string, mount *gio.Mount) {
	logger.Debugf("unmountMount id %q, mount: %v", id, mount)
	mount.UnmountWithOperation(gio.MountUnmountFlagsNone, nil, nil, gio.AsyncReadyCallback(
		func(o *gobject.Object, ret *gio.AsyncResult) {
			mount := gio.ToMount(o)
			logger.Debug("UnmountWithOperation AsyncReadyCallback")

			_, err := mount.UnmountWithOperationFinish(ret)
			if err != nil {
				m.emitError(id, err.Error())
				return
			}
			name := mount.GetName()
			gicon := mount.GetIcon()
			icon := getIconFromGIcon(gicon)
			gicon.Unref()

			go m.sendNotify(icon, "",
				fmt.Sprintf(Tr("%s removed successfully"), name))
		}))
}
