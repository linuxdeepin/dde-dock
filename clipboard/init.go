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
