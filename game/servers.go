package game

import "github.com/kettek/termfire/messages"

type Servers struct {
	game Game
}

func (s *Servers) Init(game Game) {
	// TODO: Request Metaserver
}

func (s *Servers) OnMessage(msg messages.Message) {
	// ???
}
