package mime

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/mime"
	dutils "pkg.deepin.io/lib/utils"
)

type defaultAppTable struct {
	Apps defaultAppInfos `json:"DefaultApps"`
}

type defaultAppInfo struct {
	AppId   string   `json:"AppId"`
	AppType string   `json:"AppType"`
	Types   []string `json:"SupportedType"`
}
type defaultAppInfos []*defaultAppInfo

func unmarshal(file string) (*defaultAppTable, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var table defaultAppTable
	err = json.Unmarshal(content, &table)
	if err != nil {
		return nil, err
	}

	return &table, nil
}

func genMimeAppsFile(data string) error {
	table, err := unmarshal(data)
	if err != nil {
		return err
	}

	for _, info := range table.Apps {
		for _, ty := range info.Types {
			mime.SetDefaultApp(ty, info.AppId)
		}
	}
	return nil
}

func (m *DefaultApps) initConfigData() {
	if dutils.IsFileExist(path.Join(glib.GetUserConfigDir(),
		"mimeapps.list")) {
		return
	}

	err := m.doInitConfigData()
	if err != nil {
		logger.Warning("Init mime config file failed", err)
	}
}

func (m *DefaultApps) doInitConfigData() error {
	os.Remove(path.Join(glib.GetUserConfigDir(), "mimeapps.list"))

	var data = "data.json"
	switch os.Getenv("LANGUAGE") {
	case "zh_CN", "zh_TW", "zh_HK":
		data = "data-zh_CN.json"
	}
	return genMimeAppsFile(
		findFilePath(path.Join("dde-daemon", "mime", data)))
}

func findFilePath(file string) string {
	data := path.Join(os.Getenv("HOME"), ".local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	data = path.Join("/usr/local/share", file)
	if dutils.IsFileExist(data) {
		return data
	}

	return path.Join("/usr/share", file)
}
