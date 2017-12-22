/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package power

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/xevent"
	"github.com/BurntSushi/xgbutil/xprop"
	"github.com/BurntSushi/xgbutil/xwindow"
)

func init() {
	submoduleList = append(submoduleList, newFullScreenWorkaround)
}

type fullScreenWorkaround struct {
	manager          *Manager
	keyEventAtomList []xproto.Atom
	idleId           uint32
	enable           bool
	enableMutex      sync.Mutex
	ticker           *time.Ticker
	exit             chan struct{}
	targets          []string
}

func newFullScreenWorkaround(m *Manager) (string, submodule, error) {
	name := "FullscreenWorkaround"
	wa := &fullScreenWorkaround{
		manager: m,
		enable:  true,
	}
	var atomList []xproto.Atom

	list := []string{"_NET_ACTIVE_WINDOW", "_NET_CLIENT_LIST_STACKING"}
	xu := wa.manager.helper.xu
	for _, name := range list {
		atom, err := xprop.Atm(xu, name)
		if err != nil {
			return name, nil, err
		}
		atomList = append(atomList, atom)
	}
	wa.keyEventAtomList = atomList

	wa.targets = m.settings.GetStrv("fullscreen-workaround-app-list")

	return name, wa, nil
}

func (wa *fullScreenWorkaround) detect() {
	xu := wa.manager.helper.xu
	activeWin, _ := ewmh.ActiveWindowGet(xu)
	wa.enableMutex.Lock()

	if !wa.enable {
		//logger.Debug("disabled")
		wa.enableMutex.Unlock()
		return
	}

	// 先禁止处理信号，几秒后恢复处理
	wa.enable = false
	wa.enableMutex.Unlock()

	time.AfterFunc(1*time.Second, func() {
		if wa.isFullscreenFocused(activeWin) {
			wa.tryInhibit(activeWin)
		} else {
			//logger.Debug("Try uninhibit")
			wa.uninhibit()
		}
		wa.enableMutex.Lock()
		wa.enable = true
		wa.enableMutex.Unlock()
	})
}

func (wa *fullScreenWorkaround) uninhibit() {
	screenSaver := wa.manager.helper.ScreenSaver
	if screenSaver == nil {
		logger.Warning("screenSaver is nil")
		return
	}
	if wa.idleId != 0 {
		logger.Debug("* Uninhibit:", wa.idleId)
		err := screenSaver.UnInhibit(wa.idleId)
		if err != nil {
			logger.Warning("Uninhibit failed:", wa.idleId, err)
		}
		wa.idleId = 0
	} else {
		//logger.Debug("wa.idleId == 0")
	}
}

func (wa *fullScreenWorkaround) tryInhibit(activeWin xproto.Window) {
	logger.Debug("Try inhibit")
	if wa.idleId != 0 {
		logger.Debugf("Inhibit idleId %v != 0", wa.idleId)
		return
	}

	// get active window pid and cmdline
	xu := wa.manager.helper.xu
	pid, _ := xprop.PropValNum(xprop.GetProperty(xu, activeWin, "_NET_WM_PID"))
	contents, err := ioutil.ReadFile(fmt.Sprintf("/proc/%d/cmdline", pid))
	if err != nil {
		logger.Warningf("get pid %v cmdline failed: %v", pid, err)
		return
	}
	cmdline := string(contents)
	logger.Debugf("Focused | Fullscreen, Pid %v, Cmd line: %q", pid, cmdline)
	// match process cmdline with targets
	for _, target := range wa.targets {
		if strings.Contains(cmdline, target) {
			logger.Debugf("matchs %q", target)
			wa.inhibit()
			break
		}
	}
}

func (wa *fullScreenWorkaround) inhibit() {
	screenSaver := wa.manager.helper.ScreenSaver
	if screenSaver == nil {
		logger.Warning("screenSaver is nil")
		return
	}
	id, err := screenSaver.Inhibit("idle", "Fullscreen play video")
	if err != nil {
		logger.Warning("Inhibit 'idle' failed:", err)
		return
	}
	logger.Debug("* Inhibit success:", id)
	wa.idleId = id
}

func (wa *fullScreenWorkaround) isFullscreenFocused(xid xproto.Window) bool {
	xu := wa.manager.helper.xu
	states, _ := ewmh.WmStateGet(xu, xid)
	found := 0
	//logger.Debug("window states:", states)
	for _, s := range states {
		if s == "_NET_WM_STATE_FULLSCREEN" {
			found++
		}
		if s == "_NET_WM_STATE_FOCUSED" {
			found++
		}
	}
	return found == 2
}

func (wa *fullScreenWorkaround) Start() error {
	if !wa.manager.settings.GetBoolean(settingKeyFullscreenWorkaroundEnabled) {
		logger.Info("fullscreen workaround disabled")
		return nil
	}

	// 遇到 wpp 监听信号的方式会失效，所以增加一个定时轮询
	wa.ticker = time.NewTicker(5 * time.Second)
	wa.exit = make(chan struct{})
	go func() {
		for {
			select {
			case <-wa.ticker.C:
				//logger.Debug("Loop detect tick")
				wa.detect()
			case <-wa.exit:
				wa.exit = nil
				logger.Debug("exit loop detect")
				return
			}
		}
	}()

	xu := wa.manager.helper.xu
	root := xwindow.New(xu, xu.RootWin())
	root.Listen(xproto.EventMaskPropertyChange)
	xevent.PropertyNotifyFun(func(XU *xgbutil.XUtil, ev xevent.PropertyNotifyEvent) {
		//atomName, _ := xprop.AtomName(XU, ev.Atom)
		//logger.Debugf("signal %v %s", ev.Atom, atomName)
		for _, atom := range wa.keyEventAtomList {
			if ev.Atom == atom {
				wa.detect()
			}
		}
	}).Connect(xu, root.Id)
	go xevent.Main(xu)
	return nil
}

func (wa *fullScreenWorkaround) Destroy() {
	if wa.ticker != nil {
		wa.ticker.Stop()
		wa.ticker = nil
	}
	if wa.exit != nil {
		close(wa.exit)
	}
	xevent.Quit(wa.manager.helper.xu)
}
