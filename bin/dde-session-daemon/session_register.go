package main

import "dbus/com/deepin/sessionmanager"
import "os"
import "pkg.deepin.io/lib/utils"

func ddeSessionRegister() {
	cookie := os.ExpandEnv("$DDE_SESSION_PROCESS_COOKIE_ID")
	utils.UnsetEnv("DDE_SESSION_PROCESS_COOKIE_ID")
	manager, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager")
	if err != nil {
		return
	}
	manager.Register(cookie)
}
