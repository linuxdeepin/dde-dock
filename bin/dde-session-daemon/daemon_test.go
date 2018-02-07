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

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFilterList(t *testing.T) {
	var infos = []struct {
		origin    []string
		condition []string
		ret       []string
	}{
		{
			origin:    []string{"power", "audio", "dock"},
			condition: []string{"power", "dock"},
			ret:       []string{"audio"},
		},
		{
			origin:    []string{"power", "audio", "dock"},
			condition: []string{},
			ret:       []string{"power", "audio", "dock"},
		},
		{
			origin:    []string{"power", "audio", "dock"},
			condition: []string{"power", "dock", "audio"},
			ret:       []string(nil),
		},
	}

	Convey("Test filterList", t, func() {
		for _, info := range infos {
			So(filterList(info.origin, info.condition),
				ShouldResemble, info.ret)
		}
	})
}
