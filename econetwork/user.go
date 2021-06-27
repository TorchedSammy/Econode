package econetwork

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
	"github.com/alexedwards/argon2id"
)

// A single user/player
type User struct {
	username string
	node *Node // pointer since a person won't have a node immediately on register
	op bool
}

// Someone's econode
// The idea is that we can have other people growing a single node together
type Node struct {
	name string
	balance int
}

func (e *Econetwork) register(u RegisterPayload) {
	id := 0 // TODO: make a generated snowflake
	passwordHash, _ := argon2id.CreateHash(u.Password, argon2id.DefaultParams)
	e.db.Query("INSERT INTO users (id, username, password, node) VALUES (?, ?, ?, ?)", id, u.Username, hashedPass, nil)
}
