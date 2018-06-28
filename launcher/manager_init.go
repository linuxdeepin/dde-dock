/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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

package launcher

import (
	"path/filepath"

	"pkg.deepin.io/lib/appinfo/desktopappinfo"
	"pkg.deepin.io/lib/dbusutil/proxy"
)

const (
	appsDBusDest       = "com.deepin.daemon.Apps"
	appsDBusObjectPath = "/com/deepin/daemon/Apps"
)

func (m *Manager) initItems() {
	// load items
	m.items = make(map[string]*Item)

	skipDirs := make(map[string][]string)
	skipDirs["/usr/share/applications"] = []string{"screensavers"}
	userAppsDir := getUserAppDir()
	skipDirs[userAppsDir] = []string{"menu-xdg"}
	allApps := desktopappinfo.GetAll(skipDirs)
	for _, ai := range allApps {
		if !ai.IsExecutableOk() ||
			isDeepinCustomDesktopFile(ai.GetFileName()) {
			continue
		}
		item := NewItemWithDesktopAppInfo(ai)
		m.setItemID(item)

		if m.hiddenByGSettings(item.ID) {
			continue
		}
		m.addItem(item)
	}
	logger.Debug("load items count:", len(m.items))
}

func shouldCheckDesktopFile(filename string) bool {
	dir, basename := filepath.Split(filename)
	dir = filepath.Clean(dir)
	matched, _ := filepath.Match(desktopFilePattern, basename)
	if !matched {
		return false
	}

	// ignore $HOME/.local/share/applications/menu-xdg/
	skipDir := filepath.Join(getUserAppDir(), "menu-xdg")
	if dir == skipDir {
		return false
	}
	return true
}

type popPushOp struct {
	popCount  int
	runesPush []rune
}

func (m *Manager) handlePopPushOps() {
	stack := m.searchTaskStack
	for op := range m.popPushOpChan {
		logger.Debug("op:", op)

		for i := 0; i < op.popCount; i++ {
			stack.Pop()
		}
		if len(op.runesPush) == 0 {
			// emit top result
			top := stack.topTask()
			if top == nil {
				logger.Debug("emit SearchDone []")
				m.service.Emit(m, "SearchDone", []string{})
			} else {
				top.emitResult()
			}
		}
		for _, v := range op.runesPush {
			stack.Push(v)
		}
	}
}

func (m *Manager) destroy() {
	m.appsObj.RemoveHandler(proxy.RemoveAllHandlers)
	m.sysSigLoop.Stop()
}
