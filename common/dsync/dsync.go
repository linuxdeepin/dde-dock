package dsync

import (
	"encoding/json"

	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/dbusutil/proxy"
	"pkg.deepin.io/lib/log"
	"pkg.deepin.io/lib/strv"
)

type Interface interface {
	Get() (interface{}, error)
	Set(data []byte) error
}

type Config struct {
	name       string
	core       Interface
	dbusDaemon *ofdbus.DBus
	path       dbus.ObjectPath
	sigLoop    *dbusutil.SignalLoop
	logger     *log.Logger
	methods    *struct {
		Get func() `out:"data"`
		Set func() `in:"data"`
	}
}

const (
	serviceName = "com.deepin.sync.Daemon"
	servicePath = "/com/deepin/sync/Daemon"
)

func NewConfig(name string, core Interface, sessionSigLoop *dbusutil.SignalLoop,
	path dbus.ObjectPath, logger *log.Logger) *Config {
	c := &Config{
		name:    name,
		core:    core,
		sigLoop: sessionSigLoop,
		path:    path,
		logger:  logger,
	}

	sessionBus := sessionSigLoop.Conn()
	c.dbusDaemon = ofdbus.NewDBus(sessionBus)
	c.dbusDaemon.InitSignalExt(sessionSigLoop, true)
	_, err := c.dbusDaemon.ConnectNameOwnerChanged(func(name string, oldOwner string, newOwner string) {
		if name == "com.deepin.sync.Daemon" && newOwner != "" {
			err := c.Register()
			if err != nil {
				c.logger.Warning(err)
			}
		}
	})
	if err != nil {
		logger.Warning(err)
	}
	return c
}

func (c *Config) Register() error {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	availableServices, err := ofdbus.NewDBus(sessionBus).ListNames(0)
	if err != nil {
		return err
	}
	if !strv.Strv(availableServices).Contains(serviceName) {
		c.logger.Debug("sync daemon not exists")
		return nil
	}

	obj := sessionBus.Object(serviceName, servicePath)
	err = obj.Call("com.deepin.sync.Daemon.Register", 0, c.name, c.path).Err
	return err
}

func (c *Config) Destroy() {
	c.dbusDaemon.RemoveHandler(proxy.RemoveAllHandlers)
}

func (*Config) GetInterfaceName() string {
	return "com.deepin.sync.Config"
}

func (c *Config) Get() ([]byte, *dbus.Error) {
	v, err := c.core.Get()
	if err != nil {
		return nil, dbusutil.ToError(err)
	}

	data, err := json.Marshal(v)
	if err != nil {
		return nil, dbusutil.ToError(err)
	}
	return data, nil
}

func (c *Config) Set(data []byte) *dbus.Error {
	err := c.core.Set(data)
	return dbusutil.ToError(err)
}
