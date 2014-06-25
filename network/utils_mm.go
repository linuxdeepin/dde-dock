/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
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

package network

import (
	"dbus/org/freedesktop/modemmanager1"
	"pkg.linuxdeepin.com/lib/dbus"
)

const dbusMmDest = "org.freedesktop.ModemManager1"

func mmNewModem(modemPath dbus.ObjectPath) (modem *modemmanager1.Modem, err error) {
	modem, err = modemmanager1.NewModem(dbusMmDest, modemPath)
	if err != nil {
		logger.Error(err)
	}
	return
}

func mmGetModemDeviceIdentifier(modemPath dbus.ObjectPath) (devId string, err error) {
	modem, err := mmNewModem(modemPath)
	if err != nil {
		return
	}
	devId = modem.DeviceIdentifier.Get()
	return
}
