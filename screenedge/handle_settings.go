package screenedge

import (
	"gir/gio-2.0"
)

func handleSettingsChanged() {
	setting := zoneSettings()
	setting.Connect("changed", func(s *gio.Settings, key string) {
		switch key {
		case "left-up":
			edgeActionMap[leftTopEdge] = setting.GetString(key)
		case "left-down":
			edgeActionMap[leftBottomEdge] = setting.GetString(key)
		case "right-up":
			edgeActionMap[rightTopEdge] = setting.GetString(key)
		case "right-down":
			edgeActionMap[rightBottomEdge] = setting.GetString(key)
		}
	})
	setting.GetString("left-up")
}
