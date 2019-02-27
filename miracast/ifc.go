/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
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

package miracast

import (
	"fmt"

	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	EventLinkManaged uint8 = iota + 1
	EventLinkUnmanaged
	EventSinkConnected
	EventSinkConnectedFailed
	EventSinkDisconnected
)

func (m *Miracast) ListLinks() (LinkInfos, *dbus.Error) {
	logger.Debug("call ListLinks")
	m.init()
	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	return m.links, nil
}

func (m *Miracast) ListSinks() (SinkInfos, *dbus.Error) {
	logger.Debug("call ListSinks")
	m.init()
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	return m.sinks, nil
}

func (m *Miracast) Enable(linkPath dbus.ObjectPath, enabled bool) *dbus.Error {
	logger.Debug("call Enable", linkPath, enabled)
	m.init()
	err := m.enable(linkPath, enabled)
	return dbusutil.ToError(err)
}

func (m *Miracast) enable(linkPath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(linkPath) {
		return fmt.Errorf("invalid link objPath: %v", linkPath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(linkPath)
	if link == nil {
		logger.Warning("not found the link:", linkPath)
		return fmt.Errorf("not found the link: %v", linkPath)
	}

	if enabled {
		err := m.enableWirelessManaged(link.interfaceName, false)
		if err != nil {
			logger.Warning("failed to disable manage wireless device:", err)
			return err
		}
	} else {
		err := m.enableWirelessManaged(link.interfaceName, true)
		if err != nil {
			logger.Warning("failed to enable manage wireless device:", err)
			return err
		}
	}

	managed, err := link.core.Managed().Get(0)
	if err != nil {
		return err
	}
	if managed == enabled {
		logger.Debug("work right, nothing to do:", enabled)
		return nil
	}

	link.myManaged = enabled
	return link.EnableManaged(enabled)
}

func (m *Miracast) SetLinkName(linkPath dbus.ObjectPath, name string) *dbus.Error {
	logger.Debug("call SetLinkName", linkPath, name)
	m.init()
	err := m.setLinkName(linkPath, name)
	return dbusutil.ToError(err)
}

func (m *Miracast) setLinkName(linkPath dbus.ObjectPath, name string) error {
	if !isLinkObjectPath(linkPath) {
		return fmt.Errorf("invalid link objPath: %v", linkPath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(linkPath)
	if link == nil {
		logger.Warning("not found the link:", linkPath)
		return fmt.Errorf("not found the link: %v", linkPath)
	}

	link.SetName(name)
	return nil
}

func (m *Miracast) Scanning(linkPath dbus.ObjectPath, enabled bool) *dbus.Error {
	logger.Debug("call Scanning", linkPath, enabled)
	m.init()
	err := m.scanning(linkPath, enabled)
	return dbusutil.ToError(err)
}

func (m *Miracast) scanning(linkPath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(linkPath) {
		return fmt.Errorf("invalid link objPath: %v", linkPath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(linkPath)
	if link == nil {
		logger.Warning("not found the link:", linkPath)
		return fmt.Errorf("not found the link: %v", linkPath)
	}

	if !link.myManaged {
		logger.Debug("not allow scan")
		return nil
	}

	logger.Debug("manage link", linkPath)
	err := link.EnableManaged(true)
	if err != nil {
		logger.Warning(err)
	}

	err = link.waitManaged(true)
	if err != nil {
		logger.Warning(err)
		return err
	}

	link.EnableP2PScanning(enabled)
	link.update()
	return nil
}

func (m *Miracast) Connect(sinkPath dbus.ObjectPath, x, y, w, h uint32) *dbus.Error {
	logger.Debug("call Connect", sinkPath, x, y, w, h)
	m.init()
	err := m.connect(sinkPath, x, y, w, h)
	return dbusutil.ToError(err)
}

func (m *Miracast) connect(sinkPath dbus.ObjectPath, x, y, w, h uint32) error {
	if !isSinkObjectPath(sinkPath) {
		return fmt.Errorf("invalid sink objPath: %v", sinkPath)
	}

	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	sink := m.sinks.Get(sinkPath)
	if sink == nil {
		logger.Warning("not found sink", sinkPath)
		return fmt.Errorf("not found sink %v", sinkPath)
	}

	linkPath, err := sink.peer.Link().Get(0)
	if err != nil {
		logger.Warning(err)
	}

	m.linkLocker.Lock()
	link := m.links.Get(linkPath)
	m.linkLocker.Unlock()
	if link == nil {
		logger.Warning("not found link", linkPath)
		return fmt.Errorf("not found link %v", linkPath)
	}

	err = link.EnableManaged(true)
	if err != nil {
		logger.Warning(err)
		return err
	}
	err = link.waitManaged(true)
	if err != nil {
		logger.Warning(err)
		return err
	}
	link.EnableP2PScanning(false)
	link.ConfigureForManaged()

	connected, err := sink.peer.Connected().Get(0)
	if connected {
		logger.Debug("Has connected, start session")
		m.startSession(sink, x, y, w, h)
		return nil
	}

	err = sink.peer.Connect(0, "auto", "")
	if err != nil {
		logger.Error("Failed to connect sink:", err)
		return err
	}

	go func() {
		err := waitPeerConnected(sink.peer, true)
		if err == nil {
			m.startSession(sink, x, y, w, h)
		} else {
			logger.Warning(err)
			m.emitSignalEvent(EventSinkConnectedFailed, sink.Path)
		}
	}()
	return nil
}

func (m *Miracast) Disconnect(sinkPath dbus.ObjectPath) *dbus.Error {
	logger.Debug("call Disconnect", sinkPath)
	m.init()
	err := m.disconnectSink(sinkPath)
	return dbusutil.ToError(err)
}
