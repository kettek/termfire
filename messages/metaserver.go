package messages

import (
	"strconv"
	"strings"
)

// ServerEntries is a slice of server entries, wow.
type ServerEntries []ServerEntry

// UnmarshalBinary unmarshals the binary data into a slice of server entries.
func (s *ServerEntries) UnmarshalBinary(data []byte) error {
	lines := strings.Split(string(data), "\n")
	var entry ServerEntry
	for _, line := range lines {
		if line == "START_SERVER_DATA" {
			entry = ServerEntry{}
		} else if line == "END_SERVER_DATA" {
			*s = append(*s, entry)
		} else {
			kv := strings.SplitN(line, "=", 2)
			if len(kv) != 2 {
				continue
			}
			key := kv[0]
			value := kv[1]
			switch key {
			case "hostname":
				entry.Hostname = value
			case "port":
				entry.Port, _ = strconv.Atoi(value)
			case "html_comment":
				entry.HTMLComment = value
			case "text_comment":
				entry.TextComment = value
			case "archbase":
				entry.ArchBase = value
			case "mapbase":
				entry.MapBase = value
			case "codebase":
				entry.CodeBase = value
			case "num_players":
				entry.NumPlayers, _ = strconv.Atoi(value)
			case "in_bytes":
				entry.InBytes, _ = strconv.Atoi(value)
			case "out_bytes":
				entry.OutBytes, _ = strconv.Atoi(value)
			case "uptime":
				entry.Uptime, _ = strconv.Atoi(value)
			case "version":
				entry.Version = value
			case "sc_version":
				entry.SCVersion, _ = strconv.Atoi(value)
			case "cs_version":
				entry.CSVersion, _ = strconv.Atoi(value)
			case "last_update":
				entry.LastUpdate, _ = strconv.Atoi(value)
			}
		}
	}
	return nil
}

// ServerEntry is the entry for a given server.
type ServerEntry struct {
	Hostname    string
	Port        int
	HTMLComment string
	TextComment string
	ArchBase    string
	MapBase     string
	CodeBase    string
	NumPlayers  int
	InBytes     int
	OutBytes    int
	Uptime      int
	Version     string
	SCVersion   int
	CSVersion   int
	LastUpdate  int
}
