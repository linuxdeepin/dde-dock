/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package power

//this should only use org.freedesktop.ScreenSaver interface with SimulateUserActivity

import (
	"dbus/org/freedesktop/screensaver"
	//"pkg.deepin.io/lib/logger"
	"fmt"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

type fullScreenWorkaround struct {
	xu               *xgbutil.XUtil
	targets          []string
	activeWindowAtom xproto.Atom

	ss         *screensaver.ScreenSaver
	idleId     uint32
	idleLocker sync.Mutex
}

func newFullScreenWorkaround() (*fullScreenWorkaround, error) {
	XU, err := xgbutil.NewConn()
	if err != nil {
		return nil, err
	}

	ACTIVE_WINDOW, err := xprop.Atm(XU, "_NET_ACTIVE_WINDOW")
	if err != nil {
		return nil, err
	}

	ss, err := screensaver.NewScreenSaver("org.freedesktop.ScreenSaver",
		"/org/freedesktop/ScreenSaver")
	if err != nil {
		return nil, err
	}

	return &fullScreenWorkaround{
		xu: XU,
		targets: []string{
			"libflash",
			"chrome",
			"mplayer",
			"operaplugin",
			"soffice",
			"wpp",
			"evince",
			"vlc",
			"totem",
		},
		activeWindowAtom: ACTIVE_WINDOW,
		ss:               ss,
		idleId:           0,
	}, nil
}

func (wa *fullScreenWorkaround) detectTarget(w xproto.Window) {
	pid, _ := xprop.PropValNum(xprop.GetProperty(wa.xu, w, "_NET_WM_PID"))

	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		return
	}

	if wa.isFullScreen(w) {
		for _, target := range wa.targets {
			if strings.Contains(string(contents), target) {
				wa.inhibit(target, string(contents))
				return
			}
		}
	} else {
		if wa.ss != nil && wa.idleId != 0 {
			wa.idleLocker.Lock()
			logger.Debug("[detectTarget] try to inhibit:", wa.idleId)
			err := wa.ss.UnInhibit(wa.idleId)
			if err != nil {
				logger.Warning("[detectTarget] uninhibit failed:", wa.idleId, err)
			}
			wa.idleId = 0
			wa.idleLocker.Unlock()
		}
	}
}

func (wa *fullScreenWorkaround) inhibit(target, cmdline string) {
	wa.idleLocker.Lock()
	defer wa.idleLocker.Unlock()
	if wa.ss == nil {
		return
	}

	if wa.idleId != 0 {
		logger.Debug("[inhibit] has in inhibit mode:", wa.idleId)
		return
	}

	id, err := wa.ss.Inhibit("idle", "Fullscreen play video")
	if err != nil {
		logger.Warning("Inhibit 'idle' failed:", err)
		return
	}

	logger.Debug("[inhibit] success:", id)
	wa.idleId = id
}

func (wa *fullScreenWorkaround) isFullScreen(xid xproto.Window) bool {
	states, _ := ewmh.WmStateGet(wa.xu, xid)
	found := 0
	for _, s := range states {
		if s == "_NET_WM_STATE_FULLSCREEN" {
			found++
		}
		if s == "_NET_WM_STATE_FOCUSED" {
			found++
		}
	}
	if found == 2 {
		logger.Debug("HAHAH:::::", states)
	}
	return found == 2
}

func (wa *fullScreenWorkaround) start() {
	var runner func()
	runner = func() {
		w, _ := ewmh.ActiveWindowGet(wa.xu)
		wa.detectTarget(w)
		time.AfterFunc(time.Second*5, runner)
	}
	runner()

	root := xwindow.New(wa.xu, wa.xu.RootWin())
	root.Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		if wa.activeWindowAtom == ev.Atom {
			w, _ := ewmh.ActiveWindowGet(XU)
			wa.detectTarget(w)
		}
	}).Connect(wa.xu, root.Id)
	xevent.Main(wa.xu)
}
func (wa *fullScreenWorkaround) stop() {
	xevent.Quit(wa.xu)
	if wa.ss != nil {
		screensaver.DestroyScreenSaver(wa.ss)
		wa.ss = nil
	}
}
