package main

import (
	"database/sql"
//	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/TorchedSammy/Econode/econetwork"
)

func main() {
	network := econetwork.New()
	network.Run()
}
