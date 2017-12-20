package grub2

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	"pkg.deepin.io/lib/encoding/kv"
)

const (
	grubParamsFile = "/etc/default/grub"
)

func loadGrubParams() (map[string]string, error) {
	params := make(map[string]string)
	f, err := os.Open(grubParamsFile)
	if err != nil {
		return params, err
	}
	defer f.Close()

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

const (
	grubBackground          = "GRUB_BACKGROUND"
	grubCmdlineLinuxDefault = "GRUB_CMDLINE_LINUX_DEFAULT"
	grubDefault             = "GRUB_DEFAULT"
	grubDistributor         = "GRUB_DISTRIBUTOR"
	grubGfxMode             = "GRUB_GFXMODE"
	grubTheme               = "GRUB_THEME"
	grubTimeout             = "GRUB_TIMEOUT"

	defaultGrubBackground          = "/boot/grub/themes/deepin/background.png"
	defaultGrubCmdlineLinuxDefault = "splash quiet"
	defaultGrubDefault             = "0"
	defaultGrubDefaultInt          = 0
	defaultGrubDistributor         = "`/usr/bin/lsb_release -d -s 2>/dev/null || echo Deepin`"
	defaultGrubGfxMode             = "auto"
	defaultGrubTheme               = "/boot/grub/themes/deepin/theme.txt"
	defaultGrubTimeout             = "5"
	defaultGrubTimeoutInt          = 5
)

func decodeShellValue(in string) string {
	output, err := exec.Command("/bin/sh", "-c", "echo -n "+in).Output()
	if err != nil {
		// fallback
		return strings.Trim(in, "\"")
	}
	return string(output)
}

func getTimeout(params map[string]string) int {
	timeoutStr := decodeShellValue(params[grubTimeout])
	timeoutInt, err := strconv.Atoi(timeoutStr)
	if err != nil {
		return defaultGrubTimeoutInt
	}
	return timeoutInt
}

func getGfxMode(params map[string]string) (val string) {
	val = decodeShellValue(params[grubGfxMode])
	if val == "" {
		val = defaultGrubGfxMode
	}
	return
}

func getDefaultEntry(params map[string]string) (val string) {
	val = decodeShellValue(params[grubDefault])
	if val == "" {
		val = defaultGrubDefault
	}
	return
}

func getTheme(params map[string]string) string {
	return decodeShellValue(params[grubBackground])
}

func getBackground(params map[string]string) string {
	return decodeShellValue(params[grubTheme])
}

func getGrubParamsMD5Sum(params map[string]string) string {
	return getBytesMD5Sum(getGrubParamsContent(params))
}

func getGrubParamsContent(params map[string]string) []byte {
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

func writeGrubParams(params map[string]string) (string, error) {
	logger.Debug("write grub params")
	content := getGrubParamsContent(params)

	err := ioutil.WriteFile(grubParamsFile, content, 0644)
	if err != nil {
		return "", err
	}

	return getBytesMD5Sum(content), nil
}

func getDefaultGrubParams() map[string]string {
	return map[string]string{
		grubDefault:             defaultGrubDefault,
		grubBackground:          defaultGrubBackground,
		grubTheme:               defaultGrubTheme,
		grubTimeout:             defaultGrubTimeout,
		grubGfxMode:             defaultGrubGfxMode,
		grubCmdlineLinuxDefault: quoteString(defaultGrubCmdlineLinuxDefault),
		grubDistributor:         quoteString(defaultGrubDistributor),
	}
}
