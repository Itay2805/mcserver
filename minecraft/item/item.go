package item

type Item struct {
	ID			int
	Name		string
	StackSize	int
}

func GetById(id int) *Item {
	return items[id]
}
