package game

import (
	"os"

	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/messages"
)

type Login struct {
	MessageHandler
	game Game
}

func (l *Login) Init(game Game) {
	debug.Debug("Init")
	l.game = game

	targetServer := os.Args[1]

	account := os.Args[2]
	password := os.Args[3]

	l.Once(&messages.MessageVersion{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageVersion)
		if m.SVVersion != "1030" {
			debug.Debug("Server version mismatch", m.SVVersion)
			game.SetState(&Servers{})
			return
		}

		game.SendMessage(&messages.MessageVersion{CLVersion: "1030", SVName: "termfire"})
		game.SendMessage(&messages.MessageSetup{FaceCache: true, LoginMethod: "2", ExtendedStats: true})

		l.Once(&messages.MessageSetup{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
			game.SendMessage(&messages.MessageAccountLogin{Account: account, Password: password})
			l.Once(&messages.MessageAccountPlayers{}, &messages.MessageAccountLogin{}, func(msg messages.Message, failure *messages.MessageFailure) {
				if failure != nil {
					debug.Debug("Failed to login: ", failure.Reason)
					game.SetState(&Servers{})
					return
				}
				m := msg.(*messages.MessageAccountPlayers)
				// Just join with the first one for now! :) (Play handles actual message account play command due to how we only get a unique message if a failure happens).
				game.SetState(&Play{character: m.Characters[0].Name})
			})
		})
	})

	if err := game.Connect(targetServer); err != nil {
		debug.Debug("Failed to connect to server")
		game.SetState(&Servers{})
		return
	}

}
