/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package dock

import (
	"dbus/com/deepin/daemon/display"
	"github.com/BurntSushi/xgb/xproto"
)

var displayRect xproto.Rectangle

func initDisplay() error {
	var err error
	dpy, err = display.NewDisplay(
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display",
	)
	if err != nil {
		return err
	}
	// to avoid get PrimaryRect failed
	defer func() {
		if r := recover(); r != nil {
			logger.Warning("Recovered in initDisplay:", r)
		}
	}()
	setDisplayRect(dpy.PrimaryRect.Get())
	dpy.ConnectPrimaryChanged(func(rect []interface{}) {
		setDisplayRect(rect)
		hideModemanager.updateStateWithoutDelay()
	})
	return nil
}

func setDisplayRect(rect []interface{}) {
	if len(rect) != 4 {
		return
	}
	displayRect.X, _ = rect[0].(int16)
	displayRect.Y, _ = rect[1].(int16)
	displayRect.Width, _ = rect[2].(uint16)
	displayRect.Height, _ = rect[3].(uint16)
}
