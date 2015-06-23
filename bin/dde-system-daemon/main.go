package main

import "pkg.linuxdeepin.com/lib/log"

import "pkg.linuxdeepin.com/lib"
import "pkg.linuxdeepin.com/lib/dbus"
import "os"
import _ "pkg.linuxdeepin.com/dde-daemon/accounts"
import "pkg.linuxdeepin.com/dde-daemon/loader"
import . "pkg.linuxdeepin.com/lib/gettext"

var logger = log.NewLogger("dde-daemon/dde-system-daemon")

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem("com.deepin.daemon") {
		logger.Warning("There already has an dde-daemon running.")
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
