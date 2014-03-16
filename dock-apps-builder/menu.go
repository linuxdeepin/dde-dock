package main

import (
	"fmt"
)

type MenuItem struct {
	Name    string
	Action  func()
	Enabled bool
}

type Menu struct {
	items []*MenuItem

	ids   map[int32]*MenuItem
	genID func() int32
}

func NewMenu() *Menu {
	return &Menu{
		make([]*MenuItem, 0),
		make(map[int32]*MenuItem),
		func() func() int32 {
			id := int32(0)
			return func() int32 {
				id++
				return id
			}
		}(),
	}
}

func (m *Menu) AppendItem(items ...*MenuItem) {
	m.items = append(m.items, items...)
	for _, item := range items {
		m.ids[m.genID()] = item
	}
}

func (m *Menu) HandleAction(id int32) {
	if item, ok := m.ids[id]; ok && item.Enabled {
		item.Action()
	}
}

func (m *Menu) GenerateJSON() string {
	ret := "["
	for _, item := range m.items {
		for id, _item := range m.ids {
			if _item == item {
				ret += fmt.Sprintf(`{"Id":%d, "Name": "%s", "Enabled": %v},`, id, item.Name, item.Enabled)
			}
		}
	}
	ret += "]"
	return ret
}
