package dock

import (
	"pkg.deepin.io/lib/dbus"
	"reflect"
)

type propertyHideState struct {
	*propertyBase
	val HideStateType
}

func newPropertyHideState(obj dbus.DBusObject) *propertyHideState {
	return &propertyHideState{
		propertyBase: newPropertyBase(obj, "HideState"),
		val:          HideStateUnknown,
	}
}

func (p *propertyHideState) Get() HideStateType {
	return p.val
}

func (p *propertyHideState) GetType() reflect.Type {
	return reflect.TypeOf(int32(0))
}

func (p *propertyHideState) GetValue() interface{} {
	return int32(p.val)
}

func (p *propertyHideState) SetValue(val interface{}) {
	// TODO: valid val
	newVal := HideStateType(val.(int32))
	if p.val != newVal {
		p.val = newVal
		p.Notify()
	}
}
