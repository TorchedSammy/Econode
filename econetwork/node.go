package econetwork

// Someone's econode
// The idea is that we can have other people growing a single node together
type Node struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Balance int `db:"balance"`
	Owner *Account
	Members []int
	Inventory map[string]*Item
	Multi float64 `db:"multi"`
}

func NewNode(name string, ownerAcc *Account) *Node {
	return &Node{
		Name: name,
		Owner: ownerAcc,
		Multi: 1.00,
	}
}

func (n *Node) Buy(purchase Item) {
	item, ok := n.Inventory[purchase.Name]
	if !ok {
		item = &purchase
	}
	item.Count++
	n.Inventory[purchase.Name] = item
}

