package main

import (
	"database/sql"
//	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/TorchedSammy/Econode/econetwork"
)

func main() {
	db, _ := sql.Open("sqlite3", "../../econetwork.db")
	defer db.Close()

	network := econetwork.Econetwork{
		address: "0.0.0.0",
		port: "7768",
		sessions: map[string]econetwork.User{},
		conn: nil,
		db: db,
	}
	network.Run()
}
