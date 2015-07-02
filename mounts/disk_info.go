/**
 * Copyright (c) 2011 ~ 2015 Deepin, Inc.
 *               2013 ~ 2015 jouyouyun
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

package mounts

import (
	"fmt"
	"os"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"regexp"
	"strings"
)

const (
	diskTypeNative    = "native"
	diskTypeNetwork   = "network"
	diskTypeRemovable = "removable"

	volumeKindUnix = "unix-device"
	volumeKindUUID = "uuid"

	fsAttrSize = "filesystem::size"
	fsAttrUsed = "filesystem::used"

	mtpDiskIcon = "drive-removable-media-mtp"
)

type DiskInfo struct {
	// Disk description
	Name string
	// Disk type, ex: native, removable, network...
	Type string

	CanUnmount bool
	CanEject   bool

	// The size of disk used
	Used uint64
	// The capacity of disk
	Size uint64

	Path string
	UUID string
	// The mounted path
	MountPoint string
	Icon       string
}
type DiskInfos []DiskInfo

func newDiskInfoFromMount(mount *gio.Mount) DiskInfo {
	var root = mount.GetRoot()
	defer root.Unref()
	var info = DiskInfo{
		Name:       mount.GetName(),
		MountPoint: root.GetUri(),
		CanEject:   mount.CanEject(),
		CanUnmount: mount.CanUnmount(),
		Used:       getDiskAttrUint64(root, fsAttrUsed),
		Size:       getDiskAttrUint64(root, fsAttrSize),
	}

	volume := mount.GetVolume()
	if volume != nil {
		info.Path = volume.GetIdentifier(volumeKindUnix)
		info.UUID = volume.GetIdentifier(volumeKindUUID)
		volume.Unref()
	}

	if len(info.UUID) == 0 {
		info.UUID = generateUUID()
	}

	iconObj := mount.GetIcon()
	defer iconObj.Unref()
	info.Icon = getIconFromGIcon(iconObj)

	if info.CanEject || strings.Contains(info.Icon, "usb") {
		info.Type = diskTypeRemovable
	} else if root.IsNative() {
		info.Type = diskTypeNative
	} else {
		info.Type = diskTypeNetwork
	}

	ok, _ := regexp.MatchString(`^mtp://`, info.MountPoint)
	if ok {
		info.Type = diskTypeRemovable
		info.Icon = mtpDiskIcon
	}

	return info
}

func newDiskInfoFromVolume(volume *gio.Volume) DiskInfo {
	mount := volume.GetMount()
	if mount != nil {
		defer mount.Unref()
		return newDiskInfoFromMount(mount)
	}

	var info = DiskInfo{
		Name:     volume.GetName(),
		Path:     volume.GetIdentifier(volumeKindUnix),
		UUID:     volume.GetIdentifier(volumeKindUUID),
		CanEject: volume.CanEject(),
	}

	if len(info.UUID) == 0 {
		info.UUID = generateUUID()
	}

	iconObj := volume.GetIcon()
	defer iconObj.Unref()
	info.Icon = getIconFromGIcon(iconObj)

	if info.CanEject || strings.Contains(info.Icon, "usb") {
		info.Type = diskTypeRemovable
	} else {
		info.Type = diskTypeNative
	}

	ok, _ := regexp.MatchString(`^network`, info.Path)
	if ok {
		info.Type = diskTypeNetwork
	}

	return info
}

func getDiskAttrUint64(file *gio.File, attr string) uint64 {
	info, err := file.QueryFilesystemInfo(attr, nil)
	if err != nil {
		return 0
	}
	defer info.Unref()

	return info.GetAttributeUint64(attr) / 1024
}

func getIconFromGIcon(iconObj *gio.Icon) string {
	icons := strings.Split(iconObj.ToString(), " ")
	if len(icons) > 2 {
		return icons[2]
	}

	return ""
}

func generateUUID() string {
	f, err := os.Open("/dev/urandom")
	if err != nil {
		return ""
	}

	defer f.Close()
	b := make([]byte, 16)
	f.Read(b)
	uuid := fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6],
		b[6:8], b[8:10], b[10:])

	return uuid
}
