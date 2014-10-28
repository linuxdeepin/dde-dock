package search

import (
	C "launchpad.net/gocheck"
	"sort"
)

type SearchResultListTestSuite struct {
}

var _ = C.Suite(&SearchResultListTestSuite{})

func (self *SearchResultListTestSuite) _TestSearchResultList(c *C.C) {
	res := SearchResultList{
		SearchResult{
			Id:    "chrome",
			Name:  "chrome",
			Score: 345000,
		},
		SearchResult{
			Id:    "weibo",
			Name:  "weibo",
			Score: 80000,
		},
		SearchResult{
			Id:    "music",
			Name:  "music",
			Score: 80000,
		},
	}
	c.Assert(res.Len(), C.Equals, 3)
	c.Assert(string(res[0].Id), C.Equals, "chrome")
	c.Assert(string(res[1].Id), C.Equals, "weibo")
	c.Assert(string(res[2].Id), C.Equals, "music")
	c.Assert(res.Less(0, 1), C.Equals, true)
	c.Assert(res.Less(1, 2), C.Equals, false)

	res.Swap(0, 1)
	c.Assert(string(res[0].Id), C.Equals, "weibo")

	sort.Sort(res)
	c.Assert(string(res[0].Id), C.Equals, "chrome")
	c.Assert(string(res[1].Id), C.Equals, "music")
	c.Assert(string(res[2].Id), C.Equals, "weibo")
}

func (self *SearchResultListTestSuite) TestSearchResultListReal(c *C.C) {
	list := SearchResultList{
		{"12306", "12306", 80000},
		{"google-chrome", "Google Chrome", 345000},
		{"chrome-lbfehkoinhhcknnbdgnnmjhiladcgbol-Default", "Evernote Web", 150000},
		{"chrome-kidnkfckhbdkfgbicccmdggmpgogehop-Default", "马克飞象", 150000},
		{"doit-im", "Doit.im", 80000},
		{"towerim", "Tower.im", 80000},
		{"microsoft-skydrive", "微软 SkyDrive", 80000},
		{"sina-weibo", "新浪微博", 80000},
		{"youdao-note", "有道云笔记", 80000},
		{"pirateslovedaisies", "海盗爱菊花", 80000},
		{"baidu-music", "百度音乐", 80000},
		{"xiami-music", "虾米音乐", 80000},
		{"kuwo-music", "酷我音乐网页版", 80000},
		{"kugou-music", "酷狗音乐", 80000},
		{"kingsoft-fast-docs", "金山快写", 80000},
		{"kingsoft-online-storage", "金山网盘", 80000},
	}

	sort.Sort(list)
	c.Assert(string(list[0].Id), C.Equals, "google-chrome")
}
