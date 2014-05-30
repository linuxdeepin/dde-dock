package main

import "dlib/logger"
import "dlib"
import "dlib/dbus"
import "os"
import "dde-daemon/accounts"
import "dde-daemon/systeminfo"

var Logger = logger.NewLogger("com.deepin.daemon")

func main() {
	if !dlib.UniqueOnSystem("com.deepin.daemon") {
		Logger.Warning("There already has an dde-daemon running.")
		return
	}

	go accounts.Start()
	go systeminfo.Start()

	if err := dbus.Wait(); err != nil {
		Logger.Error("dde-daemon lost dbus")
		os.Exit(-1)
	} else {
		os.Exit(0)
	}
}
