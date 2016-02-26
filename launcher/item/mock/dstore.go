/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package mock

import (
	"math/rand"
	"path/filepath"
	"pkg.deepin.io/lib/dbus"
	"time"
)

type DStore struct {
	Count           int
	disconnectCount int
	handlers        map[int]func([][]interface{})
	softs           map[string]string
}

func (m *DStore) GetPkgNameFromPath(path string) (string, error) {
	relPath, _ := filepath.Rel(".", path)
	for pkgName, path := range m.softs {
		if path == relPath {
			return pkgName, nil
		}
	}
	return "", nil
}

func (m *DStore) sendMessage(msg [][]interface{}) {
	for _, fn := range m.handlers {
		fn(msg)
	}
}

func (m *DStore) sendStartMessage(pkgName string) {
	action := []interface{}{
		"",
		dbus.MakeVariant([]interface{}{pkgName, int32(0)}),
	}
	m.sendMessage([][]interface{}{action})
}

func (m *DStore) sendUpdateMessage(pkgName string) {
	updateTime := rand.Intn(5) + 1
	for i := 0; i < updateTime; i++ {
		action := []interface{}{
			"",
			dbus.MakeVariant([]interface{}{
				pkgName,
				int32(1),
				int32(int(i+1) / updateTime),
				"update",
			}),
		}
		m.sendMessage([][]interface{}{action})
		time.Sleep(time.Duration(rand.Int31n(100)+100) * time.Millisecond)
	}
}
func (m *DStore) sendFinishedMessage(pkgName string) {
	action := []interface{}{
		"",
		dbus.MakeVariant([]interface{}{
			pkgName,
			int32(2),
			[][]interface{}{
				[]interface{}{
					pkgName,
					true,
					false,
					false,
				},
			},
		}),
	}
	m.sendMessage([][]interface{}{action})
}

func (m *DStore) sendFailedMessage(pkgName string) {
	action := []interface{}{
		"",
		dbus.MakeVariant([]interface{}{
			pkgName,
			int32(3),
			[][]interface{}{
				[]interface{}{
					pkgName,
					false,
					false,
					false,
				},
			},
			"uninstall failed",
		}),
	}
	m.sendMessage([][]interface{}{action})
}

func (m *DStore) UninstallPkg(pkgName string, purge bool) error {
	if _, ok := m.softs[pkgName]; !ok {
		m.sendFailedMessage(pkgName)
		return nil
	}
	m.sendStartMessage(pkgName)
	m.sendUpdateMessage(pkgName)
	m.sendFinishedMessage(pkgName)
	return nil
}

func (m *DStore) Connectupdate_signal(fn func([][]interface{})) func() {
	id := m.Count
	m.handlers[id] = fn
	m.Count++
	return func() {
		delete(m.handlers, id)
		m.disconnectCount++
	}
}

func NewDStore() *DStore {
	return &DStore{
		handlers: map[int]func([][]interface{}){},
		Count:    0,
		softs: map[string]string{
			"firefox":             "../testdata/firefox.desktop",
			"deepin-music-player": "../testdata/deepin-music-player.desktop",
		},
	}
}
