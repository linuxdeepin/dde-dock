package main

import "pkg.deepin.io/lib/log"

import "pkg.deepin.io/lib"
import "pkg.deepin.io/lib/dbus"
import "os"
import _ "pkg.deepin.io/dde/daemon/accounts"
import "pkg.deepin.io/dde/daemon/loader"
import . "pkg.deepin.io/lib/gettext"

var logger = log.NewLogger("daemon/dde-system-daemon")

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem("com.deepin.daemon") {
		logger.Warning("There already has an dde daemon running.")
		return
	}

	InitI18n()
	Textdomain("dde-daemon")

	logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-system-daemon")

	loader.StartAll()
	defer loader.StopAll()

	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
