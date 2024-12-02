package main

import (
	"errors"
	"net"
	"os"
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
			var buf [99999]byte
			if n, err := g.conn.Read(buf[:]); err != nil {
				debug.Debug("Failed to read from server: ", err)
			} else {
				var offset int
				for offset < n {
					length := int(buf[offset])<<8 | int(buf[offset+1])
					offset += 2
					msg := buf[offset : offset+length]
					message, err := messages.UnmarshalMessage(msg)
					if err != nil {
						debug.Debug("Failed to unmarshal message: ", err)
						return
					}
					g.state.OnMessage(message)
					offset += length
				}
			}
		}
	}()
	return nil
}

func (g *Game) Log(msg string) {
	// TODO: Add to some sort of high-level panel.
}

func (g *Game) SetState(state game.State) {
	g.state = state
	state.Init(g)
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
	if len(os.Args) < 4 {
		println("Usage: termfire <server> <account> <password>")
		os.Exit(1)
	}

	if err := debug.Start(); err != nil {
		panic(err)
	}

	var g Game

	box := tview.NewBox().SetBorder(true).SetTitle("Termfire")

	app := tview.NewApplication()
	app.SetRoot(box, true)

	app.SetAfterDrawFunc(func(s tcell.Screen) {
		g.SetState(&game.Login{})
		app.SetAfterDrawFunc(nil)
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
