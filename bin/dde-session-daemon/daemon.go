/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"errors"
	"fmt"
	"gir/gio-2.0"
	"gir/glib-2.0"
	"os"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
	"runtime/pprof"
	"sync"
)

const (
	ProfTypeCPU = "cpu"
	ProfTypeMem = "memory"
)

func runMainLoop() {
	ddeSessionRegister()
	dbus.DealWithUnhandledMessage()
	go glib.StartLoop()

	if err := dbus.Wait(); err != nil {
		logger.Errorf("Lost dbus: %v", err)
		os.Exit(-1)
	}

	logger.Info("dbus connection is closed by user")
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
	logLevel        log.Priority
	log             *log.Logger
	settings        *gio.Settings
	enabledModules  map[string]loader.Module
	disabledModules map[string]loader.Module

	cpuLocker sync.Mutex
	cpuWriter *os.File
}

// Profile: heap, goroutine, threadcreate, block
func (s *SessionDaemon) WriteProfile(profile, log string) error {
	var p = pprof.Lookup(profile)
	if p == nil {
		return fmt.Errorf("Profile '%s' not exists", profile)
	}

	fw, err := os.Create(log)
	if err != nil {
		return err
	}
	defer fw.Close()

	return p.WriteTo(fw, 1)
}

// Profile: cpu
func (s *SessionDaemon) StartCPUProfile(log string) error {
	s.cpuLocker.Lock()
	defer s.cpuLocker.Unlock()
	fw, err := os.Create(log)
	if err != nil {
		return err
	}

	err = pprof.StartCPUProfile(fw)
	if err != nil {
		fw.Close()
		return err
	}

	s.cpuWriter = fw
	return nil
}

func (s *SessionDaemon) StopCPUProfile() {
	s.cpuLocker.Lock()
	defer s.cpuLocker.Unlock()
	if s.cpuWriter == nil {
		return
	}

	pprof.StopCPUProfile()
	s.cpuWriter.Close()
	s.cpuWriter = nil
}

func (*SessionDaemon) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.Daemon",
		ObjectPath: "/com/deepin/daemon/Daemon",
		Interface:  "com.deepin.daemon.Daemon",
	}
}

func NewSessionDaemon(flags *Flags, settings *gio.Settings, logger *log.Logger) *SessionDaemon {
	session := &SessionDaemon{
		flags:           flags,
		settings:        settings,
		logLevel:        log.LevelInfo,
		log:             logger,
		enabledModules:  map[string]loader.Module{},
		disabledModules: map[string]loader.Module{},
	}

	session.initModules()

	return session
}

func (s *SessionDaemon) exitIfNotSingleton() error {
	if !lib.UniqueOnSession(s.GetDBusInfo().Dest) {
		return errors.New("There already has a dde daemon running.")
	}
	return nil
}

func (s *SessionDaemon) register() error {
	if err := s.exitIfNotSingleton(); err != nil {
		return err
	}

	if err := dbus.InstallOnSession(s); err != nil {
		s.log.Fatal(err)
	}

	return nil
}

func (s *SessionDaemon) defaultAction() {
	loader.SetLogLevel(s.logLevel)

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
			s.enabledModules[name] = module
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
	return keys(s.enabledModules)
}

func (s *SessionDaemon) EnableModules(enablingModules []string) error {
	return loader.EnableModules(enablingModules, s.getDisabledModules(), getEnableFlag(s.flags))
}

func (s *SessionDaemon) DisableModules(disableModules []string) error {
	return loader.EnableModules(s.getEnabledModules(), disableModules, getEnableFlag(s.flags))
}

func (s *SessionDaemon) ListModule(name string) error {
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
