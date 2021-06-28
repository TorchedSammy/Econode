package econetwork

import (
	"fmt"
	"errors"

	"github.com/alexedwards/argon2id"
)

var (
	ErrAccountNotExists = errors.New("account doesnt exist") // for when trying to login or get account
	ErrAccountExists = errors.New("account already exists") // trying to register an already existing username/account
)

// A client's account
type Account struct {
	Username string
	Node *Node // pointer since a person won't have a node immediately on register
	Op bool
}

// Someone's econode
// The idea is that we can have other people growing a single node together
type Node struct {
	name string
	balance int
}

func (e *Econetwork) getAccount(username string) (*Account, error) {
	rows, err := e.db.Query("SELECT * FROM users WHERE username = ?;", username)
	return nil, err
}

func (e *Econetwork) register(r RegisterPayload) error {
	_, err := e.getAccount(r.Username)
	fmt.Println(err)
	id, _ := e.sf.NextID() // TODO: make a generated snowflake
	passwordHash, _ := argon2id.CreateHash(r.Password, argon2id.DefaultParams)
	
	_, err = e.db.Exec("INSERT INTO users (id, username, password, node) VALUES (?, ?, ?, ?);", id, r.Username, passwordHash, 0)
	return err
}
