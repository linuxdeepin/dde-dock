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
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/icccm"
	"os/exec"
	"pkg.deepin.io/dde/api/soundutils"
	"time"
)

func (m *Manager) findWindow(wminstance, wmclass string) bool {
	rootWin := m.xu.RootWin()
	tree, err := xproto.QueryTree(m.xConn, rootWin).Reply()
	if err != nil {
		logger.Warning("QueryTree error:", err)
		return false
	}
	for i := int(tree.ChildrenLen) - 1; i >= 0; i-- {
		wmClass, err := icccm.WmClassGet(m.xu, tree.Children[i])
		if err == nil &&
			wmClass.Instance == wminstance &&
			wmClass.Class == wmclass {
			return true
		}
	}
	return false
}

func (m *Manager) waitWindow(wminstance, wmclass string, timeout time.Duration) {
	logger.Debug("waitWindow", wminstance, wmclass)
	ticker := time.NewTicker(time.Millisecond * 300)
	timer := time.NewTimer(timeout)
	for {
		select {
		case <-ticker.C:
			logger.Debug("waitWindow tick")
			if m.findWindow(wminstance, wmclass) {
				logger.Debug("waitWindow found")
				ticker.Stop()
				return
			}

		case <-timer.C:
			logger.Debug("waitWindow failed timeout!")
			ticker.Stop()
			return
		}
	}
}

func (m *Manager) lockWaitShow(timeout time.Duration) {
	const ddeLock = "dde-lock"
	m.doLock()
	m.waitWindow(ddeLock, ddeLock, timeout)
}

func getBatteryPowerLevelName(num uint32) string {
	if name, ok := batteryPowerLevelNameMap[num]; !ok {
		return "UnknownLevel"
	} else {
		return name
	}
}

func (m *Manager) setDPMSModeOn() {
	logger.Debug("DPMS On")
	err := dpms.ForceLevelChecked(m.xConn, dpms.DPMSModeOn).Check()
	if err != nil {
		logger.Warning("Set DPMS on error:", err)
	}
}

func (m *Manager) setDPMSModeOff() {
	logger.Debug("DPMS Off")
	err := dpms.ForceLevelChecked(m.xConn, dpms.DPMSModeOff).Check()
	if err != nil {
		logger.Warning("Set DPMS off error:", err)
	}
}

func (m *Manager) doLock() {
	logger.Debug("Lock Screen")
	if m.sessionManager != nil {
		err := m.sessionManager.RequestLock()
		if err != nil {
			logger.Error("Lock failed:", err)
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

func (m *Manager) getDisplayOutputs() []string {
	outputs := []string{}
	if m.display != nil {
		for _, objPath := range m.display.Monitors.Get() {
			monitor, err := libdisplay.NewMonitor(dbusDisplayDest, objPath)
			if err != nil {
				logger.Error("NewMonitor failed:", err)
				continue
			}
			defer libdisplay.DestroyMonitor(monitor)
			for _, name := range monitor.Outputs.Get() {
				outputs = append(outputs, name)
			}
		}
	}
	return outputs
}

func (m *Manager) setDisplayBrightness(brightnessTable map[string]float64) {
	for _, output := range m.getDisplayOutputs() {
		brightness, ok := brightnessTable[output]
		if ok {
			logger.Debugf("Change output %q brightness to %.2f", output, brightness)
			m.display.ChangeBrightness(output, brightness)
		}
	}
}

func (m *Manager) setDisplaySameBrightness(brightness float64) {
	for _, output := range m.getDisplayOutputs() {
		logger.Debugf("Change output %q brightness to %.2f", output, brightness)
		m.display.ChangeBrightness(output, brightness)
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
