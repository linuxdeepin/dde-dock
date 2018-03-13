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
	"sync"

	"github.com/BurntSushi/xgb/dpms"
	"github.com/BurntSushi/xgb/screensaver"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/screensaver")

type inhibitor struct {
	cookie uint32
	name   string
	reason string
}
type ScreenSaver struct {
	xu      *xgbutil.XUtil
	service *dbusutil.Service

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
func (ss *ScreenSaver) Inhibit(name, reason string) (uint32, *dbus.Error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ss.counter++

	ss.inhibitors[ss.counter] = inhibitor{ss.counter, name, reason}

	if len(ss.inhibitors) == 1 {
		ss.setTimeout(0, 0, false)
	}
	logger.Infof("\"%s\" want system enter inhibit, because: \"%s\"", name, reason)

	return ss.counter, nil
}

// 模拟用户操作，让系统处于使用状态，重新开始 Idle 定时器
func (ss *ScreenSaver) SimulateUserActivity() *dbus.Error {
	xproto.ForceScreenSaver(ss.xu.Conn(), 0)
	return nil
}

// 根据 id 取消对应的抑制操作
func (ss *ScreenSaver) UnInhibit(cookie uint32) *dbus.Error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	inhibitor, ok := ss.inhibitors[cookie]
	if !ok {
		logger.Warning("invalid cookie", cookie)
		return dbusutil.ToError(errors.New("invalid cookie"))
	}

	logger.Infof("%q no need inhibit.", inhibitor.name)

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

	return nil
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
		ss.blank = 1
	} else {
		ss.blank = 0
	}
	xproto.SetScreenSaver(ss.xu.Conn(), int16(seconds), int16(interval), ss.blank, 0)
	dpms.SetTimeouts(ss.xu.Conn(), 0, 0, 0)
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
	ss.service.StopExport(ss)
}

func newScreenSaver(service *dbusutil.Service) *ScreenSaver {
	s := &ScreenSaver{
		service:    service,
		inhibitors: make(map[uint32]inhibitor),
	}

	s.xu, _ = xgbutil.NewConn()
	screensaver.Init(s.xu.Conn())
	screensaver.QueryVersion(s.xu.Conn(), 1, 0)
	screensaver.SelectInput(s.xu.Conn(), xproto.Drawable(s.xu.RootWin()), screensaver.EventNotifyMask|screensaver.EventCycleMask)
	dpms.Init(s.xu.Conn())

	go s.loop()
	return s
}

var _ssaver *ScreenSaver

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("screensaver", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	service := loader.GetService()

	has, err := service.NameHasOwner(dbusServiceName)
	if err != nil {
		return err
	}
	if has {
		logger.Warning("ScreenSaver has been register, exit...")
		return nil
	}

	if _ssaver != nil {
		return nil
	}

	_ssaver = newScreenSaver(service)

	err = service.Export(dbusPath, _ssaver)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		_ssaver.destroy()
		_ssaver = nil
		return err
	}

	return nil
}

func (d *Daemon) Stop() error {
	if _ssaver == nil {
		return nil
	}

	_ssaver.destroy()
	_ssaver = nil
	return nil
}

func (ss *ScreenSaver) loop() {
	s := ss.service
	for {
		e, err := ss.xu.Conn().WaitForEvent()
		if err != nil {
			continue
		}
		switch ee := e.(type) {
		case screensaver.NotifyEvent:
			switch ee.State {
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
