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
	"bufio"
	"crypto/md5"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	dirPerm = 0755
)

type SubRecorder struct {
	root string
	// key: app, value: launched
	launchedMap        map[string]bool
	removedLaunchedMap map[string]bool
	uninstallMap       map[string]struct{}

	launchedMapMu sync.RWMutex

	statusFile      string
	statusFileOwner int

	uids     []int
	parent   *ALRecorder
	saveMu   sync.Mutex
	isSaving bool
}

func NewSubRecorder(uid int, home, root string, parent *ALRecorder) *SubRecorder {
	sr := &SubRecorder{
		root:               root,
		uids:               []int{uid}, // first uid
		parent:             parent,
		removedLaunchedMap: make(map[string]bool),
		uninstallMap:       make(map[string]struct{}),
	}

	sr.statusFile, sr.statusFileOwner = getStatusFileAndOwner(uid, home, root)
	logger.Debugf("NewSubRecorder status file: %q, owner: %d", sr.statusFile, sr.statusFileOwner)

	parent.watcher.fsWatcher.addRoot(root)
	sr.init()
	return sr
}

func (sr *SubRecorder) init() {
	subDirNames, apps := getDirsAndApps(sr.root)
	if sr.initAppLaunchedMap(apps) {
		err := MkdirAll(filepath.Dir(sr.statusFile), sr.statusFileOwner, dirPerm)
		if err != nil {
			logger.Warning(err)
		}
		sr.RequestSave()
	}

	for _, dirName := range subDirNames {
		path := filepath.Join(sr.root, dirName)
		err := sr.parent.watcher.add(path)
		if err != nil {
			logger.Warning(err)
		}
	}
}

func (sr *SubRecorder) Destroy() {
	watcher := sr.parent.watcher
	watcher.fsWatcher.removeRoot(sr.root)
	watcher.removeRecursive(sr.root)
}

const sysAppsCfgDir = "/var/lib/dde-daemon/apps"
const userAppsCfgDir = ".config/deepin/dde-daemon/apps"

// appsDir ex. $HOME/.local/share/applications
func getStatusFileAndOwner(uid int, home, appsDir string) (string, int) {
	var cfgDir string
	statusFileOwner := uid
	if home == "" {
		// system
		cfgDir = sysAppsCfgDir
		statusFileOwner = 0
	} else {
		if strings.HasPrefix(appsDir, home) {
			// user
			rel, err := filepath.Rel(home, appsDir)
			if err != nil {
				// home and appsDir all are abs path, so err should be nil
				panic(err)
			}
			appsDir = rel
			cfgDir = filepath.Join(home, userAppsCfgDir)
		} else {
			// system
			cfgDir = sysAppsCfgDir
			statusFileOwner = 0
		}
	}
	pathMd5 := md5.Sum([]byte(appsDir))
	base := fmt.Sprintf("launched-%x.csv", pathMd5)
	return filepath.Join(cfgDir, base), statusFileOwner
}

func (sr *SubRecorder) initAppLaunchedMap(apps []string) bool {
	var changed bool
	if launchedMap, err := loadStatusFromFile(sr.statusFile); err != nil {
		logger.Warning("SubRecorder.loadStatusFromFile failed", err)
		sr.resetStatus(apps)
		changed = true
	} else {
		sr.launchedMap = launchedMap
		changed = sr.checkStatus(launchedMap, apps)
	}
	return changed
}

func (sr *SubRecorder) RequestSave() {
	logger.Debug("SubRecorder.RequestSave", sr.root)
	sr.saveMu.Lock()
	if sr.isSaving {
		sr.saveMu.Unlock()
		return
	}

	sr.isSaving = true
	time.AfterFunc(2*time.Second, func() {
		err := sr.save()
		if err != nil {
			logger.Warning("SubRecorder.save error:", err)
		}
		sr.parent.emitStatusSaved(sr.root, sr.statusFile, err == nil)

		sr.saveMu.Lock()
		sr.isSaving = false
		sr.saveMu.Unlock()
	})
	sr.saveMu.Unlock()
}

func (sr *SubRecorder) writeStatus(w io.Writer) error {
	// NOTE: csv.NewWriter will new a bufio.Writer
	writer := csv.NewWriter(w)
	sr.launchedMapMu.RLock()
	err := writer.Write([]string{"# " + sr.root})
	if err != nil {
		logger.Warning(err)
	}
	for app, launched := range sr.launchedMap {
		record := make([]string, 2)
		record[0] = app
		if launched {
			record[1] = "t"
		} else {
			record[1] = "f"
		}
		if err := writer.Write(record); err != nil {
			sr.launchedMapMu.RUnlock()
			return err
		}
	}
	sr.launchedMapMu.RUnlock()

	writer.Flush()
	return writer.Error()
}

func (sr *SubRecorder) save() error {
	logger.Debug("SubRecorder.save", sr.root, sr.statusFile)
	file := sr.statusFile
	tmpFile := fmt.Sprintf("%s.new%x", file, time.Now().UnixNano())
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := sr.writeStatus(f); err != nil {
		return err
	}

	if err := os.Chown(tmpFile, sr.statusFileOwner, sr.statusFileOwner); err != nil {
		return err
	}

	if err := os.Rename(tmpFile, file); err != nil {
		return err
	}
	return nil
}

func (sr *SubRecorder) checkStatus(launchedMap map[string]bool, apps []string) bool {
	logger.Debug("SubRecorder.checkStatus", sr.root)
	var changed bool
	// apps -> to map appsMap
	appsMap := make(map[string]byte)
	for _, app := range apps {
		appsMap[app] = 0
	}
	for _, app := range apps {
		if _, ok := launchedMap[app]; !ok {
			// app added
			changed = true
			logger.Debugf("SubRecorder.checkStatus added %q", app)
			launchedMap[app] = false
		}
	}
	for app := range launchedMap {
		if _, ok := appsMap[app]; !ok {
			// app removed
			changed = true
			logger.Debugf("SubRecorder.checkStatus removed %q", app)
			delete(launchedMap, app)
		}
	}
	return changed
}

func loadStatusFromFile(dataFile string) (map[string]bool, error) {
	f, err := os.Open(dataFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	launchedMap := make(map[string]bool)
	reader := csv.NewReader(bufio.NewReader(f))
	reader.Comment = '#'
	reader.FieldsPerRecord = 2
	// parse csv file
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			logger.Warning("csv.Reader.Read error:", err)
			continue
		}
		app := record[0]
		var launched bool
		if record[1] == "t" {
			launched = true
		}
		launchedMap[app] = launched
	}
	return launchedMap, nil
}

func (sr *SubRecorder) resetStatus(apps []string) {
	launchedMap := make(map[string]bool)
	for _, app := range apps {
		launchedMap[app] = true
	}
	sr.launchedMap = launchedMap
}

func (sr *SubRecorder) uninstallHint(name string) {
	logger.Debug("SubRecorder.uninstallHint", sr.root, name)
	sr.launchedMapMu.Lock()

	sr.uninstallMap[name] = struct{}{}

	sr.launchedMapMu.Unlock()
}

func (sr *SubRecorder) handleAdded(name string) {
	logger.Debug("SubRecorder.handleAdded", sr.root, name)
	sr.launchedMapMu.Lock()

	if _, ok := sr.launchedMap[name]; !ok {
		launched, ok1 := sr.removedLaunchedMap[name]
		if ok1 {
			delete(sr.removedLaunchedMap, name)
		}
		logger.Debugf("app: %s, launched: %v", name, launched)
		sr.launchedMap[name] = launched
		sr.RequestSave()
	}
	delete(sr.uninstallMap, name)

	sr.launchedMapMu.Unlock()
}

func (sr *SubRecorder) handleRemoved(name string) {
	logger.Debug("SubRecorder.handleRemoved", sr.root, name)
	sr.launchedMapMu.Lock()
	defer sr.launchedMapMu.Unlock()

	launched, ok := sr.launchedMap[name]
	if !ok {
		return
	}
	_, uninstall := sr.uninstallMap[name]
	if !uninstall {
		sr.removedLaunchedMap[name] = launched
	} // else 卸装则不留下 launched 记录
	delete(sr.launchedMap, name)
	sr.RequestSave()
}

func (sr *SubRecorder) handleDirRemoved(name string) {
	logger.Debug("SubRecorder.handleDirRemoved", sr.root, name)
	changed := false

	sr.launchedMapMu.Lock()
	if name == "." {
		// applications dir removed
		if len(sr.launchedMap) > 0 {
			logger.Debug("SubRecorder.handleDirRemoved clear launchedMap")
			sr.launchedMap = make(map[string]bool)
			changed = true
		}
	} else {
		name = name + "/"
		// remove desktop entries under dir $name
		for app := range sr.launchedMap {
			if strings.HasPrefix(app, name) {
				logger.Debug("SubRecorder.handleDirRemoved remove", app)
				delete(sr.launchedMap, app)
				changed = true
			}
		}
	}
	sr.launchedMapMu.Unlock()

	if changed {
		sr.RequestSave()
	}
}

func (sr *SubRecorder) MarkLaunched(name string) bool {
	logger.Debug("SubRecorder.MarkLaunched", sr.root, name)
	sr.launchedMapMu.Lock()
	defer sr.launchedMapMu.Unlock()
	if launched, ok := sr.launchedMap[name]; ok {
		if !launched {
			sr.launchedMap[name] = true
			sr.RequestSave()
			return true
		}
	}
	return false
}

func (sr *SubRecorder) GetNew() (newApps []string) {
	sr.launchedMapMu.RLock()
	for app, launched := range sr.launchedMap {
		if !launched {
			newApps = append(newApps, app)
		}
	}
	sr.launchedMapMu.RUnlock()
	return
}
