/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package timedated

import (
	"fmt"

	"dbus/org/freedesktop/timedate1"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/polkit"
)

type Manager struct {
	core    *timedate1.Timedate1
	service *dbusutil.Service

	methods *struct {
		SetTime     func() `in:"usec,relative,message"`
		SetTimezone func() `in:"timezone,message"`
		SetLocalRTC func() `in:"enabled,fixSystem,message"`
		SetNTP      func() `in:"enabled,message"`
	}
}

const (
	dbusServiceName = "com.deepin.daemon.Timedated"
	dbusPath        = "/com/deepin/daemon/Timedated"
	dbusInterface   = dbusServiceName

	timedate1ActionId = "org.freedesktop.timedate1.set-time"
)

func NewManager(service *dbusutil.Service) (*Manager, error) {
	core, err := timedate1.NewTimedate1("org.freedesktop.timedate1",
		"/org/freedesktop/timedate1")
	if err != nil {
		return nil, err
	}

	polkit.Init()
	return &Manager{
		core:    core,
		service: service,
	}, nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) destroy() {
	if m.core == nil {
		return
	}
	timedate1.DestroyTimedate1(m.core)
	m.core = nil
}

func (m *Manager) checkAuthorization(method, msg string, sender dbus.Sender) error {
	pid, err := m.service.GetConnPID(string(sender))
	if err != nil {
		return err
	}

	isAuthorized, err := doAuthorized(msg, pid)
	if err != nil {
		logger.Warning("Has error occurred in doAuthorized:", err)
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
