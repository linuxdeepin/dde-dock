/*
 * Copyright (C) 2014 ~ 2017 Deepin Technology Co., Ltd.
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

package dock

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_GenerateMenuJson(t *testing.T) {
	Convey("GenerateMenuJson", t, func() {
		menu := NewMenu()
		item0 := NewMenuItem("item 0", nil, true)
		item1 := NewMenuItem("item 1", nil, true)
		item2 := NewMenuItem("item 2", nil, true)
		menu.AppendItem(item0, item1, item2)

		menuJSON := menu.GenerateJSON()
		So(menuJSON, ShouldEqual, `{"items":[{"itemId":"0","itemText":"item 0","isActive":true,"isCheckable":false,"checked":false,"itemIcon":"","itemIconHover":"","itemIconInactive":"","showCheckMark":false,"itemSubMenu":null},{"itemId":"1","itemText":"item 1","isActive":true,"isCheckable":false,"checked":false,"itemIcon":"","itemIconHover":"","itemIconInactive":"","showCheckMark":false,"itemSubMenu":null},{"itemId":"2","itemText":"item 2","isActive":true,"isCheckable":false,"checked":false,"itemIcon":"","itemIconHover":"","itemIconInactive":"","showCheckMark":false,"itemSubMenu":null}],"checkableMenu":false,"singleCheck":false}`)

		var parseResult interface{}
		err := json.Unmarshal([]byte(menuJSON), &parseResult)
		So(err, ShouldBeNil)
	})
}
