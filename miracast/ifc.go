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
	"time"

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
	m.init()
	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	return m.links, nil
}

func (m *Miracast) ListSinks() (SinkInfos, *dbus.Error) {
	m.init()
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	return m.sinks, nil
}

func (m *Miracast) Enable(linkPath dbus.ObjectPath, enabled bool) *dbus.Error {
	err := m.enable(linkPath, enabled)
	return dbusutil.ToError(err)
}

func (m *Miracast) enable(linkPath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(linkPath) {
		return fmt.Errorf("Invalid link dpath: %v", linkPath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(linkPath)
	if link == nil {
		logger.Warning("Not found the link:", linkPath)
		return fmt.Errorf("Not found the link: %v", linkPath)
	}

	if link.core.Managed.Get() == enabled {
		logger.Debug("Work right, nothing to do:", enabled)
		return nil
	}

	if v, ok := m.managingLinks[linkPath]; ok && v == enabled {
		logger.Debug("Link's managed '%s' has been setting to ", enabled)
		return nil
	}
	m.managingLinks[linkPath] = enabled

	if enabled {
		err := m.enableWirelessManaged(link.MacAddress, false)
		if err != nil {
			delete(m.managingLinks, linkPath)
			logger.Error("Failed to disable manage wireless device:", err)
			return err
		}
		time.Sleep(time.Millisecond * 500)
	}

	m.handleLinkManaged(link)
	return link.EnableManaged(enabled)
}

func (m *Miracast) SetLinkName(linkPath dbus.ObjectPath, name string) *dbus.Error {
	err := m.setLinkName(linkPath, name)
	return dbusutil.ToError(err)
}

func (m *Miracast) setLinkName(linkPath dbus.ObjectPath, name string) error {
	if !isLinkObjectPath(linkPath) {
		return fmt.Errorf("Invalid link dpath: %v", linkPath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(linkPath)
	if link == nil {
		logger.Warning("Not found the link:", linkPath)
		return fmt.Errorf("Not found the link: %v", linkPath)
	}

	link.SetName(name)
	return nil
}

func (m *Miracast) Scanning(linkPath dbus.ObjectPath, enabled bool) *dbus.Error {
	err := m.scanning(linkPath, enabled)
	return dbusutil.ToError(err)
}

func (m *Miracast) scanning(linkPath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(linkPath) {
		return fmt.Errorf("Invalid link dpath: %v", linkPath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(linkPath)
	if link == nil {
		logger.Warning("Not found the link:", linkPath)
		return fmt.Errorf("Not found the link: %v", linkPath)
	}

	link.EnableP2PScanning(enabled)
	// TODO: wait P2PScanning changed
	link.update()
	return nil
}

func (m *Miracast) Connect(sinkPath dbus.ObjectPath, x, y, w, h uint32) *dbus.Error {
	err := m.connect(sinkPath, x, y, w, h)
	return dbusutil.ToError(err)
}

func (m *Miracast) connect(sinkPath dbus.ObjectPath, x, y, w, h uint32) error {
	if !isSinkObjectPath(sinkPath) {
		return fmt.Errorf("Invalid sink dpath: %v", sinkPath)
	}

	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	sink := m.sinks.Get(sinkPath)
	if sink == nil {
		logger.Warning("Not found the sink:", sinkPath)
		return fmt.Errorf("Not found the sink: %v", sinkPath)
	}

	if sink.peer.Connected.Get() {
		logger.Debug("Has connected, start session")
		m.doConnect(sink, x, y, w, h)
		return nil
	}

	if v, ok := m.connectingSinks[sinkPath]; ok && v {
		logger.Debug("[ConnectSink] sink has connecting:", sinkPath)
		return nil
	}
	m.connectingSinks[sinkPath] = true

	m.handleSinkConnected(sink, x, y, w, h)
	err := sink.peer.Connect("auto", "")
	if err != nil {
		delete(m.connectingSinks, sinkPath)
		logger.Error("Failed to connect sink:", err)
		return err
	}

	return nil
}

func (m *Miracast) Disconnect(sinkPath dbus.ObjectPath) *dbus.Error {
	err := m.disconnect(sinkPath)
	return dbusutil.ToError(err)
}

func (m *Miracast) disconnect(sinkPath dbus.ObjectPath) error {
	if !isSinkObjectPath(sinkPath) {
		return fmt.Errorf("Invalid sink dpath: %v", sinkPath)
	}
	return m.disconnectSink(sinkPath)
}
