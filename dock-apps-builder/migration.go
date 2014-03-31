package main

// +build ignore

// this file should be rewrite

//#cgo pkg-config: glib-2.0 gio-unix-2.0 gtk+-3.0
//#include <stdlib.h>
// char* guest_app_id(long s_pid, const char* instance_name, const char* wmname, const char* wmclass, const char* icon_name);
// char* get_exe_name(int pid);
//char* icon_name_to_path(const char* name, int size);
// void init_deepin();
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

func find_exec_name_by_pid(pid uint) string {
	return C.GoString(C.get_exe_name(C.int(pid)))
}

func get_theme_icon(name string, size int) string {
	iconName := C.CString(name)
	defer func() {
		C.free(unsafe.Pointer(iconName))
	}()
	return C.GoString(C.icon_name_to_path(iconName, C.int(size)))
}

func initDeepin() {
	C.init_deepin()
}
