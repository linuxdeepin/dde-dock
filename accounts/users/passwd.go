/*
 * Copyright (C) 2013 ~ 2017 Deepin Technology Co., Ltd.
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

package users

/*
#cgo CFLAGS: -Wall -g
#cgo LDFLAGS: -lcrypt

#include <stdlib.h>
#include "passwd.h"
*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"
	"unsafe"
)

var (
	wLocker sync.Mutex
)

func EncodePasswd(words string) string {
	cwords := C.CString(words)
	defer C.free(unsafe.Pointer(cwords))

	return C.GoString(C.mkpasswd(cwords))
}

// password: has been crypt
func updatePasswd(password, username string) error {
	status := C.lock_shadow_file()
	if status != 0 {
		return fmt.Errorf("Lock shadow file failed")
	}
	defer C.unlock_shadow_file()

	content, err := ioutil.ReadFile(userFileShadow)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var datas []string
	var found bool
	for _, line := range lines {
		if len(line) == 0 {
			datas = append(datas, line)
			continue
		}

		items := strings.Split(line, ":")
		if items[0] != username {
			datas = append(datas, line)
			continue
		}

		found = true
		if items[1] == password {
			return nil
		}

		var tmp string
		for i, v := range items {
			if i != 0 {
				tmp += ":"
			}

			if i == 1 {
				tmp += password
				continue
			}

			tmp += v
		}

		datas = append(datas, tmp)
	}

	if !found {
		return fmt.Errorf("The username not exist.")
	}

	return writeStrvToFile(datas, userFileShadow, 0600)
}

func writeStrvToFile(datas []string, file string, mode os.FileMode) error {
	var content string
	for i, v := range datas {
		if i != 0 {
			content += "\n"
		}

		content += v
	}

	f, err := os.Create(file + ".bak~")
	if err != nil {
		return err
	}
	defer f.Close()

	wLocker.Lock()
	defer wLocker.Unlock()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	err = f.Sync()
	if err != nil {
		return err
	}

	os.Rename(file+".bak~", file)
	os.Chmod(file, mode)

	return nil
}
