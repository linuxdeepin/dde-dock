package desktop

import (
	. "pkg.deepin.io/lib/gettext"
)

// ComputerItem is computer item on desktop.
type ComputerItem struct {
	*Item
}

// NewComputerItem creates new computer item.
func NewComputerItem(app *Application, uri string) *ComputerItem {
	return &ComputerItem{NewItem(app, []string{uri})}
}

// GenMenuContent generates json format menu content used in DeepinMenu for ComputerItem.
func (item *ComputerItem) GenMenuContent() (*Menu, error) {
	return item.menu.AppendItem(NewMenuItem(Tr("_Open"), func() {
		item.app.displayFile("computer://")
	}, true)).AddSeparator().AppendItem(NewMenuItem(Tr("_Properties"), func() {
		item.showProperties()
	}, true)), nil
}
