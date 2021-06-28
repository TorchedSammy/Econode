package econetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/websocket"
)

type Econetwork struct {
	Address string
	sessions map[string]User
	conn *websocket.Conn
	db *sql.DB
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New() (*Econetwork, error) {
	db, err := sql.Open("sqlite3", "./econetwork.db")
	if err != nil {
		return nil, err
	}
	db.Query("CREATE TABLE users (id INTEGER PRIMARY KEY, username TEXT, password TEXT, node INTEGER);")
	db.Query("INSERT INTO users (id, username, password, node) VALUES (0, \"admin\", \"a\", 0)")

	return &Econetwork{
		Address: ":7768",
		sessions: map[string]User{},
		conn: nil,
		db: db,
	}, nil
}

func (e *Econetwork) Stop() {
	e.db.Close()
}

func (e *Econetwork) Start() {
	http.HandleFunc("/econetwork", func(w http.ResponseWriter, r *http.Request) {
		conn, _ := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity
		e.conn = conn

		go func() {
			for {
				// Read message from browser
				resp := ClientResponse{}
				err := conn.ReadJSON(&resp)
				if err != nil {
					fmt.Println(err)
				}

				// Print the message to the console
				fmt.Printf("%+v\n", resp)
				jsondata, _ := json.Marshal(resp.Data)
				switch resp.Method {
				case "register":
					registerInfo := RegisterPayload{}
					if err := json.Unmarshal(jsondata, &registerInfo); err != nil {
						e.sendMalformed("register")
						continue
					}
					e.register(registerInfo)
				}
			}
		}()
	})

	http.ListenAndServe(e.Address, nil)
}

