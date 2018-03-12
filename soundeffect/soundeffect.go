/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package soundeffect

import (
	"time"

	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/soundeffect")

func Run() {
	service, err := dbusutil.NewSessionService()
	if err != nil {
		logger.Fatal("failed to new session service:", err)
	}
	m := NewManager(service)

	err = service.Export(dbusPath, m)
	if err != nil {
		logger.Fatal("failed to export:", err)
	}

	err = gsettings.StartMonitor()
	if err != nil {
		logger.Fatal("failed to start monitor gsettings changed:", err)
	}
	err = service.RequestName(DBusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}

	service.SetAutoQuitHandler(30*time.Second, m.canQuit)
	service.Wait()
}
