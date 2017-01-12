package mime

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/mime")

func init() {
	loader.Register(NewDaemon(logger))
}

type Daemon struct {
	*loader.ModuleBase
	manager *Manager
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("mime", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	logger.BeginTracing()
	d.manager = NewManager()

	err := dbus.InstallOnSession(d.manager)
	if err != nil {
		logger.Warning("Install Manager dbus failed:", err)
		return err
	}

	media, err := NewMedia()
	if err != nil {
		logger.Error("New Media failed:", err)
		return err
	}
	d.manager.media = media

	err = dbus.InstallOnSession(media)
	if err != nil {
		logger.Warning("Install Media dbus failed:", err)
		return err
	}

	d.manager.initConfigData()
	return nil
}

func (d *Daemon) Stop() error {
	if d.manager == nil {
		return nil
	}

	d.manager.destroy()
	d.manager = nil
	return nil
}
