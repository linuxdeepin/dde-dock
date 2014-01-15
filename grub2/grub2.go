package main

import (
	"bufio"
	"dlib/dbus"
	"errors"
	"fmt"
	"io/ioutil"
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

var (
	_ENTRY_REGEXP_1 = regexp.MustCompile(`^ *(menuentry|submenu) +'(.*?)'.*$`)
	_ENTRY_REGEXP_2 = regexp.MustCompile(`^ *(menuentry|submenu) +"(.*?)".*$`)
)

type Grub2 struct {
	entries  []Entry
	settings map[string]string

	DefaultEntry string `access:"readwrite"`
	Timeout      int32  `access:"readwrite"`
	Gfxmode      string `access:"readwrite"`
	Background   string `access:"readwrite"`
	Theme        string `access:"readwrite"`
	InUpdate     bool
}

func NewGrub2() *Grub2 {
	grub := &Grub2{}
	grub.InUpdate = false
	return grub
}

func (grub *Grub2) clearEntries() {
	grub.entries = make([]Entry, 0)
}

func (grub *Grub2) clearSettings() {
	grub.settings = make(map[string]string)
}

func (grub *Grub2) readEntries() error {
	fileContent, err := ioutil.ReadFile(_GRUB_MENU)
	if err != nil {
		logError(err.Error()) // TODO
		return err
	}
	return grub.parseEntries(string(fileContent))
}

func (grub *Grub2) readSettings() error {
	fileContent, err := ioutil.ReadFile(_GRUB_CONFIG)
	if err != nil {
		logError(err.Error()) // TODO
		return err
	}
	return grub.parseSettings(string(fileContent))
}

func (grub *Grub2) writeSettings() error {
	fileContent := grub.getSettingContentToSave()
	err := ioutil.WriteFile(_GRUB_CONFIG, []byte(fileContent), 0644)
	if err != nil {
		logError(err.Error()) // TODO
		return err
	}
	return nil
}

func (grub *Grub2) generateGrubConfig() error {
	logInfo("start to generate a new grub configuration file")
	grub.InUpdate = true
	// TODO
	err := execAndWait(60, _GRUB_MKCONFIG_EXE, "-o", _GRUB_MENU)
	// execAndWait(60, _GRUB_MKCONFIG_EXE)
	grub.InUpdate = false
	logInfo("generate grub configuration finished")
	return err
}

func (grub *Grub2) parseEntries(fileContent string) error {
	grub.clearEntries()

	inMenuEntry := false
	level := 0
	numCount := make(map[int]int)
	numCount[0] = 0
	parentMenus := make([]*Entry, 0)
	parentMenus = append(parentMenus, nil)
	sl := bufio.NewScanner(strings.NewReader(fileContent))
	sl.Split(bufio.ScanLines)
	for sl.Scan() {
		line := sl.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "menuentry ") {
			if inMenuEntry {
				grub.clearEntries()
				s := "a 'menuentry' directive was detected inside the scope of a menuentry"
				logError(s)
				return errors.New(s)
			}
			title, ok := grub.parseTitle(line)
			if ok {
				entry := Entry{MENUENTRY, title, numCount[level], parentMenus[len(parentMenus)-1]}
				grub.entries = append(grub.entries, entry)
				logInfo("found entry: [%d] %s %s", level, strings.Repeat(" ", level*2), title)

				numCount[level]++
				inMenuEntry = true
				continue
			} else {
				grub.clearEntries()
				s := fmt.Sprintf("parse entry title failed from: %q", line)
				logError(s)
				return errors.New(s)
			}
		} else if strings.HasPrefix(line, "submenu ") {
			if inMenuEntry {
				grub.clearEntries()
				s := "a 'submenu' directive was detected inside the scope of a menuentry"
				logError(s)
				return errors.New(s)
			}
			title, ok := grub.parseTitle(line)
			if ok {
				entry := Entry{SUBMENU, title, numCount[level], parentMenus[len(parentMenus)-1]}
				grub.entries = append(grub.entries, entry)
				parentMenus = append(parentMenus, &entry)
				logInfo("found entry: [%d] %s %s", level, strings.Repeat(" ", level*2), title)

				level++
				numCount[level] = 0
				continue
			} else {
				grub.clearEntries()
				s := fmt.Sprintf("parse entry title failed from: %q", line)
				logError(s)
				return errors.New(s)
			}
		} else if line == "}" {
			if inMenuEntry {
				inMenuEntry = false
			} else if level > 0 {
				level--

				// delete last parent submenu
				i := len(parentMenus) - 1
				copy(parentMenus[i:], parentMenus[i+1:])
				parentMenus[len(parentMenus)-1] = nil
				parentMenus = parentMenus[:len(parentMenus)-1]
			}
		}
	}
	if err := sl.Err(); err != nil {
		logError(err.Error())
		return err
	}
	return nil
}

func (grub *Grub2) parseTitle(line string) (string, bool) {
	line = strings.TrimLeftFunc(line, unicode.IsSpace)
	if _ENTRY_REGEXP_1.MatchString(line) {
		return _ENTRY_REGEXP_1.FindStringSubmatch(line)[2], true
	} else if _ENTRY_REGEXP_2.MatchString(line) {
		return _ENTRY_REGEXP_2.FindStringSubmatch(line)[2], true
	} else {
		return "", false
	}
}

func (grub *Grub2) parseSettings(fileContent string) error {
	grub.clearSettings()

	s := bufio.NewScanner(strings.NewReader(fileContent))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "GRUB_") {
			kv := strings.SplitN(line, "=", 2)
			key, value := kv[0], kv[1]
			grub.settings[key] = unquoteString(value)
			logInfo("found setting: %s=%s", kv[0], kv[1]) // TODO
		}
	}
	if err := s.Err(); err != nil {
		logError(err.Error())
		return err
	}

	// get properties, return default value for the missing property
	grub.DefaultEntry = grub.getDefaultEntry()
	grub.Timeout = grub.getTimeout()
	grub.Gfxmode = grub.getGfxmode()
	grub.Background = grub.getBackground()
	grub.Theme = grub.getTheme()

	// reset settings, for to sync the default values
	grub.setDefaultEntry(grub.DefaultEntry)
	grub.setTimeout(grub.Timeout)
	grub.setGfxmode(grub.Gfxmode)
	grub.setBackground(grub.Background)
	grub.setTheme(grub.Theme)

	return nil
}

func (grub *Grub2) getDefaultEntry() string {
	entryTitles, _ := grub.GetEntryTitles()
	firstEntry := ""
	if len(entryTitles) > 0 {
		firstEntry = entryTitles[0]
	}
	value := grub.settings["GRUB_DEFAULT"]

	// if GRUB_DEFAULE is empty, return the first entry's title
	if len(value) == 0 {
		return firstEntry
	}

	// if GRUB_DEFAULE exist and is a valid entry name, just return it
	if stringInSlice(value, entryTitles) {
		return value
	}

	// if GRUB_DEFAULE exist and is a entry index, return its entry name
	index, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		logError(`valid value, settings["GRUB_DEFAULT"]=%s`, grub.settings["GRUB_DEFAULT"]) // TODO
		index = 0
	}
	if index >= 0 && int(index) < len(entryTitles) {
		return entryTitles[index]
	} else {
		return firstEntry
	}
}

func (grub *Grub2) getTimeout() int32 {
	if len(grub.settings["GRUB_TIMEOUT"]) == 0 {
		return 5
	}

	timeout, err := strconv.ParseInt(grub.settings["GRUB_TIMEOUT"], 10, 32)
	if err != nil {
		logError(`valid value, settings["GRUB_TIMEOUT"]=%s`, grub.settings["GRUB_TIMEOUT"]) // TODO
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
	defer func() {
		if err := recover(); err != nil {
			logFatal("grub2 failed: %v", err)
		}
	}()

	grub := NewGrub2()
	grub.Load()
	err := dbus.InstallOnSystem(grub)
	if err != nil {
		panic(err)
	}
	dbus.DealWithUnhandledMessage()

	select {}
}
