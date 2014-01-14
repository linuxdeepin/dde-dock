package main

type EntryType int

const (
	MENUENTRY EntryType = iota
	SUBMENU
)

// type EntryTitle struct{}
type Entry struct {
	entryType     EntryType
	title         string
	num           int
	parentSubMenu *Entry
}

func (entry *Entry) getFullTitle() string {
	// TODO
	if entry.parentSubMenu != nil {
		return entry.parentSubMenu.getFullTitle() + ">" + entry.title
	} else {
		return entry.title
	}
}
