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

package main

import "os"
import . "pkg.deepin.io/dde/daemon/langselector"
import "pkg.deepin.io/lib/gettext"
import "pkg.deepin.io/lib/dbus"
import "time"

func main() {
	gettext.InitI18n()
	gettext.Textdomain("dde-daemon")

	lang := Start()
	if lang == nil {
		return
	}

	dbus.DealWithUnhandledMessage()

	dbus.SetAutoDestroyHandler(time.Second*5, func() bool {
		if lang.LocaleState == LocaleStateChanging {
			return false
		} else {
			return true
		}
	})

	if err := dbus.Wait(); err != nil {
		Stop()
		os.Exit(-1)
	}

	Stop()
	os.Exit(0)
}
