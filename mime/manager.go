/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package mime

import (
	"path/filepath"
	"strings"
	"time"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/fsnotify"
	dutils "pkg.deepin.io/lib/utils"
	"pkg.deepin.io/lib/xdg/basedir"
)

const (
	AppMimeTerminal = "application/x-terminal"

	dbusServiceName = "com.deepin.daemon.Mime"
	dbusPath        = "/com/deepin/daemon/Mime"
	dbusInterface   = dbusServiceName
)

type Manager struct {
	service     *dbusutil.Service
	userManager *userAppManager
	fsWatcher   *fsnotify.Watcher
	changeTimer *time.Timer

	done     chan struct{}
	doneResp chan struct{}

	methods *struct {
		GetDefaultApp func() `in:"mime_type" out:"default_app"`
		SetDefaultApp func() `in:"mime_types,desktop_id"`
		ListApps      func() `in:"mime_type" out:"apps"`
		ListUserApps  func() `in:"mime_type" out:"user_apps"`
		AddUserApp    func() `in:"mime_types,desktop_id"`
		DeleteUserApp func() `in:"desktop_id"`
	}

	signals *struct {
		Change struct{}
	}
}

func NewManager(service *dbusutil.Service) *Manager {
	m := new(Manager)
	m.service = service
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
			if err := m.fsWatcher.Watch(dir); err != nil {
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
		case event, ok := <-watcher.Event:
			if !ok {
				logger.Error("Invalid watcher event:", event)
				return
			}

			logger.Debug("event:", event)
			base := filepath.Base(event.Name)
			if base == "mimeinfo.cache" || base == "mimeapps.list" {
				m.deferEmitChange()
			}

		case err := <-watcher.Error:
			logger.Warning("error:", err)
			return
		case <-m.done:
			return
		}
	}
}

func (m *Manager) deferEmitChange() {
	delay := 2 * time.Second
	if m.changeTimer == nil {
		m.changeTimer = time.AfterFunc(delay, func() {
			m.emitSignalChange()
		})
	} else {
		m.changeTimer.Reset(delay)
	}
}

func (m *Manager) emitSignalChange() {
	m.service.Emit(m, "Change")
}

func (m *Manager) destroy() {
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

	go func() {
		err := m.doInitConfigData()
		if err != nil {
			logger.Warning("Init mime config file failed", err)
		} else {
			logger.Info("Init mime config file successfully")
		}
	}()
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
		m.emitSignalChange()
	}()

}

// GetDefaultApp get the default app id for the special mime
// ty: the special mime
// ret0: the default app info
// ret1: error message
func (m *Manager) GetDefaultApp(ty string) (string, *dbus.Error) {
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
		return "", dbusutil.ToError(err)
	}

	defaultApp, err := marshal(info)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	return defaultApp, nil
}

// SetDefaultApp set the default app for the special mime list
// ty: the special mime
// deskId: the default app desktop id
// ret0: error message
func (m *Manager) SetDefaultApp(mimes []string, desktopDd string) *dbus.Error {
	var err error
	for _, mime := range mimes {
		if mime == AppMimeTerminal {
			err = setDefaultTerminal(desktopDd)
		} else {
			err = SetAppInfo(mime, desktopDd)
		}
		if err != nil {
			logger.Warningf("Set '%s' default app to '%s' failed: %v",
				mime, desktopDd, err)
			break
		}
	}
	return dbusutil.ToError(err)
}

// ListApps list the apps that supported the special mime
// ty: the special mime
// ret0: the app infos
func (m *Manager) ListApps(ty string) (string, *dbus.Error) {
	var infos AppInfos
	if ty == AppMimeTerminal {
		infos = getTerminalInfos()
	} else {
		infos = GetAppInfos(ty)
	}

	// filter out deepin custom desktop file
	filteredInfos := make(AppInfos, 0, len(infos))
	for _, info := range infos {
		if isDeepinCustomDesktopFile(info.fileName) {
			continue
		}

		filteredInfos = append(filteredInfos, info)
	}

	content, err := marshal(filteredInfos)
	if err != nil {
		return "", dbusutil.ToError(err)
	}

	return content, nil
}

var userAppDir string

func init() {
	userDataDir := basedir.GetUserDataDir()
	// userAppDir is $HOME/.local/share/applications
	userAppDir = filepath.Join(userDataDir, "applications")
}

// The default applications module of the DDE Control Center
// creates the desktop file with the file name beginning with
// "deepin-custom" in the applications directory under the XDG
// user data directory.
func isDeepinCustomDesktopFile(file string) bool {
	dir := filepath.Dir(file)
	base := filepath.Base(file)
	return dir == userAppDir && strings.HasPrefix(base, "deepin-custom-")
}

func (m *Manager) ListUserApps(ty string) (string, *dbus.Error) {
	apps := m.userManager.Get(ty)
	if len(apps) == 0 {
		return "", nil
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
	content, err := marshal(infos)
	if err != nil {
		return "", dbusutil.ToError(err)
	}
	return content, nil
}

func (m *Manager) AddUserApp(mimes []string, desktopId string) *dbus.Error {
	logger.Debugf("Manager.AddUserApp mimes %v desktop id: %q", mimes, desktopId)
	// check app validity
	_, err := newAppInfoById(desktopId)
	if err != nil {
		logger.Warningf("Invalid desktop id %q", desktopId)
		return dbusutil.ToError(err)
	}
	if !m.userManager.Add(mimes, desktopId) {
		return nil
	}
	err = m.userManager.Write()
	return dbusutil.ToError(err)
}

func (m *Manager) DeleteUserApp(desktopId string) *dbus.Error {
	logger.Debugf("Manager.DeleteUserApp %q", desktopId)
	err := m.userManager.Delete(desktopId)
	if err != nil {
		logger.Warningf("Delete %q failed: %v", desktopId, err)
		return dbusutil.ToError(err)
	}
	err = m.userManager.Write()
	return dbusutil.ToError(err)
}

func (*Manager) GetInterfaceName() string {
	return dbusInterface
}
