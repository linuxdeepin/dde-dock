package main

import "pkg.linuxdeepin.com/lib/log"

import "pkg.linuxdeepin.com/lib"
import "pkg.linuxdeepin.com/lib/dbus"
import "os"
import _ "pkg.linuxdeepin.com/dde-daemon/accounts"
import "pkg.linuxdeepin.com/dde-daemon"

var logger = log.NewLogger("com.deepin.daemon")

func main() {
	logger.BeginTracing()
	defer logger.EndTracing()

	if !lib.UniqueOnSystem("com.deepin.daemon") {
		logger.Warning("There already has an dde-daemon running.")
		return
	}

	logger.SetRestartCommand("/usr/lib/deepin-daemon/dde-system-daemon")

	loader.Start()
	defer loader.Stop()

	dbus.DealWithUnhandledMessage()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
