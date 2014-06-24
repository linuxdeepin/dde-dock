package loader

import (
	"pkg.linuxdeepin.com/lib/logger"
)

var Logger = logger.NewLogger("dde.daemon.loader")

type Module struct {
	Name   string
	Start  func()
	Stop   func()
	Enable bool
}

var modules = make([]*Module, 0)

func Enable(name string, enable bool) {
	for _, m := range modules {
		if m.Name == name {
			m.Enable = enable
		}
	}
}

func Register(newModule *Module) {
	for _, m := range modules {
		if m.Name == newModule.Name {
			return
		}
	}
	if newModule.Start == nil || newModule.Stop == nil {
		Logger.Fatal("can't register an incomplete module " + newModule.Name)
	}
	modules = append([]*Module{newModule}, modules...)
}

func Start() {
	for _, m := range modules {
		func() {
			defer func() {
				if err := recover(); err != nil {
					Logger.Error("Start module", m.Name, "failed:", err)
				}
			}()
			m.Start()
		}()
	}
}
func Stop() {
	for _, m := range modules {
		func() {
			defer func() {
				if err := recover(); err != nil {
					Logger.Error("Stop module", m.Name, "failed:", err)
				}
			}()
			m.Start()
		}()
	}
}
