// Package main provides ...
package dock

import (
	C "launchpad.net/gocheck"
	"os"
)

type MenuTestSuite struct{}

var _ = C.Suite(&MenuTestSuite{})

func (m *MenuTestSuite) TestGenerateMenuJson(c *C.C) {
	old := os.Getenv("LANGUAGE")
	os.Setenv("LANGUAGE", "zh_CN")
	f := NewNormalApp("./testdata/firefox.desktop")
	if f == nil {
		c.Skip("get test data failed")
		return
	}
	f.buildMenu()
	os.Setenv("LANGUAGE", old)

	c.Check(f.Menu, C.Equals, "{\"checkableMenu\":false, \"singleCheck\": false, \"items\":[{\"itemId\":\"1\", \"itemText\": \"_Run\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"2\", \"itemText\": \"新建窗口\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"3\", \"itemText\": \"新建隐私浏览窗口\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"4\", \"itemText\": \"_Undock\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}}]}")
}
