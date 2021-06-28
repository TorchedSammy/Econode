package econetwork

import "github.com/gorilla/websocket"

type Client struct {
	Account *Account
	Conn *websocket.Conn
	SessionID string
}

func (c *Client) SendSuccess(method string, data interface{}) {
	c.Conn.WriteJSON(ServerResponse{
		Code: "success",
		Method: method,
		Data: &data,
	})
}

func (c *Client) SendError(method string, data interface{}) {
	c.Conn.WriteJSON(ServerResponse{
		Code: "error",
		Method: method,
		Data: &data,
	})
}

func (c *Client) SendMalformed(method string) {
	c.Conn.WriteJSON(ServerResponse{
		Code: "malformed",
		Method: method,
	})
}

