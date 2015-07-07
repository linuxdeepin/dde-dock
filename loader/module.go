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
