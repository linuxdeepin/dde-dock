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
	"dbus/org/freedesktop/udisks2"
	"fmt"
	"gir/gio-2.0"
	"path"
	"pkg.deepin.io/lib/dbus"
	"regexp"
	"strings"
	"sync"
)

const (
	DiskTypeNative    string = "native"
	DiskTypeRemovable        = "removable"
	DiskTypeNetwork          = "network"
	DiskTypeIPhone           = "iphone"
	DiskTypePhone            = "phone"
	DiskTypeCamera           = "camera"
	DiskTypeDVD              = "dvd"
)

const (
	volumeKindUnix = "unix-device"
	volumeKindUUID = "uuid"

	fsAttrSize = "filesystem::size"
	fsAttrUsed = "filesystem::used"

	mtpDeviceIcon = "drive-removable-media-mtp"
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
}
type DiskInfos []*DiskInfo

func (m *Manager) getDiskInfos() DiskInfos {
	var infos DiskInfos
	for _, volume := range m.monitor.GetVolumes() {
		mount := volume.GetMount()
		if mount != nil {
			mount.Unref()
			volume.Unref()
			continue
		}

		info := newDiskInfoFromVolume(volume)
		volume.Unref()
		infos = infos.add(info)
	}

	for _, mount := range m.monitor.GetMounts() {
		info := newDiskInfoFromMount(mount)
		mount.Unref()
		if info == nil {
			continue
		}
		infos = infos.add(info)
	}

	return infos
}

func (m *Manager) getVolumeById(id string) *gio.Volume {
	var ret *gio.Volume = nil
	for _, volume := range m.monitor.GetVolumes() {
		if ret != nil {
			volume.Unref()
			continue
		}

		mount := volume.GetMount()
		if mount != nil && mount.Object.C != nil {
			mount.Unref()
			volume.Unref()
			continue
		}

		if getVolumeId(volume) == id {
			ret = volume
			continue
		}
	}
	return ret
}

func (m *Manager) getMountById(id string) *gio.Mount {
	var ret *gio.Mount
	for _, mount := range m.monitor.GetMounts() {
		if ret != nil {
			mount.Unref()
			continue
		}

		if getMountId(mount) == id {
			ret = mount
			continue
		}
		mount.Unref()
	}
	return ret
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

func (infos DiskInfos) getByMountPoint(point string) (*DiskInfo, error) {
	diskLocker.Lock()
	defer diskLocker.Unlock()
	for _, info := range infos {
		if info.MountPoint == point {
			return info, nil
		}
	}
	return nil, fmt.Errorf("Invalid disk mount point: %v", point)
}

func (infos DiskInfos) exists(value *DiskInfo) string {
	diskLocker.Lock()
	defer diskLocker.Unlock()
	for _, info := range infos {
		// iphone exist 2 mount struct in gio mounts
		if info.Id == value.Id ||
			((len(info.MountPoint) != 0) && (info.MountPoint == value.MountPoint)) {
			return info.Id
		}
	}
	return ""
}

func (infos DiskInfos) add(info *DiskInfo) DiskInfos {
	id := infos.exists(info)
	if len(id) != 0 {
		logger.Debug("[Add] disk exist:", id)
		infos = infos.delete(id)
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
			continue
		}
		ret = append(ret, info)
	}
	return ret
}

func (infos DiskInfos) duplicate() DiskInfos {
	diskLocker.Lock()
	defer diskLocker.Unlock()
	var newInfos DiskInfos
	for _, info := range infos {
		tmp := DiskInfo{
			Id:         info.Id,
			Name:       info.Name,
			Type:       info.Type,
			Icon:       info.Icon,
			Path:       info.Path,
			MountPoint: info.MountPoint,
			CanEject:   info.CanEject,
			CanUnmount: info.CanUnmount,
			Used:       info.Used,
			Total:      info.Total,
		}
		newInfos = append(newInfos, &tmp)
	}
	return newInfos
}

func newDiskInfoFromMount(mount *gio.Mount) *DiskInfo {
	var root = mount.GetRoot()
	var info = &DiskInfo{
		Id:         getMountId(mount),
		Name:       mount.GetName(),
		MountPoint: root.GetUri(),
		Type:       DiskTypeNative,
		CanEject:   mount.CanEject(),
		CanUnmount: mount.CanUnmount(),
		Used:       queryAttrUint64(root, fsAttrUsed),
		Total:      queryAttrUint64(root, fsAttrSize),
	}
	root.Unref()

	gicon := mount.GetIcon()
	info.Icon = getIconFromGIcon(gicon)
	gicon.Unref()

	info.correctDiskType()
	volume := mount.GetVolume()
	if volume == nil || volume.Object.C == nil {
		// All devices must have Volume objects, in addition to network devices
		if info.Type != DiskTypeNetwork {
			logger.Debugf("Invalid disk info: %#v", info)
			return nil
		}
		return info
	}

	info.Path = volume.GetIdentifier(volumeKindUnix)
	volume.Unref()
	if info.Total == 0 && info.Path != "" {
		info.Total = getTotalSizeByUDisks2(info.Path)
	}

	return info
}

func newDiskInfoFromVolume(volume *gio.Volume) *DiskInfo {
	var info = &DiskInfo{
		Id:       getVolumeId(volume),
		Name:     volume.GetName(),
		Type:     DiskTypeNative,
		Path:     volume.GetIdentifier(volumeKindUnix),
		CanEject: volume.CanEject(),
	}

	if info.Total == 0 {
		info.Total = getTotalSizeByUDisks2(info.Path)
	}

	gicon := volume.GetIcon()
	info.Icon = getIconFromGIcon(gicon)
	gicon.Unref()

	info.correctDiskType()
	return info
}

func (info *DiskInfo) correctDiskType() {
	if info.CanEject || strings.Contains(info.Icon, "usb") {
		info.Type = DiskTypeRemovable
	}

	if stringStartWith(info.MountPoint, "smb://") ||
		stringStartWith(info.MountPoint, "ftp://") ||
		stringStartWith(info.Path, "network://") {
		info.Type = DiskTypeNetwork
		return
	}

	if stringStartWith(info.MountPoint, "afc://") ||
		strings.Contains(info.Icon, "iphone") {
		info.Type = DiskTypeIPhone
		return
	}

	if stringStartWith(info.MountPoint, "mtp://") ||
		info.Icon == "phone" {
		info.Type = DiskTypePhone
		info.Icon = "phone"
		return
	}

	if stringStartWith(info.MountPoint, "gphoto2://") ||
		strings.Contains(info.Icon, "camera") {
		info.Type = DiskTypeCamera
		return
	}

	if strings.Contains(info.Icon, "dvd") {
		info.Type = DiskTypeDVD
		return
	}
}

func queryAttrUint64(file *gio.File, attr string) uint64 {
	info, err := file.QueryFilesystemInfo(attr, nil)
	if err != nil || info == nil || info.Object.C == nil {
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
	return icons[0]
}

func getTotalSizeByUDisks2(devPath string) uint64 {
	blockObj, err := udisks2.NewBlock("org.freedesktop.UDisks2",
		dbus.ObjectPath("/org/freedesktop/UDisks2/block_devices/"+path.Base(devPath)))
	if err != nil {
		logger.Debug("New udisks2 block object failed:", err)
		return 0
	}
	defer udisks2.DestroyBlock(blockObj)

	return blockObj.Size.Get() / 1024
}

func getVolumeId(volume *gio.Volume) string {
	id := volume.GetIdentifier(volumeKindUUID)
	if len(id) != 0 {
		return id
	}

	// if uuid not exist, use path as id
	return volume.GetIdentifier(volumeKindUnix)
}

func getMountId(mount *gio.Mount) string {
	root := mount.GetRoot()
	mountPoint := root.GetUri()
	// Don't unref root, root only ref once in mount
	// root.Unref()

	volume := mount.GetVolume()
	if volume == nil || volume.Object.C == nil {
		return getIdByMountPoint(mountPoint)
	}

	id := getVolumeId(volume)
	volume.Unref()
	if len(id) != 0 {
		return id
	}
	return getIdByMountPoint(mountPoint)
}

func getIdByMountPoint(mountPoint string) string {
	var id = mountPoint
	switch {
	// iphone
	case stringStartWith(mountPoint, "afc://"):
		id = regexp.MustCompile(`^afc://`).ReplaceAllString(id, "")
		id = strings.TrimRight(id, "/")
	}

	return id
}

func stringStartWith(s, key string) bool {
	if len(s) > len(key) && (string(s[:len(key)]) == key) {
		return true
	}
	return false
}
