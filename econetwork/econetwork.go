package econetwork

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

	return &Econetwork{
		Address: ":7768",
		sessions: map[string]Client{},
		sessionAsName: map[string]string{},
		highestID: 0,
		nodes: map[int]*Node{},
		conn: nil,
		db: db,
		sf: sonyflake.NewSonyflake(st),
	}, nil
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
				switch resp.Method {
				case "register":
					if c.Account != nil {
						c.SendError("register", "client already authorized")
						continue
					}

					registerInfo := AuthPayload{}
					if err := json.Unmarshal(jsondata, &registerInfo); err != nil {
						c.SendMalformed("register")
						continue
					}
					err = e.register(registerInfo)
					if err == nil {
						c.Account = &Account{Username: registerInfo.Username}
						sessionid := SessionID()
						e.sessions[sessionid] = c
						e.sessionAsName[c.Account.Username] = sessionid
						c.SessionID = sessionid
						c.SendSuccess("register", sessionid)
					} else {
						fmt.Println("Error in register method occurred\n", err)
						c.SendError("register", nil)
					}
				case "login":
					if c.Account != nil {
						c.SendError("login", "client already authorized")
						continue
					}

					loginInfo := AuthPayload{}
					if err := json.Unmarshal(jsondata, &loginInfo); err != nil {
						c.SendMalformed("login")
						continue
					}
					passMatch, err := e.login(loginInfo)
					if err == nil {
						if !passMatch {
							c.SendError("login", "password is incorrect")
							continue
						}
						loginAcc, _ := e.getAccount(loginInfo.Username)
						sessionid := SessionID()

						c.Account = loginAcc
						e.sessions[sessionid] = c
						e.sessionAsName[c.Account.Username] = sessionid

						c.SessionID = sessionid
						c.SendSuccess("login", sessionid)
					} else {
						fmt.Println("Error in login method occurred\n", err)
						c.SendError("login", nil)
					}
				case "stats":
					stats := StatsPayload{
						Nodes: len(e.nodes),
						NetworkVersion: version,
					}
					c.SendSuccess("stats", stats)
				case "newEconode":
					c, ok := e.sessions[resp.SessionID]
					if !ok {
						c.SendError("newEconode", "session not found")
						continue
					}
					if c.Account == nil {
						c.SendError("newEconode", "not authenticated")
						continue
					}
					if c.Account.Node != nil {
						c.SendError("login", "user already in a node")
						continue
					}

					econodeInfo := EconodeNewPayload{}
					if err := json.Unmarshal(jsondata, &econodeInfo); err != nil {
						c.SendMalformed("newEconode")
						continue
					}
					e.CreateNode(econodeInfo.Name, c.Account)
					c.SendSuccess("newEconode", "node created")
				case "fetchNode":
					c, ok := e.sessions[resp.SessionID]
					if !ok {
						c.SendError("fetchNode", "session not found")
						continue
					}
					if c.Account == nil {
						c.SendError("fetchNode", "not authenticated")
						continue
					}
					// check if data field is nothing (if it is, get session's node)
					node := &Node{}
					if resp.Data != nil {
						if c.Account.Node == nil {
							c.SendError("fetchNode", "user not in econode")
							continue
						}

						node = c.Account.GetNode()
					} else {
						// get from data
						var nodeName string
						if err := json.Unmarshal(jsondata, &nodeName); err != nil {
							c.SendMalformed("fetchNode")
							continue
						}

						e.GetNodeByName(nodeName)
					}
					c.SendSuccess("fetchNode", EconodeInfoPayload{
						Name: node.Name,
						Owner: node.OwnerID,
						Balance: node.Balance,
					})
				case "buyItem":
					c, ok := e.sessions[resp.SessionID]
					if !ok {
						c.SendError("buyItem", "session not found")
						continue
					}
					if c.Account == nil {
						c.SendError("buyItem", "not authenticated")
						continue
					}
					if c.Account.Node == nil {
						c.SendError("buyItem", "user not in econode")
						continue
					}

					buyInfo := ItemPurchasePayload{}
					if err := json.Unmarshal(jsondata, &buyInfo); err != nil {
						c.SendMalformed("buyItem")
						continue
					}
					fmt.Printf("%#v\n", buyInfo)
					if itemMap[buyInfo.ItemName] == nil {
						c.SendFail("buyItem", "unknown item")
						continue
					}
					err := c.Account.Node.Buy(*itemMap[buyInfo.ItemName], buyInfo.Amount)
					if err != nil {
						c.SendError("buyItem", err)
						continue
					}
					c.SendSuccess("buyItem", nil)
				case "pm":
					c, ok := e.sessions[resp.SessionID]
					if !ok {
						c.SendError("pm", "session not found")
						continue
					}
					if c.Account == nil {
						c.SendError("pm", "not authenticated")
						continue
					}
					if c.Account.Node == nil {
						c.SendError("pm", "user not in econode")
						continue
					}

					content := UserMessagePayload{}
					if err := json.Unmarshal(jsondata, &content); err != nil {
						c.SendMalformed("pm")
						continue
					}

					personSession := e.sessionAsName[content.User]
					if personSession == "" {
						c.SendFail("pm", "user not online")
						continue
					}
					personClient := e.sessions[personSession]
					msg := UserMessagePayload{
						User: c.Account.Username,
						Message: content.Message,
					}
					personClient.Outgoing("pm", msg)
				}
			}
		}()
	})

	http.ListenAndServe(e.Address, nil)
}

// Generates a session id for a user
func SessionID() string {
	idRaw := make([]byte, 24)
	rand.Read(idRaw)

	return hex.EncodeToString(idRaw)
}
