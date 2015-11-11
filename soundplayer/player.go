/**
 * Copyright (c) 2014 Deepin, Inc.
 *               2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package soundplayer

// #cgo pkg-config: glib-2.0 libcanberra
// #include <stdlib.h>
// #include "player.h"
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

var wg sync.WaitGroup

func doPlayThemeSoundASync(theme, event, device string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		doPlayThemeSound(theme, event, device)
	}()
}

func doPlayThemeSoundFileAsync(file, device string) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		doPlayThemeSoundFile(file, device)
	}()
}

func doPlayThemeSound(theme, event, device string) error {
	cTheme := C.CString(theme)
	defer C.free(unsafe.Pointer(cTheme))
	cEvent := C.CString(event)
	defer C.free(unsafe.Pointer(cEvent))
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))

	ret := C.canberra_play_system_sound(cTheme, cEvent, cDevice)
	if ret != 0 {
		msg := C.GoString(C.ca_strerror(ret))
		return fmt.Errorf("Play failed: %s", msg)
	}
	return nil
}

func doPlayThemeSoundFile(file, device string) error {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))

	ret := C.canberra_play_sound_file(cFile, cDevice)
	if ret != 0 {
		msg := C.GoString(C.ca_strerror(ret))
		return fmt.Errorf("Play failed: %s", msg)
	}
	return nil
}
