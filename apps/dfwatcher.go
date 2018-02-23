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
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/fsnotify"
)

const (
	DFWatcherDBusInterface = AppsDBusDest + ".DesktopFileWatcher"
)

// desktop file watcher
type DFWatcher struct {
	fsWatcher *fsnotify.Watcher
	sem       chan int
	eventChan chan *FileEvent
	// signal:
	Event func(name string, op uint32)
}

func NewDFWachter() (*DFWatcher, error) {
	w := new(DFWatcher)
	if fsWatcher, err := fsnotify.NewWatcher(); err != nil {
		return nil, err
	} else {
		w.fsWatcher = fsWatcher
	}

	w.sem = make(chan int, 4)
	w.eventChan = make(chan *FileEvent, 10)
	go w.listenEvents()
	return w, nil
}

func (w *DFWatcher) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       AppsDBusDest,
		ObjectPath: AppsObjectPath,
		Interface:  DFWatcherDBusInterface,
	}
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
	dbus.Emit(w, "Event", ev.Name, uint32(0))
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
				logger.Warning(err)
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
