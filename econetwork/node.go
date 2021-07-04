package econetwork

// Someone's econode
// The idea is that we can have other people growing a single node together
type Node struct {
	ID int
	Name string
	Balance int
	Owner int
	Members []int
	Inventory map[string]*Item
	Multi float64
}

func (n *Node) Buy(purchase Item) {
	item, ok := n.Inventory[purchase.Name]
	if !ok {
		item = &purchase
	}
	item.Count++
	n.Inventory[purchase.Name] = item
}
