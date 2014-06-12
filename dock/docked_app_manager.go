package dock

import (
	"container/list"
	"dlib/gio-2.0"
	"dlib/glib-2.0"
	"io/ioutil"
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

var scratchDir string = filepath.Join(os.Getenv("HOME"), ".config/dock/scratch")

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
	m.items = list.New()
	m.core = gio.NewSettings(SchemaId)
	if m.core == nil {
		return
	}

	// TODO:
	// listen changed.
	appList := m.core.GetStrv(DockedApps)
	for _, id := range appList {
		m.items.PushBack(id)
	}

	conf := glib.NewKeyFile()
	defer conf.Free()

	confFile := filepath.Join(glib.GetUserConfigDir(), "dock/apps.ini")
	_, err := conf.LoadFromFile(confFile, glib.KeyFileFlagsNone)
	if err != nil {
		logger.Debug("Open old dock config file failed:", err)
		return
	}

	inited, err := conf.GetBoolean("__Config__", "inited")
	if err == nil && inited {
		return
	}

	_, ids, err := conf.GetStringList("__Config__", "Position")
	if err != nil {
		logger.Debug("Read docked app from old config file failed:", err)
		return
	}
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
	gio.SettingsSync()
	conf.SetBoolean("__Config__", "inited", true)

	_, content, err := conf.ToData()
	if err != nil {
		return
	}

	var mode os.FileMode = 0666
	stat, err := os.Lstat(confFile)
	if err == nil {
		mode = stat.Mode()
	}

	err = ioutil.WriteFile(confFile, []byte(content), mode)
	if err != nil {
		logger.Warning("Save Config file failed:", err)
	}
}

func (m *DockedAppManager) DockedAppList() []string {
	if m.core != nil {
		appList := m.core.GetStrv(DockedApps)
		return appList
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
		logger.Info(id, "is already docked.")
		return false
	}

	logger.Info("id", id, "title", title, "cmd", cmd)
	m.items.PushBack(id)
	if app := gio.NewDesktopAppInfo(id + ".desktop"); app != nil {
		app.Unref()
	} else {
		if e := createScratchFile(id, title, icon, cmd); e != nil {
			logger.Warning(e)
			return false
		}
	}
	m.Docked(id)
	return true
}

func (m *DockedAppManager) Undock(id string) bool {
	removeItem := m.findItem(id)
	if removeItem != nil {
		logger.Info("undock", id)
		logger.Info("Remove", m.items.Remove(removeItem))
		m.core.SetStrv(DockedApps, m.toSlice())
		gio.SettingsSync()
		os.Remove(filepath.Join(scratchDir, id+".desktop"))
		os.Remove(filepath.Join(scratchDir, id+".sh"))
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
	logger.Info("sort:", items)
	for _, item := range items {
		if i := m.findItem(item); i != nil {
			m.items.PushBack(m.items.Remove(i))
		}
	}
	l := m.toSlice()
	logger.Info("sorted:", l)
	m.core.SetStrv(DockedApps, l)
	gio.SettingsSync()
}

func (m *DockedAppManager) toSlice() []string {
	appList := make([]string, 0)
	for e := m.items.Front(); e != nil; e = e.Next() {
		appList = append(appList, e.Value.(string))
	}
	return appList
}

func createScratchFile(id, title, icon, cmd string) error {
	homeDir := os.Getenv("HOME")
	path := ".config/dock/scratch"
	configDir := filepath.Join(homeDir, path)
	os.MkdirAll(configDir, 0775)
	f, err := os.OpenFile(filepath.Join(configDir, id+".desktop"),
		os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0744)
	if err != nil {
		logger.Warning(err)
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
