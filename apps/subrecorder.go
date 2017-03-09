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
	root string
	// key: app, value: launched
	launchedMap      map[string]bool
	launchedMapMutex sync.RWMutex
	statusFile       string
	uid              int

	saveCb    ALStatusSavedFun
	saveCount int32
	ticker    *time.Ticker
}

func (d *SubRecorder) initSaveTicker() {
	// checkSave every 2 seconds
	const interval = 2
	const tickerTimeout = interval + 5
	d.ticker = time.NewTicker(time.Second * interval)

	go func() {
		timer := time.NewTimer(time.Second * tickerTimeout)
		for {
			select {
			case <-d.ticker.C:
				d.checkSave()
				timer.Reset(time.Second * tickerTimeout)
			case <-timer.C:
				logger.Debug("stop check save")
				return
			}
		}
	}()
}

func NewSubRecorder(home, root string, uid int, apps []string, saveCb ALStatusSavedFun) *SubRecorder {
	d := &SubRecorder{
		root:       root,
		uid:        uid,
		statusFile: getStatusFile(home, root),
		saveCb:     saveCb,
	}
	if d.initAppLaunchedMap(apps) {
		MkdirAll(filepath.Dir(d.statusFile), d.uid, getDirPerm(d.uid))
		d.RequestSave()
	}
	d.initSaveTicker()
	logger.Debug("NewSubRecorder status file:", d.statusFile)
	logger.Debug("launchedMap len:", len(d.launchedMap))
	return d
}

func (d *SubRecorder) Destroy() {
	d.ticker.Stop()
}

func getStatusFile(home, path string) string {
	var dir string
	if home == "" {
		// system
		dir = "/var/lib/dde-daemon/apps"
	} else {
		// user
		rel, err := filepath.Rel(home, path)
		if err == nil {
			path = rel
		}
		dir = filepath.Join(home, ".config/deepin/dde-daemon/apps")
	}
	//logger.Debug("getStatusFile path", path)
	pathMd5 := md5.Sum([]byte(path))
	base := fmt.Sprintf("launched-%x.csv", pathMd5)
	return filepath.Join(dir, base)
}

func (d *SubRecorder) initAppLaunchedMap(apps []string) bool {
	var changed bool
	if launchedMap, err := loadStatusFromFile(d.statusFile); err != nil {
		logger.Warning("SubRecorder.loadStatusFromFile failed", err)
		d.resetStatus(apps)
		changed = true
	} else {
		d.launchedMap = launchedMap
		changed = d.checkStatus(launchedMap, apps)
	}
	return changed
}

func (d *SubRecorder) RequestSave() {
	atomic.AddInt32(&d.saveCount, 1)
}

func (d *SubRecorder) checkSave() {
	if atomic.SwapInt32(&d.saveCount, 0) == 0 {
		return
	}

	d.launchedMapMutex.RLock()
	err := d.save()
	if d.saveCb != nil {
		d.saveCb(d.root, d.statusFile, err == nil)
	}
	d.launchedMapMutex.RUnlock()
	if err != nil {
		logger.Warning("SubRecorder.saveDirContent error:", err)
	}
}

func (d *SubRecorder) writeStatus(w io.Writer) error {
	writer := csv.NewWriter(w)
	d.launchedMapMutex.RLock()
	writer.Write([]string{"# " + d.root})
	for app, launched := range d.launchedMap {
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
	d.launchedMapMutex.RUnlock()

	writer.Flush()
	return writer.Error()
}

func (d *SubRecorder) save() error {
	logger.Debug("SubRecorder.save", d.root, d.statusFile)
	file := d.statusFile
	uid := d.uid

	tmpFile := fmt.Sprintf("%s.new%x", file, time.Now().UnixNano())
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	if err := d.writeStatus(bufio.NewWriter(f)); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	if err := os.Chown(tmpFile, uid, uid); err != nil {
		return err
	}

	if err := os.Rename(tmpFile, file); err != nil {
		return err
	}
	return nil
}

func (d *SubRecorder) checkStatus(launchedMap map[string]bool, apps []string) bool {
	logger.Debug("SubRecorder.checkStatus", d.root)
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

func (d *SubRecorder) resetStatus(apps []string) {
	launchedMap := make(map[string]bool)
	for _, app := range apps {
		launchedMap[app] = true
	}
	d.launchedMap = launchedMap
}

func (d *SubRecorder) handleAdded(name string) {
	logger.Debug("SubRecorder.handleAdded", d.root, name)
	d.launchedMapMutex.Lock()
	defer d.launchedMapMutex.Unlock()
	launchedMap := d.launchedMap
	if _, ok := launchedMap[name]; !ok {
		launchedMap[name] = false
		d.RequestSave()
	}
}

func (d *SubRecorder) handleRemoved(name string) {
	logger.Debug("SubRecorder.handleRemoved", d.root, name)
	launchedMap := d.launchedMap
	d.launchedMapMutex.Lock()
	defer d.launchedMapMutex.Unlock()
	if _, ok := launchedMap[name]; ok {
		delete(launchedMap, name)
		d.RequestSave()
	}
}

func (d *SubRecorder) handleDirRemoved(name string) {
	logger.Debug("SubRecorder.handleDirRemoved", d.root, name)
	d.launchedMapMutex.Lock()
	defer d.launchedMapMutex.Unlock()

	changed := false
	if name == "." {
		// applications dir removed
		if len(d.launchedMap) > 0 {
			logger.Debug("SubRecorder.handleDirRemoved clear launchedMap")
			d.launchedMap = make(map[string]bool)
			changed = true
		}
	} else {
		name = name + "/"
		// remove desktop entries under dir $name
		launchedMap := d.launchedMap
		for app := range launchedMap {
			if strings.HasPrefix(app, name) {
				delete(launchedMap, app)
				logger.Debug("SubRecorder.handleDirRemoved remove", app)
				changed = true
			}
		}
	}

	if changed {
		d.RequestSave()
	}
}

func (d *SubRecorder) MarkLaunched(name string) bool {
	logger.Debug("SubRecorder.MarkLaunched", d.root, name)
	d.launchedMapMutex.Lock()
	defer d.launchedMapMutex.Unlock()
	if launched, ok := d.launchedMap[name]; ok {
		if !launched {
			d.launchedMap[name] = true
			d.RequestSave()
			return true
		}
	}
	return false
}

func (r *SubRecorder) GetNew() (newApps []string) {
	r.launchedMapMutex.RLock()
	for app, launched := range r.launchedMap {
		if !launched {
			newApps = append(newApps, app)
		}
	}
	r.launchedMapMutex.RUnlock()
	return
}
