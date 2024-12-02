package game

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/game/play"
	"github.com/kettek/termfire/messages"
	"github.com/rivo/tview"
)

var cfToW3CColor = map[messages.MessageColor]tcell.Color{
	messages.MessageColorBlack:      tcell.ColorBlack,
	messages.MessageColorWhite:      tcell.ColorWhite,
	messages.MessageColorNavy:       tcell.ColorNavy,
	messages.MessageColorRed:        tcell.ColorRed,
	messages.MessageColorOrange:     tcell.ColorOrange,
	messages.MessageColorBlue:       tcell.ColorBlue,
	messages.MessageColorDarkOrange: tcell.ColorDarkOrange,
	messages.MessageColorGreen:      tcell.ColorGreen,
	messages.MessageColorLightGreen: tcell.ColorLightGreen,
	messages.MessageColorGrey:       tcell.ColorGrey,
	messages.MessageColorBrown:      tcell.ColorBrown,
	messages.MessageColorGold:       tcell.ColorGold,
	messages.MessageColorTan:        tcell.ColorTan,
}

type Messages struct {
	view *tview.TextView
}

type Play struct {
	MessageHandler
	game      Game
	character string
	playerTag uint32
	inventory Inventory
	input     *tview.InputField
	mapp      play.Map
	messages  Messages
	topPacket uint16
}

type Object struct {
	Tag        uint32
	Name       string
	PluralName string
	Count      int
}

type Inventory struct {
	ListView *tview.List
	Items    []string
}

func (i *Inventory) AddItem(item string) {
	i.Items = append(i.Items, item)
	i.ListView.AddItem(item, "", 0, nil)
}

func (i *Inventory) RemoveItem(item string) {
}

func (i *Inventory) Clear() {
	i.Items = []string{}
	i.ListView.Clear()
}

func (p *Play) Init(game Game) (tidy func()) {
	p.game = game

	game.SendMessage(&messages.MessageAccountPlay{Character: p.character})

	p.Once(&messages.MessageAccountPlay{}, &messages.MessageAccountPlay{}, func(msg messages.Message, failure *messages.MessageFailure) {
		debug.Debug("bad character!")
		// TODO: Boot back to Login, but with a preserved login state...
	})

	// Setup our UI.
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.SetFocusFunc(func() {
		game.App().SetFocus(p.mapp.View)
	})

	left := tview.NewFlex()
	left.SetBorder(true)
	left.SetTitle("left")
	left.SetDirection(tview.FlexRow)
	middle := tview.NewFlex()
	middle.SetBorder(false)
	middle.SetTitle("middle")
	middle.SetDirection(tview.FlexRow)
	right := tview.NewFlex()
	right.SetBorder(false)
	right.SetDirection(tview.FlexRow)

	flex.AddItem(left, 0, 1, false)
	flex.AddItem(middle, 0, 2, false)
	flex.AddItem(right, 0, 1, false)

	p.messages.view = tview.NewTextView()
	p.messages.view.SetScrollable(true)
	p.messages.view.SetDynamicColors(true)
	p.messages.view.SetWrap(true)
	p.messages.view.SetWordWrap(true)

	p.input = tview.NewInputField()
	p.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			msg := &messages.MessageCommand{Command: p.input.GetText()}
			msg.Packet = p.topPacket
			msg.Repeat = 1
			p.topPacket++
			game.SendMessage(msg)
			p.input.SetText("")
		} else if key == tcell.KeyEsc {
			p.input.SetText("")
		}
		game.Pages().SwitchToPage("play")
	})

	right.AddItem(p.messages.view, 0, 1, false)
	right.AddItem(p.input, 1, 1, false)

	p.mapp.Init()
	p.mapp.SetOnResize(func(width, height int) {
		game.SendMessage(&messages.MessageSetup{
			MapSize: struct {
				Use   bool
				Value string
			}{Use: true, Value: fmt.Sprintf("%dx%d", width, height)},
		})
	})
	p.mapp.View.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		var msg *messages.MessageCommand
		if event.Key() == tcell.KeyUp {
			msg = &messages.MessageCommand{Command: "north"}
		} else if event.Key() == tcell.KeyDown {
			msg = &messages.MessageCommand{Command: "south"}
		} else if event.Key() == tcell.KeyLeft {
			msg = &messages.MessageCommand{Command: "west"}
		} else if event.Key() == tcell.KeyRight {
			msg = &messages.MessageCommand{Command: "east"}
		} else if event.Key() == tcell.KeyRune {
			if event.Rune() == 'a' {
				msg = &messages.MessageCommand{Command: "apply"}
			} else if event.Rune() == '\'' {
				game.App().SetFocus(p.input)
			} else if event.Rune() == 'i' {
				if pg, _ := game.Pages().GetFrontPage(); pg != "inventory" {
					game.Pages().SwitchToPage("inventory")
					game.App().SetFocus(p.inventory.ListView)
				}
			}
		}
		if msg != nil {
			msg.Packet = p.topPacket
			msg.Repeat = 1
			p.topPacket++
			game.SendMessage(msg)
		}
		return event
	})

	middle.AddItem(p.mapp.View, 0, 1, true)

	inventory := tview.NewFlex()
	inventory.SetTitle("Inventory")
	inventory.SetDirection(tview.FlexRow)

	p.inventory.ListView = tview.NewList()

	inventoryButtons := tview.NewFlex()
	inventoryButtons.SetDirection(tview.FlexColumn)
	inventoryButtons.AddItem(tview.NewButton("Apply").SetSelectedFunc(func() {
		item := p.inventory.Items[p.inventory.ListView.GetCurrentItem()]
		debug.Debug("apply ", item)
	}), 0, 1, false)
	inventoryButtons.AddItem(tview.NewButton("Drop").SetSelectedFunc(func() {
		item := p.inventory.Items[p.inventory.ListView.GetCurrentItem()]
		debug.Debug("drop ", item)
	}), 0, 1, false)
	inventoryButtons.AddItem(tview.NewButton("Examine").SetSelectedFunc(func() {
		item := p.inventory.Items[p.inventory.ListView.GetCurrentItem()]
		debug.Debug("examine ", item)
	}), 0, 1, false)
	inventoryButtons.AddItem(tview.NewButton("Lock").SetSelectedFunc(func() {
		item := p.inventory.Items[p.inventory.ListView.GetCurrentItem()]
		debug.Debug("lock ", item)
	}), 0, 1, false)
	inventoryButtons.AddItem(tview.NewButton("Mark").SetSelectedFunc(func() {
		item := p.inventory.Items[p.inventory.ListView.GetCurrentItem()]
		debug.Debug("mark ", item)
	}), 0, 1, false)
	inventory.AddItem(p.inventory.ListView, 0, 1, false)
	inventory.AddItem(inventoryButtons, 1, 1, false)

	game.Pages().AddPage("inventory", inventory, true, true)

	game.Pages().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			game.Pages().SwitchToPage("play")
		}
		return event
	})
	game.Pages().AddAndSwitchToPage("play", flex, true)
	game.Redraw()

	p.On(&messages.MessageSetup{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageSetup)
		if m.MapSize.Use {
			parts := strings.Split(m.MapSize.Value, "x")
			width, _ := strconv.Atoi(parts[0])
			height, _ := strconv.Atoi(parts[1])
			p.mapp.SetSize(width, height)
		}
	})

	p.On(&messages.MessageNewMap{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		p.mapp.Clear()
	})

	p.On(&messages.MessageImage2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageImage2)
		play.FaceToSizeMap[uint16(m.Face)] = play.RuneSize{Width: uint8(m.Width / 32), Height: uint8(m.Height / 32)}
	})

	p.On(&messages.MessageFace2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageFace2)

		if regexp.MustCompile(`(.*?)\.x\d\d`).MatchString(m.Name) {
			game.SendMessage(&messages.MessageAskFace{Face: uint32(m.Num)})
		}

		r := play.NameToTile(m.Name)
		play.FaceToRuneMap[uint16(m.Num)] = r
	})

	p.On(&messages.MessageSmooth{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
	})

	p.On(&messages.MessageDrawExtInfo{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageDrawExtInfo)

		colorizedText := fmt.Sprintf("[%s]%s[%s]", cfToW3CColor[m.Color], m.Message, cfToW3CColor[messages.MessageColorWhite])

		txt := p.messages.view.GetText(false)
		p.messages.view.SetText(txt + "\n" + colorizedText)
		p.messages.view.ScrollToEnd()
		game.Redraw()
	})

	p.On(&messages.MessageDeleteInventory{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		p.inventory.Clear()
		game.Redraw()
	})

	p.On(&messages.MessageItem2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageItem2)
		if m.Location == 0 {
			debug.Debug("ground items: ", m.Objects)
		} else {
			if m.Location == p.playerTag {
				for _, item := range m.Objects {
					if item.Nrof > 1 {
						p.inventory.AddItem(strconv.Itoa(int(item.Nrof)) + " " + item.PluralName)
					} else {
						p.inventory.AddItem(item.Name)
					}
				}
			} else {
				debug.Debug("some items @ ", m.Location, ": ", m.Objects)
			}
		}
		game.Redraw()
	})

	p.On(&messages.MessageStats{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageStats)
		debug.Debug("stats!", msg.Value())
		for _, stat := range m.Stats {
			switch s := stat.(type) {
			case *messages.MessageStatStr:
				debug.Debug("str: ", *s)
			}
		}
	})

	p.On(&messages.MessageMap2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageMap2)

		var setChanges []struct {
			x, y  int
			t     play.MapTile
			layer int
		}

		for _, m := range m.Coords {
			if m.Type == messages.MessageMap2CoordTypeScrollInformation {
				p.mapp.Shift(-int(m.X), -int(m.Y))
			}

			if len(m.Data) == 0 {
				// I think this is a "you are here" type message???
				if m.X == p.mapp.CenterX() && m.Y == p.mapp.CenterY() {
					//p.mapp.SetCell(m.X, m.Y, play.MapRunePlayer, tcell.ColorWhite, tcell.ColorBlack)
				}
				continue
			}
			for _, c := range m.Data {
				switch d := c.(type) {
				case *messages.MessageMap2CoordDataClear:
					p.mapp.ClearCell(m.X, m.Y)
				case *messages.MessageMap2CoordDataImage:
					t, ok := play.FaceToRuneMap[d.FaceNum]
					if !ok {
						t = play.MapTile{'?', tcell.ColorWhite, tcell.ColorBlack}
					}
					if d.FaceNum == 0 {
						p.mapp.RemoveCellLayer(m.X, m.Y, int(d.Layer))
						if size, ok := play.FaceToSizeMap[d.FaceNum]; ok {
							for x := 0; x < int(size.Width); x++ {
								for y := 0; y < int(size.Height); y++ {
									p.mapp.RemoveCellLayer(m.X-x, m.Y-y, int(d.Layer))
								}
							}
						}
					} else {
						setChanges = append(setChanges, struct {
							x, y  int
							t     play.MapTile
							layer int
						}{m.X, m.Y, t, int(d.Layer)})
						if size, ok := play.FaceToSizeMap[d.FaceNum]; ok {
							for x := 0; x < int(size.Width); x++ {
								for y := 0; y < int(size.Height); y++ {
									setChanges = append(setChanges, struct {
										x, y  int
										t     play.MapTile
										layer int
									}{m.X - x, m.Y - y, t, int(d.Layer)})
								}
							}
						}
					}
				}
			}
		}
		for _, change := range setChanges {
			p.mapp.SetCell(change.x, change.y, change.layer, change.t.R, change.t.F, change.t.B)
		}

		game.Redraw()
	})

	p.On(&messages.MessagePlayer{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessagePlayer)
		p.playerTag = m.Tag
		left.SetTitle(m.Name)
		game.Redraw()
		debug.Debug("player!", msg.Value())
	})

	p.On(&messages.MessageCommandCompleted{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageCommandCompleted)
		debug.Debug("command completed!", m)
	})

	return func() {
		game.Pages().RemovePage("play")
	}
}

func (p *Play) OnMessage(m messages.Message) {
	if !p.HasHandlerFor(m) {
		debug.Debug("Unhandled message: ", m.Kind())
		return
	}
	p.MessageHandler.OnMessage(m)
}
