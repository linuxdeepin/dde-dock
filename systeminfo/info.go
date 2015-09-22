package systeminfo

import (
	"pkg.deepin.io/dde/daemon/loader"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

type SystemInfo struct {
	// Current version, ex: "2015 Desktop"
	Version string
	// CPU information
	Processor string
	// Disk capacity
	DiskCap uint64
	// Memory size
	MemoryCap uint64
	// System architecture
	SystemType int64
}

type Daemon struct {
	info *SystemInfo
	*loader.ModuleBase
}

var logger = log.NewLogger("daemon/systeminfo")

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("systeminfo", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{}
}

func (d *Daemon) Start() error {
	if d.info != nil {
		return nil
	}

	logger.BeginTracing()

	d.info = NewSystemInfo()
	err := dbus.InstallOnSession(d.info)
	if err != nil {
		d.info = nil
		logger.Error(err)
		logger.EndTracing()
		return err
	}
	return nil
}

func (d *Daemon) Stop() error {
	if d.info == nil {
		return nil
	}

	dbus.UnInstallObject(d.info)
	d.info = nil
	logger.EndTracing()

	return nil
}

func NewSystemInfo() *SystemInfo {
	var (
		info SystemInfo
		err  error
	)

	info.Processor, err = getCPUInfo("/proc/cpuinfo")
	if err != nil {
		logger.Warning("Get cpu info failed:", err)
	}

	info.Version, err = getVersion()
	if err != nil {
		logger.Warning("Get version failed:", err)
	}

	info.MemoryCap, err = getMemoryFromFile("/proc/meminfo")
	if err != nil {
		logger.Warning("Get memory capacity failed:", err)
	}

	info.SystemType, err = getOSType()
	if err != nil {
		logger.Warning("Get os type failed:", err)
	}

	info.DiskCap, err = getDiskCap()
	if err != nil {
		logger.Warning("Get disk capacity failed:", err)
	}

	return &info
}

func (*SystemInfo) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.daemon.SystemInfo",
		ObjectPath: "/com/deepin/daemon/SystemInfo",
		Interface:  "com.deepin.daemon.SystemInfo",
	}
}
