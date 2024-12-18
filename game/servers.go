package game

import (
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/messages"
	"github.com/rivo/tview"
)

type Servers struct {
	game             Game
	Metaservers      []string
	serverEntries    serverEntries
	serversContainer *tview.List
	hostInput        *tview.InputField
}

type serverEntry struct {
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

type serverEntries []serverEntry

func (s *serverEntries) UnmarshalBinary(data []byte) error {
	lines := strings.Split(string(data), "\n")
	var entry serverEntry
	for _, line := range lines {
		if line == "START_SERVER_DATA" {
			entry = serverEntry{}
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

func (s *Servers) Init(game Game) (tidy func()) {
	s.game = game

	// Setup some UI.
	container := tview.NewFlex().SetDirection(tview.FlexRow)

	// Servers list
	box := tview.NewFlex()
	box.SetDirection(tview.FlexRow)
	box.SetBorder(true).SetTitle("Servers")

	s.serversContainer = tview.NewList()
	s.serversContainer.ShowSecondaryText(true)
	s.serversContainer.SetSelectedFunc(func(i int, _ string, _ string, _ rune) {
		server := s.serverEntries[i]
		s.hostInput.SetText(server.Hostname + ":" + strconv.Itoa(server.Port))
	})
	box.AddItem(s.serversContainer, 0, 1, true)

	container.AddItem(box, 0, 1, false)
	// Manual Join
	box2 := tview.NewFlex()
	box2.SetDirection(tview.FlexRow)

	form := tview.NewForm()
	form.SetHorizontal(true)

	s.hostInput = tview.NewInputField()
	s.hostInput.SetLabel("Host")
	s.hostInput.SetFieldWidth(40)
	form.AddFormItem(s.hostInput)
	form.AddButton("Join", func() {
		game.SetState(&Login{
			TargetServer: s.hostInput.GetText(),
		})
	})
	form.AddButton("Refresh", func() {
		s.Refresh()
	})

	box2.AddItem(form, 0, 1, false)

	container.AddItem(box2, 3, 1, false)

	go func() {
		game.Pages().AddAndSwitchToPage("servers", container, true)
		game.Redraw()
	}()

	s.Refresh()

	// TODO: Request Metaserver
	return func() {
		game.Pages().RemovePage("servers")
	}
}

func (s *Servers) OnMessage(msg messages.Message) {
	// ???
}

func (s *Servers) Refresh() {
	s.serverEntries = serverEntries{}
	for _, metaserver := range s.Metaservers {
		entries, err := s.RequestServerList(metaserver)
		if err != nil {
			debug.Debug("Failed to get server list from metaserver: ", err)
			continue
		}
		for _, e := range entries {
			found := false
			for _, server := range s.serverEntries {
				if server.Hostname == e.Hostname && server.Port == e.Port {
					found = true
					break
				}
			}
			if !found {
				s.serverEntries = append(s.serverEntries, e)
			}
		}
	}

	s.serversContainer.Clear()
	for _, server := range s.serverEntries {
		primary := server.Hostname + " " + server.TextComment
		secondary := "  Players: " + strconv.Itoa(server.NumPlayers) + " Version: " + server.Version
		s.serversContainer.AddItem(primary, secondary, 0, nil)
	}
}

func (s *Servers) RequestServerList(metaserver string) (serverEntries, error) {
	resp, err := http.Get(metaserver)
	http.DefaultClient.Timeout = 5 * time.Second
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	serverEntries := serverEntries{}

	err = serverEntries.UnmarshalBinary(body)
	if err != nil {
		return nil, err
	}

	return serverEntries, nil
}
