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
	"dbus/org/freedesktop/login1"
	"errors"
	"fmt"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"sync"
	"time"
)

const (
	AppsDBusDest            = "com.deepin.daemon.Apps"
	AppsObjectPath          = "/com/deepin/daemon/Apps"
	ALRecorderDBusInterface = AppsDBusDest + ".LaunchedRecorder"
	minUid                  = 1000
)

// app launched recorder
type ALRecorder struct {
	watcher *DFWatcher
	// key is SubRecorder.root
	subRecorders      map[string]*SubRecorder
	subRecordersMutex sync.RWMutex
	loginManager      *login1.Manager

	// signal:
	Launched         func(string)
	StatusSaved      func(root, file string, ok bool)
	ServiceRestarted func()
}

func NewALRecorder(watcher *DFWatcher) *ALRecorder {
	r := &ALRecorder{
		watcher:      watcher,
		subRecorders: make(map[string]*SubRecorder),
	}
	var err error
	r.loginManager, err = login1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		panic(err)
	}

	go r.listenEvents()

	sysDataDirs := getSystemDataDirs()
	for _, dataDir := range sysDataDirs {
		r.watchAppsDir(0, "", filepath.Join(dataDir, "applications"))
	}

	r.loginManager.ConnectUserRemoved(func(uid uint32, _ dbus.ObjectPath) {
		if uid < minUid {
			return
		}
		r.handleUserRemoved(int(uid))
	})

	go r.checkLoop()
	return r
}

func (r *ALRecorder) checkLoop() {
	// There is no need to consider stop the loop
	for {
		r.subRecordersMutex.RLock()
		for _, sr := range r.subRecorders {
			sr.doCheck()
		}
		r.subRecordersMutex.RUnlock()

		time.Sleep(time.Second * 2)
	}
}

func (r *ALRecorder) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       AppsDBusDest,
		ObjectPath: AppsObjectPath,
		Interface:  ALRecorderDBusInterface,
	}
}

func (r *ALRecorder) emitStatusSaved(root, file string, ok bool) {
	dbus.Emit(r, "StatusSaved", root, file, ok)
}

func (r *ALRecorder) emitServiceRestarted() {
	dbus.Emit(r, "ServiceRestarted")
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

func (r *ALRecorder) MarkLaunched(file string) {
	logger.Debugf("ALRecorder.MarkLaunched file: %q", file)
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	for _, sr := range r.subRecorders {
		if strings.HasPrefix(file, sr.root) {
			rel, _ := filepath.Rel(sr.root, file)
			if sr.MarkLaunched(removeDesktopExt(rel)) {
				dbus.Emit(r, "Launched", file)
			}
			return
		}
	}
	logger.Debug("MarkLaunched failed")
}

func (r *ALRecorder) GetNew(dMsg dbus.DMessage) map[string][]string {
	uid := int(dMsg.GetSenderUID())
	ret := make(map[string][]string)
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	for _, sr := range r.subRecorders {
		if intSliceContains(sr.uids, uid) {
			newApps := sr.GetNew()
			if len(newApps) > 0 {
				ret[sr.root] = newApps
			}
		}
	}
	return ret
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

func (r *ALRecorder) WatchDirs(dMsg dbus.DMessage, dataDirs []string) error {
	uid := int(dMsg.GetSenderUID())
	logger.Debugf("WatchDirs uid: %d, data dirs: %#v", uid, dataDirs)
	// check uid
	if uid < minUid {
		return errors.New("invalid uid")
	}

	// check dataDirs
	for _, dataDir := range dataDirs {
		if !filepath.IsAbs(dataDir) {
			return fmt.Errorf("%q is not absolute path", dataDir)
		}
	}

	// get home dir
	home, err := getHomeByUid(uid)
	if err != nil {
		return err
	}

	for _, dataDir := range dataDirs {
		appsDir := filepath.Join(dataDir, "applications")
		r.watchAppsDir(uid, home, appsDir)
	}
	return nil
}

func (r *ALRecorder) handleUserRemoved(uid int) {
	logger.Debug("handleUserRemoved uid:", uid)
	r.subRecordersMutex.Lock()
	defer r.subRecordersMutex.Unlock()

	for _, sr := range r.subRecorders {
		logger.Debug(sr.root, sr.uids)
		sr.uids = intSliceRemove(sr.uids, uid)
		if len(sr.uids) == 0 {
			sr.Destroy()
			delete(r.subRecorders, sr.root)
		}
	}
	logger.Debug("r.subRecorders:", r.subRecorders)
}
