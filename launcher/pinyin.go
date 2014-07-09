package launcher

import (
	pinyin "dbus/com/deepin/daemon/search"
)

type PinYinSearch struct {
	searchObj *pinyin.Search
	searchId  string
}

var pinyinSearchObj *PinYinSearch = nil

func NewPinYinSearch(data []string) (*PinYinSearch, error) {
	searchObj, err := pinyin.NewSearch("com.deepin.daemon.Search", "/com/deepin/daemon/Search")
	if err != nil {
		return nil, err
	}
	obj := &PinYinSearch{searchObj, ""}
	err = obj.Init(data)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (p *PinYinSearch) Init(data []string) error {
	var err error
	p.searchId, _, err = p.searchObj.NewSearchWithStrList(data)

	if err == nil {
		logger.Debug("search object id:", p.searchId)
	}
	return err
}

func (p *PinYinSearch) Search(key string) ([]string, error) {
	return p.searchObj.SearchString(key, p.searchId)
}

func (p *PinYinSearch) IsValid() bool {
	return p.searchId != ""
}
