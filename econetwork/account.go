package econetwork

import (
	"errors"

	"github.com/alexedwards/argon2id"
	"github.com/blockloop/scan"
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
	rows, _ := e.db.Query("SELECT * FROM users WHERE username = ?;", username)
	acc := Account{}
	scan.RowStrict(&acc, rows)

	if rows.Next() {
		return &acc, nil
	} else {
		return nil, ErrAccountNotExists
	}
}

func (e *Econetwork) register(r RegisterPayload) error {
	_, err := e.getAccount(r.Username)
	if err == nil {
		return ErrAccountExists // yes, we do check for ErrAccountNotExists and return not exists have a problem?
	}

	id, _ := e.sf.NextID()
	passwordHash, _ := argon2id.CreateHash(r.Password, argon2id.DefaultParams)
	
	_, err = e.db.Exec("INSERT INTO users (id, username, password, node) VALUES (?, ?, ?, ?);", id, r.Username, passwordHash, 0)
	return err
}
