package interfaces

type SearchId string

type SearchInterface interface {
	Search(string, []ItemInfoInterface)
	Cancel()
}

type PinYinInterface interface {
	Search(string) ([]string, error)
	IsValid() bool
}
