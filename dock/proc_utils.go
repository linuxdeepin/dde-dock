/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/
package dock

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func getProcessCmdline(pid uint) ([]string, error) {
	cmdlinePath := fmt.Sprintf("/proc/%d/cmdline", pid)
	bytes, err := ioutil.ReadFile(cmdlinePath)
	if err != nil {
		return nil, err
	}
	content := string(bytes)
	parts := strings.Split(content, "\x00")
	length := len(parts)
	if length >= 2 && parts[length-1] == "" {
		return parts[:length-1], nil
	}
	return parts, nil
}

func getProcessCwd(pid uint) (string, error) {
	cwdPath := fmt.Sprintf("/proc/%d/cwd", pid)
	cwd, err := os.Readlink(cwdPath)
	return cwd, err
}

func getProcessExe(pid uint) (string, error) {
	exePath := fmt.Sprintf("/proc/%d/exe", pid)
	exe, err := filepath.EvalSymlinks(exePath)
	return exe, err
}

func getProcessEnvVars(pid uint) (map[string]string, error) {
	envPath := fmt.Sprintf("/proc/%d/environ", pid)
	bytes, err := ioutil.ReadFile(envPath)
	if err != nil {
		return nil, err
	}
	content := string(bytes)
	lines := strings.Split(content, "\x00")
	vars := make(map[string]string, len(lines))
	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			vars[parts[0]] = parts[1]
		}
	}
	return vars, nil
}
