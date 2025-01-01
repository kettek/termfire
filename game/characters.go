package game

import (
	"fmt"

	"github.com/kettek/termfire/messages"
	"github.com/kettek/termfire/startup"
	"github.com/rivo/tview"
)

type Characters struct {
	MessageHandler
	Characters []messages.Character
}

func (c *Characters) Init(game Game) (tidy func()) {

	container := tview.NewFlex()
	container.SetDirection(tview.FlexColumn)

	containerCharacters := tview.NewFlex()
	containerCharacters.SetDirection(tview.FlexRow)
	containerCharacters.SetBorder(true)
	containerCharacters.SetTitle("Characters")
	characterList := tview.NewList()
	characterList.ShowSecondaryText(true)
	characterButtons := tview.NewForm()
	characterButtons.SetHorizontal(true)

	for _, character := range c.Characters {
		if character.Name == "" {
			continue
		}
		primary := fmt.Sprintf("%s Level %d %s", character.Name, character.Level, character.Class)
		secondary := fmt.Sprintf("  %s %s", character.Map, character.Party)
		characterList.AddItem(primary, secondary, 0, nil)
	}

	characterButtons.AddButton("Leave", func() {
		startup.Host = "" // Reset host so we don't reconnect
		game.Disconnect()
		game.SetState(&Servers{})
	})

	characterButtons.AddButton("Play", func() {
		char := c.Characters[characterList.GetCurrentItem()]
		game.SetState(&Play{character: char.Name})
	})

	containerCharacters.AddItem(characterList, 0, 1, true)
	containerCharacters.AddItem(characterButtons, 3, 1, false)

	container.AddItem(containerCharacters, 0, 1, true)

	go func() {
		if startup.Character != "" {
			game.SetState(&Play{character: startup.Character})
		} else {
			game.Pages().AddAndSwitchToPage("characters", container, true)
			game.Redraw()
		}
	}()

	return func() {
		game.Pages().RemovePage("characters")
	}
}
