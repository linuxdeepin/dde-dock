package desktop

import (
	"fmt"
	"strconv"
)

// MenuItem is menu item.
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

	id string
}

// NewMenuItem creates new menu item.
func NewMenuItem(name string, action func(), enable bool) *MenuItem {
	return &MenuItem{
		itemText:         name,
		isActive:         enable,
		isCheckable:      false,
		checked:          false,
		itemIcon:         "",
		itemIconHover:    "",
		itemIconInactive: "",
		showCheckMark:    false,
		subMenu:          nil,
		Action:           action,
	}
}

// Menu is menu.
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

// NewMenu creates new menu.
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

// SetIDGenerator changes the default id generator to the passed id generator.
func (m *Menu) SetIDGenerator(f func() string) *Menu {
	m.genID = f
	return m
}

// AddSeparator adds a separator to menu.
func (m *Menu) AddSeparator() *Menu {
	m.AppendItem(NewMenuItem("", nil, false))
	return m
}

// AppendItem appends a new menu item to menu.
func (m *Menu) AppendItem(items ...*MenuItem) *Menu {
	m.items = append(m.items, items...)
	for _, item := range items {
		if item.itemText != "" { // filter separator
			item.id = m.genID()
			m.ids[item.id] = item
		}
	}
	return m
}

func (m *Menu) handleAction(id string) bool {
	item, ok := m.ids[id]
	if ok {
		if item.isActive {
			item.Action()
		}
		return true
	}

	for _, item := range m.ids {
		if item.subMenu != nil && item.subMenu.handleAction(id) {
			return true
		}
	}

	return false
}

// HandleAction will call the action corresponding to the id.
func (m *Menu) HandleAction(id string) {
	m.handleAction(id)
}

// ToJSON generates json format menu content used in DeepinMenu.
func (m *Menu) ToJSON() string {
	ret := fmt.Sprintf(`{"checkableMenu":%v, "singleCheck":%v, "items":[`, m.checkableMenu, m.singleCheck)
	itemNumber := len(m.items)
	for i, item := range m.items {
		ret += fmt.Sprintf(`{"itemId":%q, "itemText":%q, "isActive":%v, "isCheckable":%v, "checked":%v, "itemIcon":%q, "itemIconHover":%q, "itemIconInactive":%q, "showCheckMark":%v, "itemSubMenu":`,
			item.id,
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
			ret += item.subMenu.ToJSON()
		}

		if i+1 == itemNumber {
			ret += "}"
		} else {
			ret += "},"
		}
	}
	ret += "]}"
	return ret
}
