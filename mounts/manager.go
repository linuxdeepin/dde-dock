/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mounts

import (
	libnotifications "dbus/org/freedesktop/notifications"
	"gir/gio-2.0"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/utils"
	"sync"
	"time"
)

const (
	EventTypeVolumeAdded int32 = iota + 1
	EventTypeVolumeRemoved
	EventTypeMountAdded
	EventTypeMountRemoved
	EventTypeVolumeChanged
)

const (
	mediaHandlerSchema = "org.gnome.desktop.media-handling"

	refreshDuration = time.Minute * 10
)

type Manager struct {
	DiskList DiskInfos
	Changed  func(int32, string)  // (eventType, id)
	Error    func(string, string) // (id, message)

	monitor *gio.VolumeMonitor
	setting *gio.Settings
	quit    chan struct{}

	refreshLocker sync.Mutex
	notifier      *libnotifications.Notifier
}

func newManager() *Manager {
	var m = new(Manager)
	m.monitor = gio.VolumeMonitorGet()
	m.DiskList = m.getDiskInfos()
	m.init()
	m.setting, _ = utils.CheckAndNewGSettings(mediaHandlerSchema)
	m.quit = make(chan struct{})
	return m
}

func (m *Manager) init() {
	for _, info := range m.DiskList {
		if info.CanUnmount || info.Type != DiskTypeRemovable {
			continue
		}

		logger.Debug("[init] will mount volume:", info.Name, info.Id)
		err := m.Mount(info.Id)
		if err != nil {
			logger.Warningf("Mount '%s - %s' failed: %v",
				info.Name, info.Id, err)
		}
	}
	m.refreshDiskList()
	var err error
	m.notifier, err = libnotifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if err != nil {
		logger.Warning("init notifier failed:", err)
	}
}

func (m *Manager) destroy() {
	if m.quit != nil {
		close(m.quit)
		m.quit = nil
	}

	if m.setting != nil {
		m.setting.Unref()
		m.setting = nil
	}

	m.DiskList = nil
	if m.monitor != nil {
		m.monitor.Unref()
		m.monitor = nil
	}

	if m.notifier != nil {
		libnotifications.DestroyNotifier(m.notifier)
		m.notifier = nil
	}
}

func (m *Manager) emitError(id, msg string) {
	logger.Warningf("emitError id %q msg: %v", id, msg)
	dbus.Emit(m, "Error", id, msg)
}

func (m *Manager) refreshDiskList() {
	m.refreshLocker.Lock()
	defer m.refreshLocker.Unlock()
	m.setPropDiskList(m.getDiskInfos())
}

func (m *Manager) updateDiskInfo() {
	for {
		select {
		case <-time.After(refreshDuration):
			m.refreshDiskList()
		case <-m.quit:
			return
		}
	}
}

func (m *Manager) sendNotify(icon, summary, body string) {
	notifier := m.notifier
	if notifier != nil {
		notifier.Notify(dbusDest, 0, icon, summary, body, nil, nil, 0)
		logger.Infof("send notify icon: %q, summary: %q, body: %q", icon, summary, body)
	} else {
		logger.Warning("send notify failed: notifier is nil")
	}
}
