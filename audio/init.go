/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package audio

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/audio")

func init() {
	loader.Register(NewAudioDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "audio",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}
