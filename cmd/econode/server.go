package main

import (
	"net/url"
	"strings"

	"github.com/TorchedSammy/Econode"
	"github.com/gorilla/websocket"
)

type network struct{
	ws *websocket.Conn
	session string // session ID
}

func connectToNetwork(address string) (*network, error) {
	if !strings.HasPrefix(address, "ws://") {
		address = "ws://" + address
	}

	u, err := url.Parse(address)
	if err != nil {
		return nil, err
	} else {
		u.Scheme = "ws"
		u.Path = "econetwork"

		parts := strings.Split(u.Host, ":")
		if len(parts) == 1 {
			u.Host = u.Host + ":7768" // default econode port
		}
	}

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return nil, err
	}

	return &network{
		ws: conn,
	}, nil
}

func (n *network) listenIncoming(cb func(econode.ServerResponse)) {
	go func() {
		sr := econode.ServerResponse{}
		err := n.ws.ReadJSON(&sr)
		if err != nil {
			panic(err)
		}
		
		cb(sr)
	}()
}

func (n *network) send(method string, data interface{}) {
	n.ws.WriteJSON(econode.ClientResponse{
		SessionID: n.session,
		Method: method,
		Data: &data,
	})
}
