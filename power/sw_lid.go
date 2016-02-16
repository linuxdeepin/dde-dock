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
	"io/ioutil"
	"pkg.deepin.io/dde/daemon/systeminfo"
	dutils "pkg.deepin.io/lib/utils"
	"regexp"
	"strings"
	"time"
)

const (
	swLidStateFile = "/sys/bus/platform/devices/liddev/lid_state"
	swLidOpen      = "1"
	swLidClose     = "0"
)

var swPrevState string

func (p *Power) listenSWLidState() {
	if !dutils.IsFileExist(swLidStateFile) {
		return
	}

	swPrevState = getSWLidState(swLidStateFile)
	for {
		if p.swQuit == nil {
			return
		}

		timer := time.NewTimer(time.Second * 3)
		select {
		case <-p.swQuit:
			return
		case <-timer.C:
			state := getSWLidState(swLidStateFile)
			if !isSWLidStateChanged(state) {
				continue
			}

			if state == swLidOpen {
				p.handleLidSwitch(true)
			} else if state == swLidClose {
				p.handleLidSwitch(false)
			}
		}
	}
}

// lid_state content: '1\n'
func getSWLidState(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return "1"
	}

	return strings.TrimRight(string(content), "\n")
}

func isSWLidStateChanged(state string) bool {
	if state == swPrevState {
		return false
	}
	return true
}

var isSWPlatform = func() func() bool {
	var (
		isSW bool
		cpu  string
	)

	return func() bool {
		if len(cpu) != 0 {
			return isSW
		}

		var err error
		cpu, err = systeminfo.GetCPUInfo("/proc/cpuinfo")
		if err != nil {
			cpu = ""
			return false
		}
		isSW, _ = regexp.MatchString(`^sw`, cpu)
		logger.Debug("Is SW Platform:", isSW)
		return isSW
	}
}()
