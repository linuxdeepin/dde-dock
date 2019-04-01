/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package main

import (
	"time"

	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("GreeterSetter")

func main() {
	service, err := dbusutil.NewSystemService()
	if err != nil {
		logger.Errorf("failed to new system service")
		return
	}

	var m = &Manager{
		service: service,
	}

	err = service.Export(dbusPath, m)
	if err != nil {
		logger.Errorf("failed to export:", err)
		return
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Errorf("failed to request name:", err)
		return
	}
	service.SetAutoQuitHandler(time.Second*30, nil)
	service.Wait()
	return
}
