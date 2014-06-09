/*This file is auto generate by dlib/dbus/proxyer. Don't edit it*/
package network

import "dlib/dbus"
import "dlib/dbus/property"
import "reflect"
import "sync"
import "runtime"
import "fmt"
import "errors"

/*prevent compile error*/
var _ = fmt.Println
var _ = runtime.SetFinalizer
var _ = sync.NewCond
var _ = reflect.TypeOf
var _ = property.BaseObserver{}

type SystemService struct {
	Path     dbus.ObjectPath
	DestName string
	core     *dbus.Object
}

func (obj SystemService) IsPackageSystemLocked() (arg0 bool, _err error) {
	_err = obj.core.Call("com.ubuntu.SystemService.is_package_system_locked", 0).Store(&arg0)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj SystemService) GetProxy(proxy_type string) (arg1 string, _err error) {
	_err = obj.core.Call("com.ubuntu.SystemService.get_proxy", 0, proxy_type).Store(&arg1)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj SystemService) SetProxy(proxy_type string, new_proxy string) (arg2 bool, _err error) {
	_err = obj.core.Call("com.ubuntu.SystemService.set_proxy", 0, proxy_type, new_proxy).Store(&arg2)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj SystemService) IsRebootRequired() (arg0 bool, _err error) {
	_err = obj.core.Call("com.ubuntu.SystemService.is_reboot_required", 0).Store(&arg0)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj SystemService) SetNoProxy(new_no_proxy string) (arg1 bool, _err error) {
	_err = obj.core.Call("com.ubuntu.SystemService.set_no_proxy", 0, new_no_proxy).Store(&arg1)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj SystemService) SetKeyboard(model string, layout string, variant string, options string) (arg4 bool, _err error) {
	_err = obj.core.Call("com.ubuntu.SystemService.set_keyboard", 0, model, layout, variant, options).Store(&arg4)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func (obj SystemService) GetKeyboard() (arg0 string, arg1 string, arg2 string, arg3 string, _err error) {
	_err = obj.core.Call("com.ubuntu.SystemService.get_keyboard", 0).Store(&arg0, &arg1, &arg2, &arg3)
	if _err != nil {
		fmt.Println(_err)
	}
	return
}

func NewSystemService(destName string, path dbus.ObjectPath) (*SystemService, error) {
	if !path.IsValid() {
		return nil, errors.New("The path of '" + string(path) + "' is invalid.")
	}

	core := getBus().Object(destName, path)

	obj := &SystemService{Path: path, DestName: destName, core: core}

	return obj, nil
}
