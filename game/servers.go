package game

import "github.com/kettek/termfire/messages"

type Servers struct {
	game Game
}

func (s *Servers) Init(game Game) (tidy func()) {
	// TODO: Request Metaserver
	return nil
}

func (s *Servers) OnMessage(msg messages.Message) {
	// ???
}
