package airplane_mode

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	configFile = "/var/lib/dde-daemon/airplane_mode/config.json"
)

type config struct {
	Enabled bool
}

func loadConfig(filename string, cfg *config) error {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return json.Unmarshal(content, cfg)
}

func saveConfig(filename string, cfg *config) error {
	content, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, content, 0644)
	if err != nil {
		return err
	}
	return nil
}
