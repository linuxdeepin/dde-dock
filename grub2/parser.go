package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"
	"unicode"
)

const (
	_ENTRY_REGEXP_1 = `^ *menuentry +'(.*?)'.*$`
	_ENTRY_REGEXP_2 = `^ *menuentry +"(.*?)".*$`
)

func (grub *Grub2) readEntries() {
	fileContent, err := ioutil.ReadFile(grub.grubMenuFile)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.parseEntries(string(fileContent))
}

func (grub *Grub2) readSettings() {
	fileContent, err := ioutil.ReadFile(grub.grubConfigFile)
	if err != nil {
		logError(err.Error()) // TODO
		return
	}
	grub.parseSettings(string(fileContent))
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
