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
	"errors"
	"github.com/BurntSushi/xgbutil/xrect"
)

func (m *DockManager) initDisplay() error {
	var err error
	m.dpy, err = display.NewDisplay(
		"com.deepin.daemon.Display",
		"/com/deepin/daemon/Display",
	)
	if err != nil {
		return err
	}

	err = m.setDisplayPrimaryRect(m.dpy.PrimaryRect.Get())
	if err != nil {
		return err
	}
	// set dock default width to primary width
	m.dockWidth = m.displayPrimaryRect.Width()
	m.updateDockRect()

	m.dpy.ConnectPrimaryChanged(func(rect []interface{}) {
		err := m.setDisplayPrimaryRect(rect)
		if err != nil {
			logger.Warning("error on primary changed", err)
			return
		}
		m.updateDockRect()
	})
	return nil
}

func (m *DockManager) setDisplayPrimaryRect(rect []interface{}) error {
	if len(rect) != 4 {
		return errors.New("setDisplayPrimaryRect failed: len(rect) != 4")
	}
	var primaryX, primaryY, primaryW, primaryH int

	if x, ok := rect[0].(int16); !ok {
		return errors.New("setDisplayPrimaryRect failed: convert x failed")
	} else {
		primaryX = int(x)
	}

	if y, ok := rect[1].(int16); !ok {
		return errors.New("setDisplayPrimaryRect failed: convert y failed")
	} else {
		primaryY = int(y)
	}

	if w, ok := rect[2].(uint16); !ok {
		return errors.New("setDisplayPrimaryRect failed: convert w failed")
	} else {
		primaryW = int(w)
	}

	if h, ok := rect[3].(uint16); !ok {
		return errors.New("setDisplayPrimaryRect failed: convert h failed")
	} else {
		primaryH = int(h)
	}

	if m.displayPrimaryRect == nil {
		m.displayPrimaryRect = xrect.New(primaryX, primaryY, primaryW, primaryH)
	} else {
		m.displayPrimaryRect.XSet(primaryX)
		m.displayPrimaryRect.YSet(primaryY)
		m.displayPrimaryRect.WidthSet(primaryW)
		m.displayPrimaryRect.HeightSet(primaryH)
	}
	logger.Debug("set display primary rect:", m.displayPrimaryRect)
	return nil
}
