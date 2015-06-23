package loader

import (
	"pkg.linuxdeepin.com/lib/log"
)

var logger = log.NewLogger("dde-daemon/loader")

type Module struct {
	Name   string
	Start  func()
	Stop   func()
	Enable bool
}

var modules = make([]*Module, 0)

func getModule(name string) (module *Module) {
	for _, m := range modules {
		if m.Name == name {
			module = m
			break
		}
	}
	if module == nil {
		logger.Warning("target module not found:", name)
	}
	return
}
func isModuleExist(name string) (ok bool) {
	for _, m := range modules {
		if m.Name == name {
			ok = true
			break
		}
	}
	return
}

func Start(name string) {
	m := getModule(name)
	if m != nil {
		doStart(m)
	}
}
func doStart(m *Module) {
	logger.Info("Start module:", m.Name)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Start module", m.Name, "failed:", err)
		}
	}()
	m.Start()
}

func Stop(name string) {
	m := getModule(name)
	if m != nil {
		doStop(m)
	}
}
func doStop(m *Module) {
	logger.Info("Stop module:", m.Name)
	defer func() {
		if err := recover(); err != nil {
			logger.Error("Stop module", m.Name, "failed:", err)
		}
	}()
	m.Stop()
}

func Enable(name string, enable bool) {
	m := getModule(name)
	if m != nil {
		m.Enable = enable
	}
}

func Register(newModule *Module) {
	logger.Info("Register module:", newModule.Name)
	if isModuleExist(newModule.Name) {
		logger.Warning("module already registered:", newModule.Name)
		return
	}
	if newModule.Start == nil || newModule.Stop == nil {
		logger.Error("can't register an incomplete module:", newModule.Name)
		return
	}
	modules = append([]*Module{newModule}, modules...)
}

func StartAll() {
	logger.Info("Start all modules")
	for _, m := range modules {
		if !m.Enable {
			logger.Info("skip disabled module:", m.Name)
			continue
		}
		doStart(m)
	}
}

func StopAll() {
	logger.Info("Stop all modules")
	for _, m := range modules {
		doStop(m)
	}
}
