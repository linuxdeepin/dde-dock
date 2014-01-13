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

package main

import (
	"dlib/dbus"
	"dlib/gio-2.0"
	"fmt"
)

type DiskInfo struct {
	Id         int32
	Name       string
	Type       string
	CanUnmount bool
	CanEject   bool
	UsableCap  uint32
	TotalCap   uint32
}

type ObjectInfo struct {
	Object interface{}
	Type   string
}

type Manager struct {
	DiskList []*DiskInfo
}

const (
	DEVICE_KIND = "unix-device"

	DISK_INFO_DEST = "com.deepin.daemon.DiskMount"
	DISK_INFO_PATH = "/com/deepin/daemon/DiskMount"
	DISK_INFO_IFC  = "com.deepin.daemon.DiskMount"
)

var (
	count     = 0
	monitor   = gio.VolumeMonitorGet()
	objectMap map[int32]*ObjectInfo
)

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		DISK_INFO_DEST,
		DISK_INFO_PATH,
		DISK_INFO_IFC,
	}
}

func (m *Manager) DeviceMount(id int32, mount bool) {
}

func (m *Manager) DeviceEject (id int32, eject bool) {
}

func NewDiskInfo(value interface{}, t string, id int32) *DiskInfo {
	info := &DiskInfo{}
	info.Id = id

	switch t {
	case "volume":
		{
			v := value.(*gio.Volume)
			info.Name = v.GetName()
			info.CanEject = v.CanEject()
			id := v.GetIdentifier(DEVICE_KIND)
			if ContainStart("network", id) {
				info.Type = "network"
			} else if info.CanEject {
				info.Type = "removable"
			} else {
				info.Type = "native"
			}
			break
		}
	case "driver":
		{
			v := value.(*gio.Drive)
			info.Name = v.GetName()
			info.CanEject = v.CanEject()
			id := v.GetIdentifier(DEVICE_KIND)
			if ContainStart("network", id) {
				info.Type = "network"
			} else if info.CanEject {
				info.Type = "removable"
			} else {
				info.Type = "native"
			}
			break
		}
	case "mount":
		{
			v := value.(*gio.Mount)
			info.Name = v.GetName()
			info.CanEject = v.CanEject()
			info.CanUnmount = v.CanUnmount()
			root := v.GetRoot()
			if root.IsNative() {
				info.Type = "native"
			} else if info.CanEject {
				info.Type = "removable"
			} else {
				info.Type = "network"
			}
			break
		}
	default:
		break
	}

	return info
}

func NewObjectInfo(v interface{}, t string) *ObjectInfo {
	return &ObjectInfo{Object: v, Type: t}
}

func DriverList(m *Manager) {
	drivers := monitor.GetConnectedDrives()
	for _, driver := range drivers {
		volumes := driver.GetVolumes()
		if volumes == nil {
			if driver.IsMediaRemovable() &&
				!driver.IsMediaCheckAutomatic() {
				info := NewDiskInfo(driver, "driver", int32(count))
				count += 1
				objectMap[info.Id] = NewObjectInfo(driver, "driver")
				m.DiskList = append(m.DiskList, info)
			}
			continue
		}
		for _, volume := range volumes {
			mount := volume.GetMount()
			if mount != nil {
				info := NewDiskInfo(mount, "mount", int32(count))
				count += 1
				objectMap[info.Id] = NewObjectInfo(mount, "mount")
				m.DiskList = append(m.DiskList, info)
			} else {
				info := NewDiskInfo(volume, "volume", int32(count))
				count += 1
				objectMap[info.Id] = NewObjectInfo(volume, "volume")
				m.DiskList = append(m.DiskList, info)
			}
		}
	}
}

func VolumeList(m *Manager) {
	volumes := monitor.GetVolumes()
	for _, volume := range volumes {
		driver := volume.GetDrive()
		if driver != nil {
			continue
		}
		id := volume.GetIdentifier("unix-device")
		fmt.Printf("id: %s\n", id)
		mount := volume.GetMount()
		if mount != nil {
			info := NewDiskInfo(mount, "mount", int32(count))
			count += 1
			objectMap[info.Id] = NewObjectInfo(mount, "mount")
			m.DiskList = append(m.DiskList, info)
		} else {
			info := NewDiskInfo(volume, "volume", int32(count))
			count += 1
			objectMap[info.Id] = NewObjectInfo(volume, "volume")
			m.DiskList = append(m.DiskList, info)
		}
	}
}

func MountList(m *Manager) {
	mounts := monitor.GetMounts()
	for _, mount := range mounts {
		if mount.IsShadowed() {
			continue
		}

		volume := mount.GetVolume()
		if volume != nil {
			id := volume.GetIdentifier("unix-device")
			fmt.Printf("id: %s\n", id)
			continue
		}
		info := NewDiskInfo(mount, "mount", int32(count))
		count += 1
		objectMap[info.Id] = NewObjectInfo(mount, "mount")
		m.DiskList = append(m.DiskList, info)
	}
}

func ContainStart(str1, str2 string) bool {
	for i, _ := range str1 {
		if str1[i] != str2[i] {
			return false
		}
	}

	return true
}

func main() {
	objectMap = make(map[int32]*ObjectInfo)
	m := &Manager{}
	DriverList(m)
	VolumeList(m)
	MountList(m)

	dbus.InstallOnSession(m)
	dbus.DealWithUnhandledMessage()

	select {}
}
