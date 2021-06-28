package econetwork

import (
	"github.com/alexedwards/argon2id"
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

func (e *Econetwork) register(r RegisterPayload) error {
	id, _ := e.sf.NextID() // TODO: make a generated snowflake
	passwordHash, _ := argon2id.CreateHash(r.Password, argon2id.DefaultParams)
	
	_, err := e.db.Exec("INSERT INTO users (id, username, password, node) VALUES (?, ?, ?, ?);", id, r.Username, passwordHash, 0)
	return err
}
