package main

// +build ignore

// this file should be rewrite

//#cgo pkg-config: glib-2.0 gio-unix-2.0 gtk+-3.0
//#include <stdlib.h>
// char* guest_app_id(long s_pid, const char* instance_name, const char* wmname, const char* wmclass, const char* icon_name);
import "C"
import "unsafe"

func find_app_id(pid uint, instanceName, wmName, wmClass, iconName string) string {
	iName := C.CString(instanceName)
	wName := C.CString(wmName)
	wClass := C.CString(wmClass)
	icon := C.CString(iconName)
	id := C.guest_app_id(C.long(pid), iName, wName, wClass, icon)
	defer func() {
		C.free(unsafe.Pointer(iName))
		C.free(unsafe.Pointer(wName))
		C.free(unsafe.Pointer(wClass))
		C.free(unsafe.Pointer(icon))
		C.free(unsafe.Pointer(id))
	}()
	return C.GoString(id)
}
