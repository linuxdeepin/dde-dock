package main

import (
	"container/list"
	"dlib/gio-2.0"
	"dlib/glib-2.0"
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
	if m.core != nil {
		conf := glib.NewKeyFile()
		defer conf.Free()
		confFile := filepath.Join(glib.GetUserConfigDir(), "dock/apps.ini")
		conf.LoadFromFile(confFile, glib.KeyFileFlagsNone)
		if inited, err := conf.GetBoolean("__Config__", "inited"); err != nil && !inited {
			_, ids, _ := conf.GetStringList("__Config__", "Position")
			for _, id := range ids {
				if a := gio.NewDesktopAppInfo(id + ".desktop"); a != nil {
					a.Unref()
					continue
				}

				exec, _ := conf.GetString(id, "CmdLine")
				icon, _ := conf.GetString(id, "Icon")
				title, _ := conf.GetString(id, "Name")
				createScratchFile(id, title, icon, exec)
			}
			m.core.SetStrv(DockedApps, ids)
			conf.SetBoolean("__Config__", "inited", true)
			gio.SettingsSync()
		}
	}
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

func (m *DockedAppManager) IsDocked(id string) bool {
	return m.findItem(id) != nil
}

type dockedItemInfo struct {
	Name, Icon, Exec string
}

func (m *DockedAppManager) Dock(id, title, icon, cmd string) bool {
	idElement := m.findItem(id)
	if idElement != nil {
		return false
	}

	logger.Info("id", id, "title", title, "cmd", cmd)
	m.items.PushBack(id)
	if app := gio.NewDesktopAppInfo(id + ".desktop"); app != nil {
		app.Unref()
	} else {
		if e := createScratchFile(id, title, icon, cmd); e != nil {
			logger.Error(e)
			return false
		}
	}
	m.Docked(id)
	return true
}

func (m *DockedAppManager) Undock(id string) bool {
	removeItem := m.findItem(id)
	if removeItem != nil {
		logger.Info("undocd", id)
		m.items.Remove(removeItem)
		m.core.SetStrv(DockedApps, m.toSlice())
		gio.SettingsSync()
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
	logger.Info(items)
	for _, item := range items {
		if i := m.findItem(item); i != nil {
			m.items.PushBack(m.items.Remove(i))
		}
	}
	l := m.toSlice()
	m.core.SetStrv(DockedApps, l)
	gio.SettingsSync()
}

func (m *DockedAppManager) toSlice() []string {
	list := make([]string, 0)
	for e := m.items.Front(); e != nil; e = e.Next() {
		list = append(list, e.Value.(string))
	}
	return list
}

func createScratchFile(id, title, icon, cmd string) error {
	homeDir := os.Getenv("HOME")
	path := ".config/dock/scratch"
	configDir := filepath.Join(homeDir, path)
	os.MkdirAll(configDir, 0775)
	f, err := os.Create(filepath.Join(configDir, id+".desktop"))
	if err != nil {
		return err
	}
	temp := template.Must(template.New("docked_item_temp").Parse(DockedItemTemp))
	logger.Info(title, ",", icon, ",", cmd)
	e := temp.Execute(f, dockedItemInfo{title, icon, cmd})
	if e != nil {
		return e
	}
	return nil
}
