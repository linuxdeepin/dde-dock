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

package main

import (
	"fmt"
	"os"
	"sync"

	"gir/gio-2.0"
	"gir/glib-2.0"
	"pkg.deepin.io/dde/api/session"
	"pkg.deepin.io/dde/daemon/calltrace"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/gsettings"
	"pkg.deepin.io/lib/log"
)

const (
	ProfTypeCPU = "cpu"
	ProfTypeMem = "memory"

	dbusPath        = "/com/deepin/daemon/Daemon"
	dbusServiceName = "com.deepin.daemon.Daemon"
	dbusInterface   = dbusServiceName
)

func runMainLoop() {
	err := gsettings.StartMonitor()
	if err != nil {
		logger.Fatal(err)
	}

	session.Register()
	listenDaemonSettings()

	glib.StartLoop()
	logger.Info("Loop has been terminated!")
	os.Exit(0)
}

func getEnableFlag(flag *Flags) loader.EnableFlag {
	enableFlag := loader.EnableFlagIgnoreMissingModule

	if *flag.IgnoreMissingModules {
		enableFlag = loader.EnableFlagNone
	}

	if *flag.ForceStart {
		enableFlag |= loader.EnableFlagForceStart
	}

	return enableFlag
}

type SessionDaemon struct {
	flags           *Flags
	log             *log.Logger
	settings        *gio.Settings
	enabledModules  loader.Modules
	disabledModules map[string]loader.Module

	cpuLocker sync.Mutex
	cpuWriter *os.File

	methods *struct {
		EnableModules  func() `in:"enablingModules"`
		DisableModules func() `in:"disableModules "`
		ListModule     func() `in:"name"`
		CallTrace      func() `in:"times,seconds"`
	}
}

func (*SessionDaemon) GetInterfaceName() string {
	return dbusInterface
}

func NewSessionDaemon(flags *Flags, settings *gio.Settings, logger *log.Logger) *SessionDaemon {
	session := &SessionDaemon{
		flags:           flags,
		settings:        settings,
		log:             logger,
		enabledModules:  loader.Modules{},
		disabledModules: map[string]loader.Module{},
	}

	session.initModules()

	return session
}

func (s *SessionDaemon) register(service *dbusutil.Service) error {
	err := service.Export(dbusPath, s)
	if err != nil {
		s.log.Fatal(err)
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}
	return nil
}

func (s *SessionDaemon) defaultAction() {
	err := loader.EnableModules(s.getEnabledModules(), s.getDisabledModules(), getEnableFlag(s.flags))
	if err != nil {
		fmt.Println(err)
		// TODO: define exit code.
		os.Exit(3)
	}
}

func (s *SessionDaemon) execDefaultAction() {
	s.defaultAction()
}

func (s *SessionDaemon) initModules() {
	allModules := loader.List()
	for _, module := range allModules {
		name := module.Name()
		if s.settings.GetBoolean(name) {
			s.enabledModules = append(s.enabledModules, module)
		} else {
			s.disabledModules[name] = module
		}
	}
}

func keys(m map[string]loader.Module) []string {
	keys := []string{}
	for key, _ := range m {
		keys = append(keys, key)
	}
	return keys
}

func (s *SessionDaemon) getDisabledModules() []string {
	return keys(s.disabledModules)
}

func (s *SessionDaemon) getEnabledModules() []string {
	return s.enabledModules.List()
}

func (s *SessionDaemon) enableModules(enablingModules []string) error {
	disabledModules := s.getDisabledModules()
	disabledModules = filterList(disabledModules, enablingModules)
	return loader.EnableModules(enablingModules, disabledModules, getEnableFlag(s.flags))
}

func (s *SessionDaemon) EnableModules(enablingModules []string) *dbus.Error {
	err := s.enableModules(enablingModules)
	return dbusutil.ToError(err)
}

func (s *SessionDaemon) disableModules(disableModules []string) error {
	enablingModules := s.getEnabledModules()
	enablingModules = filterList(enablingModules, disableModules)
	return loader.EnableModules(enablingModules, disableModules, getEnableFlag(s.flags))
}

func (s *SessionDaemon) DisableModules(disableModules []string) *dbus.Error {
	err := s.disableModules(disableModules)
	return dbusutil.ToError(err)
}

func (s *SessionDaemon) listModule(name string) error {
	if name == "" {
		for _, module := range loader.List() {
			fmt.Println(module.Name())
		}
		return nil
	}

	module := loader.GetModule(name)
	if module == nil {
		return fmt.Errorf("no such a module named %s", name)
	}

	for _, m := range module.GetDependencies() {
		fmt.Println(m)
	}

	return nil
}

func (s *SessionDaemon) ListModule(name string) *dbus.Error {
	err := s.listModule(name)
	return dbusutil.ToError(err)
}

func (s *SessionDaemon) CallTrace(times, seconds uint32) *dbus.Error {
	ct, err := calltrace.NewManager(seconds / times)
	if err != nil {
		logger.Warning("Failed to start calltrace:", err)
		return dbusutil.ToError(err)
	}
	ct.SetAutoDestroy(seconds)
	return nil
}

func filterList(origin, condition []string) []string {
	if len(condition) == 0 {
		return origin
	}

	var tmp = make(map[string]struct{})
	for _, v := range condition {
		tmp[v] = struct{}{}
	}

	var ret []string
	for _, v := range origin {
		_, ok := tmp[v]
		if ok {
			continue
		}
		ret = append(ret, v)
	}
	return ret
}
