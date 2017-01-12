/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mime

import (
	"github.com/fsnotify/fsnotify"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	dutils "pkg.deepin.io/lib/utils"
	"pkg.deepin.io/lib/xdg/basedir"
	"time"
)

const (
	AppMimeTerminal = "application/x-terminal"

	dbusDest = "com.deepin.daemon.Mime"
	dbusPath = "/com/deepin/daemon/Mime"
	dbusIFC  = dbusDest
)

type Manager struct {
	Change      func()
	media       *Media
	userManager *userAppManager
	fsWatcher   *fsnotify.Watcher
	changeTimer *time.Timer

	done     chan struct{}
	doneResp chan struct{}
}

func NewManager() *Manager {
	m := new(Manager)
	m.done = make(chan struct{})
	m.doneResp = make(chan struct{})
	userManager, err := newUserAppManager(userAppFile)
	if err != nil {
		userManager = &userAppManager{
			filename: userAppFile,
		}
	}
	m.userManager = userManager

	m.fsWatcher, err = fsnotify.NewWatcher()
	if err == nil {
		go m.handleFileEvents()
		dirs := getDirsNeedWatched()
		for _, dir := range dirs {
			logger.Debugf("watch dir %q", dir)
			if err := m.fsWatcher.Add(dir); err != nil {
				logger.Warning(err)
			}
		}
	} else {
		logger.Warning("new fs watcher failed:", err)
	}

	return m
}

func getDirsNeedWatched() []string {
	dirs := make([]string, 0, 5)
	appsDirs := getApplicatonsDirs()
	dirs = append(dirs, appsDirs...)

	dirs = append(dirs, basedir.GetUserConfigDir())
	sysConfDirs := basedir.GetSystemConfigDirs()
	dirs = append(dirs, sysConfDirs...)
	return dirs
}

func getApplicatonsDirs() []string {
	dirs := make([]string, 0, 3)
	dirs = append(dirs, basedir.GetUserDataDir())
	sysDirs := basedir.GetSystemDataDirs()
	dirs = append(dirs, sysDirs...)

	ret := make([]string, len(dirs))
	for i, dir := range dirs {
		ret[i] = filepath.Join(dir, "applications")
	}
	return ret
}

func (m *Manager) handleFileEvents() {
	watcher := m.fsWatcher
	defer close(m.doneResp)
	defer watcher.Close()
	for {
		select {
		case event := <-watcher.Events:
			logger.Debug("event:", event)
			base := filepath.Base(event.Name)
			if base == "mimeinfo.cache" || base == "mimeapps.list" {
				m.deferEmitChange()
			}

		case err := <-watcher.Errors:
			logger.Warning("error:", err)

		case <-m.done:
			return
		}
	}
}

func (m *Manager) deferEmitChange() {
	delay := 2 * time.Second
	if m.changeTimer == nil {
		m.changeTimer = time.AfterFunc(delay, func() {
			dbus.Emit(m, "Change")
		})
	} else {
		m.changeTimer.Reset(delay)
	}
}

func (m *Manager) destroy() {
	dbus.UnInstallObject(m)
	dbus.UnInstallObject(m.media)

	m.changeTimer.Stop()
	m.changeTimer = nil

	// send close signal to handleFileEvents goroutine
	close(m.done)
	// Wait for handleFileEvents goroutine to close
	<-m.doneResp
}

func (m *Manager) initConfigData() {
	if dutils.IsFileExist(filepath.Join(basedir.GetUserConfigDir(),
		"mimeapps.list")) {
		return
	}

	err := m.doInitConfigData()
	if err != nil {
		logger.Warning("Init mime config file failed", err)
	}
}

func (m *Manager) doInitConfigData() error {
	return genMimeAppsFile(
		findFilePath(filepath.Join("dde-daemon", "mime", "data.json")))
}

// Reset reset mimes default app
func (m *Manager) Reset() {
	resetTerminal()

	go func() {
		err := m.doInitConfigData()
		if err != nil {
			logger.Warning("Init mime config file failed", err)
		}
		dbus.Emit(m, "Change")
	}()

}

// GetDefaultApp get the default app id for the special mime
// ty: the special mime
// ret0: the default app info
// ret1: error message
func (m *Manager) GetDefaultApp(ty string) (string, error) {
	var (
		info *AppInfo
		err  error
	)
	if ty == AppMimeTerminal {
		info, err = getDefaultTerminal()
	} else {
		info, err = GetAppInfo(ty)
	}
	if err != nil {
		return "", err
	}

	return marshal(info)
}

// SetDefaultApp set the default app for the special mime list
// ty: the special mime
// deskId: the default app desktop id
// ret0: error message
func (m *Manager) SetDefaultApp(mimes []string, deskId string) error {
	var err error
	for _, mime := range mimes {
		if mime == AppMimeTerminal {
			err = setDefaultTerminal(deskId)
		} else {
			err = SetAppInfo(mime, deskId)
		}
		if err != nil {
			logger.Warningf("Set '%s' default app to '%s' failed: %v",
				mime, deskId, err)
			break
		}
	}
	return err
}

// ListApps list the apps that supported the special mime
// ty: the special mime
// ret0: the app infos
func (m *Manager) ListApps(ty string) string {
	var infos AppInfos
	if ty == AppMimeTerminal {
		infos = getTerminalInfos()
	} else {
		infos = GetAppInfos(ty)
	}
	content, _ := marshal(infos)
	return content
}

func (m *Manager) ListUserApps(ty string) string {
	apps := m.userManager.Get(ty)
	if len(apps) == 0 {
		return ""
	}
	var infos AppInfos
	for _, app := range apps {
		info, err := newAppInfoById(app.DesktopId)
		if err != nil {
			logger.Warningf("New '%s' failed: %v", app.DesktopId, err)
			continue
		}
		infos = append(infos, info)
	}
	content, _ := marshal(infos)
	return content
}

func (m *Manager) AddUserApp(mimes []string, desktopId string) error {
	// check app validity
	_, err := newAppInfoById(desktopId)
	if err != nil {
		logger.Error("Invalid desktop id:", desktopId)
		return err
	}
	if !m.userManager.Add(mimes, desktopId) {
		return nil
	}
	return m.userManager.Write()
}

func (m *Manager) DeleteUserApp(desktopId string) error {
	err := m.userManager.Delete(desktopId)
	if err != nil {
		logger.Errorf("Delete '%s' failed: %v", desktopId, err)
		return err
	}
	return m.userManager.Write()
}

func (m *Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}
