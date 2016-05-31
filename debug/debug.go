/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package debug

import (
	_ "net/http/pprof"

	"net/http"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var (
	logger = log.NewLogger("daemon/keybinding")
)

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon() *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("debug", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	go http.ListenAndServe("localhost:6969", nil)

	return nil
}

func (*Daemon) Stop() error {
	return nil
}
