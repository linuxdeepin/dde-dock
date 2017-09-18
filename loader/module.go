/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package loader

import (
	"fmt"
	"pkg.deepin.io/lib/log"
)

type Module interface {
	Name() string
	IsEnable() bool
	Enable(bool) error
	GetDependencies() []string
	SetLogLevel(log.Priority)
	LogLevel() log.Priority
	ModuleImpl
}
type Modules []Module

type ModuleImpl interface {
	Start() error
	Stop() error
}

type ModuleBase struct {
	impl    ModuleImpl
	enabled bool
	name    string
	log     *log.Logger
}

func NewModuleBase(name string, impl ModuleImpl, logger *log.Logger) *ModuleBase {
	return &ModuleBase{
		name: name,
		impl: impl,
		log:  logger,
	}
}

func (d *ModuleBase) doEnable(enable bool) error {
	if d.impl != nil {
		var fn func() error = d.impl.Stop
		if enable {
			fn = d.impl.Start
		}

		if err := fn(); err != nil {
			return err
		}
	}
	d.enabled = enable
	return nil
}

func (d *ModuleBase) Enable(enable bool) error {
	if d.enabled == enable {
		return fmt.Errorf("%s daemon is already started", d.name)
	}
	return d.doEnable(enable)
}

func (d *ModuleBase) IsEnable() bool {
	return d.enabled
}

func (d *ModuleBase) Name() string {
	return d.name
}

func (d *ModuleBase) SetLogLevel(pri log.Priority) {
	d.log.SetLogLevel(pri)
}

func (d *ModuleBase) LogLevel() log.Priority {
	return d.log.GetLogLevel()
}

func (l Modules) Get(name string) Module {
	for _, v := range l {
		if v.Name() == name {
			return v
		}
	}
	return nil
}

func (l Modules) Delete(name string) (Modules, bool) {
	var (
		tmp     Modules
		deleted bool
	)
	for _, v := range l {
		if v.Name() == name {
			deleted = true
			continue
		}
		tmp = append(tmp, v)
	}
	return tmp, deleted
}

func (l Modules) List() []string {
	var names []string
	for _, v := range l {
		names = append(names, v.Name())
	}
	return names
}

func (l Modules) Len() int {
	return len(l)
}

func (l Modules) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l Modules) Less(i, j int) bool {
	return l[i].Name() < l[j].Name()
}
