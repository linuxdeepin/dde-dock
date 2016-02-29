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
	"regexp"
	"strings"
	"sync"
)

const (
	DiskTypeNative    string = "native"
	DiskTypeRemovable        = "removable"
	DiskTypeNetwork          = "network"
)

const (
	volumeKindUnix = "unix-device"
	volumeKindUUID = "uuid"

	fsAttrSize = "filesystem::size"
	fsAttrUsed = "filesystem::used"

	mtpDeviceIcon = "drive-removable-media-mtp"
)

const (
	diskObjVolume int = iota + 1
	diskObjMount
)

var diskLocker sync.Mutex

type DiskInfo struct {
	Id         string
	Name       string
	Type       string
	Path       string
	MountPoint string
	Icon       string

	CanUnmount bool
	CanEject   bool

	Used  uint64
	Total uint64

	object *diskObject
}
type DiskInfos []*DiskInfo

type diskObject struct {
	Object interface{}
	Type   int
}

func (obj *diskObject) getVolume() *gio.Volume {
	if obj.Type != diskObjVolume {
		return nil
	}

	volume, ok := obj.Object.(*gio.Volume)
	if !ok {
		return nil
	}
	return volume
}

func (obj *diskObject) getMount() *gio.Mount {
	if obj.Type != diskObjMount {
		return nil
	}

	mount, ok := obj.Object.(*gio.Mount)
	if !ok {
		return nil
	}
	return mount
}

func (m *Manager) listDisk() DiskInfos {
	diskLocker.Lock()
	defer diskLocker.Unlock()
	var infos DiskInfos
	for _, volume := range m.monitor.GetVolumes() {
		mount := volume.GetMount()
		if mount != nil {
			mount.Unref()
			continue
		}

		info := newDiskInfoFromVolume(volume)
		infos = append(infos, info)
	}

	for _, mount := range m.monitor.GetMounts() {
		info := newDiskInfoFromMount(mount)
		infos = append(infos, info)
	}

	return infos
}

func (infos DiskInfos) get(id string) (*DiskInfo, error) {
	diskLocker.Lock()
	defer diskLocker.Unlock()
	for _, info := range infos {
		if info.Id == id {
			return info, nil
		}
	}
	return nil, fmt.Errorf("Invalid disk id: %v", id)
}

func (infos DiskInfos) add(info *DiskInfo) DiskInfos {
	tmp, _ := infos.get(info.Id)
	if tmp != nil {
		infos = infos.delete(info.Id)
	}
	diskLocker.Lock()
	infos = append(infos, info)
	diskLocker.Unlock()
	return infos
}

func (infos DiskInfos) delete(id string) DiskInfos {
	diskLocker.Lock()
	defer diskLocker.Unlock()

	var ret DiskInfos
	for _, info := range infos {
		if info.Id == id {
			info.destroy()
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos DiskInfos) destroy() {
	for _, info := range infos {
		info.destroy()
	}
}

func newDiskInfoFromMount(mount *gio.Mount) *DiskInfo {
	var root = mount.GetRoot()
	var info = &DiskInfo{
		Name:       mount.GetName(),
		MountPoint: root.GetUri(),
		Type:       DiskTypeNative,
		CanEject:   mount.CanEject(),
		CanUnmount: mount.CanUnmount(),
		Used:       queryAttrUint64(root, fsAttrUsed),
		Total:      queryAttrUint64(root, fsAttrSize),
		object: &diskObject{
			Object: mount,
			Type:   diskObjMount,
		},
	}
	root.Unref()

	volume := mount.GetVolume()
	if volume != nil {
		info.Id = volume.GetIdentifier(volumeKindUUID)
		info.Path = volume.GetIdentifier(volumeKindUnix)
		volume.Unref()
	}

	if len(info.Id) == 0 {
		if len(info.Path) != 0 {
			info.Id = info.Path
		} else {
			info.Id = info.MountPoint
		}
	}

	gicon := mount.GetIcon()
	info.Icon = getIconFromGIcon(gicon)
	gicon.Unref()

	info.correctDiskType()
	return info
}

func newDiskInfoFromVolume(volume *gio.Volume) *DiskInfo {
	var info = &DiskInfo{
		Name:     volume.GetName(),
		Id:       volume.GetIdentifier(volumeKindUUID),
		Type:     DiskTypeNative,
		Path:     volume.GetIdentifier(volumeKindUnix),
		CanEject: volume.CanEject(),
		object: &diskObject{
			Object: volume,
			Type:   diskObjVolume,
		},
	}

	if len(info.Id) == 0 {
		info.Id = info.Path
	}

	gicon := volume.GetIcon()
	info.Icon = getIconFromGIcon(gicon)
	gicon.Unref()

	info.correctDiskType()
	return info
}

func (info *DiskInfo) destroy() {
	switch info.object.Type {
	case diskObjVolume:
		volume, ok := info.object.Object.(*gio.Volume)
		if ok && volume.Object.C != nil {
			volume.Unref()
		}
	case diskObjMount:
		mount, ok := info.object.Object.(*gio.Mount)
		if ok && mount.Object.C != nil {
			mount.Unref()
		}
	}
}

var (
	mtpReg     = regexp.MustCompile(`^mtp://`)
	smbReg     = regexp.MustCompile(`^smb://`)
	ftpReg     = regexp.MustCompile(`^ftp://`)
	networkReg = regexp.MustCompile(`^network`)
)

func (info *DiskInfo) correctDiskType() {
	if info.CanEject || strings.Contains(info.Icon, "usb") {
		info.Type = DiskTypeRemovable
	}

	if smbReg.MatchString(info.MountPoint) ||
		ftpReg.MatchString(info.MountPoint) ||
		networkReg.MatchString(info.Path) {
		info.Type = DiskTypeNetwork
	}

	if mtpReg.MatchString(info.MountPoint) {
		info.Type = DiskTypeRemovable
		info.Icon = mtpDeviceIcon
	}
}

func queryAttrUint64(file *gio.File, attr string) uint64 {
	info, err := file.QueryFilesystemInfo(attr, nil)
	if err != nil {
		return 0
	}
	defer info.Unref()

	return info.GetAttributeUint64(attr) / 1024
}

func getIconFromGIcon(gicon *gio.Icon) string {
	icons := strings.Split(gicon.ToString(), " ")
	if len(icons) > 2 {
		return icons[2]
	}
	return ""
}
