/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package screenedge

import (
	"fmt"
	"io/ioutil"
	dutils "pkg.deepin.io/lib/utils"
	"strings"
)

const (
	TopLeft     = "left-up"
	TopRight    = "right-up"
	BottomLeft  = "left-down"
	BottomRight = "right-down"
)

func getProcCmdLine(pid uint32) (string, error) {
	filename := fmt.Sprintf("/proc/%v/cmdline", pid)
	if !dutils.IsFileExist(filename) {
		// TODO return a error
		return "", nil
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Warningf("ReadFile '%s' failed: %v", filename, err)
		return "", err
	}
	content := string(bytes)
	return content, nil
}

func isAppInList(pid uint32, list []string) bool {
	cmdLine, err := getProcCmdLine(pid)
	if len(cmdLine) == 0 {
		return false
	}
	if err != nil {
		return false
	}

	for _, v := range list {
		if strings.Contains(cmdLine, v) {
			logger.Debugf("cmd line %q match %q in list", cmdLine, v)
			return true
		}
	}
	return false
}
