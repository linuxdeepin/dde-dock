/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import "dbus/com/deepin/sessionmanager"
import "os"
import "pkg.deepin.io/lib/utils"

func ddeSessionRegister() {
	cookie := os.ExpandEnv("$DDE_SESSION_PROCESS_COOKIE_ID")
	utils.UnsetEnv("DDE_SESSION_PROCESS_COOKIE_ID")
	if cookie == "" {
		return
	}
	go func() {
		manager, err := sessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager")
		if err != nil {
			return
		}
		manager.Register(cookie)
	}()
}
