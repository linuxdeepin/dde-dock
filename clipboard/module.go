package clipboard

import (
	"github.com/linuxdeepin/go-x11-client"
	"github.com/linuxdeepin/go-x11-client/ext/xfixes"
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/log"
)

const dbusServiceName = "com.deepin.daemon.ClipboardManager"

var logger *log.Logger

func init() {
	logger = log.NewLogger("clipboard")
	loader.Register(newModule())
}

func newModule() *Module {
	m := new(Module)
	m.ModuleBase = loader.NewModuleBase("clipboard", m, logger)
	return m
}

type Module struct {
	*loader.ModuleBase
}

func (*Module) GetDependencies() []string {
	return nil
}

func (*Module) Start() error {
	logger.Debug("clipboard module start")

	xConn, err := x.NewConn()
	if err != nil {
		return err
	}

	initAtoms(xConn)

	_, err = xfixes.QueryVersion(xConn, xfixes.MajorVersion, xfixes.MinorVersion).Reply(xConn)
	if err != nil {
		logger.Warning(err)
	}

	m := &Manager{
		xConn: xConn,
	}

	err = m.start()
	if err != nil {
		return err
	}

	service := loader.GetService()
	err = service.Export("/com/deepin/daemon/ClipboardManager", m)
	if err != nil {
		return err
	}

	err = service.RequestName(dbusServiceName)
	if err != nil {
		return err
	}

	return nil
}

func (*Module) Stop() error {
	return nil
}
