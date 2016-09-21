/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package systeminfo

import (
	"dbus/org/freedesktop/udisks2"
	"pkg.deepin.io/lib/dbus"
)

func getDiskCap() (uint64, error) {
	udisk, err := udisks2.NewObjectManager(
		"org.freedesktop.UDisks2",
		"/org/freedesktop/UDisks2")
	if err != nil {
		return 0, err
	}

	var (
		diskCap   uint64
		driveList []dbus.ObjectPath
	)
	managers, _ := udisk.GetManagedObjects()
	for _, manager := range managers {
		block, ok := manager["org.freedesktop.UDisks2.Block"]
		if !ok {
			continue
		}

		// filter removable disk
		driveValue, _ := block["Drive"]
		drivePath := driveValue.Value().(dbus.ObjectPath)
		if len(drivePath) != 0 && drivePath != "/" {
			if isPathExists(drivePath, driveList) {
				continue
			}
			driveList = append(driveList, drivePath)
			drive := managers[drivePath]["org.freedesktop.UDisks2.Drive"]
			removable := drive["Removable"]
			if !removable.Value().(bool) {
				diskCap += drive["Size"].Value().(uint64)
			}
			continue
		}

		// if no drive exists, such as NVMe port, check whether has PartitionTable
		// TODO: parse NVMe port but no parted. If no parted, will no PartitionTable
		_, ok = manager["org.freedesktop.UDisks2.PartitionTable"]
		if !ok {
			continue
		}

		diskCap += block["Size"].Value().(uint64)
	}

	udisks2.DestroyObjectManager(udisk)
	return diskCap, nil
}

func isPathExists(item dbus.ObjectPath, list []dbus.ObjectPath) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
