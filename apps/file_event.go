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
	"github.com/fsnotify/fsnotify"
	"os"
)

type FileEvent struct {
	fsnotify.Event
	NotExist bool
	IsDir    bool
}

func NewFileCreatedEvent(name string) *FileEvent {
	return &FileEvent{
		Event: fsnotify.Event{
			Name: name,
			Op:   fsnotify.Create,
		},
	}
}

func NewFileEvent(ev fsnotify.Event) *FileEvent {
	var notExist bool
	var isDir bool
	if stat, err := os.Stat(ev.Name); os.IsNotExist(err) {
		notExist = true
	} else if err == nil {
		isDir = stat.IsDir()
	}
	return &FileEvent{
		Event:    ev,
		NotExist: notExist,
		IsDir:    isDir,
	}
}
