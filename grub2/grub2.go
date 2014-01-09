package main

import (
	"bufio"
	"dlib/dbus"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

const (
	_GRUB_MENU         = "/boot/grub/grub.cfg"
	_GRUB_CONFIG       = "/etc/default/grub"
	_GRUB_MKCONFIG_EXE = "grub-mkconfig"
)

const (
	_ENTRY_REGEXP_1 = `^ *menuentry +'(.*?)'.*$`
	_ENTRY_REGEXP_2 = `^ *menuentry +"(.*?)".*$`
)

type Grub2 struct {
	settings map[string]string

	DefaultEntry uint32 `access:readwrite`
	Timeout      int32  `access:readwrite`
	Gfxmode      string `access:readwrite`
	Background   string `access:readwrite`
	Theme        string `access:readwrite`
	Entries      []string
}

func NewGrub2() *Grub2 {
	// TODO
	grub := &Grub2{}
	return grub
}

func (grub *Grub2) readEntries() {
	fileContent, err := ioutil.ReadFile(_GRUB_MENU)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.parseEntries(string(fileContent))
}

func (grub *Grub2) readSettings() {
	fileContent, err := ioutil.ReadFile(_GRUB_CONFIG)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.parseSettings(string(fileContent))
}

func (grub *Grub2) writeSettings() {
	fileContent := grub.getSettingContentToSave()
	err := ioutil.WriteFile(_GRUB_CONFIG, []byte(fileContent), 0644)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
}

func (grub *Grub2) udpateSettings() {
	exec.Command(_GRUB_MKCONFIG_EXE + " -o " + _GRUB_MENU)
}

func (grub *Grub2) parseEntries(fileContent string) {
	// reset entries
	grub.Entries = make([]string, 0)

	s := bufio.NewScanner(strings.NewReader(fileContent))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		entry, ok := grub.parseTitle(s.Text())
		if ok {
			grub.Entries = append(grub.Entries, entry)
			logInfo(fmt.Sprintf("found entry: %s", entry)) // TODO
		}
	}
	if err := s.Err(); err != nil {
		logError(err.Error())
	}
}

func (grub *Grub2) parseTitle(line string) (string, bool) {
	line = strings.TrimLeftFunc(line, unicode.IsSpace)
	reg1 := regexp.MustCompile(_ENTRY_REGEXP_1)
	reg2 := regexp.MustCompile(_ENTRY_REGEXP_2)
	if reg1.MatchString(line) {
		return reg1.FindStringSubmatch(line)[1], true
	} else if reg2.MatchString(line) {
		return reg2.FindStringSubmatch(line)[1], true
	} else {
		return "", false
	}
}

func (grub *Grub2) parseSettings(fileContent string) {
	// reset settings
	grub.settings = make(map[string]string)

	s := bufio.NewScanner(strings.NewReader(fileContent))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "GRUB_") {
			kv := strings.SplitN(line, "=", 2)
			key, value := kv[0], kv[1]
			grub.settings[key] = unquoteString(value)
			logInfo(fmt.Sprintf("found setting: %s=%s", kv[0], kv[1])) // TODO
		}
	}
	if err := s.Err(); err != nil {
		logError(err.Error())
	}

	grub.DefaultEntry = grub.getDefaultEntry()
	grub.Timeout = grub.getTimeout()
	grub.Gfxmode = grub.getGfxmode()
	grub.Background = grub.getBackground()
	grub.Theme = grub.getTheme()
}

func (grub *Grub2) getDefaultEntry() uint32 {
	if len(grub.settings["GRUB_DEFAULT"]) == 0 {
		return 0
	}

	index, err := strconv.ParseInt(grub.settings["GRUB_DEFAULT"], 10, 32)
	if err != nil {
		logError(fmt.Sprintf(`valid value, settings["GRUB_DEFAULT"]=%s`, grub.settings["GRUB_DEFAULT"])) // TODO
		return 0
	}
	return uint32(index)
}

func (grub *Grub2) getTimeout() int32 {
	if len(grub.settings["GRUB_TIMEOUT"]) == 0 {
		return 5
	}

	timeout, err := strconv.ParseInt(grub.settings["GRUB_TIMEOUT"], 10, 32)
	if err != nil {
		logError(fmt.Sprintf(`valid value, settings["GRUB_TIMEOUT"]=%s`, grub.settings["GRUB_TIMEOUT"])) // TODO
		return 5
	}
	return int32(timeout)
}

func (grub *Grub2) getGfxmode() string {
	if len(grub.settings["GRUB_GFXMODE"]) == 0 {
		return "auto"
	}

	return grub.settings["GRUB_GFXMODE"]
}

func (grub *Grub2) getBackground() string {
	return grub.settings["GRUB_BACKGROUND"]
}

func (grub *Grub2) getTheme() string {
	return grub.settings["GRUB_THEME"]
}

func (grub *Grub2) getSettingContentToSave() string {
	fileContent := ""
	for k, v := range grub.settings {
		if len(v) > 0 {
			fileContent += k + "=" + quoteString(v) + "\n"
		}
	}
	return fileContent
}

func main() {
	grub := NewGrub2()
	err := dbus.InstallOnSession(grub)
	if err != nil {
		panic(err) // TODO
	}
	select {}
}
