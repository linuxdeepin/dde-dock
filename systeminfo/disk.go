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

	driList := getDriverList(udisk)
	if len(driList) == 0 {
		return 0, nil
	}

	var diskCap uint64
	managers, _ := udisk.GetManagedObjects()
	for _, driver := range driList {
		_, driExist := managers[driver]
		rm, _ := managers[driver]["org.freedesktop.UDisks2.Drive"]["Removable"]
		if driExist && !(rm.Value().(bool)) {
			size := managers[driver]["org.freedesktop.UDisks2.Drive"]["Size"]
			diskCap += size.Value().(uint64)
		}
	}

	udisks2.DestroyObjectManager(udisk)
	return diskCap, nil
}

func getDriverList(udisk *udisks2.ObjectManager) []dbus.ObjectPath {
	var driList []dbus.ObjectPath

	managers, _ := udisk.GetManagedObjects()
	for _, value := range managers {
		if _, ok := value["org.freedesktop.UDisks2.Block"]; ok {
			v := value["org.freedesktop.UDisks2.Block"]["Drive"]
			path := v.Value().(dbus.ObjectPath)
			if path != dbus.ObjectPath("/") {
				flag := false
				l := len(driList)
				for i := 0; i < l; i++ {
					if driList[i] == path {
						flag = true
						break
					}
				}
				if !flag {
					driList = append(driList, path)
				}
			}
		}
	}
	return driList
}
