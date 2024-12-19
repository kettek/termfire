package game

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/messages"
	"github.com/kettek/termfire/startup"
	"github.com/rivo/tview"
)

type Servers struct {
	game             Game
	Metaservers      []string
	serverEntries    messages.ServerEntries
	serversContainer *tview.List
	hostInput        *tview.InputField
}

func (s *Servers) Init(game Game) (tidy func()) {
	s.game = game

	s.Metaservers = []string{
		"http://crossfire.real-time.com/metaserver2/meta_client.php",
		"http://metaserver.eu.cross-fire.org/meta_client.php",
		"http://metaserver.us.cross-fire.org/meta_client.php",
	}

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

	container.AddItem(box, 0, 1, true)
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
		if startup.Host != "" {
			game.SetState(&Login{
				TargetServer: startup.Host,
			})
		} else {
			game.Pages().AddAndSwitchToPage("servers", container, true)
			game.Redraw()
		}
	}()

	s.Refresh()

	s.hostInput.SetText(startup.History.LastHost)

	for i, server := range s.serverEntries {
		if server.Hostname == startup.History.LastHost {
			s.serversContainer.SetCurrentItem(i)
			break
		}
	}

	// TODO: Request Metaserver
	return func() {
		game.Pages().RemovePage("servers")
	}
}

func (s *Servers) OnMessage(msg messages.Message) {
	// ???
}

func (s *Servers) Refresh() {
	s.serverEntries = messages.ServerEntries{}
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

func (s *Servers) RequestServerList(metaserver string) (messages.ServerEntries, error) {
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

	serverEntries := messages.ServerEntries{}

	err = serverEntries.UnmarshalBinary(body)
	if err != nil {
		return nil, err
	}

	return serverEntries, nil
}
