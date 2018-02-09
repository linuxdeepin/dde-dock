/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package debug

import (
	"gir/gio-2.0"
	"pkg.deepin.io/dde/daemon/calltrace"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
	"time"
)

var (
	logger = log.NewLogger("daemon/Debug")
)

type Daemon struct {
	ct   *calltrace.CallTrace
	quit chan bool
	*loader.ModuleBase
}

func NewDaemon() *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("debug", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if d.quit != nil {
		return nil
	}

	d.quit = make(chan bool)
	if d.LogLevel() != log.LevelDebug {
		loader.ToggleLogDebug(true)
	}

	go d.loop()
	return nil
}

func (d *Daemon) Stop() error {
	if d.quit == nil {
		return nil
	}

	d.quit <- true
	if d.ct != nil {
		d.ct.SetAutoDestroy(1)
	}
	if d.LogLevel() == log.LevelDebug {
		loader.ToggleLogDebug(false)
	}
	logger.Info("--------Terminate process stat loop")
	return nil
}

func (d *Daemon) loop() {
	s := gio.NewSettings("com.deepin.dde.debug")
	cpuPercentage := s.GetInt("cpu-percentage")
	memUsage := s.GetInt("mem-usage")
	duration := s.GetInt("duration")
	s.Unref()

	logger.Info("--------Start process stat loop")
	d.handleProcessStat(cpuPercentage, memUsage, duration)
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			d.handleProcessStat(cpuPercentage, memUsage, duration)
		case <-d.quit:
			ticker.Stop()
			close(d.quit)
			d.quit = nil
			return
		}
	}
}

func (d *Daemon) handleProcessStat(cpuPercentage, memUsage, duration int32) {
	cpu, _ := getCPUPercentage()
	mem, _ := getMemoryUsage()
	logger.Infof("-----------Handle process stat, cpu: %#v, mem: %#v, ct: %p", cpu, mem, d.ct)
	if cpu > float64(cpuPercentage) || mem > int64(memUsage)*1024 {
		if d.ct == nil {
			d.ct, _ = calltrace.Start(uint32(duration), logger)
		}
	} else {
		if d.ct == nil {
			return
		}
		d.ct.SetAutoDestroy(1)
	}
}
