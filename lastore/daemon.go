package lastore

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

const (
	dbusPath        = "/com/deepin/LastoreSessionHelper"
	dbusServiceName = "com.deepin.LastoreSessionHelper"
)

var logger = log.NewLogger("daemon/lastore")

func init() {
	loader.Register(newDaemon())
}

type Daemon struct {
	lastore *Lastore
	*loader.ModuleBase
}

func newDaemon() *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("lastore", daemon, logger)
	return daemon
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	service := loader.GetService()

	lastore, err := newLastore(service)
	if err != nil {
		return err
	}
	d.lastore = lastore

	err = service.Export(dbusPath, lastore, lastore.syncConfig)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	err = lastore.syncConfig.Register()
	if err != nil {
		logger.Warning("Failed to register sync service:", err)
	}
	return nil
}

func (d *Daemon) Stop() error {
	if d.lastore == nil {
		return nil
	}

	service := loader.GetService()
	err := service.ReleaseName(dbusServiceName)
	if err != nil {
		logger.Warning(err)
	}

	d.lastore.destroy()

	err = service.StopExport(d.lastore)
	if err != nil {
		logger.Warning(err)
	}

	d.lastore = nil
	return nil
}
