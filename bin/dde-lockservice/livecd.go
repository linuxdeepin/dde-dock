/**
 * Copyright (C) 2017 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

//#include <stdlib.h>
//#include <string.h>
//#include <shadow.h>
//#include <crypt.h>
//#cgo LDFLAGS: -lcrypt
//int is_livecd(const char *username)\
//{\
//    if (strcmp("deepin", username) != 0) {\
//        return 0;\
//    }\
//    struct spwd *data = getspnam(username);\
//    if (data == NULL || strlen(data->sp_pwdp) == 0) {\
//        return 0;\
//    }\
//    if (strcmp(crypt("", data->sp_pwdp), data->sp_pwdp) != 0) {\
//        return 0;\
//    }\
//    return 1;\
//}
import "C"
import "unsafe"

func isInLiveCD(username string) bool {
	cName := C.CString(username)
	ret := C.is_livecd(cName)
	C.free(unsafe.Pointer(cName))
	return (int(ret) == 1)
}
