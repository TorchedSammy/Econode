package econetwork

import (
	"fmt"

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
	_, err := e.db.Query("INSERT INTO users (id, username, password, node) VALUES (?, ?, ?, ?)", id, u.Username, passwordHash, 0)
	if err != nil {
		fmt.Println(err)
	}
}
