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
	"dlib"
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
	DiskList []DiskInfo
}

const (
	DEVICE_KIND = "unix-device"
)

var (
	count     = 0
	monitor   = gio.VolumeMonitorGet()
	objectMap map[int32]*ObjectInfo
)

func (m *Manager) DeviceEject(id int32) {
	info, ok := objectMap[id]
	if !ok {
		fmt.Printf("Eject id - %d not in objectMap.\n", id)
		return
	}

	fmt.Printf("Eject type: %s\n", info.Type)
	switch info.Type {
	case "drive":
		op := info.Object.(*gio.Drive)
		op.Eject(gio.MountUnmountFlagsNone, nil, nil)
	case "volume":
		op := info.Object.(*gio.Volume)
		op.Eject(gio.MountUnmountFlagsNone, nil, nil)
	case "mount":
		op := info.Object.(*gio.Mount)
		op.Eject(gio.MountUnmountFlagsNone, nil, nil)
	default:
		fmt.Printf("'%s' invalid type\n", info.Type)
	}
}

func (m *Manager) DeviceMount(id int32) {
	info, ok := objectMap[id]
	if !ok {
		fmt.Printf("Mount id - %d not in objectMap.\n", id)
		return
	}

	fmt.Printf("Mount type: %s\n", info.Type)
	switch info.Type {
	case "volume":
		op := info.Object.(*gio.Volume)
		op.Mount(gio.MountMountFlagsNone, nil, nil, nil)
	case "mount":
		op := info.Object.(*gio.Mount)
		op.Remount(gio.MountMountFlagsNone, nil, nil, nil)
	default:
		fmt.Printf("'%s' invalid type\n", info.Type)
	}
}

func (m *Manager) DeviceUnmount(id int32) {
	info, ok := objectMap[id]
	if !ok {
		fmt.Printf("Unmount id - %d not in objectMap.\n", id)
		return
	}

	fmt.Printf("Unmount type: %s\n", info.Type)
	switch info.Type {
	case "mount":
		op := info.Object.(*gio.Mount)
		op.Unmount(gio.MountUnmountFlagsNone, nil, nil)
	default:
		fmt.Printf("'%s' invalid type\n", info.Type)
	}
}

func newDiskInfo(value interface{}, t string, id int32) DiskInfo {
	info := DiskInfo{}
	info.Id = id

	switch t {
	case "volume":
		v := value.(*gio.Volume)
		info.Name = v.GetName()
		info.CanEject = v.CanEject()
		id := v.GetIdentifier(DEVICE_KIND)
		if containStart("network", id) {
			info.Type = "network"
		} else if info.CanEject {
			info.Type = "removable"
		} else {
			info.Type = "native"
		}
	case "drive":
		v := value.(*gio.Drive)
		info.Name = v.GetName()
		info.CanEject = v.CanEject()
		id := v.GetIdentifier(DEVICE_KIND)
		if containStart("network", id) {
			info.Type = "network"
		} else if info.CanEject {
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
		if root.IsNative() {
			info.Type = "native"
		} else if info.CanEject {
			info.Type = "removable"
		} else {
			info.Type = "network"
		}
	default:
		fmt.Printf("'%s' invalid type\n", t)
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
				info := newDiskInfo(driver, "drive", int32(count))
				count += 1
				objectMap[info.Id] = newObjectInfo(driver, "drive")
				list = append(list, info)
			}
			continue
		}
		for _, volume := range volumes {
			mount := volume.GetMount()
			if mount != nil {
				info := newDiskInfo(mount, "mount", int32(count))
				count += 1
				objectMap[info.Id] = newObjectInfo(mount, "mount")
				list = append(list, info)
			} else {
				info := newDiskInfo(volume, "volume", int32(count))
				count += 1
				objectMap[info.Id] = newObjectInfo(volume, "volume")
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
		id := volume.GetIdentifier("unix-device")
		fmt.Printf("id: %s\n", id)
		mount := volume.GetMount()
		if mount != nil {
			info := newDiskInfo(mount, "mount", int32(count))
			count += 1
			objectMap[info.Id] = newObjectInfo(mount, "mount")
			list = append(list, info)
		} else {
			info := newDiskInfo(volume, "volume", int32(count))
			count += 1
			objectMap[info.Id] = newObjectInfo(volume, "volume")
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
			id := volume.GetIdentifier("unix-device")
			fmt.Printf("id: %s\n", id)
			continue
		}
		info := newDiskInfo(mount, "mount", int32(count))
		count += 1
		objectMap[info.Id] = newObjectInfo(mount, "mount")
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
	for k, info := range objectMap {
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
		delete(objectMap, k)
	}
	count = 0
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("recover err: %s\n", err)
		}
	}()
	objectMap = make(map[int32]*ObjectInfo)
	m := &Manager{}
	m.setPropName("DiskList")
	printDiskInfo(m.DiskList)
	m.listenSignalChanged()

	dbus.InstallOnSession(m)
	dbus.DealWithUnhandledMessage()

	dlib.StartLoop()
}

func printDiskInfo(infos []DiskInfo) {
	for _, v := range infos {
		fmt.Printf("Id: %d\n", v.Id)
		fmt.Printf("Name: %s\n", v.Name)
		fmt.Printf("Type: %s\n", v.Type)
		fmt.Println("CanEject:", v.CanEject)
		fmt.Println("CanUnmount:", v.CanUnmount)
		fmt.Printf("\n")
	}
}
