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
	"os"
	"strings"
	"sync"

	"github.com/godbus/dbus"
	"github.com/linuxdeepin/go-dbus-factory/com.deepin.daemon.audio"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.miracle.wfd"
	"github.com/linuxdeepin/go-dbus-factory/org.freedesktop.miracle.wifi"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

type SinkInfo struct {
	Name      string
	P2PMac    string
	Interface string
	Connected bool
	Path      dbus.ObjectPath
	LinkPath  dbus.ObjectPath

	core   *wfd.Sink
	peer   *wifi.Peer
	locker sync.Mutex
}
type SinkInfos []*SinkInfo

func newSinkInfo(objPath dbus.ObjectPath) (*SinkInfo, error) {
	sysBus, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	core, err := wfd.NewSink(sysBus, objPath)
	if err != nil {
		return nil, err
	}

	peerPath, err := core.Peer().Get(0)
	if err != nil {
		return nil, err
	}

	peer, err := wifi.NewPeer(sysBus, peerPath)
	if err != nil {
		return nil, err
	}
	var sink = &SinkInfo{
		Path: objPath,
		core: core,
		peer: peer,
	}
	sink.update()
	return sink, nil
}

func destroySinkInfo(info *SinkInfo) {
	info.peer.RemoveHandler(proxy.RemoveAllHandlers)
}

func (sink *SinkInfo) connectSignal(m *Miracast) {
	sink.peer.InitSignalExt(m.sysSigLoop, true)
	err := sink.peer.Connected().ConnectChanged(func(hasValue bool, value bool) {
		if !hasValue {
			return
		}

		if value {
			sink.Connected = true
			m.emitSignalEvent(EventSinkConnected, sink.Path)
		} else {
			sink.Connected = false
			m.emitSignalEvent(EventSinkDisconnected, sink.Path)
		}
	})
	if err != nil {
		logger.Warning(err)
	}
}

func (sink *SinkInfo) update() {
	sink.locker.Lock()
	defer sink.locker.Unlock()
	if sink.core == nil || sink.peer == nil {
		return
	}
	sink.Name, _ = sink.peer.FriendlyName().Get(0)
	sink.P2PMac, _ = sink.peer.P2PMac().Get(0)
	sink.Interface, _ = sink.peer.Interface().Get(0)
	sink.Connected, _ = sink.peer.Connected().Get(0)
	sink.LinkPath, _ = sink.peer.Link().Get(0)
}

func (sink *SinkInfo) StartSession(x, y, w, h uint32) error {
	var (
		// format: 'x://:0'
		dpy       = "x://" + os.Getenv("DISPLAY")
		xauth     = os.Getenv("XAUTHORITY")
		audioSink = getAudioSink()
	)
	logger.Debug("[StartSession] args:", xauth, dpy, x, y, w, h, audioSink)
	sessionPath, err := sink.core.StartSession(0, xauth, dpy, x, y, w, h, audioSink)
	if err != nil {
		return err
	}
	logger.Debug("[StartSession] session path:", sessionPath)
	return nil
}

func (sink *SinkInfo) TeardownSession() error {
	p, err := sink.core.Session().Get(0)
	if err != nil {
		return err
	}
	if p == "/" {
		return fmt.Errorf("no session found")
	}

	sysBus, err := dbus.SystemBus()
	if err != nil {
		return err
	}
	session, err := wfd.NewSession(sysBus, p)
	if err != nil {
		return err
	}
	return session.Teardown(0)
}

func (sinks SinkInfos) Get(objPath dbus.ObjectPath) *SinkInfo {
	if !isSinkObjectPath(objPath) {
		return nil
	}
	for _, sink := range sinks {
		if sink.Path == objPath {
			return sink
		}
	}
	return nil
}

func (sinks SinkInfos) Remove(objPath dbus.ObjectPath) (SinkInfos, bool) {
	var (
		tmp    SinkInfos
		exists bool
	)
	for _, sink := range sinks {
		if sink.Path == objPath {
			exists = true
			continue
		}
		tmp = append(tmp, sink)
	}
	return tmp, exists
}

func isSinkObjectPath(objPath dbus.ObjectPath) bool {
	return strings.HasPrefix(string(objPath), sinkDBusPathPrefix)
}

func isSessionObjectPath(objPath dbus.ObjectPath) bool {
	return strings.HasPrefix(string(objPath), sessionDBusPathPrefix)
}

func getAudioSink() string {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return ""
	}

	obj := audio.NewAudio(sessionBus)

	defaultSinkPath, err := obj.DefaultSink().Get(0)
	if err != nil {
		return ""
	}

	sink, err := audio.NewSink(sessionBus, defaultSinkPath)
	if err != nil {
		return ""
	}

	name, err := sink.Name().Get(0)
	if err != nil {
		return ""
	}
	return name + ".monitor"
}
