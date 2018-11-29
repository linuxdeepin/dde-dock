package trayicon

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	ofdbus "github.com/linuxdeepin/go-dbus-factory/org.freedesktop.dbus"
	"pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
	"pkg.deepin.io/lib/strv"
)

type StatusNotifierWatcher struct {
	service    *dbusutil.Service
	sigLoop    *dbusutil.SignalLoop
	dbusDaemon *ofdbus.DBus

	hostServiceName string
	watchedServices strv.Strv
	PropsMu         sync.RWMutex
	// dbusutil-gen: equal=nil
	RegisteredStatusNotifierItems  strv.Strv
	IsStatusNotifierHostRegistered bool

	// dbusutil-gen: ignore
	ProtocolVersion int32

	methods *struct {
		RegisterStatusNotifierItem func() `in:"serviceName"`
		RegisterStatusNotifierHost func() `in:"serviceName"`
	}

	signals *struct {
		StatusNotifierItemRegistered struct {
			ServiceName string
		}
		StatusNotifierItemUnregistered struct {
			ServiceName string
		}
		StatusNotifierHostRegistered struct{}
	}
}

func newStatusNotifierWatcher(service *dbusutil.Service,
	sigLoop *dbusutil.SignalLoop) *StatusNotifierWatcher {
	snw := &StatusNotifierWatcher{
		service: service,
		sigLoop: sigLoop,
	}

	sessionBus := service.Conn()
	snw.dbusDaemon = ofdbus.NewDBus(sessionBus)
	return snw
}

const (
	snwDBusPath        = "/StatusNotifierWatcher"
	snwDBusServiceName = "org.kde.StatusNotifierWatcher"
)

func (*StatusNotifierWatcher) GetInterfaceName() string {
	return snwDBusServiceName
}

func (snw *StatusNotifierWatcher) isDBusServiceRegistered(serviceName string) bool {
	owner, err := snw.dbusDaemon.GetNameOwner(0, serviceName)
	if err != nil {
		logger.Warning(err)
		return false
	}
	return owner != ""
}

func (snw *StatusNotifierWatcher) RegisterStatusNotifierItem(sender dbus.Sender, serviceOrPath string) *dbus.Error {
	logger.Debug("RegisterStatusNotifierItem", serviceOrPath)

	var serviceName string
	var objPath string

	if strings.HasPrefix(serviceOrPath, "/") {
		// is path
		serviceName = string(sender)
		objPath = serviceOrPath
	} else {
		// is service name
		serviceName = serviceOrPath
		objPath = "/StatusNotifierItem"
	}

	if !snw.isDBusServiceRegistered(serviceName) {
		return dbusutil.ToError(fmt.Errorf("dbus service %q not registered", serviceName))
	}

	notifierItemId := serviceName + objPath

	snw.PropsMu.Lock()
	defer snw.PropsMu.Unlock()

	if snw.RegisteredStatusNotifierItems.Contains(notifierItemId) {
		return dbusutil.ToError(errors.New("notifier item has been registered"))
	}

	snw.watchedServices, _ = snw.watchedServices.Add(serviceName)
	newItems, _ := snw.RegisteredStatusNotifierItems.Add(notifierItemId)
	snw.setPropRegisteredStatusNotifierItems(newItems)

	snw.service.Emit(snw, "StatusNotifierItemRegistered", notifierItemId)

	return nil
}

func (snw *StatusNotifierWatcher) RegisterStatusNotifierHost(serviceName string) *dbus.Error {
	logger.Debug("RegisterStatusNotifierHost", serviceName)

	snw.PropsMu.Lock()
	defer snw.PropsMu.Unlock()

	if snw.IsStatusNotifierHostRegistered {
		return dbusutil.ToError(errors.New("host has been registered"))
	}

	snw.setPropIsStatusNotifierHostRegistered(true)
	snw.hostServiceName = serviceName

	snw.service.Emit(snw, "StatusNotifierHostRegistered")

	return nil
}

func (ss *StatusNotifierWatcher) listenDBusNameOwnerChanged() {
	ss.dbusDaemon.InitSignalExt(ss.sigLoop, true)
	ss.dbusDaemon.ConnectNameOwnerChanged(func(name string, oldOwner string, newOwner string) {
		ss.PropsMu.Lock()

		if newOwner == "" {

			if ss.hostServiceName == name {
				logger.Infof("host %s lost", name)
				ss.hostServiceName = ""
				ss.setPropIsStatusNotifierHostRegistered(false)

			} else if ss.watchedServices.Contains(name) {
				logger.Infof("item %s lost", name)

				ss.watchedServices, _ = ss.watchedServices.Delete(name)
				match := name + "/"
				newItems := make(strv.Strv, 0, len(ss.RegisteredStatusNotifierItems))
				for _, itemID := range ss.RegisteredStatusNotifierItems {
					if strings.HasPrefix(itemID, match) {
						ss.service.Emit(ss, "StatusNotifierItemUnregistered", itemID)
					} else {
						newItems = append(newItems, itemID)
					}
				}
				ss.setPropRegisteredStatusNotifierItems(newItems)
			}
		}

		ss.PropsMu.Unlock()
	})
}
