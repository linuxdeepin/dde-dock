package miracast

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("miracast", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{"network", "audio"}
}

var (
	_m     *Miracast
	logger = log.NewLogger(dbusDest)
)

func (d *Daemon) Start() error {
	if _m != nil {
		return nil
	}

	m, err := newMiracast()
	if err != nil {
		logger.Error("Failed to new manager:", err)
		return err
	}
	// fix dbus timeout
	go func() {
		m.init()
		m.ensureMiracleActive()
	}()
	_m = m

	err = dbus.InstallOnSession(m)
	if err != nil {
		logger.Error("Failed to install bus:", err)
		_m.destroy()
		_m = nil
		return err
	}
	dbus.DealWithUnhandledMessage()

	return nil
}

func (*Daemon) Stop() error {
	if _m == nil {
		return nil
	}
	_m.destroy()
	dbus.UnInstallObject(_m)
	_m = nil
	return nil
}
