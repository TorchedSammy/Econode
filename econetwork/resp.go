package econetwork

// The server response struct
// What we should send, as defined by the procotol
type ServerResponse struct {
	Code string `json:"code"`
	Method string `json:"method"`
	Data *interface{} `json:"content,omitempty"`
}

// Fields we should expect from the client
type ClientResponse struct {
	SessionID string `json:"sessionID"`
	Method string `json:"method"`
	Data *interface{} `json:"data,omitempty"`
}

