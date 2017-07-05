/**
 * Copyright (C) 2017 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"github.com/msteinert/pam"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"runtime/debug"
	"sync"
)

const (
	PromptQuestion uint32 = iota + 1
	PromptSecret
	ErrorMsg
	TextInfo
	Failure
	Successed
)

type Manager struct {
	authLocker    sync.Mutex
	authUserTable map[string]chan string // 'pid+user': 'passwd'

	// evenType, pid, username, message
	Event func(uint32, uint32, string, string)
	// username
	UserChanged func(string)
}

const (
	dbusDest = "com.deepin.dde.LockService"
	dbusPath = "/com/deepin/dde/LockService"
	dbusIFC  = "com.deepin.dde.LockService"
)

var _m *Manager

func main() {
	if !lib.UniqueOnSystem(dbusDest) {
		fmt.Println("The lock service has been running")
		return
	}

	_m = newManager()
	err := dbus.InstallOnSystem(_m)
	if err != nil {
		fmt.Println("Failed to install dbus:", err)
		return
	}
	dbus.DealWithUnhandledMessage()

	err = dbus.Wait()
	if err != nil {
		fmt.Println("Failed to wait dbus:", err)
	}
	return
}

func (*Manager) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       dbusDest,
		ObjectPath: dbusPath,
		Interface:  dbusIFC,
	}
}

func newManager() *Manager {
	var m = Manager{
		authUserTable: make(map[string]chan string),
	}

	return &m
}

func (m *Manager) CurrentUser() (string, error) {
	return getGreeterUser(greeterUserConfig)
}

func (m *Manager) IsLiveCD(username string) bool {
	return isInLiveCD(username)
}

func (m *Manager) SwitchToUser(username string) error {
	current, _ := getGreeterUser(greeterUserConfig)
	if current == username {
		return nil
	}

	err := setGreeterUser(greeterUserConfig, username)
	if err != nil {
		return err
	}
	dbus.Emit(m, "UserChanged", username)
	return nil
}

func (m *Manager) AuthenticateUser(dmsg dbus.DMessage, username string) error {
	if username == "" {
		return fmt.Errorf("No user to authenticate")
	}

	m.authLocker.Lock()
	pid := dmsg.GetSenderPID()
	id := fmt.Sprintf("%d%s", pid, username)
	_, ok := m.authUserTable[id]
	if ok {
		m.authLocker.Unlock()
		return nil
	}

	m.authUserTable[id] = make(chan string)
	m.authLocker.Unlock()
	go m.doAuthenticate(username, "", pid, true)
	return nil
}

func (m *Manager) UnlockCheck(dmsg dbus.DMessage, username, passwd string) error {
	if username == "" {
		return fmt.Errorf("No user to authenticate")
	}

	m.authLocker.Lock()
	pid := dmsg.GetSenderPID()
	id := fmt.Sprintf("%d%s", pid, username)
	v, ok := m.authUserTable[id]
	m.authLocker.Unlock()
	if ok && v != nil {
		fmt.Println("-------In authenticating:", id)
		// in authenticate
		v <- passwd
		return nil
	}

	go m.doAuthenticate(username, passwd, pid, false)
	return nil
}

func (m *Manager) doAuthenticate(username, passwd string, pid uint32, wait bool) {
	handler, err := pam.StartFunc("lightdm", username, func(style pam.Style, msg string) (string, error) {
		switch style {
		// case pam.PromptEchoOn:
		// 	if msg != "" {
		// 		fmt.Println("Echo on:", msg)
		// 		m.sendEvent(PromptQuestion, pid, username, msg)
		// 	}
		// 	// TODO: read data from input
		// 	return "", nil
		case pam.PromptEchoOff, pam.PromptEchoOn:
			if msg != "" {
				if style == pam.PromptEchoOff {
					fmt.Println("Echo off:", msg)
					m.sendEvent(PromptSecret, pid, username, msg)
				} else {
					fmt.Println("Echo on:", msg)
					m.sendEvent(PromptQuestion, pid, username, msg)
				}
			}

			if !wait {
				return passwd, nil
			}

			m.authLocker.Lock()
			id := fmt.Sprintf("%d%s", pid, username)
			v, ok := m.authUserTable[id]
			m.authLocker.Unlock()
			if !ok || v == nil {
				return "", fmt.Errorf("No passwd channel found for %s", username)
			}
			fmt.Println("-----Join select:", id)
			select {
			case tmp := <-v:
				return tmp, nil
			}
		case pam.ErrorMsg:
			if msg != "" {
				fmt.Println("ShowError:", msg)
				m.sendEvent(ErrorMsg, pid, username, msg)
			}
			return "", nil
		case pam.TextInfo:
			if msg != "" {
				fmt.Println("Text info:", msg)
				m.sendEvent(TextInfo, pid, username, msg)
			}
			return "", nil
		}
		return "", fmt.Errorf("Unexpected style: %v", style)
	})
	if err != nil {
		fmt.Println("Failed to start pam:", err)
		m.sendEvent(Failure, pid, username, err.Error())
		return
	}

	err = handler.Authenticate(pam.DisallowNullAuthtok)
	m.authLocker.Lock()
	id := fmt.Sprintf("%d%s", pid, username)
	v, ok := m.authUserTable[id]
	if ok {
		if v != nil {
			close(v)
		}
		delete(m.authUserTable, id)
	}
	m.authLocker.Unlock()
	if err != nil {
		fmt.Println("Failed to authenticate:", err)
		m.sendEvent(Failure, pid, username, err.Error())
		debug.FreeOSMemory()
		return
	}

	fmt.Println("-------Authenticate successed")
	m.sendEvent(Successed, pid, username, "Authenticated")
	debug.FreeOSMemory()
}

func (m *Manager) sendEvent(ty, pid uint32, username, msg string) {
	err := dbus.Emit(m, "Event", ty, pid, username, msg)
	if err != nil {
		fmt.Println("Failed to emit event:", ty, pid, username, msg)
	}
}
