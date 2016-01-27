package category

import (
	"errors"
	"sort"
	"strings"
	"sync"

	"pkg.deepin.io/dde/daemon/dstore"
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

	AllName         = dstore.AllName
	OthersName      = dstore.OthersName
	InternetName    = dstore.InternetName
	OfficeName      = dstore.OfficeName
	DevelopmentName = dstore.DevelopmentName
	ReadingName     = dstore.ReadingName
	GraphicsName    = dstore.GraphicsName
	GameName        = dstore.GameName
	MusicName       = dstore.MusicName
	SystemName      = dstore.SystemName
	VideoName       = dstore.VideoName
	ChatName        = dstore.ChatName
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
	lock  sync.RWMutex
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
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items[itemID] = struct{}{}
}

// RemoveItem removes a app.
func (c *Info) RemoveItem(itemID ItemID) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.items, itemID)
}

type ByItemID []ItemID

func (items ByItemID) Swap(i, j int) {
	items[i], items[j] = items[j], items[i]
}

func (items ByItemID) Len() int {
	return len(items)
}

func (items ByItemID) Less(i, j int) bool {
	return items[i] < items[j]
}

// Items returns all items belongs to this category.
func (c *Info) Items() []ItemID {
	items := []ItemID{}
	c.lock.RLock()
	for itemID := range c.items {
		items = append(items, itemID)
	}
	c.lock.RUnlock()
	sort.Sort(ByItemID(items))
	return items
}

func getCategoryID(name string) (CategoryID, error) {
	name = strings.ToLower(name)
	if id, ok := categoryNameTable[name]; ok {
		return id, nil
	}
	return OthersID, errors.New("unknown id")
}
