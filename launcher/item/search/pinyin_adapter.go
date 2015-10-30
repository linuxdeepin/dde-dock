package search

import (
	pinyin "dbus/com/deepin/daemon/search"
	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
)

// PinYinSearchAdapter is a adapter struct for pinyin.Search.
type PinYinSearchAdapter struct {
	searchObj *pinyin.Search
	searchID  SearchID
}

// NewPinYinSearchAdapter creates a new PinYinSearchAdapter object according to data.
func NewPinYinSearchAdapter(data []string) (*PinYinSearchAdapter, error) {
	searchObj, err := pinyin.NewSearch("com.deepin.daemon.Search", "/com/deepin/daemon/Search")
	if err != nil {
		return nil, err
	}
	obj := &PinYinSearchAdapter{searchObj, ""}
	err = obj.Init(data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

// Init initializes object with data.
func (p *PinYinSearchAdapter) Init(data []string) error {
	searchID, _, err := p.searchObj.NewSearchWithStrList(data)
	p.searchID = SearchID(searchID)

	return err
}

// Search executes transaction and returns found objects.
func (p *PinYinSearchAdapter) Search(key string) ([]string, error) {
	return p.searchObj.SearchString(key, string(p.searchID))
}

// IsValid returns true if this object is ok to use.
func (p *PinYinSearchAdapter) IsValid() bool {
	return p.searchID != SearchID("")
}
