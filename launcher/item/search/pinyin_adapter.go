package search

import (
	pinyin "dbus/com/deepin/daemon/search"
	. "pkg.deepin.io/dde-daemon/launcher/interfaces"
)

type PinYinSearchAdapter struct {
	searchObj *pinyin.Search
	searchId  SearchId
}

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

func (p *PinYinSearchAdapter) Init(data []string) error {
	searchId, _, err := p.searchObj.NewSearchWithStrList(data)
	p.searchId = SearchId(searchId)

	return err
}

func (p *PinYinSearchAdapter) Search(key string) ([]string, error) {
	return p.searchObj.SearchString(key, string(p.searchId))
}

func (p *PinYinSearchAdapter) IsValid() bool {
	return p.searchId != SearchId("")
}
