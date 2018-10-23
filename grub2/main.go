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

package grub2

import (
	"time"

	"pkg.deepin.io/lib/dbusutil"
)

var _g *Grub2

func RunAsDaemon() {
	allowNoCheckAuth()
	service, err := dbusutil.NewSystemService()
	if err != nil {
		logger.Fatal("failed to new system service", err)
	}
	_g = New(service)

	err = service.Export(dbusPath, _g)
	if err != nil {
		logger.Fatal("failed to export grub2:", err)
	}

	err = service.Export(themeDBusPath, _g.theme)
	if err != nil {
		logger.Fatal("failed to export grub2 theme:", err)
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Fatal("failed to request name:", err)
	}

	service.SetAutoQuitHandler(5*time.Minute, _g.canSafelyExit)
	service.Wait()
}

// write default /etc/default/grub
// generate theme background image file
// call from deepin-installer hooks/in_chroot/*_setup_bootloader_x86.job
func Setup(resolution string) error {
	return nil
}

// call from grub-themes-deepin debian/postinst
func SetupTheme() error {
	return nil
}
