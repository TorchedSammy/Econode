package econetwork

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/blockloop/scan"
)

var (
	ErrAccountNotExists = errors.New("account doesnt exist") // for when trying to login or get account
	ErrAccountExists = errors.New("account already exists") // trying to register an already existing username/account
	ErrMissingCredentials = errors.New("neither username or password were provided")
)

// A client's account
type Account struct {
	Username string `db:"username"`
	ID int `db:"id"`
	Node *Node // pointer since a person won't have a node immediately on register
	Op bool `db:"op"`
	Network *Econetwork
}

func (a *Account) GetNode() *Node {
	return a.Node
}

func (e *Econetwork) CreateNode(name string, owner *Account) {
	n := NewNode(name, owner)
	e.highestID++
	n.ID = e.highestID
	owner.Node = n
	go func(n *Node) {
		ticker := time.NewTicker(1 * time.Second)
		for {
			select {
			case <-ticker.C:
				n.Collect()
			}
		}
	}(n)
	e.nodes[n.ID] = n

	e.db.Exec("INSERT INTO nodes (id, name, owner, members, inventory, balance, gems, multi) VALUES (?, ?, ?, ?, ?, ?, ?);", n.ID, name, owner.ID, "", "", 2000.00, 0, 1.00)
}

func (e *Econetwork) getAccountByID(id int) (*Account, error) {
	rows, _ := e.db.Query("SELECT * FROM users WHERE id = ?;", id)
	acc := Account{}
	scan.RowStrict(&acc, rows)

	return e.getAccount(acc.Username)
}

func (e *Econetwork) getAccount(username string) (*Account, error) {
	if e.accountExists(username) {
		rows, _ := e.db.Query("SELECT * FROM users WHERE username = ?;", username)
		acc := Account{}
		scan.RowStrict(&acc, rows)

		nrow, _ := e.db.Query("SELECT * FROM nodes WHERE owner = ?;", acc.ID)
		node := Node{}
		scan.RowStrict(&node, nrow)
		if n, ok := e.nodes[node.ID]; !ok {
			acc.Node = &node
		} else {
			acc.Node = n
		}

		return &acc, nil
	} else {
		return nil, ErrAccountNotExists
	}
}

// Checks if an account exists in the database
func (e *Econetwork) accountExists(username string) bool {
    err := e.db.QueryRow("SELECT username FROM users WHERE username = ?", username).Scan(&username)
    if err != nil {
        if err != sql.ErrNoRows {
            fmt.Println("got another error in accountExists function", err)
            return true // TODO: ^ we should handle this properly
        }

        return false
    }

    return true
}

func (e *Econetwork) register(p AuthPayload) error {
	if p.Username == "" || p.Password == "" {
		return ErrMissingCredentials
	}

	_, err := e.getAccount(p.Username)
	if err == nil {
		return ErrAccountExists
	}

	id, _ := e.sf.NextID()
	passwordHash, _ := argon2id.CreateHash(p.Password, argon2id.DefaultParams)
	
	_, err = e.db.Exec("INSERT INTO users (id, username, password, node, op) VALUES (?, ?, ?, ?, ?);", id, p.Username, passwordHash, 0, 0)
	return err
}

func (e *Econetwork) login(p AuthPayload) (bool, error) {
	if p.Username == "" || p.Password == "" {
		return false, ErrMissingCredentials
	}

	if !e.accountExists(p.Username) {
		return false, ErrAccountNotExists
	}

	rows, _ := e.db.Query("SELECT password FROM users WHERE username = ?;", p.Username)
	var passwordHash string
	for rows.Next() {
		rows.Scan(&passwordHash)
	}

	match, _ := argon2id.ComparePasswordAndHash(p.Password, passwordHash)
	if match {
		return true, nil
	}

	return false, nil
}
