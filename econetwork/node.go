package econetwork

import "github.com/blockloop/scan"

// Someone's econode
// The idea is that we can have other people growing a single node together
type Node struct {
	ID int `db:"id"`
	Name string `db:"name"`
	Balance float64 `db:"balance"`
	OwnerID int `db:"owner"`
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

func (e *Econetwork) GetNode(id int) *Node {
	nrow, _ := e.db.Query("SELECT * FROM nodes WHERE id = ?;", id)
	node := Node{}
	scan.RowStrict(&node, nrow)
	node.Owner, _ = e.getAccountByID(node.OwnerID)

	return &node
}

func (n *Node) Buy(purchase Item) {
	item, ok := n.Inventory[purchase.Name]
	if !ok {
		item = &purchase
	}
	item.Count++
	n.Inventory[purchase.Name] = item
}

func (n *Node) CPS() float64 {
	var coins float64

	for _, itm := range n.Inventory {
		icps := itm.CPS * itm.Count // item base cps * amount in node
		coins += icps
	}

	return coins
}

func (n *Node) Collect() {
	n.Balance += n.CPS()
}

