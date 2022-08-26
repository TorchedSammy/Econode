package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/websocket"
	"github.com/sony/sonyflake"
	"github.com/blockloop/scan"
)

const version = "Econode v0.1.0-beta"
// TODO: make better motd kek
const EconetworkMOTD = "Hello!!!"

type Econetwork struct {
	Address string
	sessions map[string]Client
	sessionAsName map[string]string
	routes map[string]Route
	highestID int // highest node id
	nodes map[int]*Node
	conn *websocket.Conn
	db *sql.DB
	sf *sonyflake.Sonyflake
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New() (*Econetwork, error) {
	if _, err := os.Stat("econetwork.db"); os.IsNotExist(err) {
		os.Create("econetwork.db")
	}

	db, err := sql.Open("sqlite3", "econetwork.db")
	if err != nil {
		return nil, err
	}
	// make our tables
	// sqlite doesnt have a boolean type, so `op` is an INTEGER which'll be either 0 (false) or 1 (true)
	db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT, node INTEGER, op INTEGER);")
	db.Exec("CREATE TABLE IF NOT EXISTS nodes (id INTEGER PRIMARY KEY, name TEXT, owner INTEGER, members TEXT, inventory TEXT, balance REAL, gems INTEGER, multi REAL);") /* REAL is a float */
	var st sonyflake.Settings
	st.MachineID = func() (uint16, error) {
		return 1, nil
	}

	e := &Econetwork{
		Address: ":7768",
		sessions: map[string]Client{},
		sessionAsName: map[string]string{},
		routes: map[string]Route{},
		highestID: 0,
		nodes: map[int]*Node{},
		conn: nil,
		db: db,
		sf: sonyflake.NewSonyflake(st),
	}

	e.setupAccountRoutes()
	e.setupNodeRoutes()
	e.setupMiscRoutes()

	return e, nil
}

func (e *Econetwork) Stop() {
	e.db.Close()
}

func (e *Econetwork) Dump() {
	fmt.Println("performing db dump")
	for _, n := range e.nodes {
		fmt.Println(n.Balance, n.CPS())
		// TODO: update items, they should be stored as a json string
		// key is item name, val is amount (?)
		_, err := e.db.Exec("UPDATE nodes SET balance = ? WHERE id = ?", n.Balance, n.ID)
		fmt.Println(err)
	}
}

func (e *Econetwork) Start() {
	rows, _ := e.db.Query("SELECT id FROM nodes")
	var nodeIDs []int
	scan.RowsStrict(&nodeIDs, rows)

	// start timer to write progress to db
	ticker := time.NewTicker(1 * time.Minute) // time is low to test
	go func() {
		for {
			select {
			case <-ticker.C:
				e.Dump()
			}
		}
	}()

	for _, nID := range nodeIDs {
		n := e.GetNode(nID) 
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
		e.highestID = n.ID
	}

	http.HandleFunc("/econetwork", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		c := Client{Conn: conn}
		c.Outgoing("welcome", WelcomePayload{
			MOTD: EconetworkMOTD,
		})

		go func() {
			for {
				resp := ClientResponse{}
				err := conn.ReadJSON(&resp)
				if err != nil {
					if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
						name := "client"
						if c.Account != nil {
							name = c.Account.Username
						}
						fmt.Printf("%s disconnected\n", name)
						if c.SessionID != "" {
							delete(e.sessions, c.SessionID)
							delete(e.sessionAsName, name)
						}
						return
					}
					fmt.Println(err)
				}

				// Print the message to the console
				fmt.Printf("%+v\n", resp)
				jsondata, _ := json.Marshal(resp.Data)
				fmt.Println(resp.Data)
				rt := e.getRoute(resp.Method)
				if rt == nil {
					return
					// TODO: send error of unknown method?
				}

				data, err := rt.DataTransformer(jsondata)
				if err != nil {
					c.SendMalformed(resp.Method)
					continue
				}
				
				rt.Execute(&c, data)
			}
		}()
	})

	http.ListenAndServe(e.Address, nil)
}

func (e *Econetwork) setupMiscRoutes() {
	e.addRoutes([]Route{
		createRoute("stats", "", false, nil, func(c *Client) {
			stats := StatsPayload{
				Nodes: len(e.nodes),
				NetworkVersion: version,
			}
			c.SendSuccess("stats", stats)
		}),
	})
}
// Generates a session id for a user
func SessionID() string {
	idRaw := make([]byte, 24)
	rand.Read(idRaw)

	return hex.EncodeToString(idRaw)
}
