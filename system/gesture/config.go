package gesture

import (
	"encoding/json"
	"io/ioutil"
	"pkg.deepin.io/lib/utils"
)

type Config struct {
	LongPressDistance float64 `json:"longpress_distance"`
	Verbose           int     `json:"verbose"`
}

func loadConfig(filename string) (*Config, error) {
	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var conf Config
	err = json.Unmarshal(contents, &conf)
	if err != nil {
		return nil, err
	}
	return &conf, nil
}

func getConfigPath() string {
	suffix := "dde-daemon/gesture/conf.json"
	filename := "/etc/" + suffix
	if utils.IsFileExist(filename) {
		return filename
	}
	return "/usr/share/" + suffix
}
