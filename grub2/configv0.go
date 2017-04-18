package grub2

import (
	"strconv"
)

const configV0File = dataDir + "/grub2.json"

type ConfigV0 struct {
	NeedUpdate        bool // mark to generate grub configuration
	FixSettingsAlways bool
	EnableTheme       bool
	DefaultEntry      string
	Timeout           string
	Resolution        string
}

func NewConfigV0() *ConfigV0 {
	return new(ConfigV0)
}

func (c *ConfigV0) Load() error {
	logger.Info("load config-v0", configV0File)
	return loadJSON(configV0File, c)
}

func (c *ConfigV0) Upgrade() *Config {
	newc := NewConfig()
	newc.EnableTheme = c.EnableTheme
	newc.Resolution = c.Resolution

	// DefaultEntry str -> int
	idx, err := strconv.Atoi(c.DefaultEntry)
	if err != nil {
		idx = defaultDefaultEntry
	}
	newc.DefaultEntry = idx

	// Timeout str -> uint32
	timeout, err := strconv.ParseUint(c.Timeout, 10, 32)
	if err != nil {
		timeout = defaultTimeout
	}
	newc.Timeout = uint32(timeout)
	return newc
}
