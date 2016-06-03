/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

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
