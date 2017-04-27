/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package sessionwatcher

import (
	"dbus/org/freedesktop/login1"
	"fmt"
	"os/user"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/pulse"
)

func newLoginSession(path dbus.ObjectPath) (uint32, *login1.Session) {
	session, err := login1.NewSession(login1Dest, path)
	if err != nil {
		logger.Warning("New session '(%v)%s' failed: %v", path, err)
		return 0, nil
	}

	uinfo := session.User.Get()
	if len(uinfo) != 2 {
		logger.Warning("Invalid session user info:", path)
		login1.DestroySession(session)
		return 0, nil
	}

	uid, ok := uinfo[0].(uint32)
	if !ok {
		logger.Warning("Get session uid failed:", path)
		login1.DestroySession(session)
		return 0, nil
	}
	return uid, session
}

func suspendPulseSinks(suspend int) {
	var ctx = pulse.GetContext()
	if ctx == nil {
		logger.Error("Failed to connect pulseaudio server")
		return
	}
	for _, sink := range ctx.GetSinkList() {
		ctx.SuspendSinkById(sink.Index, suspend)
	}
}

func suspendPulseSources(suspend int) {
	var ctx = pulse.GetContext()
	if ctx == nil {
		logger.Error("Failed to connect pulseaudio server")
		return
	}
	for _, source := range ctx.GetSourceList() {
		ctx.SuspendSourceById(source.Index, suspend)
	}
}

var curUid string

func isCurrentUser(uid uint32) bool {
	if len(curUid) == 0 {
		cur, err := user.Current()
		if err != nil {
			return false
		}
		curUid = cur.Uid
	}

	return curUid == fmt.Sprint(uid)
}
