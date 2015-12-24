package category

import (
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"sort"
	"strings"
)

type XCategoryQueryIDTransaction struct {
	data map[string]string
}

func NewXCategoryQueryIDTransaction(file string) (*XCategoryQueryIDTransaction, error) {
	data, err := getXCategoryInfo(file)
	if err != nil {
		return nil, err
	}
	return &XCategoryQueryIDTransaction{data: data}, nil
}

// IDList type alias for []CategoryID, used for sorting.
type IDList []CategoryID

func (list IDList) Less(i, j int) bool {
	return list[i] < list[j]
}

func (list IDList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list IDList) Len() int {
	return len(list)
}

func getXCategory(categories []string) CategoryID {
	candidateIDs := map[CategoryID]bool{OthersID: true}
	for _, category := range categories {
		if id, err := getCategoryID(category); err == nil {
			candidateIDs[id] = true
		}
	}

	if len(candidateIDs) > 1 && candidateIDs[OthersID] {
		delete(candidateIDs, OthersID)
	}

	var ids []CategoryID
	for id := range candidateIDs {
		ids = append(ids, id)
	}

	sort.Sort(IDList(ids))

	return ids[0]
}

func (t *XCategoryQueryIDTransaction) Query(strCategories string) (CategoryID, error) {
	categories := strings.Split(strings.TrimRight(strCategories, ";"), ";")
	categoryNames := make([]string, 0, len(categories))
	for _, category := range categories {
		if name, ok := t.data[strings.ToLower(category)]; ok {
			categoryNames = append(categoryNames, name)
		}
	}
	return getXCategory(categoryNames), nil
}

func (t *XCategoryQueryIDTransaction) Free() {
}
