package econetwork

type ResponseContent struct {
	Data string `json:"data"`
	Files []string `json:"files,omitempty"`
}

// The server response struct
// What we should send, as defined by the procotol
type ServerResponse struct {
	Result string `json:"result"`
	Method string `json:"method"`
	Data *interface{} `json:"content,omitempty"`
}

// Fields we should expect from the client
type ClientResponse struct {
	SessionID string `json:"sessionID"`
	Method string `json:"method"`
	Data *interface{} `json:"data,omitempty"`
}

func (e *Econetwork) sendMalformed(method string) {
	// TODO
	// e.conn.WriteJSON()
}
