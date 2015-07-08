package power

import "pkg.deepin.io/dde/daemon/loader"
import "pkg.deepin.io/lib/dbus"
import "pkg.deepin.io/lib/log"

import libupower "dbus/org/freedesktop/upower"
import liblogin1 "dbus/org/freedesktop/login1"
import libkeybinding "dbus/com/deepin/daemon/keybinding"
import libnotifications "dbus/org/freedesktop/notifications"
import libsound "dbus/com/deepin/api/sound"

func init() {
	loader.Register(NewDaemon(logger))
}

var (
	notifier *libnotifications.Notifier
	upower   *libupower.Upower
	login1   *liblogin1.Manager
	mediaKey *libkeybinding.MediaKey
	player   *libsound.Sound

	power *Power
)

func initializeLibs() error {
	var err error
	upower, err = libupower.NewUpower(UPOWER_BUS_NAME, "/org/freedesktop/UPower")
	if err != nil {
		logger.Warning("create dbus upower failed:", err)
		return err
	}
	login1, err = liblogin1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		logger.Warning("create dbus login1 failed:", err)
		finalizeLibs()
		return err
	}
	mediaKey, err = libkeybinding.NewMediaKey("com.deepin.daemon.KeyBinding", "/com/deepin/daemon/MediaKey")
	if err != nil {
		logger.Warning("create dbus mediaKey failed:", err)
		finalizeLibs()
		return err
	}
	notifier, err = libnotifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if err != nil {
		logger.Warning("Can't build org.freedesktop.Notficaations:", err)
		finalizeLibs()
		return err
	}
	player, err = libsound.NewSound("com.deepin.api.Sound", "/com/deepin/api/Sound")
	if err != nil {
		logger.Warning("Can't build com.deepin.api.Sound:", err)
		finalizeLibs()
		return err
	}

	power = NewPower()
	return nil
}

func finalizeLibs() {
	if power != nil {
		power.batGroup.Destroy()
		power.batGroup = nil
		dbus.UnInstallObject(power)
		power = nil
	}
	if upower != nil {
		libupower.DestroyUpower(upower)
		upower = nil
	}
	if login1 != nil {
		liblogin1.DestroyManager(login1)
		login1 = nil
	}
	if mediaKey != nil {
		libkeybinding.DestroyMediaKey(mediaKey)
		mediaKey = nil
	}
	if notifier != nil {
		libnotifications.DestroyNotifier(notifier)
		notifier = nil
	}

	player = nil
}

var workaround *fullScreenWorkaround

type Daemon struct {
	*loader.ModuleBase
}

func NewDaemon(logger *log.Logger) *Daemon {
	daemon := new(Daemon)
	daemon.ModuleBase = loader.NewModuleBase("power", daemon, logger)
	return daemon
}

func (d *Daemon) GetDependencies() []string {
	return []string{"screensaver"}
}

func (d *Daemon) Start() error {
	if power != nil {
		return nil
	}

	logger.BeginTracing()

	err := initializeLibs()
	if err != nil {
		logger.Error(err)
		logger.EndTracing()
		return err
	}

	err = dbus.InstallOnSession(power)
	if err != nil {
		logger.Error("Failed InstallOnSession:", err)
		finalizeLibs()
		logger.EndTracing()
		return err
	}

	workaround = newFullScreenWorkaround()
	go workaround.start()
	return nil
}

func (d *Daemon) Stop() error {
	if power == nil {
		return nil
	}

	if workaround != nil {
		workaround.stop()
		workaround = nil
	}
	finalizeLibs()
	logger.EndTracing()
	return nil
}
