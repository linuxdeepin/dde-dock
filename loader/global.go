package loader

import (
	"pkg.linuxdeepin.com/lib/log"
	"sync"
)

var loaderInitializer sync.Once

var getLoader = func() func() *Loader {
	var loader *Loader
	return func() *Loader {
		loaderInitializer.Do(func() {
			loader = &Loader{
				modules: map[string]Module{},
				log:     log.NewLogger("dde-daemon/loader"),
			}
		})
		return loader
	}
}()

func Register(m Module) {
	loader := getLoader()
	loader.AddModule(m)
}

func List() []Module {
	return getLoader().List()
}

func GetModule(name string) Module {
	return getLoader().GetModule(name)
}

func SetLogLevel(pri log.Priority) {
	getLoader().SetLogLevel(pri)
}

func EnableModules(enablingModules []string, disableModules []string, flag EnableFlag) error {
	return getLoader().EnableModules(enablingModules, disableModules, flag)
}

func StartAll() {
	allModules := getLoader().List()
	modules := []string{}
	for _, module := range allModules {
		modules = append(modules, module.Name())
	}
	getLoader().EnableModules(modules, []string{}, EnableFlagNone)
}

// TODO: check dependencies
func StopAll() {
	modules := getLoader().List()
	for _, module := range modules {
		module.Enable(false)
	}
}
