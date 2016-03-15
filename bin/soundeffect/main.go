/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"os"
	"pkg.deepin.io/dde/daemon/soundeffect"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"time"
)

func main() {
	logger := log.NewLogger("daemon/soundeffect-runner")
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSession(soundeffect.DBusDest) {
		logger.Error("dbus not unique:", soundeffect.DBusDest)
		return
	}

	dbus.SetAutoDestroyHandler(5*time.Second, func() bool {
		return !soundeffect.IsPlaying()
	})

	soundeffect.Start()
	dbus.DealWithUnhandledMessage()
	if err := dbus.Wait(); err != nil {
		logger.Error("lost dbus session:", err)
		soundeffect.Stop()
		os.Exit(1)
	}
	soundeffect.Stop()
	os.Exit(0)
}
