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

package grub2

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"pkg.linuxdeepin.com/lib/dbus"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

var grubConfigFile = "/etc/default/grub"

func SetDefaultGrubConfigFile(file string) {
	grubConfigFile = file
}

const (
	grubMenuFile       = "/boot/grub/grub.cfg"
	grubUpdateExe      = "/usr/sbin/update-grub"
	grubTimeoutDisable = -2
)

var (
	runWithoutDBus         = false
	entryRegexpSingleQuote = regexp.MustCompile(`^ *(menuentry|submenu) +'(.*?)'.*$`)
	entryRegexpDoubleQuote = regexp.MustCompile(`^ *(menuentry|submenu) +"(.*?)".*$`)
)

// Grub2 is a dbus object, and provide properties and methods to setup
// grub2 and deepin grub2 theme.
type Grub2 struct {
	entries  []Entry
	settings map[string]string
	theme    *Theme
	config   *config

	needUpdateLock     sync.Mutex
	needUpdate         bool
	chanUpdate         chan int
	chanStopUpdateLoop chan int

	FixSettingsAlways bool   `access:"readwrite"`
	EnableTheme       bool   `access:"readwrite"`
	DefaultEntry      string `access:"readwrite"`
	Timeout           int32  `access:"readwrite"`
	Resolution        string `access:"readwrite"`
	Updating          bool
}

// NewGrub2 create a Grub2 object.
func NewGrub2() *Grub2 {
	grub := &Grub2{}
	grub.theme = NewTheme()
	grub.config = newConfig()
	grub.chanUpdate = make(chan int)
	grub.chanStopUpdateLoop = make(chan int)
	return grub
}

func DestroyGrub2(grub *Grub2) {
	grub.stopUpdateLoop()
	dbus.UnInstallObject(grub.theme)
	dbus.UnInstallObject(grub)
}

func (grub *Grub2) initGrub2() {
	grub.loadConfig()
	grub.doInitGrub2()
	grub.theme.initTheme()
	go grub.resetGfxmodeIfNeed()
	go grub.theme.regenerateBackgroundIfNeed()
	grub.startUpdateLoop()
}

func (grub *Grub2) loadConfig() {
	if grub.config.core.IsConfigFileExists() {
		grub.config.load()
	} else {
		grub.config.save()
	}
}

func (grub *Grub2) doInitGrub2() {
	err := grub.readEntries()
	if err != nil {
		logger.Error(err)
	}
	err = grub.readSettings()
	if err != nil {
		logger.Error(err)
	}

	needUpdate := false
	if grub.config.FixSettingsAlways {
		needUpdate = grub.fixSettings()
	}
	if needUpdate || grub.config.NeedUpdate {
		grub.notifyUpdate()
	}

	grub.setPropFixSettingsAlways(grub.config.FixSettingsAlways)
	grub.setPropEnableTheme(grub.config.EnableTheme)
	grub.setPropDefaultEntry(grub.getSettingDefaultEntry())
	grub.setPropTimeout(grub.getSettingTimeout())
	grub.setPropResolution(grub.getSettingGfxmode())
}

func (grub *Grub2) notifyUpdate() {
	grub.needUpdateLock.Lock()
	grub.needUpdate = true
	grub.needUpdateLock.Unlock()
	go func() {
		grub.chanUpdate <- 1
	}()
}

func (grub *Grub2) startUpdateLoop() {
	// start a goroutine to update grub configuration automatically
	go func() {
		logger.Info("update loop started")
		defer logger.Info("update loop stopped")
		for {
			select {
			case <-grub.chanStopUpdateLoop:
				break
			case <-grub.chanUpdate:
				grub.needUpdateLock.Lock()
				grub.config.NeedUpdate = grub.needUpdate
				grub.needUpdate = false
				grub.needUpdateLock.Unlock()

				if grub.config.NeedUpdate {
					grub.setPropUpdating(true)

					grub.config.save()

					logger.Info("notify to generate a new grub configuration file")
					grub2extDoGenerateGrubConfig()
					logger.Info("generate grub configuration finished")

					grub.config.NeedUpdate = false
					grub.config.save()

					// set property "Updating" to false only if don't
					// need update any more
					grub.needUpdateLock.Lock()
					if !grub.needUpdate {
						grub.setPropUpdating(false)
					}
					grub.needUpdateLock.Unlock()
				}
			}
		}
	}()
}
func (grub *Grub2) stopUpdateLoop() {
	grub.chanStopUpdateLoop <- 1
}

func (grub *Grub2) resetGfxmodeIfNeed() {
	if needUpdate := grub.resetGfxmode(); needUpdate {
		grub.notifyUpdate()

		// regenerate theme background
		screenWidth, screenHeight := getPrimaryScreenBestResolution()
		grub2extDoGenerateThemeBackground(screenWidth, screenHeight)
		grub.theme.setPropBackground(grub.theme.background)
	}
}

func (grub *Grub2) resetGfxmode() (needUpdate bool) {
	needUpdate = false
	expectedGfxmode := getPrimaryScreenBestResolutionStr()
	if expectedGfxmode != grub.getSettingGfxmode() {
		grub.setSettingGfxmode(expectedGfxmode)
		needUpdate = true
	}
	return
}

func (grub *Grub2) clearEntries() {
	grub.entries = make([]Entry, 0)
}

func (grub *Grub2) clearSettings() {
	grub.settings = make(map[string]string)
}

func (grub *Grub2) readEntries() (err error) {
	fileContent, err := ioutil.ReadFile(grubMenuFile)
	if err != nil {
		logger.Error(err)
		return
	}
	err = grub.parseEntries(string(fileContent))
	if err != nil {
		logger.Error(err)
		return
	}
	if len(grub.entries) == 0 {
		logger.Warning("there is no menu entry in %s", grubMenuFile)
	}
	return
}

func (grub *Grub2) readSettings() (err error) {
	fileContent, err := ioutil.ReadFile(grubConfigFile)
	if err != nil {
		logger.Error(err.Error())
	}
	err = grub.parseSettings(string(fileContent))

	return
}

func (grub *Grub2) fixSettings() (needUpdate bool) {
	needUpdate = false

	// reset properties, return default value for the missing property
	// default entry
	if grub.config.DefaultEntry != grub.getSettingDefaultEntry() {
		needUpdate = true
	}
	grub.setSettingDefaultEntry(grub.config.DefaultEntry)

	// timeout
	if grub.config.Timeout != grub.getSettingTimeout() {
		needUpdate = true
	}
	grub.setSettingTimeout(grub.config.Timeout)

	// gfxmode
	if grub.config.Resolution != grub.getSettingGfxmode() {
		needUpdate = true
	}
	grub.setSettingGfxmode(grub.config.Resolution)

	// disable GRUB_HIDDEN_TIMEOUT and GRUB_HIDDEN_TIMEOUT_QUIET which will conflicts with GRUB_TIMEOUT
	if len(grub.settings["GRUB_HIDDEN_TIMEOUT"]) != 0 ||
		len(grub.settings["GRUB_HIDDEN_TIMEOUT_QUIET"]) != 0 {
		grub.settings["GRUB_HIDDEN_TIMEOUT"] = ""
		grub.settings["GRUB_HIDDEN_TIMEOUT_QUIET"] = ""
		grub.writeSettings()
		needUpdate = true
	}

	// fix GRUB_DISTRIBUTOR
	grubDistroCmd := "`lsb_release -d -s 2> /dev/null || echo Debian`"
	if grub.settings["GRUB_DISTRIBUTOR"] != grubDistroCmd {
		grub.settings["GRUB_DISTRIBUTOR"] = grubDistroCmd
		grub.writeSettings()
		needUpdate = true
	}

	// disable GRUB_BACKGROUND
	if grub.settings["GRUB_BACKGROUND"] != "<none>" {
		grub.settings["GRUB_BACKGROUND"] = "<none>"
		grub.writeSettings()
		needUpdate = true
	}

	// setup deepin grub2 theme
	if grub.config.EnableTheme {
		if grub.getSettingTheme() != themeMainFile {
			grub.setSettingTheme(themeMainFile)
			needUpdate = true
		}
	} else {
		if grub.getSettingTheme() != "" {
			grub.setSettingTheme("")
			needUpdate = true
		}
	}

	return
}

func (grub *Grub2) writeSettings() {
	fileContent := grub.getSettingContentToSave()
	if runWithoutDBus {
		writeSettingsWithoutDBus(fileContent)
	} else {
		grub2extDoWriteSettings(fileContent)
	}
}

func (grub *Grub2) parseEntries(fileContent string) (err error) {
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
				err = fmt.Errorf("a 'menuentry' directive was detected inside the scope of a menuentry")
				return
			}
			title, ok := grub.parseTitle(line)
			if ok {
				entry := Entry{MENUENTRY, title, numCount[level], parentMenus[len(parentMenus)-1]}
				grub.entries = append(grub.entries, entry)
				logger.Debugf("found entry: [%d] %s %s", level, strings.Repeat(" ", level*2), title)

				numCount[level]++
				inMenuEntry = true
				continue
			} else {
				grub.clearEntries()
				err = fmt.Errorf("parse entry title failed from: %q", line)
				return
			}
		} else if strings.HasPrefix(line, "submenu ") {
			if inMenuEntry {
				grub.clearEntries()
				err = fmt.Errorf("a 'submenu' directive was detected inside the scope of a menuentry")
				return
			}
			title, ok := grub.parseTitle(line)
			if ok {
				entry := Entry{SUBMENU, title, numCount[level], parentMenus[len(parentMenus)-1]}
				grub.entries = append(grub.entries, entry)
				parentMenus = append(parentMenus, &entry)
				logger.Debugf("found entry: [%d] %s %s", level, strings.Repeat(" ", level*2), title)

				level++
				numCount[level] = 0
				continue
			} else {
				grub.clearEntries()
				err = fmt.Errorf("parse entry title failed from: %q", line)
				return
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
	err = sl.Err()
	if err != nil {
		return
	}
	return
}

func (grub *Grub2) parseTitle(line string) (string, bool) {
	line = strings.TrimLeftFunc(line, unicode.IsSpace)
	if entryRegexpSingleQuote.MatchString(line) {
		return entryRegexpSingleQuote.FindStringSubmatch(line)[2], true
	} else if entryRegexpDoubleQuote.MatchString(line) {
		return entryRegexpDoubleQuote.FindStringSubmatch(line)[2], true
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
			logger.Debugf("found setting: %s=%s", kv[0], kv[1])
		}
	}
	if err := s.Err(); err != nil {
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (grub *Grub2) getEntryTitles() (entryTitles []string, err error) {
	entryTitles = make([]string, 0)
	for _, entry := range grub.entries {
		if entry.entryType == MENUENTRY {
			entryTitles = append(entryTitles, entry.getFullTitle())
		}
	}
	if len(entryTitles) == 0 {
		err = fmt.Errorf("there is no menu entry in %s", grubMenuFile)
		return
	}
	return
}

// return setting's value or related entry(such as "Deepin 2014
// GNU/Linux") if grub.cfg exists
func (grub *Grub2) getSettingDefaultEntry() string {
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
	if isStringInArray(value, simpleEntryTitles) {
		return value
	}

	// if GRUB_DEFAULE exist and is a entry in submenu, return the first entry's title
	if isStringInArray(value, entryTitles) {
		return firstEntry
	}

	// if GRUB_DEFAULE exist and is a index number, return its entry name
	index, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		logger.Warningf(`invalid number, settings["GRUB_DEFAULT"]=%s`, grub.settings["GRUB_DEFAULT"])
		index = 0
	}
	if index >= 0 && int(index) < len(simpleEntryTitles) {
		return simpleEntryTitles[index]
	}
	return firstEntry
}

func (grub *Grub2) getSettingTimeout() int32 {
	if len(grub.settings["GRUB_TIMEOUT"]) == 0 {
		return grubTimeoutDisable
	}

	timeout, err := strconv.ParseInt(grub.settings["GRUB_TIMEOUT"], 10, 32)
	if err != nil {
		logger.Errorf(`valid value, settings["GRUB_TIMEOUT"]=%s`, grub.settings["GRUB_TIMEOUT"])
		return grubTimeoutDisable
	}
	return int32(timeout)
}

func (grub *Grub2) getSettingGfxmode() string {
	if len(grub.settings["GRUB_GFXMODE"]) == 0 {
		return "auto"
	}
	return grub.settings["GRUB_GFXMODE"]
}

func (grub *Grub2) getSettingTheme() string {
	return grub.settings["GRUB_THEME"]
}

func (grub *Grub2) setSettingDefaultEntry(title string) {
	grub.settings["GRUB_DEFAULT"] = title
	grub.config.setDefaultEntry(title)
	grub.writeSettings()
}

func (grub *Grub2) setSettingTimeout(timeout int32) {
	if timeout == grubTimeoutDisable {
		grub.settings["GRUB_TIMEOUT"] = ""
		grub.config.setTimeout(grubTimeoutDisable)
	} else {
		timeoutStr := strconv.FormatInt(int64(timeout), 10)
		grub.settings["GRUB_TIMEOUT"] = timeoutStr
		grub.config.setTimeout(timeout)
	}
	grub.writeSettings()
}

func (grub *Grub2) setSettingGfxmode(gfxmode string) {
	grub.settings["GRUB_GFXMODE"] = gfxmode
	grub.config.setResolution(gfxmode)
	grub.writeSettings()
}

func (grub *Grub2) setSettingTheme(themeFile string) {
	grub.settings["GRUB_THEME"] = themeFile
	// TODO
	// grub.config.setTheme(themeFile)
	grub.writeSettings()
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
