package network

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	VpnEnabled bool
	Devices    map[string]*DeviceConfig
}

type DeviceConfig struct {
	Enabled bool
}

const configFile = "/var/lib/dde-daemon/network/config.json"

func loadConfig(filename string, cfg *Config) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, cfg)
	return err
}

func loadConfigSafe(filename string) *Config {
	var cfg Config
	err := loadConfig(filename, &cfg)
	if err != nil && !os.IsNotExist(err) {
		logger.Warning("failed to load config:", err)
	}
	if cfg.Devices == nil {
		cfg.Devices = make(map[string]*DeviceConfig)
	}
	return &cfg
}

func saveConfig(filename string, cfg *Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	dir := filepath.Dir(filename)
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
