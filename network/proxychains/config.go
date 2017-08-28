package proxychains

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Type     string
	IP       string
	Port     uint32
	User     string
	Password string
}

func loadConfig(file string) (*Config, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg *Config) save(file string) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, data, 0600)
}
