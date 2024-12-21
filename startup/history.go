package startup

type hostHistory struct {
	Account             string
	Password            string
	RememberCredentials bool
	Character           string
}

type history struct {
	LastHost string
	Hosts    []hostHistory
}

// History contains stored information about servers, connections, etc.
var History history
