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
	"os"
	"time"
)

func (m *Miracast) handleLinkManaged(link *LinkInfo) {
	link.core.Managed.ConnectChanged(func() {
		if link.core == nil {
			return
		}
		link.update()
		logger.Debugf("Link '%s' managed changed: %v", link.Path, link.Managed)
		v, ok := m.managingLinks[link.Path]
		if !ok {
			return
		}

		if link.Managed {
			link.core.WfdSubelements.Set("000600001c440036")
			// default using 'USER' as link name
			link.SetName(os.Getenv("USER"))
			// TODO: wait WfdSubelements changed
			m.service.Emit(m, "Event", EventLinkManaged, link.Path)
		} else {
			link.core.WfdSubelements.Set("")
			m.service.Emit(m, "Event", EventLinkUnmanaged, link.Path)
			err := m.enableWirelessManaged(link.MacAddress, true)
			if err != nil {
				logger.Error("Enable nm device failed:", link.MacAddress, err)
			}
		}

		logger.Debugf("[handleLinkChanged] link '%s' managed: %v, should be: %v",
			link.Path, link.Managed, v)
		if v == link.Managed {
			delete(m.managingLinks, link.Path)
		}
	})
	link.core.P2PScanning.ConnectChanged(func() {
		if link.core == nil {
			return
		}
		link.update()
	})
}

func (m *Miracast) handleSinkConnected(sink *SinkInfo, x, y, w, h uint32) {
	var hasConnected = false
	sink.peer.Connected.ConnectChanged(func() {
		if sink.peer == nil {
			return
		}
		sink.update()
		logger.Debugf("Sink(%v)'s peer(%v) connected changed: %v", sink.Path, sink.peer.Path, sink.Connected)
		defer delete(m.connectingSinks, sink.Path)
		if !sink.Connected {
			logger.Debugf("Sink '%s' was disconnect", sink.Path)
			m.service.Emit(m, "Event", EventSinkDisconnected, sink.Path)
			return
		}

		hasConnected = true
		m.service.Emit(m, "Event", EventSinkConnected, sink.Path)
		m.doConnect(sink, x, y, w, h)
	})

	// timeout
	time.AfterFunc(time.Second*60, func() {
		if sink.core == nil || hasConnected {
			return
		}
		m.disconnect(sink.Path)
		m.service.Emit(m, "Event", EventSinkConnectedFailed, sink.Path)
	})
}

func (m *Miracast) doConnect(sink *SinkInfo, x, y, w, h uint32) {
	m.sinkLocker.Lock()
	defer m.sinkLocker.Unlock()
	if sink.core.Session.Get() != "/" {
		logger.Debug("Sink session had connected:", sink.Path, sink.core.Session.Get())
		return
	}
	err := sink.StartSession(x, y, w, h)
	if err != nil {
		logger.Error("Failed to start session:", sink.Path, err)
	}
}
