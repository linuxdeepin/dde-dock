/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package apps

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"sync"

	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.login1"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	dbusServiceName         = "com.deepin.daemon.Apps"
	dbusPath                = "/com/deepin/daemon/Apps"
	alRecorderDBusInterface = dbusServiceName + ".LaunchedRecorder"
	minUid                  = 1000
)

// app launched recorder
type ALRecorder struct {
	watcher *DFWatcher
	// key is SubRecorder.root
	subRecorders      map[string]*SubRecorder
	subRecordersMutex sync.RWMutex
	loginManager      *login1.Manager

	methods *struct {
		GetNew       func() `out:"newApps"`
		MarkLaunched func() `in:"file"`
		WatchDirs    func() `in:"dirs"`
	}

	signals *struct {
		Launched struct {
			file string
		}

		StatusSaved struct {
			root string
			file string
			ok   bool
		}

		ServiceRestarted struct{}
	}
}

func newALRecorder(watcher *DFWatcher) (*ALRecorder, error) {
	r := &ALRecorder{
		watcher:      watcher,
		subRecorders: make(map[string]*SubRecorder),
	}
	systemBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	r.loginManager = login1.NewManager(systemBus)

	go r.listenEvents()

	sysDataDirs := getSystemDataDirs()
	for _, dataDir := range sysDataDirs {
		r.watchAppsDir(0, "", filepath.Join(dataDir, "applications"))
	}

	sysSigLoop := dbusutil.NewSignalLoop(systemBus, 10)
	sysSigLoop.Start()
	r.loginManager.InitSignalExt(sysSigLoop, true)
	r.loginManager.ConnectUserRemoved(func(uid uint32, userPath dbus.ObjectPath) {
		r.handleUserRemoved(int(uid))
	})

	return r, nil
}

func (*ALRecorder) GetInterfaceName() string {
	return alRecorderDBusInterface
}

func (r *ALRecorder) Service() *dbusutil.Service {
	return r.watcher.service
}

func (r *ALRecorder) emitLaunched(file string) {
	err := r.Service().Emit(r, "Launched", file)
	if err != nil {
		logger.Warning(err)
	}
}

func (r *ALRecorder) emitStatusSaved(root, file string, ok bool) {
	err := r.Service().Emit(r, "StatusSaved", root, file, ok)
	if err != nil {
		logger.Warning(err)
	}
}

func (r *ALRecorder) emitServiceRestarted() {
	err := r.Service().Emit(r, "ServiceRestarted")
	if err != nil {
		logger.Warning(err)
	}
}

func (r *ALRecorder) listenEvents() {
	eventChan := r.watcher.eventChan
	for {
		ev := <-eventChan
		logger.Debugf("ALRecorder ev: %#v", ev)
		name := ev.Name

		if isDesktopFile(name) {
			if ev.IsFound || ev.IsCreate() || ev.IsModify() {
				// added
				r.handleAdded(name)
			} else if ev.NotExist {
				// removed
				r.handleRemoved(name)
			}
		} else if ev.NotExist && ev.IsRename() {
			// may be dir removed
			r.handleDirRemoved(name)
		}
	}
}

// handleAdded handle desktop file added
func (r *ALRecorder) handleAdded(file string) {
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	for _, sr := range r.subRecorders {
		if strings.HasPrefix(file, sr.root) {
			rel, _ := filepath.Rel(sr.root, file)
			sr.handleAdded(removeDesktopExt(rel))
			return
		}
	}
}

// handleRemoved: handle desktop file removed
func (r *ALRecorder) handleRemoved(file string) {
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	for _, sr := range r.subRecorders {
		if strings.HasPrefix(file, sr.root) {
			rel, _ := filepath.Rel(sr.root, file)
			sr.handleRemoved(removeDesktopExt(rel))
			return
		}
	}
}

func (r *ALRecorder) handleDirRemoved(file string) {
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	for _, sr := range r.subRecorders {
		if strings.HasPrefix(file, sr.root) {
			rel, _ := filepath.Rel(sr.root, file)
			sr.handleDirRemoved(rel)
			return
		}
	}
}

func (r *ALRecorder) MarkLaunched(file string) *dbus.Error {
	logger.Debugf("ALRecorder.MarkLaunched file: %q", file)
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	for _, sr := range r.subRecorders {
		if strings.HasPrefix(file, sr.root) {
			rel, _ := filepath.Rel(sr.root, file)
			if sr.MarkLaunched(removeDesktopExt(rel)) {
				r.emitLaunched(file)
			}
			return nil
		}
	}
	logger.Debug("MarkLaunched failed")
	return nil
}

func (r *ALRecorder) GetNew(sender dbus.Sender) (map[string][]string, *dbus.Error) {
	uid, err := r.Service().GetConnUID(string(sender))
	if err != nil {
		return nil, dbusutil.ToError(err)
	}

	ret := make(map[string][]string)
	r.subRecordersMutex.RLock()

	for _, sr := range r.subRecorders {
		if intSliceContains(sr.uids, int(uid)) {
			newApps := sr.GetNew()
			if len(newApps) > 0 {
				ret[sr.root] = newApps
			}
		}
	}
	r.subRecordersMutex.RUnlock()

	return ret, nil
}

func (r *ALRecorder) watchAppsDir(uid int, home, appsDir string) {
	r.subRecordersMutex.Lock()
	defer r.subRecordersMutex.Unlock()

	sr := r.subRecorders[appsDir]
	if sr != nil {
		// subRecorder exists
		logger.Debugf("subRecorder for %q exists", appsDir)
		if !intSliceContains(sr.uids, uid) {
			sr.uids = append(sr.uids, uid)
			logger.Debug("append uid", uid)
		}
		return
	}

	sr = NewSubRecorder(uid, home, appsDir, r)
	r.subRecorders[appsDir] = sr
}

func (r *ALRecorder) WatchDirs(sender dbus.Sender, dataDirs []string) *dbus.Error {
	uid, err := r.Service().GetConnUID(string(sender))
	if err != nil {
		return dbusutil.ToError(err)
	}

	logger.Debugf("WatchDirs uid: %d, data dirs: %#v", uid, dataDirs)
	// check uid
	if uid < minUid {
		return dbusutil.ToError(errors.New("invalid uid"))
	}

	// check dataDirs
	for _, dataDir := range dataDirs {
		if !filepath.IsAbs(dataDir) {
			return dbusutil.ToError(fmt.Errorf("%q is not absolute path", dataDir))
		}
	}

	// get home dir
	home, err := getHomeByUid(int(uid))
	if err != nil {
		return dbusutil.ToError(err)
	}

	for _, dataDir := range dataDirs {
		appsDir := filepath.Join(dataDir, "applications")
		r.watchAppsDir(int(uid), home, appsDir)
	}
	return nil
}

func (r *ALRecorder) handleUserRemoved(uid int) {
	logger.Debug("handleUserRemoved uid:", uid)
	r.subRecordersMutex.Lock()

	for _, sr := range r.subRecorders {
		logger.Debug(sr.root, sr.uids)
		sr.uids = intSliceRemove(sr.uids, uid)
		if len(sr.uids) == 0 {
			sr.Destroy()
			delete(r.subRecorders, sr.root)
		}
	}

	r.subRecordersMutex.Unlock()
	logger.Debug("r.subRecorders:", r.subRecorders)
}
