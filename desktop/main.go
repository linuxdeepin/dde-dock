package desktop

import (
	// "fmt"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/initializer"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("dde-daemon/desktop")

// Daemon is a wrapper for desktop daemon used in dde-sesion-daemon.
type Daemon struct {
	*loader.ModuleBase
	app *Application
}

// NewDaemon creates new daemon.
func NewDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("desktop", daemon, logger)
	return daemon
}

// GetDependencies returns the dependencies of desktop.
func (d *Daemon) GetDependencies() []string {
	return []string{}
}

// Stop stops desktop daemon.
func (d *Daemon) Stop() error {
	dbus.UnInstallObject(d.app)
	d.app = nil
	return nil
}

// Start starts desktop daemon.
func (d *Daemon) Start() error {
	initializer := initializer.NewInitializer()
	initializer.InitOnSessionBus(func(v interface{}) (interface{}, error) {
		return NewSettings()
	}).InitOnSessionBus(func(v interface{}) (interface{}, error) {
		d.app = NewApplication(v.(*Settings))
		return d.app, nil
	})

	return initializer.GetError()
}

func init() {
	loader.Register(NewDaemon())
}
