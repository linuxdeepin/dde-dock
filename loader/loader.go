package loader

import (
	"fmt"
	"pkg.linuxdeepin.com/lib/log"
	"sort"
	"sync"
)

type byName []Module

func (l byName) Len() int {
	return len(l)
}

func (l byName) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l byName) Less(i, j int) bool {
	return l[i].Name() < l[j].Name()
}

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
	modules map[string]Module
	log     *log.Logger
	lock    sync.Mutex
}

func (l *Loader) SetLogLevel(pri log.Priority) {
	l.log.SetLogLevel(pri)
}

func (l *Loader) AddModule(m Module) {
	l.lock.Lock()
	defer l.lock.Unlock()

	_, ok := l.modules[m.Name()]
	if ok {
		l.log.Debug("Register", m.Name(), "is already registered")
		return
	}

	l.log.Debug("Register module:", m.Name())
	l.modules[m.Name()] = m
}

func (l *Loader) DeleteModule(name string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.modules, name)
}

func (l *Loader) List() []Module {
	modules := []Module{}

	l.lock.Lock()
	for _, module := range l.modules {
		modules = append(modules, module)
	}
	l.lock.Unlock()

	sort.Sort(byName(modules))
	return modules
}

func (l *Loader) GetModule(name string) Module {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.modules[name]
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

	for _, node := range nodes {
		module := l.modules[node.ID]
		l.log.Debug("enable module", node.ID)
		err := module.Enable(true)
		if err != nil {
			l.log.Errorf("enable module(%s) failed", node.ID)
		}
	}

	return nil
}
