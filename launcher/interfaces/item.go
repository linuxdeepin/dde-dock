package interfaces

type ItemId string

type ItemInfoInterface interface {
	Name() string
	Icon() string
	Path() string
	Id() ItemId
	ExecCmd() string
	Description() string
	EnName() string
	GenericName() string
	Keywords() []string
	GetCategoryId() CategoryId
	SetCategoryId(CategoryId)
	GetTimeInstalled() int64
	SetTimeInstalled(int64)
}
