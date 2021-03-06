package econetwork

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/websocket"
	"github.com/sony/sonyflake"
)

const version = "Econode v0.1.0-beta"

type Econetwork struct {
	Address string
	sessions map[string]Client
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
	db.Exec("CREATE TABLE IF NOT EXISTS nodes (id INTEGER PRIMARY KEY, name TEXT, owner INTEGER, members TEXT, inventory TEXT, balance INTEGER, multi REAL);" /* REAL is a float */)
	var st sonyflake.Settings
	st.MachineID = func() (uint16, error) {
		return 1, nil
	}

	return &Econetwork{
		Address: ":7768",
		sessions: map[string]Client{},
		conn: nil,
		db: db,
		sf: sonyflake.NewSonyflake(st),
	}, nil
}

func (e *Econetwork) Stop() {
	e.db.Close()
}

func (e *Econetwork) Start() {
	http.HandleFunc("/econetwork", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		c := Client{Conn: conn}

		go func() {
			for {
				// Read message from browser
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
						}
						return
					}
					fmt.Println(err)
				}

				// Print the message to the console
				fmt.Printf("%+v\n", resp)
				jsondata, _ := json.Marshal(resp.Data)
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

						c.SessionID = sessionid
						c.SendSuccess("login", sessionid)
					} else {
						fmt.Println("Error in login method occurred\n", err)
						c.SendError("login", nil)
					}
				case "stats":
					stats := StatsPayload{
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
				case "getEconode":
					c, ok := e.sessions[resp.SessionID]
					if !ok {
						c.SendError("newEconode", "session not found")
						continue
					}
					if c.Account == nil {
						c.SendError("newEconode", "not authenticated")
						continue
					}
					if c.Account.Node == nil {
						c.SendError("newEconode", "user not in econode")
						continue
					}

					node := c.Account.GetNode()
					c.Conn.WriteJSON(EconodeInfoPayload{
						Name: node.Name,
						Owner: c.Account.ID,
						Balance: node.Balance,
					})
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
