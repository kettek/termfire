package game

import (
	"os"

	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/messages"
	"github.com/rivo/tview"
)

type Login struct {
	MessageHandler
}

func (l *Login) Init(game Game) (tidy func()) {
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
			form := tview.NewForm().
				AddInputField("Server", targetServer, 0, nil, func(text string) {
					targetServer = text
				}).
				AddInputField("Account", account, 0, nil, func(text string) {
					account = text
				}).
				AddPasswordField("Password", password, 0, '*', func(text string) {
					password = text
				}).
				AddButton("Login", func() {
					l.TryLogin(game, account, password)
				})

			game.Pages().AddAndSwitchToPage("login", form, true)
			game.Redraw()
		})
	})

	if err := game.Connect(targetServer); err != nil {
		debug.Debug("Failed to connect to server")
		game.SetState(&Login{})
	}

	return func() {
		game.Pages().RemovePage("login")
	}
}

func (l *Login) TryLogin(game Game, account, password string) {
	game.SendMessage(&messages.MessageAccountLogin{Account: account, Password: password})
	l.Once(&messages.MessageAccountPlayers{}, &messages.MessageAccountLogin{}, func(msg messages.Message, failure *messages.MessageFailure) {
		if failure != nil {
			debug.Debug("Failed to login: ", failure.Reason)
			return
		}
		m := msg.(*messages.MessageAccountPlayers)
		// Just join with the first one for now! :) (Play handles actual message account play command due to how we only get a unique message if a failure happens).
		game.SetState(&Play{character: m.Characters[0].Name})
	})

}
