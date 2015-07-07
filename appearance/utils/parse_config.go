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
	"fmt"
	"path"
	dutils "pkg.deepin.io/lib/utils"
)

const (
	regularFileMode = 0644
	dirFileMode     = 0755
)

var (
	errInvalidArgs  = fmt.Errorf("Invalid args")
	errInvalidKey   = fmt.Errorf("Invalid line key")
	errWriteKeyFile = fmt.Errorf("Write key to keyfile failed")
)

func GetUserQt4Config() string {
	dir := dutils.GetConfigDir()
	return path.Join(dir, "Trolltech.conf")
}
