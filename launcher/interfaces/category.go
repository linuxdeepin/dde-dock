package interfaces

type CategoryId int64

type CategoryInfoInterface interface {
	Id() CategoryId
	Name() string
	Items() []ItemId
	AddItem(ItemId)
	RemoveItem(ItemId)
}

type CategoryManagerInterface interface {
	AddItem(ItemId, CategoryId)
	RemoveItem(ItemId, CategoryId)
	GetAllCategory() []CategoryId
	GetCategory(id CategoryId) CategoryInfoInterface
}
