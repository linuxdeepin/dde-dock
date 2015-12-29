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
	data map[string]string
}

func NewXCategoryQueryIDTransaction(file string) (*XCategoryQueryIDTransaction, error) {
	data, err := getXCategoryInfo(file)
	if err != nil {
		data = xcategoriesFallback
	}
	return &XCategoryQueryIDTransaction{data: data}, nil
}

func (t *XCategoryQueryIDTransaction) Query(strCategories string) (string, error) {
	categories := strings.Split(strings.TrimSuffix(strCategories, ";"), ";")
	categoryNames := make([]string, 0, len(categories))
	for _, category := range categories {
		if name, ok := t.data[strings.ToLower(category)]; ok {
			categoryNames = append(categoryNames, name)
		}
	}
	possibleCategories := getUniqueCategory(categoryNames)
	return possibleCategories[0], nil
}
