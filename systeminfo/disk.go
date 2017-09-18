/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package systeminfo

import (
	"dbus/org/freedesktop/udisks2"
	"fmt"
	"pkg.deepin.io/lib/dbus"
)

type diskInfo struct {
	Drive       dbus.ObjectPath // org.freedesktop.UDisks2.Block Drive
	MountPoints []string        // org.freedesktop.UDisks2.Filesystem MountPoints
	Size        uint64          // org.freedesktop.UDisks2.Partition Size
	Table       dbus.ObjectPath // org.freedesktop.UDisks2.Partition Table
}

type diskInfoMap map[dbus.ObjectPath]diskInfo

func (set diskInfoMap) GetRootDrive() dbus.ObjectPath {
	for _, v := range set {
		if v.Drive == "" {
			continue
		}
		for _, mp := range v.MountPoints {
			if mp == "/" {
				return v.Drive
			}
		}
	}
	return ""
}

func (set diskInfoMap) GetRootTable() dbus.ObjectPath {
	for _, v := range set {
		if v.Table == "" {
			continue
		}
		for _, mp := range v.MountPoints {
			if mp == "/" {
				return v.Table
			}
		}
	}
	return ""
}

func (set diskInfoMap) GetRootSize() uint64 {
	for _, v := range set {
		if len(v.MountPoints) == 0 {
			continue
		}
		for _, mp := range v.MountPoints {
			if mp == "/" {
				return v.Size
			}
		}
	}
	return 0
}

func (set diskInfoMap) Get(key dbus.ObjectPath) *diskInfo {
	if key == "" {
		return nil
	}
	v, ok := set[key]
	if !ok {
		return nil
	}
	return &v
}

func getDiskCap() (uint64, error) {
	set, err := parseUDisksManagers()
	if err != nil {
		return 0, err
	}

	key := set.GetRootDrive()
	info := set.Get(key)
	if info != nil {
		return info.Size, nil
	}

	key = set.GetRootTable()
	info = set.Get(key)
	if info != nil {
		return info.Size, nil
	}

	// not found drive and table, try root mount point
	size := set.GetRootSize()
	err = nil
	if size == 0 {
		err = fmt.Errorf("Failed to get disk capacity: not found root mount point")
	}
	return size, err
}

func parseUDisksManagers() (diskInfoMap, error) {
	udisk, err := udisks2.NewObjectManager(
		"org.freedesktop.UDisks2",
		"/org/freedesktop/UDisks2")
	if err != nil {
		return nil, err
	}

	var set = make(diskInfoMap)
	managers, _ := udisk.GetManagedObjects()
	for k, manager := range managers {
		var info diskInfo
		block, ok := manager["org.freedesktop.UDisks2.Block"]
		if ok {
			info.Drive = block["Drive"].Value().(dbus.ObjectPath)
		}

		fs, ok := manager["org.freedesktop.UDisks2.Filesystem"]
		if ok {
			values := fs["MountPoints"].Value().([][]byte)
			for _, v := range values {
				// filter the end char '\x00'
				info.MountPoints = append(info.MountPoints, string(v[:len(v)-1]))
			}
		}

		// the manager maybe a partition, or partition table, or drive, or loop(wubi)
		partition, ok := manager["org.freedesktop.UDisks2.Partition"]
		if ok {
			info.Size = partition["Size"].Value().(uint64)
			info.Table = partition["Table"].Value().(dbus.ObjectPath)
			set[k] = info
			continue
		}

		_, ok = manager["org.freedesktop.UDisks2.PartitionTable"]
		if ok {
			info.Size = block["Size"].Value().(uint64)
			set[k] = info
			continue
		}

		drive, ok := manager["org.freedesktop.UDisks2.Drive"]
		if ok {
			info.Size = drive["Size"].Value().(uint64)
			set[k] = info
			continue
		}

		_, ok = manager["org.freedesktop.UDisks2.Loop"]
		if ok {
			info.Size = block["Size"].Value().(uint64)
			set[k] = info
			continue
		}
	}

	return set, nil
}
