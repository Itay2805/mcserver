package item

import "fmt"

type Item struct {
	ID			int
	Name		string
	StackSize	int
}

func (item *Item) String() string {
	return fmt.Sprintf("Item{ Name: \"%s\" }", item.Name)
}

func GetById(id int) *Item {
	return items[id]
}
