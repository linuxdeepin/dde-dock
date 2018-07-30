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

package power

import (
	"os/exec"
	"time"

	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/dpms"
	"github.com/linuxdeepin/go-x11-client/util/wm/icccm"
	"pkg.deepin.io/dde/api/soundutils"
)

func (m *Manager) findWindow(wmClassInstance, wmClassClass string) x.Window {
	c := m.helper.xConn
	rootWin := c.GetDefaultScreen().Root
	tree, err := x.QueryTree(c, rootWin).Reply(c)
	if err != nil {
		logger.Warning("QueryTree error:", err)
		return 0
	}

	for _, win := range tree.Children {
		wmClass, err := icccm.GetWMClass(c, win).Reply(c)
		if err == nil &&
			wmClass.Instance == wmClassInstance &&
			wmClass.Class == wmClassClass {
			return win
		}
	}
	return 0
}

func (m *Manager) waitWindowViewable(wmClassInstance, wmClassClass string, timeout time.Duration) {
	c := m.helper.xConn
	logger.Debug("waitWindowViewable", wmClassInstance, wmClassClass)
	ticker := time.NewTicker(time.Millisecond * 300)
	timer := time.NewTimer(timeout)
	for {
		select {
		case _, ok := <-ticker.C:
			if !ok {
				logger.Error("Invalid ticker event")
				return
			}

			logger.Debug("waitWindowViewable tick")
			win := m.findWindow(wmClassInstance, wmClassClass)
			if win == 0 {
				continue
			}

			winAttr, err := x.GetWindowAttributes(c, win).Reply(c)
			if err != nil {
				logger.Warning(err)
				continue
			}
			if winAttr.MapState == x.MapStateViewable {
				logger.Debug("waitWindowViewable found")
				ticker.Stop()
				return
			}

		case <-timer.C:
			logger.Debug("waitWindowViewable failed timeout!")
			ticker.Stop()
			return
		}
	}
}

func (m *Manager) lockWaitShow(timeout time.Duration) {
	const ddeLock = "dde-lock"
	m.doLock()
	m.waitWindowViewable(ddeLock, ddeLock, timeout)
}

func (m *Manager) setDPMSModeOn() {
	logger.Info("DPMS On")
	c := m.helper.xConn
	err := dpms.ForceLevelChecked(c, dpms.DPMSModeOn).Check(c)
	if err != nil {
		logger.Warning("Set DPMS on error:", err)
	}
}

func (m *Manager) setDPMSModeOff() {
	logger.Info("DPMS Off")
	c := m.helper.xConn
	err := dpms.ForceLevelChecked(c, dpms.DPMSModeOff).Check(c)
	if err != nil {
		logger.Warning("Set DPMS off error:", err)
	}
}

func (m *Manager) doLock() {
	logger.Info("Lock Screen")
	sessionManager := m.helper.SessionManager
	if sessionManager != nil {
		err := sessionManager.RequestLock(0)
		if err != nil {
			logger.Error("Lock failed:", err)
		}
	}
}

func (m *Manager) doSuspend() {
	logger.Debug("Suspend")
	sessionManager := m.helper.SessionManager
	if sessionManager != nil {
		err := sessionManager.RequestSuspend(0)
		if err != nil {
			logger.Error("Suspend failed:", err)
		}
	}
}

func (m *Manager) setDisplayBrightness(brightnessTable map[string]float64) {
	display := m.helper.Display
	for output, brightness := range brightnessTable {
		logger.Infof("Change output %q brightness to %.2f", output, brightness)
		err := display.SetBrightness(0, output, brightness)
		if err != nil {
			logger.Warningf("Change output %q brightness to %.2f failed: %v", output, brightness, err)
		} else {
			logger.Infof("Change output %q brightness to %.2f done!", output, brightness)
		}
	}
}

func (m *Manager) setAndSaveDisplayBrightness(brightnessTable map[string]float64) {
	display := m.helper.Display
	for output, brightness := range brightnessTable {
		logger.Infof("Change output %q brightness to %.2f", output, brightness)
		err := display.SetAndSaveBrightness(0, output, brightness)
		if err != nil {
			logger.Warningf("Change output %q brightness to %.2f failed: %v", output, brightness, err)
		} else {
			logger.Infof("Change output %q brightness to %.2f done!", output, brightness)
		}
	}
}

func doShowDDELowPower() {
	logger.Info("Show dde low power")
	go exec.Command(cmdDDELowPower, "--raise").Run()
}

func doCloseDDELowPower() {
	logger.Info("Close low power")
	go exec.Command(cmdDDELowPower, "--quit").Run()
}

func (m *Manager) sendNotify(icon, summary, body string) {
	n := m.helper.Notification
	n.Update(summary, body, icon)

	go func() {
		err := n.Show()
		logger.Debugf("sendNotify icon: %q, summary: %q, body: %q", icon, summary, body)
		if err != nil {
			logger.Warning("sendNotify failed:", err)
		}
	}()
}

func playSound(name string) {
	logger.Debug("play system sound", name)
	go soundutils.PlaySystemSound(name, "")
}
