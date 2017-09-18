/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package systeminfo

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

type SystemInfo struct {
	// Current deepin version, ex: "2015 Desktop"
	Version string
	// Distribution ID
	DistroID string
	// Distribution Description
	DistroDesc string
	// Distribution Version
	DistroVer string
	// CPU information
	Processor string
	// Disk capacity
	DiskCap uint64
	// Memory size
	MemoryCap uint64
	// System type: 32bit or 64bit
	SystemType int64
}

type Daemon struct {
	info *SystemInfo
	*loader.ModuleBase
}

var logger = log.NewLogger("daemon/systeminfo")

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("systeminfo", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if d.info != nil {
		return nil
	}

	logger.BeginTracing()

	d.info = NewSystemInfo()
	err := dbus.InstallOnSession(d.info)
	if err != nil {
		d.info = nil
		logger.Error(err)
		logger.EndTracing()
		return err
	}
	return nil
}

func (d *Daemon) Stop() error {
	if d.info == nil {
		return nil
	}

	dbus.UnInstallObject(d.info)
	d.info = nil
	logger.EndTracing()

	return nil
}

func NewSystemInfo() *SystemInfo {
	var info SystemInfo
	tmp, _ := doReadCache(cacheFile)
	if tmp != nil && tmp.isValidity() {
		info = *tmp
		time.AfterFunc(time.Second*10, func() {
			info.init()
			doSaveCache(&info, cacheFile)
		})
		return &info
	}

	info.init()
	go doSaveCache(&info, cacheFile)
	return &info
}

func (info *SystemInfo) init() {
	var err error
	info.Processor, err = GetCPUInfo("/proc/cpuinfo")
	if err != nil {
		logger.Warning("Get cpu info failed:", err)
	}

	info.Version, err = getVersion()
	if err != nil {
		logger.Warning("Get version failed:", err)
	}

	info.DistroID, info.DistroDesc, info.DistroVer, err = getDistro()
	if err != nil {
		logger.Warning("Get distribution failed:", err)
	}

	info.MemoryCap, err = getMemoryFromFile("/proc/meminfo")
	if err != nil {
		logger.Warning("Get memory capacity failed:", err)
	}

	if systemBit() == "64" {
		info.SystemType = 64
	} else {
		info.SystemType = 32
	}

	info.DiskCap, err = getDiskCap()
	if err != nil {
		logger.Warning("Get disk capacity failed:", err)
	}
}

func (info *SystemInfo) isValidity() bool {
	if info.Processor == "" || info.DiskCap == 0 || info.MemoryCap == 0 {
		return false
	}
	return true
}

func (*SystemInfo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.SystemInfo",
		ObjectPath: "/com/deepin/daemon/SystemInfo",
		Interface:  "com.deepin.daemon.SystemInfo",
	}
}
