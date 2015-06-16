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
	"os/exec"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"strings"
)

const (
	gsKeyAutoMount = "automount"
	gsKeyAutoOpen  = "automount-open"
)

func (m *Manager) listenDiskChanged() {
	m.monitor.Connect("mount-added", func(monitor *gio.VolumeMonitor, mount *gio.Mount) {
		if mount.CanEject() && m.isAutoOpen() {
			root := mount.GetRoot()
			var cmd = fmt.Sprintf("xdg-open %s", root.GetUri())
			root.Unref()
			go doAction(cmd)
			//err := doAction(cmd)
			//if err != nil {
			//m.logger.Warningf("Exec '%s' failed: %v",
			//cmd, err)
			//}
		}
		m.setPropDiskList(m.getDiskInfos())
	})

	m.monitor.Connect("mount-removed", func(monitor *gio.VolumeMonitor, mount *gio.Mount) {
		m.setPropDiskList(m.getDiskInfos())
	})

	m.monitor.Connect("volume-added", func(monitor *gio.VolumeMonitor, volume *gio.Volume) {
		iconObj := volume.GetIcon()
		icon := getIconFromGIcon(iconObj)
		iconObj.Unref()

		if (volume.CanEject() || strings.Contains(icon, "usb")) &&
			m.isAutoMount() {
			m.mountVolume("", volume)
		}
		m.setPropDiskList(m.getDiskInfos())
	})

	m.monitor.Connect("volume-removed", func(monitor *gio.VolumeMonitor, volume *gio.Volume) {
		m.setPropDiskList(m.getDiskInfos())
	})
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
