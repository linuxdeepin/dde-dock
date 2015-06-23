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
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/log"
	dutils "pkg.linuxdeepin.com/lib/utils"
	"sync"
	"time"
)

var logger = log.NewLogger("dde-daemon/mounts")

const (
	mediaHandlerSchema = "org.gnome.desktop.media-handling"

	diskTypeVolume int = iota + 1
	diskTypeMount

	refrashTimeout = 90
)

type diskObjectInfo struct {
	Type int
	Obj  interface{}
}

type Manager struct {
	DiskList DiskInfos

	//Error(uuid, reason)
	Error func(string, string)

	monitor *gio.VolumeMonitor
	setting *gio.Settings
	logger  *log.Logger
	endFlag chan struct{}

	cacheLocker sync.Mutex
	diskCache   map[string]*diskObjectInfo
}

func NewManager() *Manager {
	var m = Manager{}

	m.logger = logger
	m.monitor = gio.VolumeMonitorGet()
	m.setting, _ = dutils.CheckAndNewGSettings(mediaHandlerSchema)
	m.diskCache = make(map[string]*diskObjectInfo)
	m.endFlag = make(chan struct{})

	m.DiskList = m.getDiskInfos()

	return &m
}

func (m *Manager) destroy() {
	if m.diskCache != nil {
		m.clearDiskCache()
		m.diskCache = nil
	}

	if m.monitor != nil {
		m.monitor.Unref()
		m.monitor = nil
	}

	if m.logger != nil {
		m.logger.EndTracing()
		m.logger = nil
	}

	if m.endFlag != nil {
		close(m.endFlag)
		m.endFlag = nil
	}
}

func (m *Manager) refrashDiskInfos() {
	for {
		select {
		case <-time.NewTicker(time.Second * refrashTimeout).C:
			m.setPropDiskList(m.getDiskInfos())
		case <-m.endFlag:
			return
		}
	}
}

func (m *Manager) getDiskInfos() DiskInfos {
	m.clearDiskCache()
	m.diskCache = make(map[string]*diskObjectInfo)

	var infos DiskInfos
	volumes := m.monitor.GetVolumes()
	for _, volume := range volumes {
		mount := volume.GetMount()
		if mount != nil {
			mount.Unref()
			continue
		}

		info := newDiskInfoFromVolume(volume)
		m.setDiskCache(info.UUID, &diskObjectInfo{
			Type: diskTypeVolume,
			Obj:  volume,
		})
		infos = append(infos, info)
	}

	mounts := m.monitor.GetMounts()
	for _, mount := range mounts {
		info := newDiskInfoFromMount(mount)
		m.setDiskCache(info.UUID, &diskObjectInfo{
			Type: diskTypeMount,
			Obj:  mount,
		})
		infos = append(infos, info)
	}

	return infos
}

func (m *Manager) setDiskCache(key string, value *diskObjectInfo) {
	m.cacheLocker.Lock()
	defer m.cacheLocker.Unlock()
	_, ok := m.diskCache[key]
	if ok {
		m.deleteDiskCache(key)
	}

	m.diskCache[key] = value
}

func (m *Manager) getDiskCache(key string) *diskObjectInfo {
	m.cacheLocker.Lock()
	defer m.cacheLocker.Unlock()
	v, ok := m.diskCache[key]
	if !ok {
		return nil
	}
	return v
}

func (m *Manager) deleteDiskCache(key string) {
	m.cacheLocker.Lock()
	defer m.cacheLocker.Unlock()
	v, ok := m.diskCache[key]
	if !ok {
		return
	}

	switch v.Type {
	case diskTypeVolume:
		volume := v.Obj.(*gio.Volume)
		volume.Unref()
	case diskTypeMount:
		mount := v.Obj.(*gio.Mount)
		mount.Unref()
	}
	delete(m.diskCache, key)
}

func (m *Manager) clearDiskCache() {
	m.cacheLocker.Lock()
	defer m.cacheLocker.Unlock()
	for _, v := range m.diskCache {
		switch v.Type {
		case diskTypeVolume:
			volume := v.Obj.(*gio.Volume)
			volume.Unref()
		case diskTypeMount:
			mount := v.Obj.(*gio.Mount)
			mount.Unref()
		}
	}
}
