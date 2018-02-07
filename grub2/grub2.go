/*
 * Copyright (C) 2017 ~ 2018 Deepin Technology Co., Ltd.
 *
 * Author:     jouyouyun <jouyouwen717@gmail.com>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package grub2

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"unicode"

	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/log"
)

var logger *log.Logger

func SetLogger(v *log.Logger) {
	logger = v
}

type Grub2 struct {
	modifyFuncChan  chan ModifyFunc
	mkconfigManager *MkconfigManager
	entries         []Entry
	theme           *Theme
	setPropMu       sync.Mutex

	// props:
	DefaultEntry string
	EnableTheme  bool
	Resolution   string
	Timeout      uint32

	Updating bool
}

const (
	propNameDefaultEntry = "DefaultEntry"
	propNameEnableTheme  = "EnableTheme"
	propNameResolution   = "Resolution"
	propNameTimeout      = "Timeout"
	propNameUpdating     = "Updating"
)

// return -1 for failed
func (g *Grub2) defaultEntryStr2Idx(str string) int {
	entriesLv1 := g.getEntryTitlesLv1()
	return getStringIndexInArray(str, entriesLv1)
}

func (g *Grub2) defaultEntryIdx2Str(idx int) (string, error) {
	entriesLv1 := g.getEntryTitlesLv1()
	length := len(entriesLv1)
	if length == 0 {
		return "", errors.New("no entry")
	}
	if 0 <= idx && idx < length {
		return entriesLv1[idx], nil
	} else {
		return "", errors.New("index out of range")
	}
}

func (g *Grub2) applyParams(params map[string]string) {
	//timeout
	timeout := getTimeout(params)
	if timeout < 0 {
		timeout = 999
	}
	g.Timeout = uint32(timeout)

	// enable theme
	var enableTheme bool
	theme := getTheme(params)
	if theme != "" {
		enableTheme = true
	}
	g.EnableTheme = enableTheme

	// resolution
	g.Resolution = getGfxMode(params)

	// default entry
	defaultEntry := getDefaultEntry(params)

	defaultEntryIdx, err := strconv.Atoi(defaultEntry)
	if err == nil {
		// is a num
		g.DefaultEntry, _ = g.defaultEntryIdx2Str(defaultEntryIdx)
	} else {
		// not a num
		if defaultEntry == "saved" {
			// TODO saved
			g.DefaultEntry, _ = g.defaultEntryIdx2Str(0)
		} else {
			g.DefaultEntry = defaultEntry
		}
	}
}

type ModifyFunc func(map[string]string)

func getModifyFuncEnableTheme(enable bool) ModifyFunc {
	if enable {
		return func(params map[string]string) {
			params[grubTheme] = quoteString(defaultGrubTheme)
			params[grubBackground] = quoteString(defaultGrubBackground)
		}
	} else {
		return func(params map[string]string) {
			delete(params, grubTheme)
			delete(params, grubBackground)
		}
	}
}

func getModifyFuncTimeout(timeout uint32) ModifyFunc {
	return func(params map[string]string) {
		params[grubTimeout] = strconv.Itoa(int(timeout))
	}
}

func getModifyFuncResolution(val string) ModifyFunc {
	return func(params map[string]string) {
		params[grubGfxMode] = quoteString(val)
	}
}

func getModifyFuncDefaultEntry(idx int) ModifyFunc {
	return func(params map[string]string) {
		params[grubDefault] = strconv.Itoa(idx)
	}
}

func New() *Grub2 {
	g := &Grub2{}

	g.readEntries()

	params, err := loadGrubParams()
	if err != nil {
		logger.Warning(err)
	}
	paramsMD5Sum := getGrubParamsMD5Sum(params)
	logger.Debug("paramsHash:", paramsMD5Sum)

	var needCallMkconfig bool
	if mkconfigLog, err := loadLog(); err != nil {
		needCallMkconfig = true
	} else {
		ok, _ := mkconfigLog.Verify(paramsMD5Sum)
		if !ok {
			needCallMkconfig = true
		}
	}

	g.applyParams(params)
	g.modifyFuncChan = make(chan ModifyFunc)
	g.mkconfigManager = newMkconfigManager(g.modifyFuncChan, func(running bool) {
		// state change callback
		if g.Updating != running {
			g.Updating = running
			dbus.NotifyChange(g, propNameUpdating)
		}
	})
	go g.mkconfigManager.loop()

	// init theme
	g.theme = NewTheme(g)
	g.theme.initTheme()
	go g.theme.regenerateBackgroundIfNeed()

	if needCallMkconfig {
		g.modifyFuncChan <- func(_ map[string]string) {
			// NoOp
		}
	}

	return g
}

func (grub *Grub2) readEntries() (err error) {
	fileContent, err := ioutil.ReadFile(grubScriptFile)
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
		logger.Warning("there is no menu entry in %s", grubScriptFile)
	}
	return
}

func (grub *Grub2) resetEntries() {
	grub.entries = make([]Entry, 0)
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
			title, ok := parseTitle(line)
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
			title, ok := parseTitle(line)
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

var (
	entryRegexpSingleQuote = regexp.MustCompile(`^ *(menuentry|submenu) +'(.*?)'.*$`)
	entryRegexpDoubleQuote = regexp.MustCompile(`^ *(menuentry|submenu) +"(.*?)".*$`)
)

func parseTitle(line string) (string, bool) {
	line = strings.TrimLeftFunc(line, unicode.IsSpace)
	if entryRegexpSingleQuote.MatchString(line) {
		return entryRegexpSingleQuote.FindStringSubmatch(line)[2], true
	} else if entryRegexpDoubleQuote.MatchString(line) {
		return entryRegexpDoubleQuote.FindStringSubmatch(line)[2], true
	} else {
		return "", false
	}
}

func (g *Grub2) getScreenWidthHeight() (w, h uint16, err error) {
	return parseResolution(g.Resolution)
}

func (g *Grub2) canSafelyExit() bool {
	logger.Debug("call canSafelyExit")
	if g.mkconfigManager.IsRunning() || g.theme.Updating {
		return false
	}
	return true
}
