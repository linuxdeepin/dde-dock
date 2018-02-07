/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package dock

import (
	"errors"
	"fmt"
	"path/filepath"
	"pkg.deepin.io/lib/procfs"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	process procfs.Process
	cmdline []string
	args    []string
	exe     string
	cwd     string
	environ procfs.EnvVars
	hasPid  bool
	ppid    uint
}

func NewProcessInfoWithCmdline(cmd []string) *ProcessInfo {
	if len(cmd) == 0 {
		return nil
	}
	return &ProcessInfo{
		cmdline: cmd,
		args:    cmd[1:],
		exe:     cmd[0],
	}
}

func NewProcessInfo(pid uint) (*ProcessInfo, error) {
	if pid == 0 {
		return nil, errors.New("pid is 0")
	}

	process := procfs.Process(pid)
	pInfo := &ProcessInfo{
		process: process,
		hasPid:  true,
	}
	var err error

	// exe
	pInfo.exe, err = process.Exe()
	if err != nil {
		return nil, err
	}

	// cwd
	pInfo.cwd, err = process.Cwd()
	if err != nil {
		return nil, err
	}

	// cmdline
	pInfo.cmdline, err = process.Cmdline()
	if err != nil {
		return nil, err
	}

	// args
	pInfo.args = getCmdlineArgs(pInfo.exe, pInfo.cwd, pInfo.cmdline)
	if err != nil {
		return nil, err
	}

	// environ
	pInfo.environ, _ = process.Environ()

	// ppid
	if status, err := process.Status(); err == nil {
		pInfo.ppid, _ = status.PPid()
	}

	return pInfo, nil
}

func getCmdlineArgs(exe, cwd string, cmdline []string) []string {
	ok := verifyExe(exe, cwd, cmdline[0])
	if !ok {
		logger.Debug("first arg is not exe file, contains arguments")
		// try again
		parts := strings.Split(cmdline[0], " ")
		ok = verifyExe(exe, cwd, parts[0])
		if !ok {
			logger.Warningf("failed to find right exe, exe: %q, cwd: %q, cmdline: %#v", exe, cwd, cmdline)
			return nil
		} else {
			return append(parts[1:], cmdline[1:]...)
		}
	} else {
		return cmdline[1:]
	}
}

func verifyExe(exe, cwd, firstArg string) bool {
	if filepath.Base(firstArg) == firstArg {
		logger.Debug("basename equal")
		return true
	}

	if !filepath.IsAbs(firstArg) {
		firstArg = filepath.Join(cwd, firstArg)
	}
	// firstArg is abs path
	logger.Debugf("firstArg: %q", firstArg)
	firstArgPath, err := filepath.EvalSymlinks(firstArg)
	if err != nil {
		logger.Warning(err)
		// first arg is not exe file, contains arguments
		return false
	}
	logger.Debugf("firstArgPath: %q", firstArgPath)
	return exe == firstArgPath
}

func (p *ProcessInfo) getJoinedExeArgs() string {
	var cmdline string
	cmdline = strconv.Quote(p.exe)
	for _, arg := range p.args {
		cmdline += (" " + strconv.Quote(arg))
	}
	return cmdline + " $@"
}

func (p *ProcessInfo) GetShellScriptLines() string {
	cmdline := p.getJoinedExeArgs()
	return fmt.Sprintf("#!/bin/sh\ncd %q\nexec %s\n", p.cwd, cmdline)
}

func (p *ProcessInfo) GetOneCommandLine() string {
	cmdline := p.getJoinedExeArgs()
	return fmt.Sprintf("sh -c 'cd %q;exec %s;'", p.cwd, cmdline)
}
