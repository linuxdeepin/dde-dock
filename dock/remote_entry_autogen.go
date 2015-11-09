/*This file is auto generate by pkg.deepin.io/dbus-generator. Don't edit it*/
package dock

import "pkg.deepin.io/lib/dbus"
import "pkg.deepin.io/lib/dbus/property"
import "reflect"
import "sync"
import "runtime"
import "fmt"
import "errors"

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

	signals       map[<-chan *dbus.Signal]struct{}
	signalsLocker sync.Mutex

	Id   *dbusPropertyRemoteEntryId
	Type *dbusPropertyRemoteEntryType
	Data *dbusPropertyRemoteEntryData
}

func (obj *RemoteEntry) _createSignalChan() <-chan *dbus.Signal {
	obj.signalsLocker.Lock()
	ch := getBus().Signal()
	obj.signals[ch] = struct{}{}
	obj.signalsLocker.Unlock()
	return ch
}
func (obj *RemoteEntry) _deleteSignalChan(ch <-chan *dbus.Signal) {
	obj.signalsLocker.Lock()
	delete(obj.signals, ch)
	getBus().DetachSignal(ch)
	obj.signalsLocker.Unlock()
}
func DestroyRemoteEntry(obj *RemoteEntry) {
	obj.signalsLocker.Lock()
	for ch, _ := range obj.signals {
		getBus().DetachSignal(ch)
	}
	obj.signalsLocker.Unlock()

	obj.Id.Reset()
	obj.Type.Reset()
	obj.Data.Reset()
}

func (obj *RemoteEntry) Activate(arg1 int32, arg2 int32, arg3 uint32) (arg0 bool, _err error) {
	_err = obj.core.Call("dde.dock.Entry.Activate", 0, arg1, arg2, arg3).Store(&arg0)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) ContextMenu(arg0 int32, arg1 int32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.ContextMenu", 0, arg0, arg1).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) HandleMenuItem(arg0 string, arg1 uint32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.HandleMenuItem", 0, arg0, arg1).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) HandleDragDrop(arg0 int32, arg1 int32, arg2 string, arg3 uint32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.HandleDragDrop", 0, arg0, arg1, arg2, arg3).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) HandleDragEnter(arg0 int32, arg1 int32, arg2 string, arg3 uint32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.HandleDragEnter", 0, arg0, arg1, arg2, arg3).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) HandleDragLeave(arg0 int32, arg1 int32, arg2 string, arg3 uint32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.HandleDragLeave", 0, arg0, arg1, arg2, arg3).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) HandleDragOver(arg0 int32, arg1 int32, arg2 string, arg3 uint32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.HandleDragOver", 0, arg0, arg1, arg2, arg3).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) HandleMouseWheel(arg0 int32, arg1 int32, arg2 int32, arg3 uint32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.HandleMouseWheel", 0, arg0, arg1, arg2, arg3).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) SecondaryActivate(arg0 int32, arg1 int32, arg2 uint32) (_err error) {
	_err = obj.core.Call("dde.dock.Entry.SecondaryActivate", 0, arg0, arg1, arg2).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) ShowQuickWindow() (_err error) {
	_err = obj.core.Call("dde.dock.Entry.ShowQuickWindow", 0).Store()
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj *RemoteEntry) ConnectDataChanged(callback func(arg0 string, arg1 string)) func() {
	__conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0,
		"type='signal',path='"+string(obj.Path)+"', interface='dde.dock.Entry',sender='"+obj.DestName+"',member='DataChanged'")
	sigChan := obj._createSignalChan()
	go func() {
		for v := range sigChan {
			if v.Path != obj.Path || v.Name != "dde.dock.Entry.DataChanged" || 2 != len(v.Body) {
				continue
			}
			if reflect.TypeOf(v.Body[0]) != reflect.TypeOf((*string)(nil)).Elem() {
				continue
			}
			if reflect.TypeOf(v.Body[1]) != reflect.TypeOf((*string)(nil)).Elem() {
				continue
			}

			callback(v.Body[0].(string), v.Body[1].(string))
		}
	}()
	return func() {
		obj._deleteSignalChan(sigChan)
	}
}

type dbusPropertyRemoteEntryId struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryId) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.Id is not writable")
}

func (this *dbusPropertyRemoteEntryId) Get() string {
	return this.GetValue().(string)
}
func (this *dbusPropertyRemoteEntryId) GetValue() interface{} /*string*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "Id").Store(&r)
	if err == nil && r.Signature().String() == "s" {
		return r.Value().(string)
	} else {
		fmt.Println("dbusProperty:Id error:", err, "at dde.dock.Entry")
		return *new(string)
	}
}
func (this *dbusPropertyRemoteEntryId) GetType() reflect.Type {
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

type dbusPropertyRemoteEntryData struct {
	*property.BaseObserver
	core *dbus.Object
}

func (this *dbusPropertyRemoteEntryData) SetValue(notwritable interface{}) {
	fmt.Println("dde.dock.Entry.Data is not writable")
}

func (this *dbusPropertyRemoteEntryData) Get() map[string]string {
	return this.GetValue().(map[string]string)
}
func (this *dbusPropertyRemoteEntryData) GetValue() interface{} /*map[string]string*/ {
	var r dbus.Variant
	err := this.core.Call("org.freedesktop.DBus.Properties.Get", 0, "dde.dock.Entry", "Data").Store(&r)
	if err == nil && r.Signature().String() == "a{ss}" {
		return r.Value().(map[string]string)
	} else {
		fmt.Println("dbusProperty:Data error:", err, "at dde.dock.Entry")
		return *new(map[string]string)
	}
}
func (this *dbusPropertyRemoteEntryData) GetType() reflect.Type {
	return reflect.TypeOf((*map[string]string)(nil)).Elem()
}

func NewRemoteEntry(destName string, path dbus.ObjectPath) (*RemoteEntry, error) {
	if !path.IsValid() {
		return nil, errors.New("The path of '" + string(path) + "' is invalid.")
	}

	core := getBus().Object(destName, path)

	obj := &RemoteEntry{Path: path, DestName: destName, core: core, signals: make(map[<-chan *dbus.Signal]struct{})}

	obj.Id = &dbusPropertyRemoteEntryId{&property.BaseObserver{}, core}
	obj.Type = &dbusPropertyRemoteEntryType{&property.BaseObserver{}, core}
	obj.Data = &dbusPropertyRemoteEntryData{&property.BaseObserver{}, core}

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
				v.Body[0].(string) == "dde.dock.Entry" {
				props := v.Body[1].(map[string]dbus.Variant)
				for key, _ := range props {
					if false {
					} else if key == "Id" {
						obj.Id.Notify()

					} else if key == "Type" {
						obj.Type.Notify()

					} else if key == "Data" {
						obj.Data.Notify()
					}
				}
			} else if v.Name == "dde.dock.Entry.PropertiesChanged" && len(v.Body) == 1 && reflect.TypeOf(v.Body[0]) == typeKeyValues {
				for key, _ := range v.Body[0].(map[string]dbus.Variant) {
					if false {
					} else if key == "Id" {
						obj.Id.Notify()

					} else if key == "Type" {
						obj.Type.Notify()

					} else if key == "Data" {
						obj.Data.Notify()
					}
				}
			}
		}
	}()

	runtime.SetFinalizer(obj, func(_obj *RemoteEntry) { DestroyRemoteEntry(_obj) })
	return obj, nil
}
