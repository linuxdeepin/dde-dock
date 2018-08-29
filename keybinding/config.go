package keybinding

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Config struct {
	HandleTouchPadToggle bool
}

var globalConfig Config

func loadConfig() {
	data, err := ioutil.ReadFile("/var/lib/dde-daemon/keybinding/config.json")
	if err != nil {
		if !os.IsNotExist(err) {
			logger.Warning(err)
		}
		return
	}

	err = json.Unmarshal(data, &globalConfig)
	if err != nil {
		logger.Warning(err)
	}

	logger.Debugf("loadConfig %#v", globalConfig)
}
