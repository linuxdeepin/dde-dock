package proxychains

import (
	dbus "github.com/godbus/dbus"
	notifications "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.notifications"
	. "pkg.deepin.io/lib/gettext"
)

var (
	notification            *notifications.Notifications
	notifyIconProxyEnabled  = "notification-network-proxy-enabled"
	notifyIconProxyDisabled = "notification-network-proxy-disabled"
)

func init() {
	sessionBus, err := dbus.SessionBus()
	if err != nil {
		notification = nil
		return
	}
	notification = notifications.NewNotifications(sessionBus)
}

func createNotify(appName string) func(string, string, string) {
	var nid uint32 = 0
	return func(icon, summary, body string) {
		if notification == nil {
			logger.Warning("notification is nil")
			logger.Debugf("%s %s %s", icon, summary, body)
			return
		}
		var err error
		nid, err = notification.Notify(0, appName, nid,
			icon, summary, body, nil, nil, -1)
		if err != nil {
			logger.Warning(err)
			return
		}
	}
}

var notify = createNotify("dde-control-center")

func notifyAppProxyEnabled() {
	notify(notifyIconProxyEnabled, Tr("Network"), Tr("Application proxy is set successfully"))
}
func notifyAppProxyEnableFailed() {
	notify(notifyIconProxyDisabled, Tr("Network"), Tr("Failed to set the application proxy"))
}
