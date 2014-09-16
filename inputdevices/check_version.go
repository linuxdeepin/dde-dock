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

package inputdevices

import (
	"os"
	"path"
	dutils "pkg.linuxdeepin.com/lib/utils"
)

const (
	_VERSION     = "0.1"
	_VERSION_DIR = ".config/dde-daemon/inputdevices"
)

func (m *Manager) isVersionRight() bool {
	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		logger.Debug("GetHomeDir Failed")
		return false
	}

	versionFile := path.Join(homeDir, _VERSION_DIR, "version")
	if !dutils.IsFileExist(versionFile) {
		m.newVersionFile()
		return false
	}

	return true
}

func (m *Manager) newVersionFile() {
	homeDir := dutils.GetHomeDir()
	if len(homeDir) < 1 {
		logger.Debug("GetHomeDir Failed")
		return
	}

	vDir := path.Join(homeDir, _VERSION_DIR)
	if !dutils.IsFileExist(vDir) {
		if err := os.MkdirAll(vDir, 0755); err != nil {
			logger.Warningf("MkdirAll '%s' failed: %v", vDir, err)
			return
		}
	}

	vFile := path.Join(vDir, "version")
	fp, err := os.Create(vFile)
	if err != nil {
		logger.Warningf("Create '%s' failed: %v", vFile, err)
		return
	}
	defer fp.Close()

	fp.WriteString(_VERSION)
	fp.Sync()

	return
}
