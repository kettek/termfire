package game

import (
	"github.com/kettek/termfire/messages"
	"github.com/rivo/tview"
)

type State interface {
	Init(Game) (tidy func())
	OnMessage(messages.Message)
}

type Game interface {
	App() *tview.Application
	Redraw()
	Connect(addr string) error
	Disconnect()
	SetState(state State)
	SendMessage(msg messages.Message) error
	Log(string)
	Pages() *tview.Pages
}
