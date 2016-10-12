/**
 * Copyright (C) 2016 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package launcher

import (
	"fmt"
	"sort"
)

type MatchResult struct {
	score SearchScore
	item  *Item
}

func (r *MatchResult) String() string {
	if r == nil {
		return "<nil>"
	}
	return fmt.Sprintf("<MatchResult item=%v score=%v>", r.item.ID, r.score)
}

type MatchResults []*MatchResult

// impl sort interface
func (p MatchResults) Len() int           { return len(p) }
func (p MatchResults) Less(i, j int) bool { return p[i].score < p[j].score }
func (p MatchResults) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (results MatchResults) GetTruncatedOrderedIDs() []string {
	sort.Sort(sort.Reverse(results))
	const idsMaxLen = 42

	idsLen := len(results)
	if idsLen > idsMaxLen {
		idsLen = idsMaxLen
	}
	ids := make([]string, idsLen)

	for i := 0; i < idsLen; i++ {
		ids[i] = results[i].item.ID
	}
	return ids
}

func (results MatchResults) Copy() MatchResults {
	resultsCopy := make(MatchResults, len(results))
	copy(resultsCopy, results)
	return resultsCopy
}
