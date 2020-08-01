/*
 * Copyright (C) 2013 ~ 2018 Deepin Technology Co., Ltd.
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

type Manager struct {
	service    *dbusutil.Service
	writeStart bool
	writeEnd   chan bool

	methods *struct { //nolint
		NewSearchWithStrList  func() `in:"list" out:"md5sum,ok"`
		NewSearchWithStrDict  func() `in:"dict" out:"md5sum,ok"`
		SearchString          func() `in:"str,md5sum" out:"result"`
		SearchStartWithString func() `in:"str,md5sum" out:"result"`
	}
}

const (
	dbusServiceName = "com.deepin.daemon.Search"
	dbusPath        = "/com/deepin/daemon/Search"
	dbusInterface   = "com.deepin.daemon.Search"
)

var (
	logger = log.NewLogger("daemon/search")
)

func newManager(service *dbusutil.Service) *Manager {
	m := Manager{
		service: service,
	}

	m.writeStart = false

	return &m
}

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()
	logger.SetRestartCommand("/usr/lib/deepin-daemon/search")

	service, err := dbusutil.NewSessionService()
	if err != nil {
		logger.Fatal("failed to new session service:", err)
	}

	m := newManager(service)
	err = service.Export(dbusPath, m)
	if err != nil {
		logger.Fatal("failed to export:", err)
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}

	service.SetAutoQuitHandler(time.Second*5, func() bool {
		if m.writeStart {
			select { //nolint
			case <-m.writeEnd:
				return true
			}
		}
		return true
	})
	service.Wait()
}
