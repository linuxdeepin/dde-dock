/*
 * Copyright (C) 2015 ~ 2018 Deepin Technology Co., Ltd.
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

package dstore

import (
	"sort"
	"strings"
)

// IDList type alias for []string, used for sorting.
type IDList []string

func (list IDList) Less(i, j int) bool {
	return list[i] < list[j]
}

func (list IDList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list IDList) Len() int {
	return len(list)
}

func getUniqueCategory(categories []string) []string {
	candidateIDs := map[string]bool{OthersName: true}
	for _, category := range categories {
		candidateIDs[category] = true
	}

	if len(candidateIDs) > 1 {
		delete(candidateIDs, OthersName)
	}

	var ids []string
	for id := range candidateIDs {
		ids = append(ids, id)
	}

	sort.Sort(IDList(ids))

	return ids
}

type XCategoryQueryIDTransaction struct {
	data             map[string]string
	deepinCategories map[string]string
}

func NewXCategoryQueryIDTransaction(file, categoriesFile string) (*XCategoryQueryIDTransaction, error) {
	data, err := getXCategoryInfo(file)
	if err != nil {
		data = xcategoriesFallback
	}
	deepinCategories := map[string]string{}
	for _, c := range GetAllInfos(categoriesFile) {
		id := c.ID
		deepinCategories[id] = id
	}
	return &XCategoryQueryIDTransaction{
		data:             data,
		deepinCategories: deepinCategories,
	}, nil
}

func (t *XCategoryQueryIDTransaction) Query(strCategories string) (string, error) {
	categories := strings.Split(strings.TrimSuffix(strCategories, ";"), ";")
	categoryNames := make([]string, 0, len(categories))
	for _, category := range categories {
		if name, ok := t.data[strings.ToLower(category)]; ok {
			categoryNames = append(categoryNames, name)
		}
	}

	if len(categoryNames) != 0 {
		if c, ok := t.deepinCategories[categoryNames[0]]; ok {
			return c, nil
		}
	}

	possibleCategories := getUniqueCategory(categoryNames)
	return possibleCategories[0], nil
}
