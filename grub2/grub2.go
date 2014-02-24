/**
 * Copyright (c) 2013 ~ 2014 Deepin, Inc.
 *               2013 ~ 2014 Xu FaSheng
 *
 * Author:      Xu FaSheng <fasheng.xu@gmail.com>
 * Maintainer:  Xu FaSheng <fasheng.xu@gmail.com>
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 **/

package main

import (
	"bufio"
	"dlib/dbus"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

const (
	_GRUB_MENU            = "/boot/grub/grub.cfg"
	_GRUB_CONFIG          = "/etc/default/grub"
	_GRUB_MKCONFIG_EXE    = "/usr/sbin/grub-mkconfig"
	_GRUB_TIMEOUT_DISABLE = -2
	_GRUB_CACHE_FILE      = "/var/cache/dde-daemon/grub2.json"
)

var (
	_ENTRY_REGEXP_1 = regexp.MustCompile(`^ *(menuentry|submenu) +'(.*?)'.*$`)
	_ENTRY_REGEXP_2 = regexp.MustCompile(`^ *(menuentry|submenu) +"(.*?)".*$`)
)

type CacheConfig struct {
	LastScreenWidth, LastScreenHeight int32
	NeedUpdate                        bool // mark to generate grub configuration
}

type Grub2 struct {
	entries  []Entry
	settings map[string]string
	theme    *Theme
	config   CacheConfig

	needUpdateLock sync.Mutex
	needUpdate     bool
	chanUpdate     chan int

	DefaultEntry string `access:"readwrite"`
	Timeout      int32  `access:"readwrite"`
}

func NewGrub2() *Grub2 {
	grub := &Grub2{}
	grub.theme = NewTheme()
	grub.config = CacheConfig{1024, 768, true} // default value
	grub.chanUpdate = make(chan int)
	return grub
}

func (grub *Grub2) notifyUpdate() {
	go func() {
		grub.chanUpdate <- 1
	}()
	grub.needUpdateLock.Lock()
	grub.needUpdate = true
	grub.needUpdateLock.Unlock()
}

func (grub *Grub2) load() {
	err := grub.readEntries()
	if err != nil {
		panic(err)
	}
	err = grub.readSettings()
	if err != nil {
		panic(err)
	}

	if isFileExists(_GRUB_CACHE_FILE) {
		err = grub.readCacheConfig()
		if err != nil {
			panic(err)
		}
	} else {
		grub.writeCacheConfig()
	}
	if grub.config.NeedUpdate {
		grub.notifyUpdate()
	}

	grub.resetGfxmodeIfNeed()

	// start a goroutine to update grub configuration automatically
	go func() {
		for {
			<-grub.chanUpdate
			grub.needUpdateLock.Lock()
			grub.config.NeedUpdate = grub.needUpdate
			grub.needUpdate = false
			grub.needUpdateLock.Unlock()

			if grub.config.NeedUpdate {
				grub.writeCacheConfig()

				grub.generateGrubConfig()

				grub.config.NeedUpdate = false
				grub.writeCacheConfig()
			}
		}
	}()
}

func (grub *Grub2) resetGfxmodeIfNeed() {
	w, h := getPrimaryScreenBestResolution()
	gfxmode := fmt.Sprintf("%dx%d;auto", w, h)
	if gfxmode != grub.getGfxmode() || w != grub.config.LastScreenWidth || h != grub.config.LastScreenHeight {
		grub.setGfxmode(gfxmode)
		grub.writeSettings()

		grub.config.LastScreenWidth = w
		grub.config.LastScreenHeight = h

		grub.notifyUpdate()

		grub.theme.generateBackground()
	}
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
		logError(err.Error())
		return err
	}
	return grub.parseEntries(string(fileContent))
}

func (grub *Grub2) readSettings() error {
	fileContent, err := ioutil.ReadFile(_GRUB_CONFIG)
	if err != nil {
		logError(err.Error())
		return err
	}
	return grub.parseSettings(string(fileContent))
}

func (grub *Grub2) writeSettings() error {
	grub.setTheme(grub.theme.mainFile) // enable deepin grub2 theme
	fileContent := grub.getSettingContentToSave()
	err := ioutil.WriteFile(_GRUB_CONFIG, []byte(fileContent), 0664)
	if err != nil {
		logError(err.Error())
		return err
	}
	return nil
}

func (grub *Grub2) readCacheConfig() (err error) {
	fileContent, err := ioutil.ReadFile(_GRUB_CACHE_FILE)
	if err != nil {
		logError(err.Error())
		return
	}
	err = json.Unmarshal(fileContent, &grub.config)
	if err != nil {
		logError(err.Error())
		return
	}
	return
}

func (grub *Grub2) writeCacheConfig() (err error) {
	// ensure parent directory exists
	if !isFileExists(_GRUB_CACHE_FILE) {
		os.MkdirAll(path.Dir(_GRUB_CACHE_FILE), 0755)
	}
	fileContent, err := json.Marshal(grub.config)
	if err != nil {
		logError(err.Error())
		return
	}
	err = ioutil.WriteFile(_GRUB_CACHE_FILE, fileContent, 0644)
	if err != nil {
		logError(err.Error())
		return
	}
	return
}

func (grub *Grub2) generateGrubConfig() (err error) {
	logInfo("start to generate a new grub configuration file")
	err = execAndWait(30, _GRUB_MKCONFIG_EXE, "-o", _GRUB_MENU)
	return err
	if err != nil {
		logError("generate grub configuration failed")
	} else {
		logInfo("generate grub configuration finished")
	}
	return
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
			logInfo("found setting: %s=%s", kv[0], kv[1])
		}
	}
	if err := s.Err(); err != nil {
		logError(err.Error())
		return err
	}

	// reset properties, return default value for the missing property
	grub.DefaultEntry = grub.getDefaultEntry()
	grub.Timeout = grub.getTimeout()
	dbus.NotifyChange(grub, "DefaultEntry")
	dbus.NotifyChange(grub, "Timeout")

	// reset settings to sync the default values
	grub.setDefaultEntry(grub.DefaultEntry)
	grub.setTimeout(grub.Timeout)

	return nil
}

func (grub *Grub2) getEntryTitles() ([]string, error) {
	entryTitles := make([]string, 0)
	for _, entry := range grub.entries {
		if entry.entryType == MENUENTRY {
			entryTitles = append(entryTitles, entry.getFullTitle())
		}
	}
	if len(entryTitles) == 0 {
		s := fmt.Sprintf("there is no menu entry in %s", _GRUB_MENU)
		logError(s)
		return entryTitles, errors.New(s)
	}
	return entryTitles, nil
}

func (grub *Grub2) getDefaultEntry() string {
	entryTitles, _ := grub.getEntryTitles()
	simpleEntryTitles, _ := grub.GetSimpleEntryTitles()
	firstEntry := ""
	if len(simpleEntryTitles) > 0 {
		firstEntry = simpleEntryTitles[0]
	}
	value := grub.settings["GRUB_DEFAULT"]

	// if GRUB_DEFAULE is empty, return the first entry's title
	if len(value) == 0 {
		return firstEntry
	}

	// if GRUB_DEFAULE exist and is a valid entry name, just return it
	if stringInSlice(value, simpleEntryTitles) {
		return value
	}

	// if GRUB_DEFAULE exist and is a entry in submenu, return the first entry's title
	if stringInSlice(value, entryTitles) {
		return firstEntry
	}

	// if GRUB_DEFAULE exist and is a index number, return its entry name
	index, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		logError(`invalid number, settings["GRUB_DEFAULT"]=%s`, grub.settings["GRUB_DEFAULT"])
		index = 0
	}
	if index >= 0 && int(index) < len(simpleEntryTitles) {
		return simpleEntryTitles[index]
	} else {
		return firstEntry
	}
}

func (grub *Grub2) getTimeout() int32 {
	if len(grub.settings["GRUB_TIMEOUT"]) == 0 {
		return _GRUB_TIMEOUT_DISABLE
	}

	timeout, err := strconv.ParseInt(grub.settings["GRUB_TIMEOUT"], 10, 32)
	if err != nil {
		logError(`valid value, settings["GRUB_TIMEOUT"]=%s`, grub.settings["GRUB_TIMEOUT"])
		return _GRUB_TIMEOUT_DISABLE
	}
	return int32(timeout)
}

func (grub *Grub2) getGfxmode() string {
	if len(grub.settings["GRUB_GFXMODE"]) == 0 {
		return "auto"
	}

	return grub.settings["GRUB_GFXMODE"]
}

func (grub *Grub2) getTheme() string {
	return grub.settings["GRUB_THEME"]
}

func (grub *Grub2) setDefaultEntry(title string) {
	grub.settings["GRUB_DEFAULT"] = title
}

func (grub *Grub2) setTimeout(timeout int32) {
	if timeout == _GRUB_TIMEOUT_DISABLE {
		grub.settings["GRUB_TIMEOUT"] = ""
	} else {
		timeoutStr := strconv.FormatInt(int64(timeout), 10)
		grub.settings["GRUB_TIMEOUT"] = timeoutStr
	}
}

func (grub *Grub2) setGfxmode(gfxmode string) {
	grub.settings["GRUB_GFXMODE"] = gfxmode
}

func (grub *Grub2) setTheme(themeFile string) {
	grub.settings["GRUB_THEME"] = themeFile
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
	err := dbus.InstallOnSystem(grub)
	if err != nil {
		panic(err)
	}
	err = dbus.InstallOnSystem(grub.theme)
	if err != nil {
		panic(err)
	}

	// load after dbus service installed to ensure property changed
	// signal send success
	grub.load()
	grub.theme.load()

	dbus.DealWithUnhandledMessage()

	select {}
}
