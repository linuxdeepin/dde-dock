/*
 * Copyright (C) 2016 ~ 2017 Deepin Technology Co., Ltd.
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
	"sync/atomic"
	"time"
)

type SubRecorder struct {
	root   string
	rootOk bool
	// key: app, value: launched
	launchedMap      map[string]bool
	launchedMapMutex sync.RWMutex

	statusFile      string
	statusFileOwner int

	uids     []int
	parent   *ALRecorder
	needSave int32 // if value == 1, need save status file
}

func NewSubRecorder(uid int, home, root string, parent *ALRecorder) *SubRecorder {
	sr := &SubRecorder{
		root:   root,
		uids:   []int{uid}, // first uid
		parent: parent,
	}

	sr.statusFile, sr.statusFileOwner = getStatusFileAndOwner(uid, home, root)
	logger.Debugf("NewSubRecorder status file: %q, owner: %d", sr.statusFile, sr.statusFileOwner)
	sr.initRoot()
	return sr
}

func (sr *SubRecorder) initRoot() {
	sr.rootOk = sr.getRootOk()
	if !sr.rootOk {
		return
	}

	subDirNames, apps := getDirsAndApps(sr.root)
	if sr.initAppLaunchedMap(apps) {
		MkdirAll(filepath.Dir(sr.statusFile), sr.statusFileOwner, getDirPerm(sr.statusFileOwner))
		sr.RequestSave()
	}

	for _, dirName := range subDirNames {
		path := filepath.Join(sr.root, dirName)
		sr.parent.watcher.add(path)
	}
}

func (sr *SubRecorder) doCheck() {
	rootOkChanged := sr.checkRoot()
	if rootOkChanged {
		if sr.rootOk {
			logger.Debugf("sr root %q Ok false => true", sr.root)
			// rootOk false => true
			// do init
			subDirNames, apps := getDirsAndApps(sr.root)
			if sr.initAppLaunchedMap(apps) {
				MkdirAll(filepath.Dir(sr.statusFile), sr.statusFileOwner, getDirPerm(sr.statusFileOwner))
				sr.RequestSave()
			}

			for _, dirName := range subDirNames {
				path := filepath.Join(sr.root, dirName)
				sr.parent.watcher.addRecursive(path, true)
			}
		} else {
			// rootOk true => false
			logger.Debugf("sr root %q Ok true => false", sr.root)
		}
	}

	sr.checkSave()
}

// return true if sr.rootOk changed
func (sr *SubRecorder) checkRoot() bool {
	oldRootOk := sr.rootOk
	sr.rootOk = sr.getRootOk()
	return oldRootOk != sr.rootOk
}

func (sr *SubRecorder) getRootOk() bool {
	// rootOk: root exist and is dir
	fileInfo, err := os.Stat(sr.root)
	if err != nil {
		return false
	}

	if fileInfo.IsDir() {
		return true
	}

	logger.Warning(sr.root, "is not a direcotry")
	return false
}

func (sr *SubRecorder) Destroy() {
	sr.parent.watcher.removeRecursive(sr.root)
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
	atomic.StoreInt32(&sr.needSave, 1)
}

func (sr *SubRecorder) checkSave() {
	if atomic.SwapInt32(&sr.needSave, 0) == 0 {
		return
	}

	sr.launchedMapMutex.RLock()
	err := sr.save()
	sr.launchedMapMutex.RUnlock()
	if err != nil {
		logger.Warning("SubRecorder.saveDirContent error:", err)
	}
	sr.parent.emitStatusSaved(sr.root, sr.statusFile, err == nil)
}

func (sr *SubRecorder) writeStatus(w io.Writer) error {
	writer := csv.NewWriter(w)
	sr.launchedMapMutex.RLock()
	writer.Write([]string{"# " + sr.root})
	for app, launched := range sr.launchedMap {
		record := make([]string, 2)
		record[0] = app
		if launched {
			record[1] = "t"
		} else {
			record[1] = "f"
		}
		if err := writer.Write(record); err != nil {
			return err
		}
	}
	sr.launchedMapMutex.RUnlock()

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
	if err := sr.writeStatus(bufio.NewWriter(f)); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
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

func (sr *SubRecorder) handleAdded(name string) {
	logger.Debug("SubRecorder.handleAdded", sr.root, name)
	sr.launchedMapMutex.Lock()

	if _, ok := sr.launchedMap[name]; !ok {
		sr.launchedMap[name] = false
		sr.RequestSave()
	}

	sr.launchedMapMutex.Unlock()
}

func (sr *SubRecorder) handleRemoved(name string) {
	logger.Debug("SubRecorder.handleRemoved", sr.root, name)
	sr.launchedMapMutex.Lock()

	if _, ok := sr.launchedMap[name]; ok {
		delete(sr.launchedMap, name)
		sr.RequestSave()
	}

	sr.launchedMapMutex.Unlock()
}

func (sr *SubRecorder) handleDirRemoved(name string) {
	logger.Debug("SubRecorder.handleDirRemoved", sr.root, name)
	changed := false

	sr.launchedMapMutex.Lock()
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
	sr.launchedMapMutex.Unlock()

	if changed {
		sr.RequestSave()
	}
}

func (sr *SubRecorder) MarkLaunched(name string) bool {
	logger.Debug("SubRecorder.MarkLaunched", sr.root, name)
	sr.launchedMapMutex.Lock()
	defer sr.launchedMapMutex.Unlock()
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
	sr.launchedMapMutex.RLock()
	for app, launched := range sr.launchedMap {
		if !launched {
			newApps = append(newApps, app)
		}
	}
	sr.launchedMapMutex.RUnlock()
	return
}
