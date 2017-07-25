/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

import (
	"os/exec"
	"pkg.deepin.io/dde/api/drandr"
	"strings"
)

func init() {
	submoduleList = append(submoduleList, newLidSwitchHandler)
}

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
	power.ConnectLidClosed(h.onLidClosed)
	power.ConnectLidOpened(h.onLidOpened)
	return nil
}

func (h *LidSwitchHandler) onLidClosed() {
	logger.Info("Lid closed")
	m := h.manager
	if !m.LidClosedSleep.Get() {
		return
	}
	outputs, err := getWorkingOutputNames(m.helper)
	if err != nil {
		logger.Warning("getWorkingOutputNames failed:", err)
		return
	}
	logger.Debug("working outputs:", outputs)
	outputs = outputsAfterLidClosed(outputs)
	logger.Debug("working outputs after lid closed:", outputs)
	if len(outputs) > 0 {
		if err := h.startAskUser(); err != nil {
			logger.Warning("LidSwitchHandler.startAskUser failed", err)
		}
	} else {
		m.doSuspend()
	}
}

func (h *LidSwitchHandler) onLidOpened() {
	logger.Info("Lid opened")
	if err := h.stopAskUser(); err != nil {
		logger.Warning("stopAskUser error:", err)
	}
}

func (h *LidSwitchHandler) startAskUser() error {
	h.cmd = exec.Command("/usr/lib/deepin-daemon/dde-suspend-dialog")
	logger.Debug("startAskUser cmd.Start")
	if err := h.cmd.Start(); err != nil {
		return err
	}
	go func() {
		err := h.cmd.Wait()
		if err == nil {
			h.manager.doSuspend()
		} else {
			logger.Debug("cmd exit with", err)
		}
	}()
	return nil
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

// guess the working outputs after user close the lid
func outputsAfterLidClosed(outputs []string) []string {
	ret := make([]string, 0, len(outputs))
	// found built ouput
	var found bool
	for _, output := range outputs {
		if !found && isBuiltinOuput(output) {
			// skip built output
			continue
			found = true
		}
		ret = append(ret, output)
	}
	return ret
}

// copy from display module of project startdde
func isBuiltinOuput(name string) bool {
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

// get the names of the working display outputs
func getWorkingOutputNames(helper *Helper) ([]string, error) {
	conn := helper.xu.Conn()
	screenInfo, err := drandr.GetScreenInfo(conn)
	if err != nil {
		return nil, err
	}

	outputs := screenInfo.Outputs.ListConnectionOutputs()
	workingOutputs := make(drandr.OutputInfos, 0, len(outputs))
	for _, output := range outputs {
		crtc := &output.Crtc
		if crtc.Width == 0 || crtc.Height == 0 {
			continue
		}
		workingOutputs = append(workingOutputs, output)
	}
	return workingOutputs.ListNames(), nil
}

func (h *LidSwitchHandler) Destroy() {
}
