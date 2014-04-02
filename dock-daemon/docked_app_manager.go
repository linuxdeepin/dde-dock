package main

import (
	"container/list"
	"dlib/gio-2.0"
	"os"
	"path/filepath"
	"text/template"
)

const (
	SchemaId       string = "com.deepin.dde.dock"
	DockedApps     string = "docked-apps"
	DockedItemTemp string = `[Desktop Entry]
Name={{ .Name }}
Exec={{ .Exec }}
Icon={{ .Icon }}
Type=Application
Terminal=false
StartupNotify=false
`
)

type DockedAppManager struct {
	core  *gio.Settings
	items *list.List

	Docked   func(id string) // find indicator on front-end.
	Undocked func(id string)
}

func NewDockedAppManager() *DockedAppManager {
	m := &DockedAppManager{}
	m.init()
	return m
}

func (m *DockedAppManager) init() {
	m.core = gio.NewSettings(SchemaId)
	m.items = list.New()
}

func (m *DockedAppManager) DockedAppList() []string {
	if m.core != nil {
		list := m.core.GetStrv(DockedApps)
		for _, id := range list {
			m.items.PushBack(id)
		}
		return list
	}
	return make([]string, 0)
}

type dockedItemInfo struct {
	Name, Icon, Exec string
}

func (m *DockedAppManager) Dock(id, title, icon, cmd string) bool {
	idElement := m.findItem(id)
	if idElement != nil {
		return false
	}

	m.items.PushBack(id)
	if app := gio.NewDesktopAppInfo(id + ".desktop"); app != nil {
		app.Unref()
	} else {
		homeDir := os.Getenv("HOME")
		path := ".config/dock/scratch"
		configDir := filepath.Join(homeDir, path)
		os.MkdirAll(configDir, 0775)
		f, err := os.Create(filepath.Join(configDir, id+".desktop"))
		if err != nil {
			return false
		}
		temp := template.Must(template.New("docked_item_temp").Parse(DockedItemTemp))
		logger.Info(title, ",", icon, ",", cmd)
		e := temp.Execute(f, dockedItemInfo{title, icon, cmd})
		if e != nil {
			logger.Error(e)
		}
	}
	m.Docked(id)
	return true
}

func (m *DockedAppManager) Undock(id string) bool {
	removeItem := m.findItem(id)
	if removeItem != nil {
		m.items.Remove(removeItem)
		m.core.SetStrv(DockedApps, m.toSlice())
		os.Remove(filepath.Join(
			os.Getenv("HOME"),
			".config/dock/scratch",
			id+".desktop",
		))
		m.Undocked(id)
		return true
	} else {
		return false
	}
}

func (m *DockedAppManager) findItem(id string) *list.Element {
	for e := m.items.Front(); e != nil; e = e.Next() {
		if e.Value.(string) == id {
			return e
		}
	}
	return nil
}

func (m *DockedAppManager) Sort(items []string) {
	stack := list.New()
	for _, item := range items {
		if i := m.findItem(item); i != nil {
			m.items.Remove(i)
			if stack.Len() != 0 {
				m.items.InsertBefore(i, stack.Back())
				stack.Remove(stack.Back())
			}
			stack.PushBack(i)
		}
	}
	if stack.Len() != 0 {
		m.items.PushBackList(stack)
	}
	m.core.SetStrv(DockedApps, m.toSlice())
}

func (m *DockedAppManager) toSlice() []string {
	list := make([]string, 0)
	for e := m.items.Front(); e != nil; e = e.Next() {
		list = append(list, e.Value.(string))
	}
	return list
}
