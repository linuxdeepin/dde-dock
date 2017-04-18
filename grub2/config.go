package grub2

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"pkg.deepin.io/lib/encoding/kv"
	"sort"
	"strconv"
)

const (
	dataDir    = "/var/cache/deepin"
	configFile = dataDir + "/grub2-v1.json"

	defaultDefaultEntry = 0
	defaultEnableTheme  = true
	defaultResolution   = "auto"
	defaultTimeout      = 5
)

type Config struct {
	DefaultEntry int
	EnableTheme  bool
	Resolution   string
	Timeout      uint32
}

func NewConfig() *Config {
	return new(Config)
}

func (c *Config) UseDefault() {
	c.DefaultEntry = defaultDefaultEntry
	c.EnableTheme = defaultEnableTheme
	c.Resolution = defaultResolution
	c.Timeout = defaultTimeout
}

func (c *Config) String() string {
	if c == nil {
		return "<nil>"
	}
	return fmt.Sprintf("entry: %d, theme: %v, timeout: %d, resolution: %s",
		c.DefaultEntry, c.EnableTheme, c.Timeout, c.Resolution)
}

func (c *Config) Hash() string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("DefaultEntry:%v", c.DefaultEntry))
	io.WriteString(h, fmt.Sprintf("EnableTheme:%v", c.EnableTheme))
	io.WriteString(h, fmt.Sprintf("Resolution:%v", c.Resolution))
	io.WriteString(h, fmt.Sprintf("Timeout:%v", c.Timeout))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (c1 *Config) Equal(c2 *Config) bool {
	return c1.DefaultEntry == c2.DefaultEntry &&
		c1.EnableTheme == c2.EnableTheme &&
		c1.Resolution == c2.Resolution &&
		c1.Timeout == c2.Timeout
}

const (
	GRUB_THEME                 = "/boot/grub/themes/deepin/theme.txt"
	GRUB_BACKGROUND            = "/boot/grub/themes/deepin/background.png"
	GRUB_DISTRIBUTOR           = "`/usr/bin/lsb_release -d -s 2>/dev/null || echo Deepin`"
	GRUB_CMDLINE_LINUX_DEFAULT = "splash quiet "
)

func (c *Config) GetGrubParamsContent() []byte {
	params, err := loadGrubParams()
	if err != nil {
		logger.Warning("loadGrubParams failed:", err)
		params = make(map[string]string)
	}

	// modify params
	// keep boot option
	if params["GRUB_CMDLINE_LINUX_DEFAULT"] == "" {
		params["GRUB_CMDLINE_LINUX_DEFAULT"] = quoteString(GRUB_CMDLINE_LINUX_DEFAULT)
	}
	params["GRUB_DISTRIBUTOR"] = quoteString(GRUB_DISTRIBUTOR)
	params["GRUB_DEFAULT"] = quoteString(strconv.Itoa(c.DefaultEntry))
	params["GRUB_GFXMODE"] = quoteString(c.Resolution)
	params["GRUB_TIMEOUT"] = quoteString(strconv.FormatUint(uint64(c.Timeout), 10))
	// disable GRUB_HIDDEN_TIMEOUT and GRUB_HIDDEN_TIMEOUT_QUIET which will conflicts with GRUB_TIMEOUT
	delete(params, "GRUB_HIDDEN_TIMEOUT")
	delete(params, "GRUB_HIDDEN_TIMEOUT_QUIET")

	if c.EnableTheme {
		params["GRUB_THEME"] = quoteString(GRUB_THEME)
		params["GRUB_BACKGROUND"] = quoteString(GRUB_BACKGROUND)
	} else {
		delete(params, "GRUB_THEME")
		delete(params, "GRUB_BACKGROUND")
	}

	keys := make(sort.StringSlice, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	keys.Sort()

	// write buf
	var buf bytes.Buffer
	buf.WriteString("# Written by " + DBusDest + "\n")
	for _, k := range keys {
		buf.WriteString(k + "=" + params[k] + "\n")
	}
	// if you want let the grub-mkconfig exit with error code,
	// uncomment the next line.
	//buf.WriteString("=\n")
	return buf.Bytes()
}

func loadGrubParams() (map[string]string, error) {
	f, err := os.Open(grubParamsFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	params := make(map[string]string)
	r := kv.NewReader(f)
	r.TrimSpace = kv.TrimLeadingTailingSpace
	r.Comment = '#'
	for {
		pair, err := r.Read()
		if err != nil {
			break
		}
		if pair.Key == "" {
			continue
		}
		params[pair.Key] = pair.Value
	}

	return params, nil
}

func (c *Config) Load() error {
	logger.Info("load config", configFile)
	return loadJSON(configFile, c)
}

func (c *Config) Save() error {
	return saveJSON(configFile, c)
}
