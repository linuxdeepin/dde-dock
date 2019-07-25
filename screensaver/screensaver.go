/*
 * Copyright (C) 2014 ~ 2018 Deepin Technology Co., Ltd.
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

package screensaver

import (
	"errors"
	"strings"
	"sync"

	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/dpms"
	"github.com/linuxdeepin/go-x11-client/ext/screensaver"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/screensaver")

type inhibitor struct {
	sender dbus.Sender
	cookie uint32
	name   string
	reason string
}

type ScreenSaver struct {
	xConn      *x.Conn
	service    *dbusutil.Service
	sigLoop    *dbusutil.SignalLoop
	dbusDaemon *ofdbus.DBus

	blank        byte
	idleTime     uint32
	idleInterval uint32

	inhibitors map[uint32]inhibitor
	counter    uint32
	mu         sync.Mutex

	//Inhibit state, we need save the SetTimeout value,
	//so we can recover the correct state when enter UnInhibit state.
	lastVals *timeoutVals

	signals *struct {
		// Idle 定时器超时信号，当系统在给定时间内未被使用时发送
		IdleOn struct{}

		// Idle 超时时，如果设置了壁纸切换，则发送此信号
		CycleActive struct{}

		// Idle 超时后，如果系统被使用就发送此信号，重新开始 Idle 计时器
		IdleOff struct{}
	}

	methods *struct {
		Inhibit    func() `in:"name,reason" out:"cookie"`
		UnInhibit  func() `in:"cookie"`
		SetTimeout func() `in:"seconds,interval,blank"`
	}
}

type timeoutVals struct {
	seconds, interval uint32
	blank             bool
}

// 抑制 Idle 计时器，不再检测系统是否空闲，然后返回一个 id，用来取消此操作。
//
// name: 抑制 Idle 计时器的程序名称
//
// reason: 抑制原因
//
// ret0: 此次操作对应的 id，用来取消抑制
func (ss *ScreenSaver) Inhibit(sender dbus.Sender, name, reason string) (uint32,
	*dbus.Error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.counter++

	ss.inhibitors[ss.counter] = inhibitor{
		cookie: ss.counter,
		name:   name,
		reason: reason,
		sender: sender,
	}

	if len(ss.inhibitors) == 1 {
		ss.setTimeout(0, 0, false)
	}
	logger.Infof("sender %s %q want system enter inhibit, because: %q",
		sender, name, reason)

	return ss.counter, nil
}

// 模拟用户操作，让系统处于使用状态，重新开始 Idle 定时器
func (ss *ScreenSaver) SimulateUserActivity() *dbus.Error {
	err := x.ForceScreenSaverChecked(ss.xConn, x.ScreenSaverReset).Check(ss.xConn)
	return dbusutil.ToError(err)
}

// 根据 id 取消对应的抑制操作
func (ss *ScreenSaver) UnInhibit(sender dbus.Sender, cookie uint32) *dbus.Error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	inhibitor, ok := ss.inhibitors[cookie]
	if !ok {
		logger.Warning("invalid cookie", cookie)
		return dbusutil.ToError(errors.New("invalid cookie"))
	}

	if inhibitor.sender != sender {
		return dbusutil.ToError(errors.New("sender not match"))
	}

	logger.Infof("%q no need inhibit.", inhibitor.name)
	ss.unInhibit(cookie)

	return nil
}

func (ss *ScreenSaver) unInhibit(cookie uint32) {
	delete(ss.inhibitors, cookie)
	if len(ss.inhibitors) == 0 {
		logger.Info("Enter un-inhibit state")
		if ss.lastVals != nil {
			logger.Info("recover from ", ss.lastVals)
			ss.setTimeout(ss.lastVals.seconds, ss.lastVals.interval, ss.lastVals.blank)
			ss.lastVals = nil
		} else {
			ss.setTimeout(ss.idleTime, ss.idleInterval, ss.blank == 1)
		}
	}
}

// 设置 Idle 的定时器超时时间
//
// seconds: 超时时间，以秒为单位
//
// interval: 屏保模式下，背景更换的间隔时间
//
// blank: 是否黑屏，此参数暂时无效
func (ss *ScreenSaver) SetTimeout(seconds, interval uint32, blank bool) *dbus.Error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if len(ss.inhibitors) > 0 {
		ss.lastVals = &timeoutVals{seconds, interval, blank}
		logger.Info("Current is inhibit state, the value", ss.lastVals, "will apply when in unhibit state")
	} else {
		ss.setTimeout(seconds, interval, blank)

		ss.idleTime = seconds
		ss.idleInterval = interval
		if blank {
			ss.blank = 1
		} else {
			ss.blank = 0
		}
	}
	return nil
}

func (ss *ScreenSaver) setTimeout(seconds, interval uint32, blank bool) {
	if blank {
		ss.blank = x.BlankingPreferred
	} else {
		ss.blank = x.BlankingNotPreferred
	}

	err := x.SetScreenSaverChecked(ss.xConn, int16(seconds), int16(interval), ss.blank,
		x.ExposuresNotAllowed).Check(ss.xConn)
	if err != nil {
		logger.Warning(err)
	}

	err = dpms.SetTimeoutsChecked(ss.xConn, 0, 0,
		0).Check(ss.xConn)
	if err != nil {
		logger.Warning(err)
	}

	logger.Info("SetTimeout to ", seconds, interval, blank)
}

const (
	dbusServiceName = "org.freedesktop.ScreenSaver"
	dbusPath        = "/org/freedesktop/ScreenSaver"
	dbusInterface   = dbusServiceName
)

func (*ScreenSaver) GetInterfaceName() string {
	return dbusInterface
}

func (ss *ScreenSaver) destroy() {
	ss.sigLoop.Stop()
	ss.dbusDaemon.RemoveHandler(proxy.RemoveAllHandlers)
	if ss.xConn != nil {
		ss.xConn.Close()
	}
}

func newScreenSaver(service *dbusutil.Service) (*ScreenSaver, error) {
	s := &ScreenSaver{
		service:    service,
		inhibitors: make(map[uint32]inhibitor),
	}

	sessionBus := service.Conn()
	s.dbusDaemon = ofdbus.NewDBus(sessionBus)
	s.sigLoop = dbusutil.NewSignalLoop(sessionBus, 10)

	var err error
	s.xConn, err = x.NewConn()
	if err != nil {
		return nil, err
	}

	// query screensaver ext version
	ssVersion, err := screensaver.QueryVersion(s.xConn, screensaver.MajorVersion, screensaver.MinorVersion).Reply(s.xConn)
	if err != nil {
		return nil, err
	}
	logger.Debugf("screensaver ext version %d.%d", ssVersion.ServerMajorVersion,
		ssVersion.ServerMinorVersion)

	// query dpms ext version
	dpmsVersion, err := dpms.GetVersion(s.xConn, 1, 1).Reply(s.xConn)
	if err != nil {
		logger.Warning("failed to get dpms ext version:", err)
	} else {
		logger.Debugf("dpms ext version %d.%d", dpmsVersion.ServerMajorVersion,
			dpmsVersion.ServerMinorVersion)
	}

	root := s.xConn.GetDefaultScreen().Root
	err = screensaver.SelectInputChecked(s.xConn, x.Drawable(root), screensaver.EventNotifyMask|
		screensaver.EventCycleMask).Check(s.xConn)
	if err != nil {
		logger.Warning(err)
	}

	s.listenDBusNameOwnerChanged()
	s.sigLoop.Start()
	go s.loop()
	return s, nil
}

func (ss *ScreenSaver) listenDBusNameOwnerChanged() {
	ss.dbusDaemon.InitSignalExt(ss.sigLoop, true)
	ss.dbusDaemon.ConnectNameOwnerChanged(func(name string, oldOwner string, newOwner string) {
		if strings.HasPrefix(name, ":") &&
			name == oldOwner && newOwner == "" {

			ss.mu.Lock()
			for cookie, inhibitor := range ss.inhibitors {
				if string(inhibitor.sender) == name {
					logger.Infof("app %s %q disconnect from DBus",
						name, inhibitor.name)
					ss.unInhibit(cookie)
					break
				}
			}
			ss.mu.Unlock()
		}
	})
}

var _ssaver *ScreenSaver

func (ss *ScreenSaver) loop() {
	ssExtData := ss.xConn.GetExtensionData(screensaver.Ext())

	eventChan := make(chan x.GenericEvent, 10)
	ss.xConn.AddEventChan(eventChan)

	for ev := range eventChan {
		switch ev.GetEventCode() {
		case screensaver.NotifyEventCode + ssExtData.FirstEvent:
			event, _ := screensaver.NewNotifyEvent(ev)
			s := ss.service
			switch event.State {
			case screensaver.StateCycle:
				s.Emit(ss, "CycleActive")
			case screensaver.StateOn:
				s.Emit(ss, "IdleOn")
			case screensaver.StateOff:
				s.Emit(ss, "IdleOff")
			}
		}
	}
}
