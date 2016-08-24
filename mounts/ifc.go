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
	"os/exec"
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
		info := newDiskInfoFromMount(mount)
		if info != nil && len(info.MountPoint) != 0 {
			m.ejectMount(info)
			return nil
		}
	}

	volume := m.getVolumeById(id)
	if volume != nil {
		info := newDiskInfoFromVolume(volume)
		if info != nil {
			m.ejectVolume(info)
			return nil
		}
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
		info := newDiskInfoFromVolume(volume)
		if info != nil {
			m.mountVolume(info)
			return nil
		}
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
		info := newDiskInfoFromMount(mount)
		if info != nil {
			m.unmountMount(info)
			return nil
		}
	}

	err := fmt.Errorf("Not found GMount by '%s'", id)
	m.emitError(id, err.Error())
	return err
}

func (m *Manager) ejectVolume(info *DiskInfo) {
	logger.Debugf("ejectVolume info: %#v", info)
	go func() {
		err := doDiskOperation("eject", info.Path)
		if err != nil {
			logger.Warning("[ejectVolume] failed:", info.Path, err)
			m.emitError(info.Id, err.Error())
		}
	}()
}

func (m *Manager) ejectMount(info *DiskInfo) {
	logger.Debugf("ejectMount info: %#v", info)
	go func() {
		err := doDiskOperation("eject", info.MountPoint)
		if err != nil {
			logger.Warning("[ejectMount] failed:", info.MountPoint, err)
			m.emitError(info.Id, err.Error())
		}
	}()
}

func (m *Manager) mountVolume(info *DiskInfo) {
	logger.Debugf("mountVolume info: %#v", info)
	go func() {
		err := doDiskOperation("mount", info.Path)
		if err != nil {
			logger.Warning("[mountVolume] failed:", info.Path, err)
			m.emitError(info.Id, err.Error())
		}
	}()
}

func (m *Manager) unmountMount(info *DiskInfo) {
	logger.Debugf("unmountMount info: %#v", info)
	go func() {
		err := doDiskOperation("unmount", info.MountPoint)
		if err != nil {
			logger.Warning("[unmountMount] failed:", info.MountPoint, err)
			m.emitError(info.Id, err.Error())
			return
		}
		go m.sendNotify(info.Icon, "",
			fmt.Sprintf(Tr("%s removed successfully"), info.MountPoint))
	}()
}

func doDiskOperation(ty, path string) error {
	var args []string
	switch ty {
	case "eject":
		args = append(args, "-e")
	case "mount":
		args = append(args, []string{"-m", "-d"}...)
	case "unmount":
		args = append(args, "-u")
	}
	args = append(args, path)
	out, err := exec.Command("gvfs-mount", args...).CombinedOutput()
	if err != nil {
		if len(out) != 0 {
			return fmt.Errorf("%s", string(out))
		}
		return err
	}
	return nil
}
