package power

import (
	libdisplay "dbus/com/deepin/daemon/display"
	libsessionwatcher "dbus/com/deepin/daemon/sessionwatcher"
	liblockfront "dbus/com/deepin/dde/lockfront"
	libsessionmanager "dbus/com/deepin/sessionmanager"
	libpower "dbus/com/deepin/system/power"
	liblogin1 "dbus/org/freedesktop/login1"
	libnotifications "dbus/org/freedesktop/notifications"
	libscreensaver "dbus/org/freedesktop/screensaver"

	"github.com/BurntSushi/xgb/dpms"
	"github.com/BurntSushi/xgbutil"
)

type Helper struct {
	Power          *libpower.Power
	Notifier       *libnotifications.Notifier
	SessionManager *libsessionmanager.SessionManager
	SessionWatcher *libsessionwatcher.SessionWatcher
	ScreenSaver    *libscreensaver.ScreenSaver
	Display        *libdisplay.Display
	LockFront      *liblockfront.LockFront
	Login1Manager  *liblogin1.Manager

	xu *xgbutil.XUtil
}

func NewHelper() (*Helper, error) {
	h := &Helper{}
	err := h.init()
	if err != nil {
		return nil, err
	}
	return h, nil
}

func (h *Helper) init() error {
	var err error
	h.Power, err = libpower.NewPower("com.deepin.system.Power", "/com/deepin/system/Power")
	if err != nil {
		logger.Warning("init Power failed:", err)
		return err
	}

	h.Notifier, err = libnotifications.NewNotifier("org.freedesktop.Notifications", "/org/freedesktop/Notifications")
	if err != nil {
		logger.Warning("init Notifier failed:", err)
		return err
	}

	h.SessionManager, err = libsessionmanager.NewSessionManager("com.deepin.SessionManager", "/com/deepin/SessionManager")
	if err != nil {
		logger.Warning("init SessionManager failed:", err)
		return err
	}

	h.ScreenSaver, err = libscreensaver.NewScreenSaver("org.freedesktop.ScreenSaver", "/org/freedesktop/ScreenSaver")
	if err != nil {
		logger.Warning("init ScreenSaver failed:", err)
		return err
	}

	h.Display, err = libdisplay.NewDisplay(dbusDisplayDest, dbusDisplayPath)
	if err != nil {
		logger.Warning("init Display failed:", err)
		return err
	}

	h.LockFront, err = liblockfront.NewLockFront("com.deepin.dde.lockFront", "/com/deepin/dde/lockFront")
	if err != nil {
		logger.Warning("init LockFront failed:", err)
		return err
	}

	h.SessionWatcher, err = libsessionwatcher.NewSessionWatcher("com.deepin.daemon.SessionWatcher", "/com/deepin/daemon/SessionWatcher")
	if err != nil {
		logger.Warning("init SessionWatcher failed:", err)
		return err
	}

	h.Login1Manager, err = liblogin1.NewManager("org.freedesktop.login1", "/org/freedesktop/login1")
	if err != nil {
		logger.Warning("init login1 manager failed:", err)
		return err
	}

	// init X conn
	h.xu, err = xgbutil.NewConn()
	if err != nil {
		return err
	}
	dpms.Init(h.xu.Conn())
	return nil
}

func (h *Helper) Destroy() {
	if h.Power != nil {
		libpower.DestroyPower(h.Power)
		h.Power = nil
	}

	if h.Notifier != nil {
		libnotifications.DestroyNotifier(h.Notifier)
		h.Notifier = nil
	}

	if h.SessionManager != nil {
		libsessionmanager.DestroySessionManager(h.SessionManager)
		h.SessionManager = nil
	}

	if h.ScreenSaver != nil {
		libscreensaver.DestroyScreenSaver(h.ScreenSaver)
		h.ScreenSaver = nil
	}

	if h.Display != nil {
		libdisplay.DestroyDisplay(h.Display)
		h.Display = nil
	}

	if h.LockFront != nil {
		h.LockFront = nil
	}

	if h.SessionWatcher != nil {
		libsessionwatcher.DestroySessionWatcher(h.SessionWatcher)
		h.SessionWatcher = nil
	}

	if h.Login1Manager != nil {
		liblogin1.DestroyManager(h.Login1Manager)
		h.Login1Manager = nil
	}

	// NOTE: Don't close x conn, because the bug of lib xgbutil.
	// [xgbutil] eventloop.go:27: BUG: Could not read an event or an error.
	if h.xu != nil {
		//h.xu.Conn().Close()
		h.xu = nil
	}
}
