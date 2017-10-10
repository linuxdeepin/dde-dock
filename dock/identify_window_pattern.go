/*
 * Copyright (C) 2017 ~ 2017 Deepin Technology Co., Ltd.
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

package dock

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type WindowPatterns []WindowPattern

type WindowPattern struct {
	Rules       []WindowRule `json:"rules"`
	Result      string       `json:"ret"`
	ParsedRules []*WindowRuleParsed
}

type WindowRule [2]string
type WindowRuleParsed struct {
	Key         string
	ValueParsed *RuleValueParsed
}

func (rule *WindowRule) Parse() *WindowRuleParsed {
	key, value := rule[0], rule[1]
	return &WindowRuleParsed{
		Key:         key,
		ValueParsed: parseRuleValue(value),
	}
}

func loadWindowPatterns(file string) (WindowPatterns, error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var patterns WindowPatterns
	err = json.Unmarshal(content, &patterns)
	if err != nil {
		return nil, err
	}
	logger.Debugf("loadWindowPatterns: ok count %d", len(patterns))

	// parse pattterns
	for i := range patterns {
		pattern := &patterns[i]
		rules := pattern.Rules
		// parse rules in pattern
		pattern.ParsedRules = make([]*WindowRuleParsed, len(rules))
		for j := range rules {
			rule := &rules[j]
			pattern.ParsedRules[j] = rule.Parse()
		}
	}

	return patterns, nil
}

func (patterns WindowPatterns) Match(winInfo *WindowInfo) string {
	for i := range patterns {
		pattern := &patterns[i]
		rules := pattern.ParsedRules
		patternOk := true

		logger.Debugf("try pattern %d", i)
		for j := range rules {
			rule := rules[j]
			ok := rule.Match(winInfo)
			if !ok {
				// pattern match fail
				patternOk = false
				break
			}
		}

		if patternOk {
			// pattern match success
			logger.Debugf("pattern match success")
			return pattern.Result
		}
	}
	// fail
	return ""
}

func parseRuleKey(winInfo *WindowInfo, key string) string {
	switch key {
	case "hasPid":
		if winInfo.process != nil && winInfo.process.hasPid {
			return "t"
		}
		return "f"
	case "exec":
		// executable file base name
		if winInfo.process != nil {
			return filepath.Base(winInfo.process.exe)
		}
	case "arg":
		// command line arguments
		if winInfo.process != nil {
			return strings.Join(winInfo.process.args, " ")
		}
	// xprop example:
	// WM_CLASS(STRING) = "dman", "DManual"
	// wmClass.Instance is dman
	// wmClass.Class is DManual
	case "wmi":
		// wmClass.Instance
		if winInfo.wmClass != nil {
			return winInfo.wmClass.Instance
		}
	case "wmc":
		// wmClass.Class
		if winInfo.wmClass != nil {
			return winInfo.wmClass.Class
		}
	case "wmn":
		// WM_NAME
		return winInfo.wmName
	case "wmrole":
		// WM_ROLE
		return winInfo.wmRole

	default:
		const envPrefix = "env."
		if strings.HasPrefix(key, envPrefix) {
			envName := key[len(envPrefix):]
			if winInfo.process != nil {
				return winInfo.process.environ.Get(envName)
			}
		}
	}
	return ""
}

func (rule *WindowRuleParsed) Match(winInfo *WindowInfo) bool {
	keyParsed := parseRuleKey(winInfo, rule.Key)
	fn := rule.ValueParsed.Fn
	if fn == nil {
		logger.Warningf("WindowRule.Match: badRuleValue %q", rule.ValueParsed.Original)
		return false
	}
	result := fn(keyParsed)
	logger.Debugf("%s %q %v ? %v", rule.Key, keyParsed, rule.ValueParsed, result)
	return result
}

type RuleValueParsed struct {
	Fn       RuleMatchFunc
	Type     byte
	Flags    RuleValueParsedFlag
	Original string
	Value    string
}
type RuleMatchFunc func(string) bool

type RuleValueParsedFlag uint

const (
	RuleValueParsedFlagNone RuleValueParsedFlag = 1 << iota
	RuleValueParsedFlagNegative
	RuleValueParsedFlagIgnoreCase
)

func (p *RuleValueParsed) String() string {
	var buf bytes.Buffer
	if p.Fn == nil {
		return "bad rule"
	}

	if p.Flags&RuleValueParsedFlagNegative != 0 {
		buf.WriteString("not ")
	}

	var typeDesc string
	switch p.Type {
	case '=', 'e', 'E':
		typeDesc = "equal"
	case 'c', 'C':
		typeDesc = "contains"
	case 'r', 'R':
		typeDesc = "match regexp "
	default:
		typeDesc = "<unknown type>"
	}
	buf.WriteString(typeDesc)

	if p.Flags&RuleValueParsedFlagIgnoreCase != 0 {
		buf.WriteString(" (ignore case)")
	}
	// "$k not equal (ignore case) $p.value"
	fmt.Fprintf(&buf, " %q", p.Value)
	return buf.String()
}

func negativeRule(fn RuleMatchFunc) RuleMatchFunc {
	return func(k string) bool {
		return !fn(k)
	}
}

var regexpCache map[string]*regexp.Regexp = make(map[string]*regexp.Regexp)
var regexpCacheMutex sync.Mutex

func getRegexp(expr string) *regexp.Regexp {
	regexpCacheMutex.Lock()
	defer regexpCacheMutex.Unlock()

	reg, ok := regexpCache[expr]
	if ok {
		return reg
	}
	reg, err := regexp.Compile(expr)
	if err != nil {
		logger.Warning(err)
	}
	regexpCache[expr] = reg
	return reg
}

// "=:XXX" equal XXX
// "=!XXX" not equal XXX

// "c:XXX" contains XXX
// "c!XXX" not contains XXX

// "r:XXX" match regexp XXX
// "r!XXX" not match regexp XXX

// e c r ignore case
// = E C R not ignore case
func parseRuleValue(val string) *RuleValueParsed {
	var ret = &RuleValueParsed{
		Original: val,
	}
	if len(val) < 2 {
		return ret
	}
	var negative bool
	switch val[1] {
	case ':':
	case '!':
		ret.Flags |= RuleValueParsedFlagNegative
		negative = true
	default:
		return ret
	}
	// type
	value := val[2:]
	ret.Value = value
	ret.Type = val[0]

	var fn RuleMatchFunc
	switch val[0] {
	case 'C':
		fn = func(k string) bool {
			return strings.Contains(k, value)
		}
	case 'c':
		ret.Flags |= RuleValueParsedFlagIgnoreCase
		fn = func(k string) bool {
			return strings.Contains(
				strings.ToLower(k),
				strings.ToLower(value))
		}
	case '=', 'E':
		fn = func(k string) bool {
			return k == value
		}
	case 'e':
		ret.Flags |= RuleValueParsedFlagIgnoreCase
		fn = func(k string) bool {
			return strings.EqualFold(k, value)
		}

	case 'R':
		fn = func(k string) bool {
			reg := getRegexp(value)
			if reg == nil {
				return false
			}
			return reg.MatchString(k)
		}
	case 'r':
		ret.Flags |= RuleValueParsedFlagIgnoreCase
		fn = func(k string) bool {
			reg := getRegexp("(?i)" + value)
			if reg == nil {
				return false
			}
			return reg.MatchString(k)
		}

	default:
		return ret
	}
	if negative {
		ret.Fn = negativeRule(fn)
	} else {
		ret.Fn = fn
	}
	return ret
}
