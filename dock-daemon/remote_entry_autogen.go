/*This file is auto generate by dlib/dbus/proxyer. Don't edit it*/
package main

import "dlib/dbus"
import "dlib/dbus/property"
import "reflect"
import "sync"
import "runtime"
import "fmt"
import "errors"
import "strings"

/*prevent compile error*/
var _ = fmt.Println
var _ = runtime.SetFinalizer
var _ = sync.NewCond
var _ = reflect.TypeOf
var _ = property.BaseObserver{}

type RemoteEntry struct {
	Path     dbus.ObjectPath
	DestName string
	core     *dbus.Object

	signals       map[chan *dbus.Signal]bool
	signalsLocker sync.Mutex

	ID                 *dbusPropertyRemoteEntryID
	Type               *dbusPropertyRemoteEntryType
	Tooltip            *dbusPropertyRemoteEntryTooltip
	Icon               *dbusPropertyRemoteEntryIcon
	Status             *dbusPropertyRemoteEntryStatus
	QuickWindowVieable *dbusPropertyRemoteEntryQuickWindowVieable
	Allocation         *dbusPropertyRemoteEntryAllocation
}

func (obj RemoteEntry) _createSignalChan() chan *dbus.Signal {
	obj.signalsLocker.Lock()
	ch := make(chan *dbus.Signal, 30)
	getBus().Signal(ch)
	obj.signals[ch] = false
	obj.signalsLocker.Unlock()
	return ch
}
func (obj RemoteEntry) _deleteSignalChan(ch chan *dbus.Signal) {
	obj.signalsLocker.Lock()
	delete(obj.signals, ch)
	getBus().DetachSignal(ch)
	close(ch)
	obj.signalsLocker.Unlock()
}
func DestroyRemoteEntry(obj *RemoteEntry) {
	obj.signalsLocker.Lock()
	for ch, _ := range obj.signals {
		getBus().DetachSignal(ch)
		close(ch)
	}
	obj.signals = make(map[chan *dbus.Signal]bool)
	obj.signalsLocker.Unlock()
}

func (obj RemoteEntry) Activate(arg0 int32, arg1 int32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.Activate", 0, arg0, arg1).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj RemoteEntry) ContextMenu(arg0 int32, arg1 int32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.ContextMenu", 0, arg0, arg1).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj RemoteEntry) OnDragDrop(arg0 int32, arg1 int32, arg2 string) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.OnDragDrop", 0, arg0, arg1, arg2).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj RemoteEntry) OnDragEnter(arg0 int32, arg1 int32, arg2 string) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.OnDragEnter", 0, arg0, arg1, arg2).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj RemoteEntry) OnDragLeave(arg0 int32, arg1 int32, arg2 string) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.OnDragLeave", 0, arg0, arg1, arg2).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj RemoteEntry) OnDragOver(arg0 int32, arg1 int32, arg2 string) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.OnDragOver", 0, arg0, arg1, arg2).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj RemoteEntry) QuickWindow(arg0 int32, arg1 int32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.QuickWindow", 0, arg0, arg1).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj RemoteEntry) SecondaryActivate(arg0 int32, arg1 int32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.SecondaryActivate", 0, arg0, arg1).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

type dbusPropertyRemoteEntryID struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryID) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.ID is not writable")
}

func (this *dbusPropertyRemoteEntryID) Get() string {
	return this.GetValue().(string)
}
func (this *dbusPropertyRemoteEntryID) GetValue() interface{} /*string*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "ID").Store(&r)
	if err == nil && r.Signature().String() == "s" {
		return r.Value().(string)
	} else {
		fmt.Println("dbusProperty:ID error:", err, "at dde.dock.Entry")
		return *new(string)
	}
}
func (this *dbusPropertyRemoteEntryID) GetType() reflect.Type {
	return reflect.TypeOf((*string)(nil)).Elem()
}

type dbusPropertyRemoteEntryType struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryType) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.Type is not writable")
}

func (this *dbusPropertyRemoteEntryType) Get() string {
	return this.GetValue().(string)
}
func (this *dbusPropertyRemoteEntryType) GetValue() interface{} /*string*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "Type").Store(&r)
	if err == nil && r.Signature().String() == "s" {
		return r.Value().(string)
	} else {
		fmt.Println("dbusProperty:Type error:", err, "at dde.dock.Entry")
		return *new(string)
	}
}
func (this *dbusPropertyRemoteEntryType) GetType() reflect.Type {
	return reflect.TypeOf((*string)(nil)).Elem()
}

type dbusPropertyRemoteEntryTooltip struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryTooltip) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.Tooltip is not writable")
}

func (this *dbusPropertyRemoteEntryTooltip) Get() string {
	return this.GetValue().(string)
}
func (this *dbusPropertyRemoteEntryTooltip) GetValue() interface{} /*string*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "Tooltip").Store(&r)
	if err == nil && r.Signature().String() == "s" {
		return r.Value().(string)
	} else {
		fmt.Println("dbusProperty:Tooltip error:", err, "at dde.dock.Entry")
		return *new(string)
	}
}
func (this *dbusPropertyRemoteEntryTooltip) GetType() reflect.Type {
	return reflect.TypeOf((*string)(nil)).Elem()
}

type dbusPropertyRemoteEntryIcon struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryIcon) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.Icon is not writable")
}

func (this *dbusPropertyRemoteEntryIcon) Get() string {
	return this.GetValue().(string)
}
func (this *dbusPropertyRemoteEntryIcon) GetValue() interface{} /*string*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "Icon").Store(&r)
	if err == nil && r.Signature().String() == "s" {
		return r.Value().(string)
	} else {
		fmt.Println("dbusProperty:Icon error:", err, "at dde.dock.Entry")
		return *new(string)
	}
}
func (this *dbusPropertyRemoteEntryIcon) GetType() reflect.Type {
	return reflect.TypeOf((*string)(nil)).Elem()
}

type dbusPropertyRemoteEntryStatus struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryStatus) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.Status is not writable")
}

func (this *dbusPropertyRemoteEntryStatus) Get() int32 {
	return this.GetValue().(int32)
}
func (this *dbusPropertyRemoteEntryStatus) GetValue() interface{} /*int32*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "Status").Store(&r)
	if err == nil && r.Signature().String() == "i" {
		return r.Value().(int32)
	} else {
		fmt.Println("dbusProperty:Status error:", err, "at dde.dock.Entry")
		return *new(int32)
	}
}
func (this *dbusPropertyRemoteEntryStatus) GetType() reflect.Type {
	return reflect.TypeOf((*int32)(nil)).Elem()
}

type dbusPropertyRemoteEntryQuickWindowVieable struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryQuickWindowVieable) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.QuickWindowVieable is not writable")
}

func (this *dbusPropertyRemoteEntryQuickWindowVieable) Get() bool {
	return this.GetValue().(bool)
}
func (this *dbusPropertyRemoteEntryQuickWindowVieable) GetValue() interface{} /*bool*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "QuickWindowVieable").Store(&r)
	if err == nil && r.Signature().String() == "b" {
		return r.Value().(bool)
	} else {
		fmt.Println("dbusProperty:QuickWindowVieable error:", err, "at dde.dock.Entry")
		return *new(bool)
	}
}
func (this *dbusPropertyRemoteEntryQuickWindowVieable) GetType() reflect.Type {
	return reflect.TypeOf((*bool)(nil)).Elem()
}

type dbusPropertyRemoteEntryAllocation struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryAllocation) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.Allocation is not writable")
}

func (this *dbusPropertyRemoteEntryAllocation) Get() []interface{} {
	return this.GetValue().([]interface{})
}
func (this *dbusPropertyRemoteEntryAllocation) GetValue() interface{} /*[]interface {}*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "Allocation").Store(&r)
	if err == nil && r.Signature().String() == "(nnqq)" {
		return r.Value().([]interface{})
	} else {
		fmt.Println("dbusProperty:Allocation error:", err, "at dde.dock.Entry")
		return *new([]interface{})
	}
}
func (this *dbusPropertyRemoteEntryAllocation) GetType() reflect.Type {
	return reflect.TypeOf((*[]interface{})(nil)).Elem()
}

func NewRemoteEntry(destName string, path dbus.ObjectPath) (*RemoteEntry, error) {
	if !path.IsValid() {
		return nil, errors.New("The path of '" + string(path) + "' is invalid.")
	}

	core := getBus().Object(destName, path)
	var v string
	core.Call("org.freedesktop.DBus.Introspectable.Introspect", 0).Store(&v)
	if strings.Index(v, "dde.dock.Entry") == -1 {
		return nil, errors.New("'" + string(path) + "' hasn't interface 'dde.dock.Entry'.")
	}

	obj := &RemoteEntry{Path: path, DestName: destName, core: core, signals: make(map[chan *dbus.Signal]bool)}

	obj.ID = &dbusPropertyRemoteEntryID{&property.BaseObserver{}, core}
	obj.Type = &dbusPropertyRemoteEntryType{&property.BaseObserver{}, core}
	obj.Tooltip = &dbusPropertyRemoteEntryTooltip{&property.BaseObserver{}, core}
	obj.Icon = &dbusPropertyRemoteEntryIcon{&property.BaseObserver{}, core}
	obj.Status = &dbusPropertyRemoteEntryStatus{&property.BaseObserver{}, core}
	obj.QuickWindowVieable = &dbusPropertyRemoteEntryQuickWindowVieable{&property.BaseObserver{}, core}
	obj.Allocation = &dbusPropertyRemoteEntryAllocation{&property.BaseObserver{}, core}

	getBus().BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "type='signal',path='"+string(path)+"',interface='org.freedesktop.DBus.Properties',sender='"+destName+"',member='PropertiesChanged'")
	getBus().BusObject().Call("org.freedesktop.DBus.AddMatch", 0, "type='signal',path='"+string(path)+"',interface='dde.dock.Entry',sender='"+destName+"',member='PropertiesChanged'")
	sigChan := obj._createSignalChan()
	go func() {
		typeString := reflect.TypeOf("")
		typeKeyValues := reflect.TypeOf(map[string]dbus.Variant{})
		typeArrayValues := reflect.TypeOf([]string{})
		for v := range sigChan {
			if v.Name == "org.freedesktop.DBus.Properties.PropertiesChanged" &&
				len(v.Body) == 3 &&
				reflect.TypeOf(v.Body[0]) == typeString &&
				reflect.TypeOf(v.Body[1]) == typeKeyValues &&
				reflect.TypeOf(v.Body[2]) == typeArrayValues &&
				v.Body[0].(string) != "dde.dock.Entry" {
				props := v.Body[1].(map[string]dbus.Variant)
				for key, _ := range props {
					if false {
					} else if key == "ID" {
						obj.ID.Notify()

					} else if key == "Type" {
						obj.Type.Notify()

					} else if key == "Tooltip" {
						obj.Tooltip.Notify()

					} else if key == "Icon" {
						obj.Icon.Notify()

					} else if key == "Status" {
						obj.Status.Notify()

					} else if key == "QuickWindowVieable" {
						obj.QuickWindowVieable.Notify()

					} else if key == "Allocation" {
						obj.Allocation.Notify()
					}
				}
			} else if v.Name == "dde.dock.Entry.PropertiesChanged" && len(v.Body) == 1 && reflect.TypeOf(v.Body[0]) == typeKeyValues {
				for key, _ := range v.Body[0].(map[string]dbus.Variant) {
					if false {
					} else if key == "ID" {
						obj.ID.Notify()

					} else if key == "Type" {
						obj.Type.Notify()

					} else if key == "Tooltip" {
						obj.Tooltip.Notify()

					} else if key == "Icon" {
						obj.Icon.Notify()

					} else if key == "Status" {
						obj.Status.Notify()

					} else if key == "QuickWindowVieable" {
						obj.QuickWindowVieable.Notify()

					} else if key == "Allocation" {
						obj.Allocation.Notify()
					}
				}
			}
		}
	}()

	runtime.SetFinalizer(obj, func(_obj *RemoteEntry) { DestroyRemoteEntry(_obj) })
	return obj, nil
}

var __conn *dbus.Conn = nil

func getBus() *dbus.Conn {
	if __conn == nil {
		var err error
		__conn, err = dbus.SessionBus()
		if err != nil {
			panic(err)
		}
	}
	return __conn
}
