package econetwork

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Econetwork struct {
	address string
	port int
	sessions map[string]User
	conn *websocket.Conn
	db sql.DB
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (e *Econetwork) Run() {
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
				switch resp.Method {
				case "register":
					registerInfo := RegisterPayload{}
					if err := json.Unmarshal(resp.Data, &registerInfo); err != nil {
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

