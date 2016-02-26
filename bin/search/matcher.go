/**
 * Copyright (C) 2013 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Result score && Sorted by it
const (
	POOR          uint32 = 50000
	BELOW_AVERAGE        = 60000
	AVERAGE              = 70000
	ABOVE_AVERAGE        = 75000
	GOOD                 = 80000
	VERY_GOOD            = 85000
	EXCELLENT            = 90000
	HIGHEST              = 100000
)

func addMatcher(template, key string,
	score uint32, m map[*regexp.Regexp]uint32) error {
	reg, err := regexp.Compile(fmt.Sprintf(template, key))
	if err != nil {
		return err
	}

	m[reg] = score
	return nil
}

// learn from synapse
func getMatchers(key string) map[*regexp.Regexp]uint32 {
	// * create a couple of regexes and try to help with matching
	// * match with these regular expressions (with descending score):
	// * 1) ^query$
	// * 2) ^query
	// * 3) \bquery
	// * 4) split to words and seach \bword1.+\bword2 (if there are 2+ words)
	// * 5) query
	// * 6) split to characters and search \bq.+\bu.+\be.+\br.+\by
	// * 7) split to characters and search \bq.*u.*e.*r.*y
	m := make(map[*regexp.Regexp]uint32)

	addMatcher(`(?i)^(%s)$`, key, HIGHEST, m)
	addMatcher(`(?i)^(%s)`, key, EXCELLENT, m)
	addMatcher(`(?i)\b(%s)`, key, VERY_GOOD, m)

	words := strings.Fields(key)
	if len(words) > 1 {
		addMatcher(`(?i)\b(%s)`, strings.Join(words, `).+\b(`),
			GOOD, m)
	}

	addMatcher(`(?i)(%s)`, key, BELOW_AVERAGE, m)

	charSpliter, err := regexp.Compile(`\s*`)
	if err != nil {
		logger.Warning("Get char spliter failed:", err)
		return m
	}

	chars := charSpliter.Split(key, -1)
	if len(words) == 1 && len(chars) <= 5 {
		addMatcher(`(?i)\b(%s)`, strings.Join(chars, `).+\b(`),
			ABOVE_AVERAGE, m)
	}

	addMatcher(`(?i)\b(%s)`, strings.Join(chars, `).*(`), POOR, m)

	return m
}
