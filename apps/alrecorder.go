/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package apps

import (
	"dbus/org/freedesktop/login1"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"strings"
	"sync"
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
	eventChan         chan *FileEvent
	loginManager      *login1.Manager

	// signal:
	Launched    func(string)
	StatusSaved ALStatusSavedFun
}

type ALStatusSavedFun func(root, file string, ok bool)

func NewALRecorder(watcher *DFWatcher) *ALRecorder {
	eventChan := make(chan *FileEvent)
	r := &ALRecorder{
		watcher:      watcher,
		eventChan:    eventChan,
		subRecorders: make(map[string]*SubRecorder),
	}
	var err error
	r.loginManager, err = login1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		panic(err)
	}

	go r.listenEvents()
	watcher.eventChan = eventChan

	appDirs := getSystemAppDirs()
	for _, appDir := range appDirs {
		r.addSystemAppDir(appDir)
	}

	r.loginManager.ConnectUserNew(func(uid uint32, _ dbus.ObjectPath) {
		if uid < minUid {
			return
		}
		r.addUserAppDir(int(uid))
	})
	r.loginManager.ConnectUserRemoved(func(uid uint32, _ dbus.ObjectPath) {
		if uid < minUid {
			return
		}
		r.removeUserAppDir(int(uid))
	})

	for _, uid := range r.listUsers() {
		r.addUserAppDir(uid)
	}
	return r
}

func (r *ALRecorder) listUsers() (uids []int) {
	users, err := r.loginManager.ListUsers()
	if err != nil {
		logger.Warning(err)
		return
	}

	for _, user := range users {
		// user struct {uid, name, object_path}
		if len(user) > 0 {
			uid := int(user[0].(uint32))

			if uid < minUid {
				continue
			}
			uids = append(uids, uid)
		}
	}
	return
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

func (r *ALRecorder) listenEvents() {
	for {
		ev := <-r.eventChan
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

func (r *ALRecorder) addAppDir(home, path string, uid int) {
	logger.Debugf("ALRecorder.addAppDir %q %q %d", home, path, uid)
	if r.subRecorderExists(path) {
		logger.Debug("ALRecorder.addAppDir faield, subRecorder exists")
		return
	}

	MkdirAll(path, uid, getDirPerm(uid))
	dirNames, appNames := getDirsAndApps(path)
	sr := NewSubRecorder(home, path, uid, appNames, r.emitStatusSaved)
	r.addSubRecorder(sr)
	for _, dirName := range dirNames {
		r.watcher.add(filepath.Join(path, dirName))
	}
}

func (r *ALRecorder) addSystemAppDir(path string) {
	r.addAppDir("", path, 0)
}

func (r *ALRecorder) addUserAppDir(uid int) {
	home, appDir, err := getUserDir(uid)
	if err != nil {
		logger.Warning(err)
		return
	}
	r.addAppDir(home, appDir, uid)
}

func (r *ALRecorder) removeAppDir(path string) {
	r.removeSubRecorder(path)
	r.watcher.removeRecursive(path)
}

func (r *ALRecorder) removeUserAppDir(uid int) {
	_, appDir, err := getUserDir(uid)
	if err != nil {
		logger.Warning(err)
		return
	}
	r.removeAppDir(appDir)
}

func (r *ALRecorder) addSubRecorder(sr *SubRecorder) {
	r.subRecordersMutex.Lock()
	defer r.subRecordersMutex.Unlock()

	r.subRecorders[sr.root] = sr
}

func (r *ALRecorder) subRecorderExists(root string) bool {
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	_, ok := r.subRecorders[root]
	return ok
}

func (r *ALRecorder) removeSubRecorder(root string) {
	r.subRecordersMutex.Lock()
	defer r.subRecordersMutex.Unlock()

	if sr, ok := r.subRecorders[root]; ok {
		delete(r.subRecorders, root)
		sr.Destroy()
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

func (r *ALRecorder) GetNew(dMsg dbus.DMessage) (map[string][]string, error) {
	_, appDir, err := getUserDir(int(dMsg.GetSenderUID()))
	if err != nil {
		return nil, err
	}
	ret := make(map[string][]string)
	r.subRecordersMutex.RLock()
	defer r.subRecordersMutex.RUnlock()

	for _, sr := range r.subRecorders {
		if sr.root == appDir || isSystemAppDir(sr.root) {
			newApps := sr.GetNew()
			if len(newApps) > 0 {
				ret[sr.root] = newApps
			}
		}
	}
	return ret, nil
}
