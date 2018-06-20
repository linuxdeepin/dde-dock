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

package power

import (
	"fmt"
)

var submoduleList []func(*Manager) (string, submodule, error)

type submodule interface {
	Start() error
	Destroy()
}

func (m *Manager) initSubmodules() {
	m.submodules = make(map[string]submodule, len(submoduleList))
	// new all submodule
	for _, newMethod := range submoduleList {
		name, submoduleInstance, err := newMethod(m)
		logger.Debug("New submodule:", name)
		if err != nil {
			logger.Warningf("New submodule %v failed: %v", name, err)
			continue
		}
		m.submodules[name] = submoduleInstance
	}
}

func (m *Manager) _startSubmodule(name string) error {
	submodule, ok := m.submodules[name]
	if !ok {
		return fmt.Errorf("%v not exist", name)
	}
	return submodule.Start()
}

func (m *Manager) startSubmodules() {
	startOrder := []string{"PowerSavePlan", "LidSwitchHandler"}
	for _, name := range startOrder {
		logger.Infof("submodule %v start", name)
		err := m._startSubmodule(name)
		if err != nil {
			logger.Warningf("submodule %v start failed: %v", name, err)
		}
	}
}

func (m *Manager) destroySubmodules() {
	if m.submodules != nil {
		for name, submodule := range m.submodules {
			logger.Debug("destroy submodule:", name)
			submodule.Destroy()
		}
		m.submodules = nil
	}
}
