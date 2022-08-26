package main

import (
	"errors"
	"math"
	"strconv"
	"strings"

	"github.com/blockloop/scan"
)

var (
	ErrNotEnoughMoney = errors.New("not enough money to do this")
)

// Someone's econode
// The idea is that we can have other people growing a single node together
type Node struct {
	ID int `db:"id"`
	Name string `db:"name"`
	OwnerID int `db:"owner"`
	Owner *Account
	membersRaw string `db:"members"`
	Members []int
	invRaw string `db:"inventory"`
	Inventory map[string]Item
	Balance float64 `db:"balance"`
	Gems int `db:"gems"`
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
	if len(node.membersRaw) != 0 { // yes this is kinda stupid
		for _, mIDstr := range strings.Split(node.membersRaw, ",") {
			mID, _ := strconv.Atoi(mIDstr)
			node.Members = append(node.Members, mID)
		}
	}
	node.Inventory = make(map[string]Item)

	return &node
}

func (e *Econetwork) GetNodeByName(name string) *Node {
	row, _ := e.db.Query("SELECT id FROM nodes WHERE name = ?;", name)
	var nodeID int
	scan.RowsStrict(&nodeID, row)

	return e.GetNode(nodeID)
}

func (e *Econetwork) setupNodeRoutes() {
	e.addRoutes([]Route{
		createRoute("fetchNode", "", true, nil, func(c *Client) {
			node := c.Account.GetNode()
			c.SendSuccess("fetchNode", EconodeInfoPayload{
				Name: node.Name,
				Owner: node.OwnerID,
				Balance: node.Balance,
			})
		}),
	})
}

func (n *Node) Buy(purchase Item, amount float64) error {
	item, ok := n.Inventory[purchase.Name]
	if !ok {
		item = purchase
	}
	item.Count += amount
	price := (item.Price * (math.Pow(1.2, item.Count) - 1)) / 0.2
	if price > n.Balance {
		return ErrNotEnoughMoney
	}

	n.Inventory[purchase.Name] = item
	n.Balance -= price
	return nil
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
