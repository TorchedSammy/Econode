package main

// Here we define individual payloads/`data` types for different methods

// For register and login methods
type AuthPayload struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type UserMessagePayload struct {
	User string `json:"username"`
	Message string `json:"message"`
}

type StatsPayload struct {
	Nodes int `json:"nodes"`
	Accounts int `json:"accounts"`
	NetworkVersion string `json:"networkVersion"`
}

type EconodeInfoPayload struct {
	Name string `json:"name"`
	Owner int `json:"ownerID"`
	Balance float64 `json:"balance"`
}

type AccountInfoPayload struct {
	Username string `json:"username"`
	ID int `json:"id"`
	Node int `json:"nodeID"`
	Op bool `json:"op"`
}

type ItemPurchasePayload struct {
	ItemName string `json:"itemName"`
	Amount float64 `json:"amount"`
}

type WelcomePayload struct {
	MOTD string `json:"motd"`
}

