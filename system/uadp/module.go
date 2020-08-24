package uadp

import (
	"pkg.deepin.io/dde/daemon/loader"
)

func init() {
	loader.Register(newModule())
}

type Module struct {
	uadp *Uadp
	*loader.ModuleBase
}

func newModule() *Module {
	m := new(Module)
	m.ModuleBase = loader.NewModuleBase("Uadp", m, logger)
	return m
}

func (m *Module) GetDependencies() []string {
	return []string{}
}

func (m *Module) Start() error {
	service := loader.GetService()

	if m.uadp != nil {
		return nil
	}
	var err error
	m.uadp, err = newUadp(service)
	if err != nil {
		logger.Warning("failed to newUadp:", err)
		return err
	}

	err = service.Export(dbusPath, m.uadp)
	if err != nil {
		logger.Warning("failed to Export uadp:", err)
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		logger.Warning("failed to RequestName:", err)
		return err
	}

	return nil
}

func (m *Module) Stop() error {
	if m.uadp == nil {
		return nil
	}

	service := loader.GetService()
	err := service.ReleaseName(dbusServiceName)
	if err != nil {
		logger.Warning("failed to releaseName:", err)
	}

	err = service.StopExport(m.uadp)
	if err != nil {
		logger.Warning("failed to stopExport:", err)
	}
	m.uadp = nil

	return nil
}
