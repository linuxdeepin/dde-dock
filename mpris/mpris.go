package mpris

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

var _m *Manager

type Daemon struct {
	*loader.ModuleBase
}

func init() {
	loader.Register(NewDaemon(logger))
}

func NewDaemon(log *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("mpris", daemon, log)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if _m != nil {
		return nil
	}

	logger.BeginTracing()

	var err error
	_m, err = NewManager()
	if err != nil {
		logger.Error("Create mpris manager failed:", err)
		logger.EndTracing()
		return err
	}
	_m.listenMediakey()

	return nil
}

func (d *Daemon) Stop() error {
	if _m == nil {
		return nil
	}

	_m.destroy()
	_m = nil

	return nil
}
