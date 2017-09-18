/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	"path/filepath"
	"pkg.deepin.io/lib/procfs"
	"strconv"
	"strings"
)

type IdentifyWindowFunc struct {
	Name string
	Fn   _IdentifyWindowFunc
}

type _IdentifyWindowFunc func(*DockManager, *WindowInfo) (string, *AppInfo)

func (m *DockManager) registerIdentifyWindowFuncs() {
	m.registerIdentifyWindowFunc("PidEnv", identifyWindowByPidEnv)
	m.registerIdentifyWindowFunc("Cmdline-XWalk", identifyWindowByCmdlineXWalk)
	m.registerIdentifyWindowFunc("Rule", identifyWindowByRule)
	m.registerIdentifyWindowFunc("Bamf", identifyWindowByBamf)
	m.registerIdentifyWindowFunc("Pid", identifyWindowByPid)
	m.registerIdentifyWindowFunc("Scratch", identifyWindowByScratch)
	m.registerIdentifyWindowFunc("GtkAppId", identifyWindowByGtkAppId)
	m.registerIdentifyWindowFunc("WmClass", identifyWindowByWmClass)
}

func (m *DockManager) registerIdentifyWindowFunc(name string, fn _IdentifyWindowFunc) {
	m.identifyWindowFuns = append(m.identifyWindowFuns, &IdentifyWindowFunc{
		Name: name,
		Fn:   fn,
	})
}

func (m *DockManager) identifyWindow(winInfo *WindowInfo) (string, *AppInfo) {
	logger.Debugf("identifyWindow: window id: %v, window innerId: %q", winInfo.window, winInfo.innerId)
	if winInfo.innerId == "" {
		logger.Debug("identifyWindow: failed winInfo no innerId")
		return "", nil
	}

	for idx, item := range m.identifyWindowFuns {
		name := item.Name
		logger.Debugf("identifyWindow: try %s:%d", name, idx)
		innerId, appInfo := item.Fn(m, winInfo)
		if innerId != "" {
			// success
			logger.Debugf("identifyWindow by %s success, innerId: %q, appInfo: %v", name, innerId, appInfo)
			// NOTE: if name == "Pid", appInfo may be nil
			if appInfo != nil {
				fixedAppInfo := fixAutostartAppInfo(appInfo)
				if fixedAppInfo != nil {
					appInfo = fixedAppInfo
					appInfo.identifyMethod = name + "+FixAutostart"
					innerId = fixedAppInfo.innerId
				} else {
					appInfo.identifyMethod = name
				}
			}
			return innerId, appInfo
		}
	}
	// fail
	logger.Debugf("identifyWindow: failed")
	return winInfo.innerId, nil
}

func fixAutostartAppInfo(appInfo *AppInfo) *AppInfo {
	file := appInfo.GetFileName()
	if isInAutostartDir(file) {
		logger.Debug("file is in autostart dir")
		base := filepath.Base(file)
		return NewAppInfo(base)
	}
	return nil
}

func identifyWindowByScratch(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	desktopFile := filepath.Join(scratchDir, addDesktopExt(winInfo.innerId))
	logger.Debugf("try scratch desktop file: %q", desktopFile)
	appInfo := NewAppInfoFromFile(desktopFile)
	if appInfo != nil {
		// success
		return appInfo.innerId, appInfo
	}
	// fail
	return "", nil
}

func identifyWindowByPid(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	if winInfo.pid != 0 {
		logger.Debugf("identifyWindowByPid: pid: %d", winInfo.pid)
		entry := m.Entries.GetByWindowPid(winInfo.pid)
		if entry != nil {
			// success
			return entry.innerId, entry.appInfo
		}
	}
	// fail
	return "", nil
}

func identifyWindowByGtkAppId(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	gtkAppId := winInfo.gtkAppId
	logger.Debugf("identifyWindowByGtkAppId: gtkAppId: %q", gtkAppId)
	if gtkAppId != "" {
		appInfo := NewAppInfo(gtkAppId)
		if appInfo != nil {
			// success
			return appInfo.innerId, appInfo
		}
	}
	// fail
	return "", nil
}

func identifyWindowByPidEnv(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	pid := winInfo.pid
	process := winInfo.process
	if process != nil && pid != 0 {
		launchedDesktopFile := process.environ.Get("GIO_LAUNCHED_DESKTOP_FILE")
		launchedDesktopFilePid, _ := strconv.ParseUint(
			process.environ.Get("GIO_LAUNCHED_DESKTOP_FILE_PID"), 10, 32)

		logger.Debugf("identifyWindowByPidEnv: launchedDesktopFile: %q, pid: %d",
			launchedDesktopFile, launchedDesktopFilePid)

		var try bool
		if uint(launchedDesktopFilePid) == pid {
			try = true
		} else if uint(launchedDesktopFilePid) == process.ppid && process.ppid != 0 {
			logger.Debug("ppid equal")
			parentProcess := procfs.Process(process.ppid)
			cmdline, err := parentProcess.Cmdline()
			if err == nil && len(cmdline) > 0 {
				logger.Debugf("parent process cmdline: %#v", cmdline)
				base := filepath.Base(cmdline[0])
				if base == "sh" || base == "bash" {
					try = true
				}
			}
		}

		if try {
			appInfo := NewAppInfoFromFile(launchedDesktopFile)
			if appInfo != nil {
				// success
				return appInfo.innerId, appInfo
			}
		}
	}
	// fail
	return "", nil
}

func identifyWindowByRule(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	ret := m.windowPatterns.Match(winInfo)
	if ret == "" {
		return "", nil
	}
	logger.Debug("identifyWindowByRule ret:", ret)
	// parse ret
	// id=$appId or env
	var appInfo *AppInfo
	if len(ret) > 4 && strings.HasPrefix(ret, "id=") {
		appInfo = NewAppInfo(ret[3:])
	} else if ret == "env" {
		process := winInfo.process
		if process != nil {
			launchedDesktopFile := process.environ.Get("GIO_LAUNCHED_DESKTOP_FILE")
			if launchedDesktopFile != "" {
				appInfo = NewAppInfoFromFile(launchedDesktopFile)
			}
		}
	} else {
		logger.Warningf("bad ret: %q", ret)
	}

	if appInfo != nil {
		return appInfo.innerId, appInfo
	}
	return "", nil
}

func identifyWindowByWmClass(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	if winInfo.wmClass != nil {
		instance := winInfo.wmClass.Instance
		if instance != "" {
			appInfo := NewAppInfo(instance)
			if appInfo != nil {
				return appInfo.innerId, appInfo
			}
		}

		class := winInfo.wmClass.Class
		if class != "" {
			appInfo := NewAppInfo(class)
			if appInfo != nil {
				return appInfo.innerId, appInfo
			}
		}
	}
	// fail
	return "", nil
}

func identifyWindowByBamf(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	// bamf
	win := winInfo.window
	desktop := getDesktopFromWindowByBamf(win)
	if desktop != "" {
		appInfo := NewAppInfoFromFile(desktop)
		if appInfo != nil {
			// success
			return appInfo.innerId, appInfo
		}
	}
	return "", nil
}

func identifyWindowByCmdlineXWalk(m *DockManager, winInfo *WindowInfo) (string, *AppInfo) {
	process := winInfo.process
	if process == nil || winInfo.pid == 0 {
		return "", nil
	}

	exeBase := filepath.Base(process.exe)
	args := process.args
	if exeBase != "xwalk" || len(args) == 0 {
		return "", nil
	}
	lastArg := args[len(args)-1]
	logger.Debugf("lastArg: %q", lastArg)

	if filepath.Base(lastArg) == "manifest.json" {
		appId := filepath.Base(filepath.Dir(lastArg))
		appInfo := NewAppInfo(appId)
		if appInfo != nil {
			// success
			return appInfo.innerId, appInfo
		}
	}
	// failed
	return "", nil
}
