/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package clipboard

// #cgo pkg-config: gtk+-3.0 x11 glib-2.0
// #cgo CFLAGS: -Wall -g
// #include "gsd-clipboard-manager.h"
import "C"

import . "pkg.deepin.io/dde/daemon/loader"
import "pkg.deepin.io/lib/log"

var logger = log.NewLogger("daemon/clipboard")

func init() {
	Register(NewClipboardDaemon(logger))
	//loader.Register(&loader.Module{
	//Name:   "clipboard",
	//Start:  Start,
	//Stop:   Stop,
	//Enable: true,
	//})
}

type Daemon struct {
	*ModuleBase
}

func NewClipboardDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = NewModuleBase("clipboard", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	C.start_clip_manager()
	return nil
}

func (*Daemon) Stop() error {
	C.stop_clip_manager()
	return nil
}
