package main

import "dbus/com/deepin/sessionmanager"
import "os"
import "pkg.deepin.io/lib/utils"

func ddeSessionRegister() {
	cookie := os.ExpandEnv("$DDE_SESSION_PROCESS_COOKIE_ID")
	utils.UnsetEnv("DDE_SESSION_PROCESS_COOKIE_ID")
	if cookie == "" {
		logger.Warning("get DDE_SESSION_PROCESS_COOKIE_ID failed")
		return
	}
	go func() {
		manager, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager")
		if err != nil {
			logger.Warning("register failed:", err)
			return
		}
		manager.Register(cookie)
	}()
}
