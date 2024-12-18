package main

import (
	"errors"
	"net"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/game"
	"github.com/kettek/termfire/messages"
	"github.com/rivo/tview"
)

type Game struct {
	conn  net.Conn
	state game.State
	app   *tview.Application
	tidy  func()
	pages *tview.Pages
}

func (g *Game) App() *tview.Application {
	return g.app
}

func (g *Game) Redraw() {
	// What the hell is this... Why can't g.app propagate children updates/changes...
	go func() {
		g.app.QueueUpdateDraw(func() {})
	}()
}

func (g *Game) Pages() *tview.Pages {
	return g.pages
}

func (g *Game) Connect(targetServer string) error {
	if !strings.Contains(targetServer, ":") {
		targetServer += ":13327"
	}

	conn, err := net.Dial("tcp", targetServer)
	if err != nil {
		g.Log("Failed to connect to server")
		return err
	} else {
		g.Log("Connected to server")
	}
	g.conn = conn

	go func() {
		for {
			var buf [32767]byte
			if n, err := g.conn.Read(buf[:]); err != nil {
				debug.Debug("Failed to read from server: ", err)
				// FIXME: Various game states need to know the connection exploded... add some sort on OnDisconnected-type func to each state.
				return
			} else {
				var offset int
				for offset < n {
					length := int(buf[offset])<<8 | int(buf[offset+1])

					// TODO: Turn into a loop?
					if offset+length > n {
						// FIXME: Something is really wrong here -- we get message too long too often... are we failing to read some messages somehow...?
						debug.Debug("message too long: ", offset, length, n, string(buf[offset:n]))

						if n2, err := g.conn.Read(buf[n:]); err != nil {
							debug.Debug("Failed to read from server: ", err)
							return
						} else {
							debug.Debug("Read more data: ", n2)
							n += n2
						}
						continue
					}

					offset += 2
					msg := buf[offset : offset+length]
					message, err := messages.UnmarshalMessage(msg)
					if err != nil {
						debug.Debug("Failed to unmarshal message: ", err)
						return
					}
					offset += length
					g.app.QueueUpdate(func() {
						g.state.OnMessage(message)
					})
				}
			}
		}
	}()
	return nil
}

func (g *Game) Disconnect() {
	if g.conn != nil {
		g.conn.Close()
		g.conn = nil
	}
}

func (g *Game) Log(msg string) {
	// TODO: Add to some sort of high-level panel.
}

func (g *Game) SetState(state game.State) {
	if g.tidy != nil {
		g.tidy()
	}
	g.state = state
	g.tidy = state.Init(g)
}

func (g *Game) SendMessage(msg messages.Message) error {
	bytes := msg.Bytes()
	if len(bytes) > 0 {
		// TODO
		//g.logPanel.Add("C->S" + string(bytes))
		g.conn.Write([]byte{byte(len(bytes) >> 8), byte(len(bytes))})
		g.conn.Write(bytes)
		return nil
	}
	return errors.New("empty message")
}

func bytesToStringAndHex(b []byte) string {
	result := ""
	for _, c := range b {
		result += string(c) + "  "
	}
	result += "\n"
	for _, c := range b {
		result += strconv.FormatInt(int64(c), 16) + " "
	}
	return result
}

func main() {
	if err := debug.Start(); err != nil {
		panic(err)
	}

	var g Game

	g.pages = tview.NewPages()

	g.app = tview.NewApplication()
	g.app.EnableMouse(true)
	g.app.SetRoot(g.pages, true)
	g.app.SetFocus(g.pages)

	g.app.SetAfterDrawFunc(func(s tcell.Screen) {
		g.app.SetAfterDrawFunc(nil)
		g.SetState(&game.Servers{})
	})

	if err := g.app.Run(); err != nil {
		panic(err)
	}
}
