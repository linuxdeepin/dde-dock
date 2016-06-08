package dock

import (
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/dbus/property"
)

type propertyBase struct {
	property.BaseObserver
	obj  dbus.DBusObject
	name string
}

func newPropertyBase(obj dbus.DBusObject, name string) *propertyBase {
	return &propertyBase{
		BaseObserver: property.BaseObserver{},
		obj:          obj,
		name:         name,
	}
}

func (p *propertyBase) Notify() {
	dbus.NotifyChange(p.obj, p.name)
	p.BaseObserver.Notify()
}
