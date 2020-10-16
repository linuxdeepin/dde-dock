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

package power

import (
	"os/exec"
	"strings"
)

func init() {
	submoduleList = append(submoduleList, newLidSwitchHandler)
}

// nolint
const (
	lidSwitchStateUnknown = iota
	lidSwitchStateOpen
	lidSwitchStateClose
)

type LidSwitchHandler struct {
	manager *Manager
	cmd     *exec.Cmd
}

func newLidSwitchHandler(m *Manager) (string, submodule, error) {
	h := &LidSwitchHandler{
		manager: m,
	}
	return "LidSwitchHandler", h, nil
}

func (h *LidSwitchHandler) Start() error {
	power := h.manager.helper.Power
	_, err := power.ConnectLidClosed(h.onLidClosed)
	if err != nil {
		return err
	}
	_, err = power.ConnectLidOpened(h.onLidOpened)
	if err != nil {
		return err
	}
	return nil
}

func (h *LidSwitchHandler) onLidClosed() {
	logger.Info("Lid closed")
	var onBattery bool
	m := h.manager
	m.setPrepareSuspend(suspendStateLidClose)
	m.PropsMu.Lock()
	m.lidSwitchState = lidSwitchStateClose
	onBattery = h.manager.OnBattery
	m.PropsMu.Unlock()
	m.claimOrReleaseAmbientLight()
	var lidCloseAction int32
	if onBattery {
		lidCloseAction = m.BatteryLidClosedAction.Get() // 获取合盖操作
	} else {
		lidCloseAction = m.LinePowerLidClosedAction.Get() // 获取合盖操作
	}
	switch lidCloseAction {
	case powerActionShutdown:
		m.doShutdown()
	case powerActionSuspend:
		m.doSuspend()
	case powerActionHibernate:
		m.doHibernate()
	case powerActionTurnOffScreen:
		m.doTurnOffScreen()
	case powerActionDoNothing:
		return
	}
}

func (h *LidSwitchHandler) onLidOpened() {
	logger.Info("Lid opened")
	m := h.manager
	m.setPrepareSuspend(suspendStateLidOpen)
	m.PropsMu.Lock()
	m.lidSwitchState = lidSwitchStateOpen
	m.PropsMu.Unlock()
	m.claimOrReleaseAmbientLight()

	if err := h.stopAskUser(); err != nil {
		logger.Warning("stopAskUser error:", err)
	}

	err := m.helper.ScreenSaver.SimulateUserActivity(0)
	if err != nil {
		logger.Warning(err)
	}
}

func (h *LidSwitchHandler) stopAskUser() error {
	if h.cmd == nil {
		return nil
	}

	if h.cmd.ProcessState == nil {
		// h.cmd is running
		logger.Debug("stopAskUser: kill process")
		err := h.cmd.Process.Kill()
		if err != nil {
			return err
		}
	} else {
		logger.Debug("stopAskUser: cmd exited")
	}
	h.cmd = nil
	return nil
}

// copy from display module of project startdde
func isBuiltinOutput(name string) bool {
	name = strings.ToLower(name)
	switch {
	case strings.Contains(name, "lvds"):
		// Most drivers use an "LVDS" prefix
		fallthrough
	case strings.Contains(name, "lcd"):
		// fglrx uses "LCD" in some versions
		fallthrough
	case strings.Contains(name, "edp"):
		// eDP is for internal built-in panel connections
		fallthrough
	case strings.Contains(name, "dsi"):
		return true
	case name == "default":
		// now sunway notebook has only one output named default
		return true
	}
	return false
}

func (h *LidSwitchHandler) Destroy() {
}
