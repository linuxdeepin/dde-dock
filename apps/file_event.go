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
	"pkg.deepin.io/lib/fsnotify"
)

type FileEvent struct {
	*fsnotify.FileEvent
	NotExist bool
	IsDir    bool
	IsFound  bool
}

func NewFileFoundEvent(name string) *FileEvent {
	ev := &fsnotify.FileEvent{
		Name: name,
	}
	return &FileEvent{
		FileEvent: ev,
		IsFound:   true,
	}
}

func NewFileEvent(ev *fsnotify.FileEvent) *FileEvent {
	var notExist bool
	var isDir bool
	if stat, err := os.Stat(ev.Name); os.IsNotExist(err) {
		notExist = true
	} else if err == nil {
		isDir = stat.IsDir()
	}
	return &FileEvent{
		FileEvent: ev,
		NotExist:  notExist,
		IsDir:     isDir,
	}
}
