/**
 * Copyright (C) 2014 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package search

import (
	"fmt"
	"regexp"
)

// item score
const (
	Poor         uint32 = 50000
	BelowAverage uint32 = 60000
	Average      uint32 = 70000
	AboveAverage uint32 = 75000
	Good         uint32 = 80000
	VeryGood     uint32 = 85000
	Excellent    uint32 = 90000
	Highest      uint32 = 100000
)

func addMatcher(template string, key string, score uint32, m map[*regexp.Regexp]uint32) error {
	r, err := regexp.Compile(fmt.Sprintf(template, key))
	if err != nil {
		return err
	}

	m[r] = score
	return nil
}

// learn from synpase
// TODO:
// 1. analyse the code of synapse much deeply.
// 2. add a weight for frequency.
func getMatchers(key string) map[*regexp.Regexp]uint32 {
	// * create a couple of regexes and try to help with matching
	// * match with these regular expressions (with descending score):
	// * 1) ^query$
	// * 2) ^query
	// * 3) \bquery
	// * 4) query
	m := make(map[*regexp.Regexp]uint32, 0)
	// ^query$
	addMatcher(`(?i)^(%s)$`, key, Highest, m)
	// ^query
	addMatcher(`(?i)^(%s)`, key, Excellent, m)
	// \bquery
	addMatcher(`(?i)\b(%s)`, key, VeryGood, m)
	// query
	addMatcher("(?i)(%s)", key, BelowAverage, m)

	return m
}
