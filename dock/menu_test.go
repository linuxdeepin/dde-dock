// Package main provides ...
package dock

import (
	"testing"
)

func TestGenerateMenuJson(t *testing.T) {
	initDeepin()
	f := NewNormalApp("./test-data/firefox.desktop")
	f.buildMenu()
	if f.Menu != "{\"checkableMenu\":false, \"singleCheck\": false, \"items\":[{\"itemId\":\"1\", \"itemText\": \"_Run\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"2\", \"itemText\": \"新建窗口\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"3\", \"itemText\": \"新建隐私浏览窗口\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}},{\"itemId\":\"4\", \"itemText\": \"_Undock\", \"isActive\": true, \"isCheckable\":false, \"checked\":false, \"itemIcon\":\"\", \"itemIconHover\":\"\", \"itemIconInactive\":\"\", \"showCheckMark\":false, \"itemSubMenu\":{\"checkableMenu\":false, \"singleCheck\":false, \"items\": []}}]}" {
		t.Fail()
	}
}
