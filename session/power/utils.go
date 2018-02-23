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
	"github.com/BurntSushi/xgb/dpms"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil/icccm"
	"os/exec"
	"pkg.deepin.io/dde/api/soundutils"
	"time"
)

func (m *Manager) findWindow(wminstance, wmclass string) xproto.Window {
	xu := m.helper.xu
	rootWin := xu.RootWin()
	tree, err := xproto.QueryTree(xu.Conn(), rootWin).Reply()
	if err != nil {
		logger.Warning("QueryTree error:", err)
		return 0
	}
	for i := int(tree.ChildrenLen) - 1; i >= 0; i-- {
		win := tree.Children[i]
		wmClass, err := icccm.WmClassGet(xu, win)
		if err == nil &&
			wmClass.Instance == wminstance &&
			wmClass.Class == wmclass {
			return win
		}
	}
	return 0
}

func (m *Manager) waitWindowViewable(wminstance, wmclass string, timeout time.Duration) {
	xu := m.helper.xu
	logger.Debug("waitWindowViewable", wminstance, wmclass)
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
			win := m.findWindow(wminstance, wmclass)
			if win == 0 {
				continue
			}

			winAttr, err := xproto.GetWindowAttributes(xu.Conn(), win).Reply()
			if err != nil {
				logger.Warning(err)
				continue
			}
			if winAttr.MapState == xproto.MapStateViewable {
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
	xu := m.helper.xu
	err := dpms.ForceLevelChecked(xu.Conn(), dpms.DPMSModeOn).Check()
	if err != nil {
		logger.Warning("Set DPMS on error:", err)
	}
}

func (m *Manager) setDPMSModeOff() {
	logger.Info("DPMS Off")
	xu := m.helper.xu
	err := dpms.ForceLevelChecked(xu.Conn(), dpms.DPMSModeOff).Check()
	if err != nil {
		logger.Warning("Set DPMS off error:", err)
	}
}

func (m *Manager) doLock() {
	logger.Info("Lock Screen")
	sessionManager := m.helper.SessionManager
	if sessionManager != nil {
		err := sessionManager.RequestLock()
		if err != nil {
			logger.Error("Lock failed:", err)
		}
	}
}

func (m *Manager) doSuspend() {
	logger.Debug("Suspend")
	sessionManager := m.helper.SessionManager
	if sessionManager != nil {
		err := sessionManager.RequestSuspend()
		if err != nil {
			logger.Error("Suspend failed:", err)
		}
	}
}

func (m *Manager) setDisplayBrightness(brightnessTable map[string]float64) {
	display := m.helper.Display
	for output, brightness := range brightnessTable {
		logger.Infof("Change output %q brightness to %.2f", output, brightness)
		err := display.SetBrightness(output, brightness)
		if err != nil {
			logger.Warningf("Change output %q brightness to %.2f failed: %v", output, brightness, err)
		} else {
			logger.Infof("Change output %q brightness to %.2f done!", output, brightness)
		}
	}
}

func doShowLowpower() {
	logger.Info("Show low power")
	go exec.Command(cmdLowPower, "--raise").Run()
}

func execCommand(cmd string) {
	logger.Infof("exec %q", cmd)
	go exec.Command("sh", "-c", cmd).Run()
}

func doCloseLowpower() {
	logger.Info("Close low power")
	go exec.Command(cmdLowPower, "--quit").Run()
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
