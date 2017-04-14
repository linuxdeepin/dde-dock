/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package loader

import (
	"fmt"
	"pkg.deepin.io/lib/log"
	"sync"
	"time"
)

type EnableFlag int

const (
	EnableFlagNone EnableFlag = 1 << iota
	EnableFlagIgnoreMissingModule
	EnableFlagForceStart
)

func (flags EnableFlag) HasFlag(flag EnableFlag) bool {
	return flags&flag != 0
}

const (
	ErrorNoDependencies int = iota
	ErrorCircleDependencies
	ErrorMissingModule
	ErrorInternalError
	ErrorConflict
)

type EnableError struct {
	ModuleName string
	Code       int
	detail     string
}

func (e *EnableError) Error() string {
	switch e.Code {
	case ErrorNoDependencies:
		return fmt.Sprintf("%s's dependencies is not meet, %s is need", e.ModuleName, e.detail)
	case ErrorCircleDependencies:
		return "dependency circle"
		// return fmt.Sprintf("%s and %s dependency each other.", e.ModuleName, e.detail)
	case ErrorMissingModule:
		return fmt.Sprintf("%s is missing", e.ModuleName)
	case ErrorInternalError:
		return fmt.Sprintf("%s started failed: %s", e.ModuleName, e.detail)
	case ErrorConflict:
		return fmt.Sprintf("tring to enable disabled module(%s)", e.ModuleName)
	}
	panic("EnableError: Unknown Error, Should not be reached")
}

type Loader struct {
	modules Modules
	log     *log.Logger
	lock    sync.Mutex
}

func (l *Loader) SetLogLevel(pri log.Priority) {
	l.log.SetLogLevel(pri)

	l.lock.Lock()
	defer l.lock.Unlock()

	for _, module := range l.modules {
		module.SetLogLevel(pri)
	}
}

func (l *Loader) AddModule(m Module) {
	l.lock.Lock()
	defer l.lock.Unlock()

	tmp := l.modules.Get(m.Name())
	if tmp != nil {
		l.log.Debug("Register", m.Name(), "is already registered")
		return
	}

	l.log.Debug("Register module:", m.Name())
	l.modules = append(l.modules, m)
}

func (l *Loader) DeleteModule(name string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.modules, _ = l.modules.Delete(name)
}

func (l *Loader) List() []Module {
	var modules Modules

	l.lock.Lock()
	for _, module := range l.modules {
		modules = append(modules, module)
	}
	l.lock.Unlock()

	return modules
}

func (l *Loader) GetModule(name string) Module {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.modules.Get(name)
}

func (l *Loader) EnableModules(enablingModules []string, disableModules []string, flag EnableFlag) error {
	l.lock.Lock()
	defer l.lock.Unlock()

	builder := NewDAGBuilder(l, enablingModules, disableModules, flag)
	dag, err := builder.Execute()
	if err != nil {
		return err
	}

	nodes, ok := dag.TopologicalDag()
	if !ok {
		return &EnableError{Code: ErrorCircleDependencies}
	}

	for _, name := range enablingModules {
		node := nodes.Get(name)
		if node == nil {
			continue
		}
		module := l.modules.Get(node.ID)
		l.log.Info("enable module", node.ID)
		startTime := time.Now()
		err := module.Enable(true)
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		if err != nil {
			l.log.Fatalf("enable module %s failed: %s, cost %s", node.ID, err, duration)
		}
		l.log.Info("enable module", node.ID, "done, cost", duration)
	}

	return nil
}
