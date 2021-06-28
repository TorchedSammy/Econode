package econetwork

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/blockloop/scan"
)

var (
	ErrAccountNotExists = errors.New("account doesnt exist") // for when trying to login or get account
	ErrAccountExists = errors.New("account already exists") // trying to register an already existing username/account
)

// A client's account
type Account struct {
	Username string `db:"username"`
	Node *Node // pointer since a person won't have a node immediately on register
	Op bool `db:"op"`
}

// Someone's econode
// The idea is that we can have other people growing a single node together
type Node struct {
	name string
	balance int
}

func (e *Econetwork) getAccount(username string) (*Account, error) {
	if e.accountExists(username) {
		rows, _ := e.db.Query("SELECT * FROM users WHERE username = ?;", username)
		acc := Account{}
		scan.RowStrict(&acc, rows)

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
            fmt.Println("got another error in checkAccount function", err)
            return true // TODO: ^ we should handle this properly
        }

        return false
    }

    return true
}

func (e *Econetwork) register(r RegisterPayload) error {
	_, err := e.getAccount(r.Username)
	if err != nil {
		return ErrAccountExists // yes, we do check for ErrAccountNotExists and return not exists have a problem?
	}

	id, _ := e.sf.NextID()
	passwordHash, _ := argon2id.CreateHash(r.Password, argon2id.DefaultParams)
	
	_, err = e.db.Exec("INSERT INTO users (id, username, password, node, op) VALUES (?, ?, ?, ?, ?);", id, r.Username, passwordHash, 0, 0)
	return err
}
