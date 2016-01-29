/**
 * Copyright (c) 2015 Deepin, Inc.
 *               2015 Xu FaSheng
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

/*
Package bluetooth is a high level dbus warp for bluez5. You SHOULD always use dbus to call
the interface of Bluetooth.

The dbus interface of bluetooth to operate bluez is designed to asynchronous, and all result is returned by signal.
Other interface to get adapter/device informations will return immediately, because it was cached.
*/
package bluetooth
