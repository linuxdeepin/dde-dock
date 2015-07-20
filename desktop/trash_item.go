package desktop

import (
	"fmt"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/gio-2.0"
)

// TrashItem is TrashItem.
type TrashItem struct {
	*Item
}

// NewTrashItem creates new trash item.
func NewTrashItem(app *Application, uri string) *TrashItem {
	return &TrashItem{NewItem(app, []string{uri})}
}

// GenMenu generates json format menu content used in DeepinMenu for TrashItem.
func (item *TrashItem) GenMenu() (*Menu, error) {
	item.menu = NewMenu()
	clearMenuItemText := Tr("_Clear")

	trash := gio.FileNewForUri("trash://")
	info, err := trash.QueryInfo(gio.FileAttributeTrashItemCount, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		return nil, err
	}
	defer info.Unref()

	trashedItemCount := info.GetAttributeUint32(gio.FileAttributeTrashItemCount)
	if item.app.settings.ShowTrashedItemCountIsEnable() {
		clearMenuItemText = fmt.Sprintf(NTr("_Clear %d Item", "_Clear %d Items", int(trashedItemCount)), trashedItemCount)
	}

	return item.menu.AppendItem(NewMenuItem(Tr("_Open"), func() {
		item.app.displayFile("trash://")
	}, true)).AddSeparator().AppendItem(NewMenuItem(clearMenuItemText, func() {
		item.emitRequestEmptyTrash()
	}, trashedItemCount != 0)), nil
}
