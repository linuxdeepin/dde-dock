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

package trayicon

import (
	x "github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/composite"
	"github.com/linuxdeepin/go-x11-client/ext/damage"
	"github.com/linuxdeepin/go-x11-client/util/atom"
	"github.com/linuxdeepin/go-x11-client/util/wm/ewmh"
	"github.com/linuxdeepin/go-x11-client/util/wm/icccm"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

func init() {
	loader.Register(NewDaemon(logger))
}

var (
	logger = log.NewLogger("daemon/trayicon")

	XConn     *x.Conn
	ewmhConn  *ewmh.Conn
	icccmConn *icccm.Conn

	XA_NET_SYSTEM_TRAY_S0         x.Atom
	XA_NET_SYSTEM_TRAY_OPCODE     x.Atom
	XA_NET_SYSTEM_TRAY_VISUAL     x.Atom
	XA_NET_SYSTEM_TRAY_ORIENTAION x.Atom
	XA_MANAGER                    x.Atom
)

func initX() {
	damage.QueryVersion(XConn, damage.MajorVersion, damage.MinorVersion).Reply(XConn)
	composite.QueryVersion(XConn, composite.MajorVersion, composite.MinorVersion).Reply(XConn)

	XA_NET_SYSTEM_TRAY_S0, _ = atom.GetVal(XConn, "_NET_SYSTEM_TRAY_S0")
	XA_NET_SYSTEM_TRAY_OPCODE, _ = atom.GetVal(XConn, "_NET_SYSTEM_TRAY_OPCODE")
	XA_NET_SYSTEM_TRAY_VISUAL, _ = atom.GetVal(XConn, "_NET_SYSTEM_TRAY_VISUAL")
	XA_NET_SYSTEM_TRAY_ORIENTAION, _ = atom.GetVal(XConn, "NET_SYSTEM_TRAY_ORIENTAION")
	XA_MANAGER, _ = atom.GetVal(XConn, "MANAGER")
}
