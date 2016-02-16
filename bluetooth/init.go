/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package bluetooth

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/bluetooth")

func init() {
	loader.Register(newBluetoothDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "bluetooth",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
