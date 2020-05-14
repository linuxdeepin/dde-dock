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
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"syscall"

	polkit "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.policykit1"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.timedate1"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/keyfile"
)

//go:generate dbusutil-gen -type Manager manager.go

type Manager struct {
	core      *timedate1.Timedate
	service   *dbusutil.Service
	PropsMu   sync.RWMutex
	NTPServer string

	methods *struct {
		SetTime      func() `in:"usec,relative,message"`
		SetTimezone  func() `in:"timezone,message"`
		SetLocalRTC  func() `in:"enabled,fixSystem,message"`
		SetNTP       func() `in:"enabled,message"`
		SetNTPServer func() `in:"server,message"`
	}
}

const (
	dbusServiceName = "com.deepin.daemon.Timedated"
	dbusPath        = "/com/deepin/daemon/Timedated"
	dbusInterface   = dbusServiceName

	timedate1ActionId = "org.freedesktop.timedate1.set-time"

	timeSyncCfgFile = "/etc/systemd/timesyncd.conf.d/deepin.conf"
)

func NewManager(service *dbusutil.Service) (*Manager, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	core := timedate1.NewTimedate(systemBus)

	server, err := getNTPServer()
	if err != nil {
		logger.Warning(err)
	}

	err = startNTPbyFirstBoot(core)
	if err != nil {
		logger.Error(err)
	}

	return &Manager{
		core:      core,
		service:   service,
		NTPServer: server,
	}, nil
}

func startNTPbyFirstBoot(core *timedate1.Timedate) error {
	filePath := "/var/lib/dde-daemon/firstBootFile"
	err := syscall.Access(filePath, syscall.F_OK)
	if err != nil {
		err = core.SetNTP(0, true, false)
		if err != nil {
			logger.Error(err)
		}
		logger.Error(err)
		file, e := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
		if e != nil {
			logger.Error("create file error ! !", e)
			return e
		}
		writer := bufio.NewWriter(file)
		_, e = writer.Write([]byte("This system first boot complete, do not delete this file\n"))
		if e != nil {
			logger.Error("write content failed！！", e)
			return e
		}
		e = writer.Flush()
		if e != nil {
			logger.Error("flush content error！！", e)
			return e
		}

	} else {
		logger.Info("file exist")
	}
	return nil
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}

func (m *Manager) destroy() {
	if m.core == nil {
		return
	}
	m.core = nil
}

func (m *Manager) checkAuthorization(method, msg string, sender dbus.Sender) error {
	isAuthorized, err := doAuthorized(msg, string(sender))
	if err != nil {
		logger.Warning("Has error occurred in doAuthorized:", err)
		return err
	}
	if !isAuthorized {
		logger.Warning("Failed to authorize")
		return fmt.Errorf("[%s] Failed to authorize for %v", method, sender)
	}
	return nil
}

func doAuthorized(msg, sysBusName string) (bool, error) {
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return false, err
	}
	authority := polkit.NewAuthority(systemBus)
	subject := polkit.MakeSubject(polkit.SubjectKindSystemBusName)
	subject.SetDetail("name", sysBusName)
	detail := map[string]string{
		"polkit.message": msg,
	}
	ret, err := authority.CheckAuthorization(0, subject, timedate1ActionId,
		detail, polkit.CheckAuthorizationFlagsAllowUserInteraction, "")
	if err != nil {
		return false, err
	}
	return ret.IsAuthorized, nil
}

func setNTPServer(server string) error {
	kf := keyfile.NewKeyFile()
	err := kf.LoadFromFile(timeSyncCfgFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	kf.SetString("Time", "NTP", server)

	dir := filepath.Dir(timeSyncCfgFile)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = kf.SaveToFile(timeSyncCfgFile)
	return err
}

func getNTPServer() (string, error) {
	kf := keyfile.NewKeyFile()
	err := kf.LoadFromFile(timeSyncCfgFile)
	if err != nil && !os.IsNotExist(err) {
		return "", err
	}

	server, _ := kf.GetString("Time", "NTP")
	return server, nil
}

func restartSystemdService(service, mode string) error {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	systemdObj := sysBus.Object("org.freedesktop.systemd1", "/org/freedesktop/systemd1")
	var jobPath dbus.ObjectPath
	err = systemdObj.Call("org.freedesktop.systemd1.Manager.RestartUnit", 0, service, mode).Store(&jobPath)
	return err
}
