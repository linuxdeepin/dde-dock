package dock

import (
	"pkg.linuxdeepin.com/lib/dbus"
	"sync"
)

type coordinateRange struct {
	X0 int32
	Y0 int32
	X1 int32
	Y1 int32
}

type XMouseAreaInterface interface {
	ConnectCursorInto(func(int32, int32, string)) func()
	ConnectCursorOut(func(int32, int32, string)) func()
	UnregisterArea(string) error
	RegisterAreas(interface{}, int32) (string, error)
	RegisterFullScreen() (string, error)
}

type XMouseAreaProxyer struct {
	lock    sync.RWMutex
	area    XMouseAreaInterface
	areaId  string
	idValid bool
}

func (a *XMouseAreaProxyer) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		"com.deepin.daemon.Dock",
		"/dde/dock/XMouseAreaProxyer",
		"dde.dock.XMouseAreaProxyer",
	}
}

func NewXMouseAreaProxyer(area XMouseAreaInterface, err error) (*XMouseAreaProxyer, error) {
	if err != nil {
		return nil, err
	}
	return &XMouseAreaProxyer{area: area, idValid: false}, nil
}

func (a *XMouseAreaProxyer) connectMotionInto(callback func(int32, int32, string)) func() {
	return a.area.ConnectCursorInto(func(x, y int32, id string) {
		a.lock.Lock()
		defer a.lock.Unlock()
		if !a.idValid || id != a.areaId {
			logger.Warningf("valid: %v, event id: %v, areaId: %v", a.idValid, id, a.areaId)
			return
		}
		callback(x, y, id)
	})
}

func (a *XMouseAreaProxyer) connectMotionOut(callback func(int32, int32, string)) func() {
	return a.area.ConnectCursorOut(func(x, y int32, id string) {
		a.lock.Lock()
		defer a.lock.Unlock()
		if !a.idValid || id != a.areaId {
			logger.Warningf("valid: %v, event id: %v, areaId: %v", a.idValid, id, a.areaId)
			return
		}
		callback(x, y, id)
	})
}

func (a *XMouseAreaProxyer) unregister() {
	if a.idValid {
		a.area.UnregisterArea(a.areaId)
		a.idValid = false
	}
}

func (a *XMouseAreaProxyer) registerArea(registerHandler func() (string, error)) {
	a.lock.Lock()
	defer a.lock.Unlock()

	newAreaId, err := registerHandler()
	if err != nil {
		logger.Warning("register mousearea failed:", err)
		return
	}

	if a.areaId != newAreaId {
		a.unregister()
	}
	a.idValid = true
	a.areaId = newAreaId
}

func (a *XMouseAreaProxyer) RegisterAreas(areas []coordinateRange, eventMask int32) {
	a.registerArea(func() (string, error) {
		return a.area.RegisterAreas(areas, eventMask)
	})
}

func (a *XMouseAreaProxyer) RegisterFullScreen() {
	a.registerArea(a.area.RegisterFullScreen)
}
