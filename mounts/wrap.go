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

// #cgo pkg-config: gio-2.0
// #include <stdlib.h>
// #include "disk_listener.h"
import "C"

import (
	"fmt"
	"os/exec"
	"strings"
	"unsafe"

	"pkg.deepin.io/dde/api/soundutils"
	"pkg.deepin.io/lib/gio-2.0"
)

const (
	gsKeyAutoMount = "automount"
	gsKeyAutoOpen  = "automount-open"
)

func startDiskListener() {
	C.start_disk_listener()
}

//export handleDiskChanged
func handleDiskChanged(event, uuid *C.char) {
	var ev = C.GoString(event)
	var id = C.GoString(uuid)
	logger.Debug("Disk event:", ev, id)
	if len(id) != 0 {
		defer C.free(unsafe.Pointer(uuid))
	}
	if _manager == nil {
		return
	}

	_manager.setPropDiskList(_manager.getDiskInfos())
	switch ev {
	case "volume-added":
		soundutils.PlaySystemSound(soundutils.KeyDevicePlug, "", false)
		handleVolumeAdded(id)
	case "volume-removed":
		soundutils.PlaySystemSound(soundutils.KeyDeviceUnplug, "", false)
	case "mount-added":
		handleMountAdded(id)
	}
}

func handleVolumeAdded(id string) {
	info := _manager.getDiskCache(id)
	if info == nil || info.Type != diskTypeVolume {
		return
	}

	volume := info.Obj.(*gio.Volume)
	iconObj := volume.GetIcon()
	icon := getIconFromGIcon(iconObj)
	iconObj.Unref()

	if (volume.CanEject() || strings.Contains(icon, "usb")) &&
		_manager.isAutoMount() {
		_manager.mountVolume(id, volume)
	}
}

func handleMountAdded(id string) {
	info := _manager.getDiskCache(id)
	if info == nil || info.Type != diskTypeMount {
		return
	}

	mount := info.Obj.(*gio.Mount)
	if mount.CanEject() && _manager.isAutoOpen() {
		root := mount.GetRoot()
		var cmd = fmt.Sprintf("gvfs-open %s", root.GetUri())
		root.Unref()
		go doAction(cmd)
	}
}

func (m *Manager) isAutoMount() bool {
	if m.setting == nil {
		return false
	}

	return m.setting.GetBoolean(gsKeyAutoMount)
}

func (m *Manager) isAutoOpen() bool {
	if m.setting == nil {
		return false
	}

	return m.setting.GetBoolean(gsKeyAutoOpen)
}

func doAction(cmd string) error {
	out, err := exec.Command("/bin/sh", "-c",
		cmd).CombinedOutput()
	if err != nil {
		return fmt.Errorf(string(out))
	}

	return nil
}
