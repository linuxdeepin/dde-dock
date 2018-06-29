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
	"os"
	"path/filepath"
	"time"

	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/fsnotify"
	"pkg.deepin.io/lib/log"
)

const (
	dfWatcherDBusInterface = dbusServiceName + ".DesktopFileWatcher"
)

// desktop file watcher
type DFWatcher struct {
	service   *dbusutil.Service
	fsWatcher *fsWatcher
	sem       chan int
	eventChan chan *FileEvent

	signals *struct {
		Event struct {
			name string
			op   uint32
		}
	}
}

func newDFWatcher(service *dbusutil.Service) (*DFWatcher, error) {
	w := new(DFWatcher)

	interval := 6 * time.Second
	if logger.GetLogLevel() == log.LevelDebug {
		interval = 3 * time.Second
	}
	fsWatcher, err := newFsWatcher(interval)
	if err != nil {
		return nil, err
	}
	fsWatcher.trySuccessCb = func(file string) {
		w.addRecursive(file, true)
	}
	w.fsWatcher = fsWatcher
	w.service = service
	w.sem = make(chan int, 4)
	w.eventChan = make(chan *FileEvent, 10)
	go w.listenEvents()
	return w, nil
}

func (*DFWatcher) GetInterfaceName() string {
	return dfWatcherDBusInterface
}

func (w *DFWatcher) listenEvents() {
	for {
		select {
		case ev, ok := <-w.fsWatcher.Event:
			if !ok {
				logger.Error("Invalid event:", ev)
				return
			}
			logger.Debug("event", ev)
			w.handleEvent(ev)

		case err := <-w.fsWatcher.Error:
			logger.Warning("error", err)
			return
		}
	}
}

func (w *DFWatcher) handleEvent(event *fsnotify.FileEvent) {
	w.fsWatcher.handleEvent(event)
	ev := NewFileEvent(event)
	file := ev.Name
	if (ev.IsCreate() || ev.IsRename()) && ev.IsDir {
		// it exist and is dir
		w.addRecursive(file, true)
		return
	}
	w.notifyEvent(ev)
}

func (w *DFWatcher) notifyEvent(ev *FileEvent) {
	logger.Debugf("notifyEvent %q", ev.Name)
	err := w.service.Emit(w, "Event", ev.Name, uint32(0))
	if err != nil {
		logger.Warning(err)
	}
	w.eventChan <- ev
}

func (w *DFWatcher) add(path string) error {
	logger.Debug("DFWatcher.add", path)
	return w.fsWatcher.Watch(path)
}

func (w *DFWatcher) remove(path string) error {
	logger.Debug("DFWatcher.remove", path)
	return w.fsWatcher.RemoveWatch(path)
}

func (w *DFWatcher) addRecursive(path string, loadExisted bool) {
	logger.Debug("DFWatcher.addRecursive", path, loadExisted)
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logger.Warning(err)
			return nil
		}
		if info.IsDir() {
			logger.Debug("DFWatcher.addRecursive watch", path)
			err := w.add(path)
			if err != nil {
				logger.Warning(err)
			}
		} else if loadExisted {
			if isDesktopFile(path) {
				w.notifyEvent(NewFileFoundEvent(path))
			}
		}
		return nil
	})
}

func (w *DFWatcher) removeRecursive(path string) {
	logger.Debug("DFWatcher.removeRecursive", path)
	filepath.Walk(path,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				// ignore file not exist error
				if !os.IsNotExist(err) {
					logger.Warning(err)
				}
				return nil
			}

			if info.IsDir() {
				err := w.remove(path)
				if err != nil {
					logger.Warning(err)
				}
			}
			return nil
		})
}
