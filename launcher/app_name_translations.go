package launcher

import (
	"bufio"
	"encoding/json"
	"os"
	. "pkg.deepin.io/dde/daemon/launcher/log"
	"pkg.deepin.io/lib/gettext"
)

const (
	dataDir                 = "/usr/share/dde/data/"
	appNameTranslationsFile = dataDir + "app_name_translations.json"
)

func loadNameMap() map[string]string {
	file, err := os.Open(appNameTranslationsFile)
	if err != nil {
		Log.Warning(err)
		return nil
	}
	reader := bufio.NewReader(file)
	dec := json.NewDecoder(reader)
	var data map[string](map[string]string)
	err = dec.Decode(&data)
	if err != nil {
		Log.Warning(err)
		return nil
	}

	lang := gettext.QueryLang()
	nameMap := data[lang]
	Log.Debug("loadNameMap:", nameMap)
	return nameMap
}
