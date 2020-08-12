/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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
	"time"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

const (
	dbusServiceName = "com.deepin.daemon.SystemInfo"
	dbusPath        = "/com/deepin/daemon/SystemInfo"
	dbusInterface   = dbusServiceName
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
	// CPU max MHz
	CPUMaxMHz float64
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
	service := loader.GetService()

	d.info = NewSystemInfo()

	err := service.Export(dbusPath, d.info)
	if err != nil {
		d.info = nil
		logger.Error(err)
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		d.info = nil
		logger.Error(err)
		return err
	}
	return nil
}

func (d *Daemon) Stop() error {
	if d.info == nil {
		return nil
	}

	service := loader.GetService()
	_ = service.StopExport(d.info)
	d.info = nil

	return nil
}

func NewSystemInfo() *SystemInfo {
	var info SystemInfo
	tmp, _ := doReadCache(cacheFile)
	if tmp != nil && tmp.isValidity() {
		info = *tmp
		time.AfterFunc(time.Second*10, func() {
			info.init()
			_ = doSaveCache(&info, cacheFile)
		})
		return &info
	}

	info.init()
	go func() {
		_ = doSaveCache(&info, cacheFile)
	}()
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

	lscpuRes, err := runLscpu()
	if err != nil {
		logger.Warning("run lscpu failed:", err)
		return
	}

	info.CPUMaxMHz, err = getCPUMaxMHzByLscpu(lscpuRes)
	if err != nil {
		logger.Warning(err)
	}

	if info.Processor == "" {
		info.Processor, err = getProcessorByLscpu(lscpuRes)
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (info *SystemInfo) isValidity() bool {
	if info.Processor == "" || info.DiskCap == 0 || info.MemoryCap == 0 {
		return false
	}
	return true
}

func (*SystemInfo) GetInterfaceName() string {
	return dbusInterface
}
