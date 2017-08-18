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
	"pkg.deepin.io/lib/fsnotify"
	"os"
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
