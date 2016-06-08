package dock

import (
	"github.com/BurntSushi/xgb/xproto"
	"pkg.deepin.io/lib/dbus"
	"reflect"
)

type propertyFrontendWindow struct {
	*propertyBase
	val xproto.Window
}

func newPropertyFrontendWindow(obj dbus.DBusObject) *propertyFrontendWindow {
	return &propertyFrontendWindow{
		propertyBase: newPropertyBase(obj, "FrontendWindow"),
	}
}

func (p *propertyFrontendWindow) Get() xproto.Window {
	return p.val
}

func (p *propertyFrontendWindow) GetType() reflect.Type {
	return reflect.TypeOf(uint32(0))
}

func (p *propertyFrontendWindow) GetValue() interface{} {
	return uint32(p.val)
}

func (p *propertyFrontendWindow) SetValue(val interface{}) {
	newVal := xproto.Window(val.(uint32))
	if p.val != newVal {
		p.val = newVal
		p.Notify()
	}
}
