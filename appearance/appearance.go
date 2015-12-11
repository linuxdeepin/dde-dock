// Manage desktop appearance
package appearance

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var (
	_m     *Manager
	logger = log.NewLogger("daemon/appearance")
)

type Daemon struct {
	*loader.ModuleBase
}

func init() {
	loader.Register(NewAppearanceDaemon(logger))
}

func NewAppearanceDaemon(logger *log.Logger) *Daemon {
	var d = new(Daemon)
	d.ModuleBase = loader.NewModuleBase("appearance", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (*Daemon) Start() error {
	if _m != nil {
		return nil
	}

	logger.BeginTracing()
	_m = NewManager()
	err := dbus.InstallOnSession(_m)
	if err != nil {
		logger.Error("Install dbus failed:", err)
		_m.destroy()
		logger.EndTracing()
		return err
	}
	go _m.listenCursorChanged()
	_m.listenGSettingChanged()

	return nil
}

func (*Daemon) Stop() error {
	if _m == nil {
		return nil
	}

	_m.destroy()
	logger.EndTracing()
	_m = nil
	return nil
}
