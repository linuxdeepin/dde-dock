package main

import (
	"bufio"
	"dlib/dbus"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"unicode"
)

const (
	GRUB_MENU         = "/boot/grub/grub.cfg"
	GRUB_CONFIG       = "/etc/default/grub"
	GRUB_MKCONFIG_EXE = "grub-mkconfig"
)

const (
	_ENTRY_REGEXP_1 = `^ *menuentry +'(.*?)'.*$`
	_ENTRY_REGEXP_2 = `^ *menuentry +"(.*?)".*$`
)

type Grub2 struct {
	entries  []string
	settings map[string]string

	DefaultEntry uint32 `access:readwrite`
	Timeout      uint32 `access:readwrite`
	Gfxmode      string `access:readwrite`
	Background   string `access:readwrite`
	Theme        string `access:readwrite`
}

func NewGrub2() *Grub2 {
	// TODO
	grub := &Grub2{}
	return grub
}

func (grub *Grub2) readEntries() {
	fileContent, err := ioutil.ReadFile(GRUB_MENU)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.parseEntries(string(fileContent))
}

func (grub *Grub2) readSettings() {
	fileContent, err := ioutil.ReadFile(GRUB_CONFIG)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.parseSettings(string(fileContent))
}

func (grub *Grub2) writeSettings() {
	fileContent := grub.getSettingContentToSave()
	err := ioutil.WriteFile(GRUB_CONFIG, []byte(fileContent), 0644)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
}

func (grub *Grub2) udpateSettings() {
	exec.Command(GRUB_MKCONFIG_EXE + " -o " + GRUB_MENU)
}

func (grub *Grub2) parseEntries(fileContent string) {
	// reset entries
	grub.entries = make([]string, 0)

	s := bufio.NewScanner(strings.NewReader(fileContent))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		entry, ok := grub.parseTitle(s.Text())
		if ok {
			grub.entries = append(grub.entries, entry)
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
