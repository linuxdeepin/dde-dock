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
	"os/exec"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/gobject-2.0"
	"strings"
	"time"
)

const (
	TIME_DURATION = 30

	MEDIA_HAND_AUTO_MOUNT = "automount"
	MEDIA_HAND_AUTO_OPEN  = "automount-open"
)

var (
	mediaHandSetting = gio.NewSettings("org.gnome.desktop.media-handling")
)

func (m *Manager) refrashDiskInfoList() {
	for {
		select {
		case <-time.NewTimer(time.Second * TIME_DURATION).C:
			logger.Debug("Refrash Disk Info List")
			m.setPropName("DiskList")
			//logger.Infof("Disk List: %v", m.DiskList)
		case <-m.quitFlag:
			return
		}
	}
}

func (m *Manager) endDiskrefrash() {
	close(m.quitFlag)
}

func (m *Manager) listenSignalChanged() {
	monitor.Connect("mount-added", func(volumeMonitor *gio.VolumeMonitor, mount *gio.Mount) {
		// Judge whether the property 'mount_and_open' set true
		// if true, open the device use exec.Command("xdg-open", "device").Run()
		logger.Info("EVENT: mount added")
		if mount.CanUnmount() &&
			mediaHandSetting.GetBoolean(MEDIA_HAND_AUTO_MOUNT) &&
			mediaHandSetting.GetBoolean(MEDIA_HAND_AUTO_OPEN) {
			uri := mount.GetRoot().GetUri()
			go exec.Command("/usr/bin/xdg-open", uri).Run()
		}
		m.setPropName("DiskList")
	})
	monitor.Connect("mount-removed", func(volumeMonitor *gio.VolumeMonitor, mount *gio.Mount) {
		logger.Info("EVENT: mount removed")
		m.setPropName("DiskList")
	})
	monitor.Connect("mount-changed", func(volumeMonitor *gio.VolumeMonitor, mount *gio.Mount) {
		m.setPropName("DiskList")
	})

	monitor.Connect("volume-added", func(volumeMonitor *gio.VolumeMonitor, volume *gio.Volume) {
		icons := volume.GetIcon().ToString()
		as := strings.Split(icons, " ")
		iconName := ""
		if len(as) > 2 {
			iconName = as[2]
		}
		if (volume.CanEject() || strings.Contains(iconName, "usb")) &&
			mediaHandSetting.GetBoolean(MEDIA_HAND_AUTO_MOUNT) {
			volume.Mount(gio.MountMountFlagsNone, nil, nil, gio.AsyncReadyCallback(func(o *gobject.Object, res *gio.AsyncResult) {
				_, err := volume.MountFinish(res)
				if err != nil {
					logger.Warningf("volume mount failed: %s", err)
					m.setPropName("DiskList")
				}
			}))
		} else {
			m.setPropName("DiskList")
		}
	})
	monitor.Connect("volume-removed", func(volumeMonitor *gio.VolumeMonitor, volume *gio.Volume) {
		m.setPropName("DiskList")
	})
	monitor.Connect("volume-changed", func(volumeMonitor *gio.VolumeMonitor, volume *gio.Volume) {
		m.setPropName("DiskList")
	})

	monitor.Connect("drive-disconnected", func(volumeMonitor *gio.VolumeMonitor, drive *gio.Drive) {
		m.setPropName("DiskList")
	})
	monitor.Connect("drive-connected", func(volumeMonitor *gio.VolumeMonitor, drive *gio.Drive) {
		m.setPropName("DiskList")
	})
	monitor.Connect("drive-changed", func(volumeMonitor *gio.VolumeMonitor, drive *gio.Drive) {
		m.setPropName("DiskList")
	})
}
