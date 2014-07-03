/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
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
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/gobject-2.0"
)

func (m *Manager) DeviceEject(uuid string) (bool, string) {
	info, ok := objectMap[uuid]
	if !ok {
		Logger.Warning("Eject id - %s not in objectMap.", uuid)
		return false, fmt.Sprintf("Invalid Id: %s\n", uuid)
	}

	Logger.Infof("Eject type: %s", info.Type)
	switch info.Type {
	case "drive":
		op := info.Object.(*gio.Drive)
		op.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
			_, err := op.EjectFinish(res)
			if err != nil {
				m.Error(uuid, err.Error())
				Logger.Warningf("drive eject failed: %s, %s", uuid, err)
			}
		}))
	case "volume":
		op := info.Object.(*gio.Volume)
		op.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
			_, err := op.EjectFinish(res)
			if err != nil {
				m.Error(uuid, err.Error())
				Logger.Warningf("volume eject failed: %s, %s", uuid, err)
			}
		}))
	case "mount":
		op := info.Object.(*gio.Mount)
		op.Eject(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
			_, err := op.EjectFinish(res)
			if err != nil {
				m.Error(uuid, err.Error())
				Logger.Warningf("mount eject failed: %s, %s", uuid, err)
			}
		}))
	default:
		Logger.Errorf("'%s' invalid type", info.Type)
		return false, fmt.Sprintf("Invalid type: '%s'\n", info.Type)
	}

	return true, ""
}

func (m *Manager) DeviceMount(uuid string) (bool, string) {
	info, ok := objectMap[uuid]
	if !ok {
		Logger.Warning("Mount id - %s not in objectMap.", uuid)
		return false, fmt.Sprintf("Invalid Id: %s\n", uuid)
	}

	Logger.Infof("Mount type: %s", info.Type)
	switch info.Type {
	case "volume":
		op := info.Object.(*gio.Volume)
		op.Mount(gio.MountMountFlagsNone, nil, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
			_, err := op.MountFinish(res)
			if err != nil {
				m.Error(uuid, err.Error())
				Logger.Warningf("volume mount failed: %s, %s", uuid, err)
			}
		}))
	case "mount":
		op := info.Object.(*gio.Mount)
		op.Remount(gio.MountMountFlagsNone, nil, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
			_, err := op.RemountFinish(res)
			if err != nil {
				m.Error(uuid, err.Error())
				Logger.Warningf("mount remount failed: %s, %s", uuid, err)
			}
		}))
	default:
		Logger.Errorf("'%s' invalid type", info.Type)
		return false, fmt.Sprintf("Invalid type: '%s'\n", info.Type)
	}

	return true, ""
}

func (m *Manager) DeviceUnmount(uuid string) (bool, string) {
	info, ok := objectMap[uuid]
	if !ok {
		Logger.Warningf("Unmount id - %s not in objectMap.", uuid)
		return false, fmt.Sprintf("Invalid Id: %s\n", uuid)
	}

	Logger.Infof("Unmount type: %s", info.Type)
	switch info.Type {
	case "mount":
		op := info.Object.(*gio.Mount)
		op.Unmount(gio.MountUnmountFlagsNone, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
			_, err := op.UnmountFinish(res)
			if err != nil {
				m.Error(uuid, err.Error())
				Logger.Warningf("mount unmount failed: %s, %s", uuid, err)
			}
		}))
	default:
		Logger.Errorf("'%s' invalid type", info.Type)
		return false, fmt.Sprintf("Invalid type: '%s'\n", info.Type)
	}

	return true, ""
}
