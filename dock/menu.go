package dock

import (
	"fmt"
	"strconv"
)

type MenuItem struct {
	// Name    string
	// Enabled bool

	itemText         string
	isActive         bool
	isCheckable      bool
	checked          bool
	itemIcon         string
	itemIconHover    string
	itemIconInactive string
	showCheckMark    bool
	subMenu          *Menu

	Action func()
}

// TODO: set properties.

func NewMenuItem(name string, action func(), enable bool) *MenuItem {
	return &MenuItem{
		name,
		enable,
		false,
		false,
		"",
		"",
		"",
		false,
		nil,
		action,
	}
}

type Menu struct {
	// Dock handle these.
	// x, y            int32
	// isDockMenu      bool
	// cornerDirection Direction
	// content         *MenuContent
	items []*MenuItem

	ids map[string]*MenuItem

	checkableMenu bool
	singleCheck   bool

	genID func() string
}

func NewMenu() *Menu {
	return &Menu{
		make([]*MenuItem, 0),
		make(map[string]*MenuItem),
		false,
		false,
		func() func() string {
			id := int64(0)
			return func() string {
				id++
				return strconv.FormatInt(id, 10)
			}
		}(),
	}
}

func (m *Menu) AddSeparator() *Menu {
	m.AppendItem(NewMenuItem("", nil, false))
	return m
}

func (m *Menu) AppendItem(items ...*MenuItem) {
	m.items = append(m.items, items...)
	for _, item := range items {
		if item.itemText != "" { // filter separator
			m.ids[m.genID()] = item
		}
	}
}

func (m *Menu) HandleAction(id string) {
	if item, ok := m.ids[id]; ok && item.isActive {
		item.Action()
	}
}

func (m *Menu) GenerateJSON() string {
	ret := fmt.Sprintf(`{"checkableMenu":%v, "singleCheck": %v, "items":[`, m.checkableMenu, m.singleCheck)
	itemNumber := len(m.items)
	for i, item := range m.items {
		for id, _item := range m.ids {
			if _item == item {
				ret += fmt.Sprintf(`{"itemId":"%s", "itemText": "%s", "isActive": %v, "isCheckable":%v, "checked":%v, "itemIcon":"%s", "itemIconHover":"%s", "itemIconInactive":"%s", "showCheckMark":%v, "itemSubMenu":`,
					id,
					item.itemText,
					item.isActive,
					item.isCheckable,
					item.checked,
					item.itemIcon,
					item.itemIconHover,
					item.itemIconInactive,
					item.showCheckMark,
				)

				if item.subMenu == nil {
					ret += `{"checkableMenu":false, "singleCheck":false, "items": []}`
				} else {
					ret += item.subMenu.GenerateJSON()
				}

				if i+1 == itemNumber {
					ret += "}"
				} else {
					ret += "},"
				}
			}
		}
	}
	ret += "]}"
	return ret
}
