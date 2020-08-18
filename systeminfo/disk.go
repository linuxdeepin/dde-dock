/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"fmt"

	"github.com/godbus/dbus"
)
//nolint
type diskInfo struct {
	Drive       dbus.ObjectPath // org.freedesktop.UDisks2.Block Drive
	MountPoints []string        // org.freedesktop.UDisks2.Filesystem MountPoints
	Size        uint64          // org.freedesktop.UDisks2.Partition Size
	Table       dbus.ObjectPath // org.freedesktop.UDisks2.Partition Table
}

type diskInfoMap map[dbus.ObjectPath]diskInfo //nolint

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
	dlist, err := GetDiskList()
	if err != nil {
		fmt.Println("Failed to get disk list:", err)
		return 0, err
	}

	rdisk := dlist.GetRoot()
	if rdisk == nil {
		fmt.Println("Failed to get root disk")
		return 0, err
	}
	return uint64(rdisk.Size), err
}
