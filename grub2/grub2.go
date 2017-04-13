/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package grub2

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"pkg.deepin.io/lib/dbus"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"unicode"
)

const (
	DefaultGrubSettingFile = "/etc/default/grub"
	DefaultGrubMenuFile    = "/boot/grub/grub.cfg"
)

var grubSettingFile = DefaultGrubSettingFile

func SetDefaultGrubSettingFile(file string) {
	grubSettingFile = file
}

const (
	grubUpdateCmd                 = "/usr/sbin/update-grub"
	lsbReleaseCmd                 = "/usr/bin/lsb_release"
	defaultGrubDefaultEntry       = "0"
	defaultGrubGfxmode            = "auto"
	defaultGrubTimeout            = "5"
	defaultGrubTimeoutInt   int32 = 5
)

var (
	runWithoutDbus         = false
	grubMenuFile           = DefaultGrubMenuFile
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

	// DefaultEntry should always is the entry title instead of index,
	// and will convert to index when saving settings
	DefaultEntry      string `access:"readwrite"`
	FixSettingsAlways bool   `access:"readwrite"`
	EnableTheme       bool   `access:"readwrite"`
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
	grub.resetEntries()
	grub.resetSettings()
	return grub
}

func DestroyGrub2(grub *Grub2) {
	grub.stopUpdateLoop()
	dbus.UnInstallObject(grub.theme)
	dbus.UnInstallObject(grub)
}

func (grub *Grub2) initGrub2() {
	grub.config.Load()
	grub.doInitGrub2()
	grub.theme.initTheme()
	go grub.theme.regenerateBackgroundIfNeed()
	grub.startUpdateLoop()
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
	grub.DefaultEntry = grub.getSettingDefaultEntry()
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
					grub2extDoGenerateGrubMenu()
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

func (grub *Grub2) resetEntries() {
	grub.entries = make([]Entry, 0)
}

func (grub *Grub2) resetSettings() {
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
	fileContent, err := ioutil.ReadFile(grubSettingFile)
	if err != nil {
		logger.Error(err)
	}
	err = grub.parseSettings(string(fileContent))

	return
}

func (grub *Grub2) fixSettings() (needUpdate bool) {
	needUpdate = grub.doFixSettings()
	if needUpdate {
		grub.writeSettings()
		grub.config.save()
	}
	return
}

func (grub *Grub2) doFixSettings() (needUpdate bool) {
	needUpdate = false

	// reset properties, return default value for the missing property
	// default entry
	if grub.config.DefaultEntry != grub.doGetSettingDefaultEntry() {
		logger.Infof("fix setting DefaultEntry %s->%s", grub.doGetSettingDefaultEntry(), grub.config.DefaultEntry)
		needUpdate = true
	}
	grub.doSetSettingDefaultEntry(grub.config.DefaultEntry)

	// timeout
	if grub.config.Timeout != grub.doGetSettingTimeout() {
		logger.Infof("fix setting Timeout %s->%s", grub.doGetSettingTimeout(), grub.config.Timeout)
		needUpdate = true
	}
	grub.doSetSettingTimeout(grub.config.Timeout)

	// gfxmode
	if grub.config.Resolution != grub.doGetSettingGfxmode() {
		logger.Infof("fix setting Resolution %s->%s", grub.doGetSettingGfxmode(), grub.config.Resolution)
		needUpdate = true
	}
	grub.doSetSettingGfxmode(grub.config.Resolution)

	// disable GRUB_HIDDEN_TIMEOUT and GRUB_HIDDEN_TIMEOUT_QUIET which will conflicts with GRUB_TIMEOUT
	if len(grub.settings["GRUB_HIDDEN_TIMEOUT"]) != 0 ||
		len(grub.settings["GRUB_HIDDEN_TIMEOUT_QUIET"]) != 0 {
		logger.Info("fix setting GRUB_HIDDEN_TIMEOUT")
		grub.settings["GRUB_HIDDEN_TIMEOUT"] = ""
		grub.settings["GRUB_HIDDEN_TIMEOUT_QUIET"] = ""
		needUpdate = true
	}

	// setup deepin grub2 theme
	if grub.config.EnableTheme {
		if grub.doGetSettingTheme() != themeMainFile {
			logger.Infof("fix setting theme %s->%s", grub.doGetSettingTheme(), themeMainFile)
			grub.doSetSettingTheme(themeMainFile)
			needUpdate = true
		}
		// make GRUB_BACKGROUND same with the theme background
		if grub.settings["GRUB_BACKGROUND"] != themeBgFile {
			logger.Info("fix setting GRUB_BACKGROUND")
			grub.settings["GRUB_BACKGROUND"] = themeBgFile
			needUpdate = true
		}
	} else {
		if grub.doGetSettingTheme() != "" {
			logger.Infof("fix setting theme %s->%s", grub.doGetSettingTheme(), "<disabled>")
			grub.doSetSettingTheme("")
			needUpdate = true
		}
		if grub.settings["GRUB_BACKGROUND"] != "" {
			logger.Info("fix setting GRUB_BACKGROUND to empty")
			grub.settings["GRUB_BACKGROUND"] = ""
			needUpdate = true
		}
	}

	return
}

func (grub *Grub2) fixSettingDistro() (needUpdate bool) {
	needUpdate = grub.doFixSettingDistro()
	if needUpdate {
		grub.writeSettings()
	}
	return
}
func (grub *Grub2) doFixSettingDistro() (needUpdate bool) {
	// fix GRUB_DISTRIBUTOR
	wantGrubDistroCmd := fmt.Sprintf("`%s -d -s 2>/dev/null || echo Debian`", lsbReleaseCmd)
	if grub.doGetSettingDistributor() != wantGrubDistroCmd {
		needUpdate = true
		grub.doSetSettingDistributor(wantGrubDistroCmd)
	}
	return
}

func (grub *Grub2) writeSettings() {
	fileContent := grub.getSettingContentToSave()
	if runWithoutDbus {
		doWriteGrubSettings(fileContent)
	} else {
		grub2extDoWriteGrubSettings(fileContent)
	}
}

func (grub *Grub2) parseEntries(fileContent string) (err error) {
	grub.resetEntries()

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
				grub.resetEntries()
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
				grub.resetEntries()
				err = fmt.Errorf("parse entry title failed from: %q", line)
				return
			}
		} else if strings.HasPrefix(line, "submenu ") {
			if inMenuEntry {
				grub.resetEntries()
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
				grub.resetEntries()
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
	grub.resetSettings()

	s := bufio.NewScanner(strings.NewReader(fileContent))
	s.Split(bufio.ScanLines)
	for s.Scan() {
		line := s.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "GRUB_") {
			a := strings.SplitN(line, "=", 2)
			if len(a) != 2 {
				continue
			}
			key, value := a[0], a[1]
			grub.settings[key] = unquoteString(value)
			logger.Debugf("found setting: %s=%s", a[0], a[1])
		}
	}
	if err := s.Err(); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// getAllEntriesLv1 return all entires titles in level one.
func (grub *Grub2) getEntryTitlesLv1() (entryTitles []string) {
	for _, entry := range grub.entries {
		if entry.parentSubMenu == nil {
			entryTitles = append(entryTitles, entry.getFullTitle())
		}
	}
	return
}

// return default entry or related entry title(such as "Deepin 2014
// GNU/Linux") if possible
func (grub *Grub2) getSettingDefaultEntry() (entry string) {
	entry = grub.doGetSettingDefaultEntry()
	if len(entry) == 0 {
		entry = defaultGrubDefaultEntry
	}

	// convert to simple stype
	entry = convertToSimpleEntry(entry)

	// if there is no entry titles, just return origin value
	entryTitles := grub.getEntryTitlesLv1()
	if len(entryTitles) == 0 {
		return
	}

	// if entry titles exists and the origin value is a valid title,
	// just return it
	if isStringInArray(entry, entryTitles) {
		return
	}

	// if entry titles exists and the origin value is an index number,
	// return it related title
	if i, err := strconv.ParseInt(entry, 10, 32); err == nil {
		if i >= 0 && int(i) < len(entryTitles) {
			entry = convertToSimpleEntry(entryTitles[i])
		}
	}
	return
}

func (grub *Grub2) setSettingDefaultEntry(title string) {
	// Convert the default entry value to index if possible for that
	// the disk number may be different between live-cd and normal
	// mode and the entry value will be invalid then.
	entriesLv1 := grub.getEntryTitlesLv1()
	if i := getStringIndexInArray(title, entriesLv1); i != -1 {
		title = strconv.Itoa(i)
	}
	grub.doSetSettingDefaultEntry(title)
	grub.writeSettings()
	grub.config.save()
}
func (grub *Grub2) doGetSettingDefaultEntry() string {
	return grub.settings["GRUB_DEFAULT"]
}
func (grub *Grub2) doSetSettingDefaultEntry(value string) {
	grub.settings["GRUB_DEFAULT"] = value
	grub.config.doSetDefaultEntry(value)
}

func (grub *Grub2) getSettingTimeout() (timeout int32) {
	timeout = defaultGrubTimeoutInt // default timeout
	timeoutStr := grub.doGetSettingTimeout()
	if len(timeoutStr) == 0 {
		return
	}
	timeout64, err := strconv.ParseInt(timeoutStr, 10, 32)
	if err != nil {
		logger.Errorf(`valid value, settings["GRUB_TIMEOUT"]=%s`, timeoutStr)
		return
	}
	timeout = int32(timeout64)
	return
}
func (grub *Grub2) setSettingTimeout(timeout int32) {
	grub.doSetSettingTimeoutLogic(timeout)
	grub.writeSettings()
	grub.config.save()
}
func (grub *Grub2) doSetSettingTimeoutLogic(timeout int32) {
	timeoutStr := strconv.FormatInt(int64(timeout), 10)
	grub.doSetSettingTimeout(timeoutStr)
}
func (grub *Grub2) doGetSettingTimeout() string {
	return grub.settings["GRUB_TIMEOUT"]
}
func (grub *Grub2) doSetSettingTimeout(value string) {
	grub.settings["GRUB_TIMEOUT"] = value
	grub.config.doSetTimeout(value)
}

func (grub *Grub2) getSettingGfxmode() string {
	gfxmode := grub.doGetSettingGfxmode()
	if len(gfxmode) == 0 {
		return defaultGrubGfxmode
	}
	return gfxmode
}
func (grub *Grub2) setSettingGfxmode(gfxmode string) {
	grub.doSetSettingGfxmode(gfxmode)
	grub.writeSettings()
	grub.config.save()
}
func (grub *Grub2) doGetSettingGfxmode() string {
	return grub.settings["GRUB_GFXMODE"]
}
func (grub *Grub2) doSetSettingGfxmode(value string) {
	grub.settings["GRUB_GFXMODE"] = value
	grub.config.doSetResolution(value)
}

func (grub *Grub2) getSettingTheme() string {
	return grub.doGetSettingTheme()
}

func (grub *Grub2) setEnableTheme(enable bool) {
	if enable {
		grub.setSettingTheme(themeMainFile)
	} else {
		grub.setSettingTheme("")
	}
}

func (grub *Grub2) setSettingTheme(themeFile string) {
	grub.doSetSettingTheme(themeFile)
	if themeFile == "" {
		grub.settings["GRUB_BACKGROUND"] = ""
	} else {
		grub.settings["GRUB_BACKGROUND"] = themeBgFile
	}
	grub.writeSettings()
}
func (grub *Grub2) doGetSettingTheme() string {
	return grub.settings["GRUB_THEME"]
}
func (grub *Grub2) doSetSettingTheme(value string) {
	grub.settings["GRUB_THEME"] = value
}

func (grub *Grub2) doGetSettingDistributor() string {
	return grub.settings["GRUB_DISTRIBUTOR"]
}
func (grub *Grub2) doSetSettingDistributor(value string) {
	grub.settings["GRUB_DISTRIBUTOR"] = value
}

func (grub *Grub2) getSettingContentToSave() string {
	// sort lines before saving
	lines := make(sort.StringSlice, 0)
	for k, v := range grub.settings {
		if len(v) > 0 {
			l := k + "=" + quoteString(v)
			lines = append(lines, l)
		}
	}
	lines.Sort()
	fileContent := ""
	for _, l := range lines {
		fileContent += l + "\n"
	}
	return fileContent
}
