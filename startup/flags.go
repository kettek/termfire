package startup

import "flag"

// Host is the host to connect to.
var Host string

// Account is the username to connect with.
var Account string

// Password is the password to connect with.
var Password string

// Character is the character to connect with.
var Character string

func init() {
	flag.StringVar(&Host, "host", "", "The host to connect to")
	flag.StringVar(&Account, "account", "", "The account to connect with")
	flag.StringVar(&Password, "password", "", "The password for the account")
	flag.StringVar(&Character, "character", "", "The character to connect with")

	flag.Parse()
}
