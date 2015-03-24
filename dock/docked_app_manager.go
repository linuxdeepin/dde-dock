package dock

import (
	"container/list"
	"io/ioutil"
	"os"
	"path/filepath"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	"strings"
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
		m.items.PushBack(strings.ToLower(strings.Replace(id, "_", "-", -1)))
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
		if a := NewDesktopAppInfo(id + ".desktop"); a != nil {
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
	item := m.findItem(id)
	if item != nil {
		return true
	}

	if id = guess_desktop_id(id); id != "" {
		item = m.findItem(id)
	}
	// logger.Info("IsDocked:", item, item != nil)
	return item != nil
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

	id = strings.ToLower(id)
	idElement = m.findItem(id)
	if idElement != nil {
		logger.Info(id, "is already docked.")
		return false
	}

	logger.Debug("id", id, "title", title, "cmd", cmd)
	m.items.PushBack(id)
	if guess_desktop_id(id) == "" {
		// if app := NewDesktopAppInfo(id + ".desktop"); app != nil {
		// 	app.Unref()
		// } else {
		if e := createScratchFile(id, title, icon, cmd); e != nil {
			logger.Warning("create scratch file failed:", e)
			return false
		}
	}
	dbus.Emit(m, "Docked", id)
	app := ENTRY_MANAGER.runtimeApps[id]
	if app != nil {
		app.buildMenu()
	}
	return true
}

func (m *DockedAppManager) doUndock(id string) bool {
	removeItem := m.findItem(id)
	if removeItem == nil {
		logger.Warning("not find docked app:", id)
		return false
	}

	logger.Info("Undock", id)
	logger.Info("Remove", m.items.Remove(removeItem))
	m.core.SetStrv(DockedApps, m.toSlice())
	gio.SettingsSync()
	os.Remove(filepath.Join(scratchDir, id+".desktop"))
	os.Remove(filepath.Join(scratchDir, id+".sh"))
	os.Remove(filepath.Join(scratchDir, id+".png"))
	dbus.Emit(m, "Undocked", removeItem.Value.(string))
	app := ENTRY_MANAGER.runtimeApps[id]
	if app != nil {
		app.buildMenu()
	}

	return true
}

func (m *DockedAppManager) Undock(id string) bool {
	id = strings.ToLower(id)
	logger.Debug("undock lower id:", id)
	if m.doUndock(id) {
		return true
	}

	tmpId := ""
	if tmpId = guess_desktop_id(id); tmpId != "" {
		logger.Debug("undock guess desktop id:", tmpId)
		m.doUndock(tmpId)
		return true
	}

	tmpId = strings.Replace(id, "-", "_", -1)
	if m.doUndock(tmpId) {
		logger.Debug("undock replace - to _:", tmpId)
		return true
	}

	return false
}

func (m *DockedAppManager) findItem(id string) *list.Element {
	for e := m.items.Front(); e != nil; e = e.Next() {
		if strings.ToLower(e.Value.(string)) == strings.ToLower(id) {
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
	defer f.Close()
	temp := template.Must(template.New("docked_item_temp").Parse(DockedItemTemp))
	logger.Debug(title, ",", icon, ",", cmd)
	e := temp.Execute(f, dockedItemInfo{title, icon, cmd})
	if e != nil {
		return e
	}
	return nil
}
