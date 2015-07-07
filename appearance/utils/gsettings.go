/**
 * Copyright (c) 2011 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 jouyouyun
 *
 * Author:      jouyouyun <jouyouwen717@gmail.com>
 * Maintainer:  jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package utils

import (
	"pkg.deepin.io/lib/gio-2.0"
	dutils "pkg.deepin.io/lib/utils"
	"sync"
)

var (
	locker sync.Mutex

	schemaGSettingsMap = make(map[string]*gio.Settings)
	gsettingsCountMap  = make(map[*gio.Settings]int)
)

func NewGSettings(schema string) *gio.Settings {
	s, ok := schemaGSettingsMap[schema]
	if !ok {
		return newGsettings(schema, false)
	}

	locker.Lock()
	gsettingsCountMap[s]++
	locker.Unlock()

	return s
}

func CheckAndNewGSettings(schema string) *gio.Settings {
	s, ok := schemaGSettingsMap[schema]
	if !ok {
		return newGsettings(schema, true)
	}

	locker.Lock()
	gsettingsCountMap[s]++
	locker.Unlock()

	return s
}

func Unref(s *gio.Settings) {
	if s == nil {
		return
	}

	_, ok := gsettingsCountMap[s]
	if !ok {
		return
	}

	locker.Lock()
	gsettingsCountMap[s]--
	if gsettingsCountMap[s] == 0 {
		delete(gsettingsCountMap, s)
		deleteSchemaByGSettings(s)
		s.Unref()
		s = nil
	}
	locker.Unlock()

	return
}

func newGsettings(schema string, check bool) *gio.Settings {
	if check {
		if !dutils.IsGSchemaExist(schema) {
			return nil
		}
	}

	s := gio.NewSettings(schema)

	locker.Lock()
	schemaGSettingsMap[schema] = s
	gsettingsCountMap[s]++
	locker.Unlock()

	return s
}

func deleteSchemaByGSettings(s *gio.Settings) {
	var schema string
	for k, v := range schemaGSettingsMap {
		if v == s {
			schema = k
			break
		}
	}

	delete(schemaGSettingsMap, schema)
}
