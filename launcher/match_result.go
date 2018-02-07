/*
 * Copyright (C) 2016 ~ 2018 Deepin Technology Co., Ltd.
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
