package game

import (
	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/game/play"
	"github.com/kettek/termfire/messages"
	"github.com/kettek/termfire/startup"
	"github.com/rivo/tview"
)

type Login struct {
	MessageHandler
	TargetServer string
}

func (l *Login) Init(game Game) (tidy func()) {
	targetServer := l.TargetServer

	var account, password string
	if startup.Account != "" {
		account = startup.Account
	}
	if startup.Password != "" {
		password = startup.Password
	}

	// Clear out our face cache.
	play.GlobalObjectMapper.Reset()

	l.Once(&messages.MessageVersion{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageVersion)
		if m.SVVersion != "1030" {
			debug.Debug("Server version mismatch", m.SVVersion)
			game.SetState(&Servers{})
			return
		}

		game.SendMessage(&messages.MessageVersion{CLVersion: "1030", SVName: "termfire"})
		game.SendMessage(&messages.MessageSetup{
			FaceCache: struct {
				Use   bool
				Value bool
			}{Use: true, Value: true},
			LoginMethod: struct {
				Use   bool
				Value string
			}{Use: true, Value: "2"},
			ExtendedStats: struct {
				Use   bool
				Value bool
			}{Use: true, Value: true},
			Sound2: struct {
				Use   bool
				Value uint8
			}{Use: true, Value: 1},
		})

		// We need to add any messages here to our face to rune map cache, as otherwise they'll be lost to the void.
		l.On(&messages.MessageFace2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
			m := msg.(*messages.MessageFace2)

			r, fg, bg := play.GlobalObjectMapper.GetRuneAndColors(m.Name)
			play.GlobalObjectMapper.FaceToRune[uint16(m.Num)] = play.MapTile{R: play.MapRune(r), F: tcell.GetColor(fg), B: tcell.GetColor(bg)}

			debug.Debug("face2!", msg.Value())
		})

		l.Once(&messages.MessageSetup{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
			form := tview.NewForm().
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

			if account != "" && password != "" {
				l.TryLogin(game, account, password)
			}
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
		game.SetState(&Characters{
			Characters: m.Characters,
		})
	})

}
