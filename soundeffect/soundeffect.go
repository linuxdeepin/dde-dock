package soundeffect

import (
	. "pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger = log.NewLogger("daemon/soundeffect")

type Daemon struct {
	*ModuleBase
}

func init() {
	Register(NewSoundEffectDaemon(logger))
}

func NewSoundEffectDaemon(logger *log.Logger) *Daemon{
	var d = new(Daemon)
	d.ModuleBase = NewModuleBase("soundeffect", d, logger)
	return d
}

func (*Daemon) GetDependencies() []string {
	return []string{}
}

var _manager *Manager

func (*Daemon) Start() error {
	if _manager != nil {
		return nil
	}

	logger.BeginTracing()
	var err error
	_manager, err = NewManager()
	if err != nil {
		logger.Error("New Manager failed:", err)
		return err
	}

	err = dbus.InstallOnSession(_manager)
	if err != nil {
		logger.Error("Install session bus failed:", err)
		return err
	}
	_manager.handleGSetting()

	return nil
}

func (*Daemon) Stop() error{
	if _manager == nil {
		return nil
	}

	_manager.setting.Unref()
	_manager = nil
	logger.EndTracing()
	return nil
}
