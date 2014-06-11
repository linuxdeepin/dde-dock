/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package bluetooth

import (
	"dbus/org/bluez"
	"dlib/dbus"
)

type adapter struct {
	bluezAdapter *bluez.Adapter1

	Path                dbus.ObjectPath
	Alias               string
	Powered             bool
	Discoverable        bool
	DiscoverableTimeout uint32
}

func (b *Bluetooth) newAdapter(apath dbus.ObjectPath) (a *adapter) {
	a = &adapter{Path: apath}
	a.bluezAdapter, _ = bluezNewAdapter(apath)
	a.Alias = a.bluezAdapter.Alias.Get()
	a.Powered = a.bluezAdapter.Powered.Get()
	a.Discoverable = a.bluezAdapter.Discoverable.Get()
	a.DiscoverableTimeout = a.bluezAdapter.DiscoverableTimeout.Get()

	// TODO connect properties
	a.bluezAdapter.Alias.ConnectChanged(func() {
		a.Alias = a.bluezAdapter.Alias.Get()
		b.updatePropAdapters()
	})
	a.bluezAdapter.Powered.ConnectChanged(func() {
		a.Powered = a.bluezAdapter.Powered.Get()
		b.updatePropAdapters()
	})
	a.bluezAdapter.Discoverable.ConnectChanged(func() {
		a.Discoverable = a.bluezAdapter.Discoverable.Get()
		b.updatePropAdapters()
	})
	a.bluezAdapter.DiscoverableTimeout.ConnectChanged(func() {
		a.DiscoverableTimeout = a.bluezAdapter.DiscoverableTimeout.Get()
		b.updatePropAdapters()
	})

	return
}

func (b *Bluetooth) addAdapter(apath dbus.ObjectPath) {
	if b.isAdapterExists(apath) {
		logger.Warning("repeat add adapter:", apath)
		return
	}
	a := b.newAdapter(apath)
	b.adapters = append(b.adapters, a)
	b.updatePropAdapters()
}

func (b *Bluetooth) removeAdapter(apath dbus.ObjectPath) {
	i := b.getAdapterIndex(apath)
	if i < 0 {
		logger.Warning("repeat remove adapter:", apath)
		return
	}
	copy(b.adapters[i:], b.adapters[i+1:])
	b.adapters[len(b.adapters)-1] = nil
	b.adapters = b.adapters[:len(b.adapters)-1]
	b.updatePropAdapters()
}

func (b *Bluetooth) isAdapterExists(apath dbus.ObjectPath) bool {
	if b.getAdapterIndex(apath) >= 0 {
		return true
	}
	return false
}

func (b *Bluetooth) getAdapterIndex(apath dbus.ObjectPath) int {
	for i, a := range b.adapters {
		if a.Path == apath {
			return i
		}
	}
	return -1
}
