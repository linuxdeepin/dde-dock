package appearance

import (
	"io/ioutil"
	"os"
	"path"
)

const _fontConfVersion = "1.1"

var _fontVersionConf = os.Getenv("HOME") + "/.config/fontconfig/conf.d/deepin_conf.version"

func (m *Manager) checkFontConfVersion() bool {
	if isVersionRight(_fontConfVersion, _fontVersionConf) {
		return true
	}

	logger.Debug("Font config version not same, will delete config and create")
	err := os.Remove(_fontVersionConf)
	if err != nil {
		logger.Warning("Failed to remove font version:", err)
	}

	err = os.MkdirAll(path.Dir(_fontVersionConf), 0755)
	if err != nil {
		logger.Warning("Failed to create font version directory:", err)
		return false
	}

	err = ioutil.WriteFile(_fontVersionConf,
		[]byte(_fontConfVersion), 0644)
	if err != nil {
		logger.Warning("Failed to write font version:", err)
		return false
	}
	return false
}

func isVersionRight(version, file string) bool {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return false
	}

	return string(data) == version
}
