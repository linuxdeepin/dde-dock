/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	libdisplay "dbus/com/deepin/daemon/display"
	"github.com/BurntSushi/xgb/dpms"
	"os/exec"
	"pkg.deepin.io/dde/api/soundutils"
	"time"
)

func getBatteryPowerLevelName(num uint32) string {
	if name, ok := batteryPowerLevelNameMap[num]; !ok {
		return "UnknownLevel"
	} else {
		return name
	}
}

func (m *Manager) setDPMSModeOn() {
	logger.Debug("DPMS On")
	dpms.ForceLevel(m.xConn, dpms.DPMSModeOn)
}

func (m *Manager) setDPMSModeOff() {
	logger.Debug("DPMS Off")
	dpms.ForceLevel(m.xConn, dpms.DPMSModeOff)
}

func (m *Manager) doLock() {
	logger.Debug("Lock Screen")
	if m.sessionManager != nil {
		err := m.sessionManager.RequestLock()
		if err != nil {
			logger.Error("Lock failed:", err)
		}

		// wait dde-lock show
		if m.lockFront != nil {
			lockWaiter := &struct {
				ch chan int
			}{make(chan int)}
			go func() {
				for lockWaiter.ch != nil {
					time.Sleep(300 * time.Millisecond)
					logger.Debug("check lock result")
					_, err = m.lockFront.LockResult()
					if err == nil {
						if lockWaiter.ch != nil {
							lockWaiter.ch <- 1
						}
						break
					}
				}
			}()

			select {
			case <-time.After(3 * time.Second):
				logger.Warning("lock timeout")
			case <-lockWaiter.ch:
				logger.Debug("lock done")
			}
			close(lockWaiter.ch)
			lockWaiter.ch = nil
		}
	}
}

func (m *Manager) doShutdown() {
	if m.sessionManager != nil {
		err := m.sessionManager.RequestShutdown()
		if err != nil {
			logger.Error("Shutdown failed:", err)
		}
	}
}

func (m *Manager) doSuspend() {
	logger.Debug("Suspend")
	if m.sessionManager != nil {
		err := m.sessionManager.RequestSuspend()
		if err != nil {
			logger.Error("Suspend failed:", err)
		}
	}
}

func (m *Manager) setDisplayBrightness(brightness float64) {
	outputNames := []string{}
	if m.display != nil {
		for _, objPath := range m.display.Monitors.Get() {
			monitor, err := libdisplay.NewMonitor(dbusDisplayDest, objPath)
			if err != nil {
				logger.Error("NewMonitor failed:", err)
				continue
			}
			defer libdisplay.DestroyMonitor(monitor)
			for _, name := range monitor.Outputs.Get() {
				outputNames = append(outputNames, name)
			}
		}
		for _, name := range outputNames {
			logger.Debugf("Change output %q brightness to %.2f", name, brightness)
			m.display.ChangeBrightness(name, brightness)
		}
	}
}

func doShowLowpower() {
	logger.Debug("Show low power")
	go exec.Command(cmdLowPower, "--raise").Run()
}

func execCommand(cmd string) {
	logger.Debugf("exec %q", cmd)
	go exec.Command("sh", "-c", cmd).Run()
}

func doCloseLowpower() {
	logger.Debug("Close low power")
	go exec.Command(cmdLowPower, "--quit").Run()
}

func doShutDownInteractive() {
	go exec.Command("dde-shutdown").Run()
}

func (m *Manager) isMultiScreen() bool {
	if m.display != nil {
		monitorObjPaths := m.display.Monitors.Get()
		if len(monitorObjPaths) > 1 {
			return true
		} else if len(monitorObjPaths) == 1 {
			monitor, err := libdisplay.NewMonitor(dbusDisplayDest, monitorObjPaths[0])
			if err == nil {
				// NOTE: 复制屏时只有一个合一的 monitor， IsComposited 属性为 true
				if monitor.IsComposited.Get() {
					return true
				}
				return false
			}
			logger.Warning(err)
			return false
		}
	}
	return false
}

func (m *Manager) sendNotify(icon, summary, body string) {
	go func() {
		if m.notifier != nil {
			m.notifier.Notify(dbusDest, 0, icon, summary, body, nil, nil, 0)
			logger.Debug("send notify ", icon, summary, body)
		}
	}()
}

func playSound(name string) {
	logger.Debug("play sound", name)
	soundutils.PlaySystemSound(name, "", false)
}
