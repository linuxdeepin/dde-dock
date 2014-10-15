// Package main provides ...
package dock

import (
	C "launchpad.net/gocheck"
)

type MenuTestSuite struct{}

var _ = C.Suite(&MenuTestSuite{})

func (m *MenuTestSuite) TestGenerateMenuJson(c *C.C) {
	initDeepin()
	f := NewNormalApp("./test-data/firefox.desktop")
	f.buildMenu()

	c.Check(f.Menu, C.Equals, "{\"checkableMenu\":false, \"singleCheck\": false, \"items\":[{\"itemId\":\"1\", \"itemText\": \"_Run\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"2\", \"itemText\": \"新建窗口\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"3\", \"itemText\": \"新建隐私浏览窗口\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"4\", \"itemText\": \"_Undock\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}}]}")
}
