/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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
	"pkg.deepin.io/lib/dbus"
	"time"
)

const (
	EventLinkManaged uint8 = iota + 1
	EventLinkUnmanaged
	EventSinkConnected
	EventSinkConnectedFailed
	EventSinkDisconnected
)

func (m *Miracast) ListLinks() LinkInfos {
	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	return m.links
}

func (m *Miracast) ListSinks() SinkInfos {
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	return m.sinks
}

func (m *Miracast) Enable(dpath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(dpath) {
		return fmt.Errorf("Invalid link dpath: %v", dpath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(dpath)
	if link == nil {
		logger.Warning("Not found the link:", dpath)
		return fmt.Errorf("Not found the link: %v", dpath)
	}

	if link.core.Managed.Get() == enabled {
		logger.Debug("Work right, nothing to do:", enabled)
		return nil
	}

	if v, ok := m.managingLinks[dpath]; ok && v == enabled {
		logger.Debug("Link's managed '%s' has been setting to ", enabled)
		return nil
	}
	m.managingLinks[dpath] = enabled

	if enabled {
		err := m.enableWirelessManaged(link.MacAddress, false)
		if err != nil {
			delete(m.managingLinks, dpath)
			logger.Error("Failed to disable manage wireless device:", err)
			return err
		}
		time.Sleep(time.Millisecond * 500)
	}

	m.handleLinkManaged(link)
	return link.EnableManaged(enabled)
}

func (m *Miracast) SetLinkName(dpath dbus.ObjectPath, name string) error {
	if !isLinkObjectPath(dpath) {
		return fmt.Errorf("Invalid link dpath: %v", dpath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(dpath)
	if link == nil {
		logger.Warning("Not found the link:", dpath)
		return fmt.Errorf("Not found the link: %v", dpath)
	}

	link.SetName(name)
	return nil
}

func (m *Miracast) Scanning(dpath dbus.ObjectPath, enabled bool) error {
	if !isLinkObjectPath(dpath) {
		return fmt.Errorf("Invalid link dpath: %v", dpath)
	}

	m.linkLocker.Lock()
	defer m.linkLocker.Unlock()
	link := m.links.Get(dpath)
	if link == nil {
		logger.Warning("Not found the link:", dpath)
		return fmt.Errorf("Not found the link: %v", dpath)
	}

	link.EnableP2PScanning(enabled)
	// TODO: wait P2PScanning changed
	link.update()
	return nil
}

func (m *Miracast) Connect(dpath dbus.ObjectPath, x, y, w, h uint32) error {
	if !isSinkObjectPath(dpath) {
		return fmt.Errorf("Invalid sink dpath: %v", dpath)
	}

	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	sink := m.sinks.Get(dpath)
	if sink == nil {
		logger.Warning("Not found the sink:", dpath)
		return fmt.Errorf("Not found the sink: %v", dpath)
	}

	if sink.peer.Connected.Get() {
		logger.Debug("Has connected, start session")
		m.doConnect(sink, x, y, w, h)
		return nil
	}

	if v, ok := m.connectingSinks[dpath]; ok && v {
		logger.Debug("[ConnectSink] sink has connecting:", dpath)
		return nil
	}
	m.connectingSinks[dpath] = true

	m.handleSinkConnected(sink, x, y, w, h)
	err := sink.peer.Connect("auto", "")
	if err != nil {
		delete(m.connectingSinks, dpath)
		logger.Error("Failed to connect sink:", err)
		return err
	}

	return nil
}

func (m *Miracast) Disconnect(dpath dbus.ObjectPath) error {
	if !isSinkObjectPath(dpath) {
		return fmt.Errorf("Invalid sink dpath: %v", dpath)
	}
	return m.disconnectSink(dpath)
}
