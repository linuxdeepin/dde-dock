/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

/*
Package bluetooth is a high level dbus warp for bluez5. You SHOULD always use dbus to call
the interface of Bluetooth.

The dbus interface of bluetooth to operate bluez is designed to asynchronous, and all result is returned by signal.
Other interface to get adapter/device informations will return immediately, because it was cached.
*/
package bluetooth
