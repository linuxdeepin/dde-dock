/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
