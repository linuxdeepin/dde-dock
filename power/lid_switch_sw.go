// +build sw

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
	"fmt"
	"io/ioutil"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
	"time"
)

const (
	swLidOpen  = "1"
	swLidClose = "0"
)

var swLidStateFile = "/sys/bus/platform/devices/liddev/lid_state"

func init() {
	submoduleList = append(submoduleList, newLidStateChangeListener)
}

type lidStateChangeListener struct {
	prevState string
	ticker    *time.Ticker
	manager   *Manager
}

func newLidStateChangeListener(m *Manager) (string, submodule, error) {
	const name = "LidStateChangeListener"
	if !dutils.IsFileExist(swLidStateFile) {
		return name, nil, fmt.Errorf("file not exist %q", swLidStateFile)
	}
	s := &lidStateChangeListener{
		manager: m,
	}
	return name, s, nil
}

func (s *lidStateChangeListener) Start() error {
	s.ticker = time.NewTicker(3 * time.Second)
	s.prevState = getSWLidState()
	go func() {
		for range s.ticker.C {
			logger.Debug("lid state check tick")
			newState := getSWLidState()
			if s.prevState != newState {
				// handle state change
				var closed bool
				if newState == swLidOpen {
					closed = false
				} else if newState == swLidClose {
					closed = true
				}
				s.manager.handleLidSwitch(closed)
				s.prevState = newState
			}
		}
	}()
	return nil
}

// lid_state content: '1\n'
func getSWLidState() string {
	content, err := ioutil.ReadFile(swLidStateFile)
	if err != nil {
		logger.Warning(err)
		return swLidOpen
	}
	return strings.TrimRight(string(content), "\n")
}

func (s *lidStateChangeListener) Destroy() {
	if s.ticker != nil {
		s.ticker.Stop()
		s.ticker = nil
	}
}
