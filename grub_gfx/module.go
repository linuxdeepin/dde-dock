package grub_gfx

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

const moduleName = "grub-gfx"

var logger = log.NewLogger(moduleName)

type module struct {
	*loader.ModuleBase
}

func (*module) GetDependencies() []string {
	return nil
}

func (d *module) Start() error {
	logger.Debug("module start")
	detectChange()
	return nil
}

func (d *module) Stop() error {
	return nil
}

func newModule() *module {
	d := new(module)
	d.ModuleBase = loader.NewModuleBase(moduleName, d, logger)
	return d
}

func init() {
	loader.Register(newModule())
}
