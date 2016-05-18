package dock

import (
	"fmt"
	"path/filepath"
	"strings"
)

type ProcessInfo struct {
	pid     uint
	cmdline []string
	args    []string
	exe     string
	cwd     string
}

func NewProcessInfo(pid uint) (*ProcessInfo, error) {
	pInfo := &ProcessInfo{
		pid: pid,
	}
	var err error

	// exe
	pInfo.exe, err = getProcessExe(pid)
	if err != nil {
		return nil, err
	}

	// cwd
	pInfo.cwd, err = getProcessCwd(pid)
	if err != nil {
		return nil, err
	}

	// cmdline
	pInfo.cmdline, err = getProcessCmdline(pid)
	if err != nil {
		return nil, err
	}

	// args
	pInfo.args = getCmdlineArgs(pInfo.exe, pInfo.cwd, pInfo.cmdline)
	if err != nil {
		return nil, err
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
		logger.Error(err)
		// first arg is not exe file, contains arguments
		return false
	}
	logger.Debugf("firstArgPath: %q", firstArgPath)
	return exe == firstArgPath
}

func (p *ProcessInfo) GetShellScript() string {
	cmdlineJoined := strings.Join(p.cmdline, " ")
	return fmt.Sprintf("cd %q\n%s\n", p.cwd, cmdlineJoined)
}
