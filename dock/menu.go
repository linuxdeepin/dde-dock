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
	"fmt"
	"strconv"
)

// json sample
// {
//    "checkableMenu" : false,
//    "items" : [
//       {
//          "itemText" : "item 1",
//          "isActive" : true,
//          "itemSubMenu" : nil,
//          "itemId" : "2",
//          "itemIconInactive" : "",
//          "checked" : false,
//          "itemIconHover" : "",
//          "itemIcon" : "",
//          "showCheckMark" : false,
//          "isCheckable" : false
//       },
//    ],
//    "singleCheck" : false
// }

type MenuItem struct {
	Id            string `json:"itemId"`
	Text          string `json:"itemText"`
	IsActive      bool   `json:"isActive"`
	IsCheckable   bool   `json:"isCheckable"`
	Checked       bool   `json:"checked"`
	Icon          string `json:"itemIcon"`
	IconHover     string `json:"itemIconHover"`
	IconInactive  string `json:"itemIconInactive"`
	ShowCheckMark bool   `json:"showCheckMark"`
	SubMenu       *Menu  `json:"itemSubMenu"`

	action func(uint32)
}

func NewMenuItem(name string, action func(uint32), enable bool) *MenuItem {
	return &MenuItem{
		Text:     name,
		IsActive: enable,
		action:   action,
	}
}

type Menu struct {
	Items         []*MenuItem `json:"items"`
	CheckableMenu bool        `json:"checkableMenu"`
	SingleCheck   bool        `json:"singleCheck"`

	itemCount int64
}

func NewMenu() *Menu {
	return &Menu{}
}

func (menu *Menu) allocId() string {
	idStr := strconv.FormatInt(menu.itemCount, 10)
	menu.itemCount++
	return idStr
}

func (m *Menu) AppendItem(items ...*MenuItem) {
	for _, item := range items {
		if item.Text != "" {
			item.Id = m.allocId()
			m.Items = append(m.Items, item)
		}
	}
}

func (m *Menu) HandleAction(id string, timestamp uint32) {
	for _, item := range m.Items {
		if id == item.Id && item.IsActive {
			fmt.Println(id)
			item.action(timestamp)
		}
	}
}

func (m *Menu) GenerateJSON() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		logger.Warning(err)
		return ""
	}
	return string(bytes)
}
