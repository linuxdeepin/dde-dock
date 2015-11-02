package category

import (
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/lib/gio-2.0"
	"sort"
	"strings"
)

type XCategoryQueryIDTransition struct {
}

func NewXCategoryQueryIDTransition(string) (*XCategoryQueryIDTransition, error) {
	return &XCategoryQueryIDTransition{}, nil
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

func (transition *XCategoryQueryIDTransition) Query(app *gio.DesktopAppInfo) (CategoryID, error) {
	categories := strings.Split(strings.TrimRight(app.GetCategories(), ";"), ";")
	return getXCategory(categories), nil
}

func (transition *XCategoryQueryIDTransition) Free() {
}
