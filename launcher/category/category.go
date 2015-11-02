package category

import (
	"errors"
	"strings"

	. "pkg.deepin.io/dde/daemon/launcher/interfaces"
	"pkg.deepin.io/lib/gettext"
)

// category id and name.
const (
	OthersID CategoryID = iota - 2
	AllID
	InternetID
	OfficeID
	DevelopmentID
	ReadingID
	GraphicsID
	GameID
	MusicID
	SystemID
	VideoID
	ChatID

	AllName         = "all"
	OthersName      = "others"
	InternetName    = "internet"
	OfficeName      = "office"
	DevelopmentName = "development"
	ReadingName     = "reading"
	GraphicsName    = "graphics"
	GameName        = "game"
	MusicName       = "music"
	SystemName      = "system"
	VideoName       = "video"
	ChatName        = "chat"
)

var (
	categoryNameTable = map[string]CategoryID{
		OthersName:      OthersID,
		AllName:         AllID,
		InternetName:    InternetID,
		OfficeName:      OfficeID,
		DevelopmentName: DevelopmentID,
		ReadingName:     ReadingID,
		GraphicsName:    GraphicsID,
		GameName:        GameID,
		MusicName:       MusicID,
		SystemName:      SystemID,
		VideoName:       VideoID,
		ChatName:        ChatID,
	}
)

// Info for category.
type Info struct {
	id    CategoryID
	name  string
	items map[ItemID]struct{}
}

func NewInfo(id CategoryID, name string) *Info {
	return &Info{
		id:    id,
		name:  name,
		items: map[ItemID]struct{}{},
	}
}

// ID returns category id.
func (c *Info) ID() CategoryID {
	return c.id
}

// Name returns category english name.
func (c *Info) Name() string {
	return c.name
}

// LocaleName returns category's locale name.
func (c *Info) LocaleName() string {
	return gettext.Tr(c.name)
}

// AddItem adds a new app.
func (c *Info) AddItem(itemID ItemID) {
	c.items[itemID] = struct{}{}
}

// RemoveItem removes a app.
func (c *Info) RemoveItem(itemID ItemID) {
	delete(c.items, itemID)
}

// Items returns all items belongs to this category.
func (c *Info) Items() []ItemID {
	items := []ItemID{}
	for itemID := range c.items {
		items = append(items, itemID)
	}
	return items
}

func getCategoryID(name string) (CategoryID, error) {
	name = strings.ToLower(name)
	if id, ok := categoryNameTable[name]; ok {
		return id, nil
	}
	return OthersID, errors.New("unknown id")
}
