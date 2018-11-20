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
	"errors"
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
	flags                *Flags
	log                  *log.Logger
	settings             *gio.Settings
	part1EnabledModules  []string
	part1DisabledModules []string
	part2EnabledModules  []string
	part2DisabledModules []string

	cpuLocker sync.Mutex
	cpuWriter *os.File

	methods *struct {
		CallTrace func() `in:"times,seconds"`
	}
}

func (*SessionDaemon) GetInterfaceName() string {
	return dbusInterface
}

func NewSessionDaemon(flags *Flags, settings *gio.Settings, logger *log.Logger) *SessionDaemon {
	session := &SessionDaemon{
		flags:    flags,
		settings: settings,
		log:      logger,
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

func (s *SessionDaemon) initModules() {
	part1ModuleNames := []string{
		"dock",
		"trayicon",
		"x-event-monitor",
	}

	part2ModuleNames := []string{
		"network",
		"audio",
		"screensaver",
		"sessionwatcher",
		"power", // need screensaver and sessionwatcher
		"launcher",
		"service-trigger",
		"clipboard",
		"keybinding",
		"appearance",
		"inputdevices",
		"gesture",
		"housekeeping",
		"timedate",
		"bluetooth",
		"screenedge",
		"fprintd",
		"mime",
		"miracast", // need network
		"systeminfo",
		"lastore",
		"grub-gfx",
		"calltrace",
		"debug",
	}

	allModules := loader.List()
	if len(part1ModuleNames)+len(part2ModuleNames) != len(allModules) {
		panic("module names len not equal")
	}

	for _, moduleName := range part1ModuleNames {
		if s.isModuleDefaultEnabled(moduleName) {
			s.part1EnabledModules = append(s.part1EnabledModules, moduleName)
		} else {
			s.part1DisabledModules = append(s.part1DisabledModules, moduleName)
		}
	}

	for _, moduleName := range part2ModuleNames {
		if s.isModuleDefaultEnabled(moduleName) {
			s.part2EnabledModules = append(s.part2EnabledModules, moduleName)
		} else {
			s.part2DisabledModules = append(s.part2DisabledModules, moduleName)
		}
	}
}

func (s *SessionDaemon) isModuleDefaultEnabled(moduleName string) bool {
	mod := loader.GetModule(moduleName)
	if mod == nil {
		panic(fmt.Errorf("not found module %q", moduleName))
	}

	return s.settings.GetBoolean(moduleName)
}

func (s *SessionDaemon) getAllDefaultEnabledModules() []string {
	result := make([]string, len(s.part1EnabledModules)+len(s.part2EnabledModules))
	n := copy(result, s.part1EnabledModules)
	copy(result[n:], s.part2EnabledModules)
	return result
}

func (s *SessionDaemon) getAllDefaultDisabledModules() []string {
	result := make([]string, len(s.part1DisabledModules)+len(s.part2DisabledModules))
	n := copy(result, s.part1DisabledModules)
	copy(result[n:], s.part2DisabledModules)
	return result
}

func (s *SessionDaemon) execDefaultAction() {
	var err error
	if hasDDECookie {
		// start part1
		err = loader.EnableModules(s.part1EnabledModules, s.part1DisabledModules, 0)
		session.Register()

	} else {
		err = loader.EnableModules(s.getAllDefaultEnabledModules(),
			s.getAllDefaultDisabledModules(), getEnableFlag(s.flags))
	}

	if err != nil {
		fmt.Println(err)
		os.Exit(3)
	}
}

func (s *SessionDaemon) enableModules(enablingModules []string) error {
	disabledModules := filterList(s.getAllDefaultDisabledModules(), enablingModules)
	return loader.EnableModules(enablingModules, disabledModules, getEnableFlag(s.flags))
}

func (s *SessionDaemon) disableModules(disableModules []string) error {
	enablingModules := filterList(s.getAllDefaultEnabledModules(), disableModules)
	return loader.EnableModules(enablingModules, disableModules, getEnableFlag(s.flags))
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

func (s *SessionDaemon) CallTrace(times, seconds uint32) *dbus.Error {
	ct, err := calltrace.NewManager(seconds / times)
	if err != nil {
		logger.Warning("Failed to start calltrace:", err)
		return dbusutil.ToError(err)
	}
	ct.SetAutoDestroy(seconds)
	return nil
}

func (s *SessionDaemon) StartPart2() *dbus.Error {
	if !hasDDECookie {
		return dbusutil.ToError(errors.New("env DDE_SESSION_PROCESS_COOKIE_ID is empty"))
	}
	// start part2
	err := loader.EnableModules(s.part2EnabledModules, s.part2DisabledModules, 0)
	return dbusutil.ToError(err)
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
