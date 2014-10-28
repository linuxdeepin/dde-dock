package setting

type CategoryDisplayMode int64

const (
	CategoryDisplayModeUnknown CategoryDisplayMode = iota - 1
	CategoryDisplayModeIcon
	CategoryDisplayModeText

	CategoryDisplayModeKey string = "category-display-mode"
)

func (c CategoryDisplayMode) String() string {
	switch c {
	case CategoryDisplayModeUnknown:
		return "unknown category display mode"
	case CategoryDisplayModeText:
		return "display text mode"
	case CategoryDisplayModeIcon:
		return "display icon mode"
	default:
		return "unknown mode"
	}
}
