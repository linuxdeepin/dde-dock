/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     Hualet Wang <mr.asianwang@gmail.com>
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

package shortcuts

import (
	"testing"

	"github.com/linuxdeepin/go-x11-client/util/keysyms"
)

func TestKey2Mode(t *testing.T) {
	mod, res := key2Mod("super_l")
	if mod != keysyms.ModMaskSuper || res != true {
		t.Fatalf("key2Mod on super_l failed")
	}

	mod, res = key2Mod("super_r")
	if mod != keysyms.ModMaskSuper || res != true {
		t.Fatalf("key2Mod on super_l failed")
	}

	mod, res = key2Mod("r")
	if mod != 0 || res != false {
		t.Fatalf("key2Mod on r failed")
	}
}

func BenchmarkKey2Mode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		key2Mod("super_l")
	}
}
