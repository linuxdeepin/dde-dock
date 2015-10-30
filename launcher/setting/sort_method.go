package setting

// SortMethod type for sort method.
type SortMethod int64

// sort method.
const (
	SortMethodUnknown SortMethod = iota - 1
	SortMethodByName
	SortMethodByCategory
	SortMethodByTimeInstalled
	SortMethodByFrequency

	SortMethodkey string = "sort-method"
)

func (s SortMethod) String() string {
	switch s {
	case SortMethodUnknown:
		return "unknown sort method"
	case SortMethodByName:
		return "sort by name"
	case SortMethodByCategory:
		return "sort by category"
	case SortMethodByTimeInstalled:
		return "sort by time installed"
	case SortMethodByFrequency:
		return "sort by frequency"
	default:
		return "unknown sort method"
	}
}
