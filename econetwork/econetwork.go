package econetwork

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	_ "github.com/mattn/go-sqlite3"
	"github.com/gorilla/websocket"
	"github.com/sony/sonyflake"
)

type Econetwork struct {
	Address string
	sessions map[string]User
	conn *websocket.Conn
	db *sql.DB
	sf *sonyflake.Sonyflake
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func New() (*Econetwork, error) {
	os.Create("econetwork.db")
	db, err := sql.Open("sqlite3", "econetwork.db")
	if err != nil {
		return nil, err
	}
	// make our tables
	db.Exec("CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY, username TEXT, password TEXT, node INTEGER);")
	var st sonyflake.Settings
	st.MachineID = 1

	return &Econetwork{
		Address: ":7768",
		sessions: map[string]User{},
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
		e.conn = conn

		go func() {
			for {
				// Read message from browser
				resp := ClientResponse{}
				err := conn.ReadJSON(&resp)
				if err != nil {
					if websocket.IsCloseError(err) || websocket.IsUnexpectedCloseError(err) {
						fmt.Println("client disconnected")
						return
					}
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

