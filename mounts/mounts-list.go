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

package mounts

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/logger"
	"regexp"
	"strings"
)

type DiskInfo struct {
	Name       string
	Type       string
	CanUnmount bool
	CanEject   bool
	UsableCap  int64
	TotalCap   int64
	Path       string
	UUID       string
	MountURI   string
	IconName   string
}

type ObjectInfo struct {
	Object interface{}
	Type   string
}

type Manager struct {
	DiskList []DiskInfo
	Error    func(string, string)

	quitFlag chan bool
}

const (
	DEVICE_KIND = "unix-device"
)

var (
	monitor   = gio.VolumeMonitorGet()
	objectMap = make(map[string]*ObjectInfo)
	Logger    = logger.NewLogger(DISK_INFO_DEST)
)

func newDiskInfo(value interface{}, t string) DiskInfo {
	defer func() {
		if err := recover(); err != nil {
			Logger.Error("Received Error: ", err)
		}
	}()

	info := DiskInfo{}

	switch t {
	case "volume":
		v := value.(*gio.Volume)
		info.Name = v.GetName()
		info.CanEject = v.CanEject()
		id := v.GetIdentifier(DEVICE_KIND)
		//info.TotalCap, info.UsableCap = getDiskCap(id)
		info.Path = v.GetIdentifier(gio.VolumeIdentifierKindUnixDevice)
		info.UUID = v.GetIdentifier(gio.VolumeIdentifierKindUuid)
		Logger.Debugf("VOLUME Name: %s, UUID: %v", info.Name, info.UUID)
		if len(info.UUID) < 1 {
			info.UUID = generateUUID()
			Logger.Debugf("VOLUME Name: %s, Generate UUID: %v", info.Name, info.UUID)
		}

		if mount := v.GetMount(); mount != nil {
			info.CanUnmount = mount.CanUnmount()
		}

		icons := v.GetIcon().ToString()
		as := strings.Split(icons, " ")
		if len(as) > 2 {
			info.IconName = as[2]
		}

		if containStart("network", id) {
			info.Type = "network"
		} else if info.CanEject ||
			strings.Contains(info.IconName, "usb") {
			info.Type = "removable"
		} else {
			info.Type = "native"
		}
	case "drive":
		v := value.(*gio.Drive)
		info.Name = v.GetName()
		info.CanEject = v.CanEject()
		id := v.GetIdentifier(DEVICE_KIND)
		//info.TotalCap, info.UsableCap = getDiskCap(id)
		info.Path = v.GetIdentifier(gio.VolumeIdentifierKindUnixDevice)
		info.UUID = v.GetIdentifier(gio.VolumeIdentifierKindUuid)
		Logger.Infof("DRIVE Name: %s, UUID: %v", info.Name, info.UUID)
		if len(info.UUID) < 1 {
			info.UUID = generateUUID()
			Logger.Infof("DRIVE Name: %s, Generate UUID: %v", info.Name, info.UUID)
		}

		icons := v.GetIcon().ToString()
		as := strings.Split(icons, " ")
		if len(as) > 2 {
			info.IconName = as[2]
		}
		if containStart("network", id) {
			info.Type = "network"
		} else if info.CanEject ||
			strings.Contains(info.IconName, "usb") {
			info.Type = "removable"
		} else {
			info.Type = "native"
		}
	case "mount":
		v := value.(*gio.Mount)
		info.Name = v.GetName()
		info.CanEject = v.CanEject()
		info.CanUnmount = v.CanUnmount()
		root := v.GetRoot()
		info.MountURI = root.GetUri()
		info.TotalCap, info.UsableCap = getDiskCap(root.GetPath())

		icons := v.GetIcon().ToString()
		as := strings.Split(icons, " ")
		if len(as) > 2 {
			info.IconName = as[2]
		}

		if info.CanEject ||
			strings.Contains(info.IconName, "usb") {
			info.Type = "removable"
		} else if root.IsNative() {
			info.Type = "native"
		} else {
			info.Type = "network"
		}
		if volume := v.GetVolume(); volume != nil {
			info.Path = volume.GetIdentifier(gio.VolumeIdentifierKindUnixDevice)
			info.UUID = volume.GetIdentifier(gio.VolumeIdentifierKindUuid)

		}
		Logger.Infof("MOUNT Name: %s, UUID: %v", info.Name, info.UUID)
		if len(info.UUID) < 1 {
			info.UUID = generateUUID()
			Logger.Infof("MOUNT Name: %s, Generate UUID: %v", info.Name, info.UUID)
		}

		if ok, _ := regexp.MatchString(`^mtp://`, info.MountURI); ok {
			info.Type = "removable"
			info.IconName = "drive-removable-media-mtp"
		}

		if info.TotalCap == 0 {
			info.TotalCap, info.UsableCap = getDiskCap(info.Path)
		}
	default:
		Logger.Errorf("'%s' invalid type", t)
	}

	return info
}

func newObjectInfo(v interface{}, t string) *ObjectInfo {
	return &ObjectInfo{Object: v, Type: t}
}

func driverList() []DiskInfo {
	list := []DiskInfo{}
	drivers := monitor.GetConnectedDrives()
	for _, driver := range drivers {
		volumes := driver.GetVolumes()
		if volumes == nil {
			if driver.IsMediaRemovable() &&
				!driver.IsMediaCheckAutomatic() {
				info := newDiskInfo(driver, "drive")
				objectMap[info.UUID] = newObjectInfo(driver, "drive")
				list = append(list, info)
			}
			continue
		}
		for _, volume := range volumes {
			mount := volume.GetMount()
			if mount != nil {
				info := newDiskInfo(mount, "mount")
				objectMap[info.UUID] = newObjectInfo(mount, "mount")
				list = append(list, info)
			} else {
				info := newDiskInfo(volume, "volume")
				objectMap[info.UUID] = newObjectInfo(volume, "volume")
				list = append(list, info)
			}
		}
	}

	return list
}

func volumeList() []DiskInfo {
	list := []DiskInfo{}
	volumes := monitor.GetVolumes()
	for _, volume := range volumes {
		driver := volume.GetDrive()
		if driver != nil {
			continue
		}
		//id := volume.GetIdentifier("unix-device")
		mount := volume.GetMount()
		if mount != nil {
			info := newDiskInfo(mount, "mount")
			objectMap[info.UUID] = newObjectInfo(mount, "mount")
			list = append(list, info)
		} else {
			info := newDiskInfo(volume, "volume")
			objectMap[info.UUID] = newObjectInfo(volume, "volume")
			list = append(list, info)
		}
	}
	return list
}

func mountList() []DiskInfo {
	list := []DiskInfo{}
	mounts := monitor.GetMounts()
	for _, mount := range mounts {
		if mount.IsShadowed() {
			continue
		}

		volume := mount.GetVolume()
		if volume != nil {
			continue
		}
		info := newDiskInfo(mount, "mount")
		objectMap[info.UUID] = newObjectInfo(mount, "mount")
		list = append(list, info)
	}
	return list
}

func containStart(str1, str2 string) bool {
	for i, _ := range str1 {
		if str1[i] != str2[i] {
			return false
		}
	}

	return true
}

func getDiskInfoList() []DiskInfo {
	list := []DiskInfo{}

	destroyObjectMap()
	l1 := driverList()
	l2 := volumeList()
	l3 := mountList()
	list = append(list, l1...)
	list = append(list, l2...)
	list = append(list, l3...)

	return list
}

func destroyObjectMap() {
	for _, info := range objectMap {
		switch info.Type {
		case "drive":
			op := info.Object.(*gio.Drive)
			op.Unref()
		case "volume":
			op := info.Object.(*gio.Volume)
			op.Unref()
		case "mount":
			op := info.Object.(*gio.Mount)
			op.Unref()
		}
	}
	objectMap = make(map[string]*ObjectInfo)
}

func NewManager() *Manager {
	m := &Manager{}
	m.setPropName("DiskList")
	m.listenSignalChanged()
	m.quitFlag = make(chan bool)

	return m
}

var _manager *Manager

func GetManager() *Manager {
	if _manager == nil {
		_manager = NewManager()
	}

	return _manager
}

func Start() {
	Logger.BeginTracing()

	m := GetManager()
	err := dbus.InstallOnSession(m)
	if err != nil {
		Logger.Fatal("Install DBus Session Failed:", err)
	}

	go m.refrashDiskInfoList()
}

func Stop() {
	GetManager().quitFlag <- true
	Logger.EndTracing()
}
