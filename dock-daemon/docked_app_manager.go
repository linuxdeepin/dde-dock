package main

import (
	"dlib/gio-2.0"
)

const (
	SchemaId   string = "com.deepin.dde.dock"
	DockedApps string = "docked-apps"
)

type DockedAppManager struct {
	core  *gio.Settings
	items []string

	Docked   func(id string, indicator string)
	Undocked func(id string)
}

func NewDockedAppManager() *DockedAppManager {
	m := &DockedAppManager{}
	m.init()
	return m
}

func (m *DockedAppManager) init() {
	m.core = gio.NewSettings(SchemaId)
}

func (m *DockedAppManager) DockedAppList() []string {
	if m.core != nil {
		return m.core.GetStrv(DockedApps)
	}
	return make([]string, 0)
}

func (m *DockedAppManager) Dock(id string, indicator string) bool {
	m.Docked(id, indicator)
	return true
}

func (m *DockedAppManager) Undock(id string) bool {
	m.Undocked(id)
	return true
}
