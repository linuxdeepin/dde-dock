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
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	swLidOpen  = "1"
	swLidClose = "0"
)

const swLidStateFile = "/sys/bus/platform/devices/liddev/lid_state"

func (m *Manager) initLidSwitchSW() {
	_, err := os.Stat(swLidStateFile)
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Warning(err)
			return
		}
		// else err is Not Exist Error, ignore it
	} else {
		m.HasLidSwitch = true
		go m.swLidSwitchCheckLoop()
	}
}

func (m *Manager) swLidSwitchCheckLoop() {
	prevState := getLidStateSW()
	for {
		time.Sleep(time.Second * 3)
		newState := getLidStateSW()
		if prevState != newState {
			prevState = newState

			var closed bool
			switch newState {
			case swLidClose:
				closed = true
			case swLidOpen:
				closed = false
			default:
				logger.Warningf("unknown lid state %q", newState)
				continue
			}
			m.handleLidSwitchEvent(closed)
		}
	}
}

// lid_state content: '1\n'
func getLidStateSW() string {
	content, err := ioutil.ReadFile(swLidStateFile)
	if err != nil {
		logger.Warning(err)
		return swLidOpen
	}
	return strings.TrimRight(string(content), "\n")
}
