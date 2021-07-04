package econetwork

// Here we define individual payloads/`data` types for different methods

// For register and login methods
type AuthPayload struct {
	Password string
	Username string
}

type StatsPayload struct {
	Nodes int
	Accounts int
	NetworkVersion string
}

type EconodeNewPayload struct {
	Name string
}

type EconodeInfoPayload struct {
	Name string
	Owner string
	Balance int
}

