/**
 * Copyright (C) 2017 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package timedated

import (
	"dbus/org/freedesktop/timedate1"
	"fmt"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/polkit"
)

type Manager struct {
	core *timedate1.Timedate1
}

const (
	dbusDest = "com.deepin.daemon.Timedated"
	dbusPath = "/com/deepin/daemon/Timedated"
	dbusIFC  = dbusDest

	timedate1ActionId = "org.freedesktop.timedate1.set-time"
)

func NewManager() (*Manager, error) {
	core, err := timedate1.NewTimedate1("org.freedesktop.timedate1",
		"/org/freedesktop/timedate1")
	if err != nil {
		return nil, err
	}

	polkit.Init()
	return &Manager{
		core: core,
	}, nil
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func (m *Manager) destroy() {
	if m.core == nil {
		return
	}
	timedate1.DestroyTimedate1(m.core)
	m.core = nil
}

func (m *Manager) checkAuthorization(method, msg string, pid uint32) error {
	isAuthorized, err := doAuthorized(msg, pid)
	if err != nil {
		logger.Warning("Has error occured in doAuthorized:", err)
		return err
	}
	if !isAuthorized {
		logger.Warning("Failed to authorize")
		return fmt.Errorf("[%s] Failed to authorize for %v", method, pid)
	}
	return nil
}

func doAuthorized(msg string, pid uint32) (bool, error) {
	var subject = polkit.NewSubject(polkit.SubjectKindUnixProcess)
	subject.SetDetail("pid", pid)
	var t = uint64(0)
	subject.SetDetail("start-time", t)
	var detail = make(map[string]string)
	detail["polkit.gettext_domain"] = "dde-daemon"
	detail["polkit.message"] = msg
	var cancelId string
	ret, err := polkit.CheckAuthorization(subject, timedate1ActionId,
		detail, polkit.CheckAuthorizationFlagsAllowUserInteraction, cancelId)
	subject = nil
	if err != nil {
		return false, err
	}
	return ret.IsAuthorized, nil
}
