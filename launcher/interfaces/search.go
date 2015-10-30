package interfaces

// SearchID is type for pinyin search.
type SearchID string

// Search is interface for search transaction.
type Search interface {
	Search(string, []ItemInfo)
	Cancel()
}

// PinYin is interface for pinyin search transaction.
type PinYin interface {
	Search(string) ([]string, error)
	IsValid() bool
}
