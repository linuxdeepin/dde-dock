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

package dock

import (
	"path/filepath"

	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/xdg/basedir"

	x "github.com/linuxdeepin/go-x11-client"
)

func init() {
	loader.Register(NewDaemon(logger))
}

var (
	logger      = log.NewLogger("daemon/dock")
	homeDir     string
	scratchDir  string
	dockManager *Manager

	globalXConn *x.Conn

	atomNetShowingDesktop       x.Atom
	atomNetClientList           x.Atom
	atomNetActiveWindow         x.Atom
	atomNetWMIcon               x.Atom
	atomNetWMName               x.Atom
	atomNetWMState              x.Atom
	atomNetWMWindowType         x.Atom
	atomXEmbedInfo              x.Atom
	atomNetFrameExtents         x.Atom
	atomGtkFrameExtents         x.Atom
	atomNetWmStateHidden        x.Atom
	atomWmWindowRole            x.Atom
	atomUTF8String              x.Atom
	atomFlatpakAppId            x.Atom
	atomGtkApplicationId        x.Atom
	atomNetWmWindowOpacity      x.Atom
	atomWmClientLeader          x.Atom
	atomWmCommand               x.Atom
	atomNetWmStateFocused       x.Atom
	atomNetWmWindowTypeDesktop  x.Atom
	atomNetWmActionMinimize     x.Atom
	atomWmStateDemandsAttention x.Atom
	atomNetWmStateSkipTaskbar   x.Atom
	atomNetWmStateModal         x.Atom
	atomNetWmWindowTypeDialog   x.Atom
	atomNetWmStateMaximizedVert x.Atom
	atomNetWmStateMaximizedHorz x.Atom
	atomNetWmStateAbove         x.Atom
	atomNetWmActionClose        x.Atom
	atomNetWmAllowedActions     x.Atom
	atomNetWmPid                x.Atom
	atomMotifWmHints            x.Atom
)

func initDir() {
	homeDir = basedir.GetUserHomeDir()
	scratchDir = filepath.Join(basedir.GetUserConfigDir(), "dock/scratch")
	logger.Debugf("scratch dir: %q", scratchDir)
}

func initAtom() {
	atomNetShowingDesktop, _ = getAtom("_NET_SHOWING_DESKTOP")
	atomNetClientList, _ = getAtom("_NET_CLIENT_LIST")
	atomNetActiveWindow, _ = getAtom("_NET_ACTIVE_WINDOW")
	atomNetWMIcon, _ = getAtom("_NET_WM_ICON")
	atomNetWMName, _ = getAtom("_NET_WM_NAME")
	atomNetWMState, _ = getAtom("_NET_WM_STATE")
	atomNetWMWindowType, _ = getAtom("_NET_WM_WINDOW_TYPE")
	atomXEmbedInfo, _ = getAtom("_XEMBED_INFO")
	atomNetFrameExtents, _ = getAtom("_NET_FRAME_EXTENTS")
	atomGtkFrameExtents, _ = getAtom("_GTK_FRAME_EXTENTS")
	atomNetWmStateHidden, _ = getAtom("_NET_WM_STATE_HIDDEN")
	atomWmWindowRole, _ = getAtom("WM_WINDOW_ROLE")
	atomUTF8String, _ = getAtom("UTF8_STRING")
	atomFlatpakAppId, _ = getAtom("FLATPAK_APPID")
	atomGtkApplicationId, _ = getAtom("_GTK_APPLICATION_ID")
	atomNetWmWindowOpacity, _ = getAtom("_NET_WM_WINDOW_OPACITY")
	atomWmClientLeader, _ = getAtom("WM_CLIENT_LEADER")
	atomWmCommand, _ = getAtom("WM_COMMAND")
	atomNetWmStateFocused, _ = getAtom("_NET_WM_STATE_FOCUSED")
	atomNetWmWindowTypeDesktop, _ = getAtom("_NET_WM_WINDOW_TYPE_DESKTOP")
	atomNetWmActionMinimize, _ = getAtom("_NET_WM_ACTION_MINIMIZE")
	atomWmStateDemandsAttention, _ = getAtom("_NET_WM_STATE_DEMANDS_ATTENTION")
	atomNetWmStateSkipTaskbar, _ = getAtom("_NET_WM_STATE_SKIP_TASKBAR")
	atomNetWmStateModal, _ = getAtom("_NET_WM_STATE_MODAL")
	atomNetWmWindowTypeDialog, _ = getAtom("_NET_WM_WINDOW_TYPE_DIALOG")
	atomNetWmStateMaximizedVert, _ = getAtom("_NET_WM_STATE_MAXIMIZED_VERT")
	atomNetWmStateMaximizedHorz, _ = getAtom("_NET_WM_STATE_MAXIMIZED_HORZ")
	atomNetWmStateAbove, _ = getAtom("_NET_WM_STATE_ABOVE")
	atomNetWmActionClose, _ = getAtom("_NET_WM_ACTION_CLOSE")
	atomNetWmAllowedActions, _ = getAtom("_NET_WM_ALLOWED_ACTIONS")
	atomNetWmPid, _ = getAtom("_NET_WM_PID")
	atomMotifWmHints, _ = getAtom("_MOTIF_WM_HINTS")
}
