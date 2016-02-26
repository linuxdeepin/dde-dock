/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mounts

import (
	"sync"
	"time"

	"gir/gio-2.0"
	"pkg.deepin.io/lib/log"
	dutils "pkg.deepin.io/lib/utils"
)

var logger = log.NewLogger("daemon/mounts")

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
	// All disk info list in system
	DiskList DiskInfos

	// Error(uuid, reason) signal. It will be emited if operation failure
	//
	// uuid: the disk uuid
	// reason: detail info about the failure
	Error func(string, string)

	setting *gio.Settings
	logger  *log.Logger
	endFlag chan struct{}

	locker      sync.Mutex
	cacheLocker sync.Mutex
	diskCache   map[string]*diskObjectInfo
}

func NewManager() *Manager {
	var m = Manager{}

	m.logger = logger
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
	m.locker.Lock()
	defer m.locker.Unlock()

	m.clearDiskCache()
	m.diskCache = make(map[string]*diskObjectInfo)

	var infos DiskInfos
	var monitor = gio.VolumeMonitorGet()
	defer monitor.Unref()
	volumes := monitor.GetVolumes()
	for _, volume := range volumes {
		mount := volume.GetMount()
		if mount != nil {
			continue
		}

		info := newDiskInfoFromVolume(volume)
		m.setDiskCache(info.UUID, &diskObjectInfo{
			Type: diskTypeVolume,
			Obj:  volume,
		})
		infos = append(infos, info)
	}

	mounts := monitor.GetMounts()
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
	tmp, ok := m.diskCache[key]
	if ok {
		freeDiskInfoObj(tmp)
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

	freeDiskInfoObj(v)
	delete(m.diskCache, key)
}

func (m *Manager) clearDiskCache() {
	m.cacheLocker.Lock()
	defer m.cacheLocker.Unlock()
	for _, v := range m.diskCache {
		freeDiskInfoObj(v)
	}
	m.diskCache = nil
}

func freeDiskInfoObj(v *diskObjectInfo) {
	switch v.Type {
	case diskTypeVolume:
		volume := v.Obj.(*gio.Volume)
		volume.Unref()
	case diskTypeMount:
		mount := v.Obj.(*gio.Mount)
		mount.Unref()
	}
}
