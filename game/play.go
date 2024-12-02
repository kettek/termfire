package game

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/debug"
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

type mapRune rune

const (
	mapRuneWall   mapRune = '█'
	mapRuneWindow         = '▓'
	mapRuneStones         = '·'
	mapRuneDirt           = '░'
	mapRuneDoor           = '+'
	mapRuneGate           = '‡'
	mapRuneWater          = '~'
	mapRunePlayer         = '@'
	mapRuneCoin           = '¢'
	mapRuneBed            = '&'
	mapRuneTable          = 'T'
	mapRuneChair          = 'h'
	mapRuneScroll         = '!'
	mapRuneLever          = '/'
	mapHouse              = '#'
	mapShop               = '#'
	mapTower              = '#'
	mapPlant              = '♣'
	mapTree               = '♠'
	mapLight              = '☼'
	mapPond               = '≈'
	mapFountain           = '⌂'
	mapSign               = '☺'
	mapStatue             = '☻'
	mapWell               = 'O'
	mapEmpty              = ' '
	mapClock              = '♦'
)

type mapTile struct {
	r mapRune
	f tcell.Color
	b tcell.Color
}

var nameToMapTile = map[string]mapTile{
	"wall":      {mapRuneWall, tcell.ColorWhite, tcell.ColorBlack},
	"window":    {mapRuneWindow, tcell.ColorWhite, tcell.ColorBlack},
	"floor":     {mapRuneDirt, tcell.ColorWhite, tcell.ColorBlack},
	"stones":    {mapRuneStones, tcell.ColorGray, tcell.ColorBlack},
	"dirt":      {mapRuneDirt, tcell.ColorBrown, tcell.ColorBlack},
	"grass":     {mapRuneDirt, tcell.ColorGreen, tcell.ColorBlack},
	"ground":    {mapRuneDirt, tcell.ColorWhite, tcell.ColorBlack},
	"cobble":    {mapRuneDirt, tcell.ColorGray, tcell.ColorBlack},
	"door":      {mapRuneDoor, tcell.ColorWhite, tcell.ColorBlack},
	"gate":      {mapRuneGate, tcell.ColorGray, tcell.ColorBlack},
	"water":     {mapRuneWater, tcell.ColorBlue, tcell.ColorBlack},
	"player":    {mapRunePlayer, tcell.ColorWhite, tcell.ColorBlack},
	"coin":      {mapRuneCoin, tcell.ColorYellow, tcell.ColorBlack},
	"bed":       {mapRuneBed, tcell.ColorRed, tcell.ColorBlack},
	"table":     {mapRuneTable, tcell.ColorBeige, tcell.ColorBlack},
	"chair":     {mapRuneChair, tcell.ColorBeige, tcell.ColorBlack},
	"scroll":    {mapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"card":      {mapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"book":      {mapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"lever":     {mapRuneLever, tcell.ColorGray, tcell.ColorBlack},
	"house":     {mapHouse, tcell.ColorBlack, tcell.ColorWhite},
	"barrack":   {mapHouse, tcell.ColorBrown, tcell.ColorWhite},
	"tavern":    {mapHouse, tcell.ColorBeige, tcell.ColorWhite},
	"guild":     {mapHouse, tcell.ColorDarkGray, tcell.ColorWhite},
	"fort":      {mapHouse, tcell.ColorBlack, tcell.ColorWhite},
	"tower":     {mapTower, tcell.ColorBlack, tcell.ColorWhite},
	"shop":      {mapShop, tcell.ColorBlack, tcell.ColorYellow},
	"store":     {mapShop, tcell.ColorBlack, tcell.ColorYellow},
	"market":    {mapShop, tcell.ColorBlack, tcell.ColorYellow},
	"bank":      {mapShop, tcell.ColorBlack, tcell.ColorYellow},
	"shrine":    {mapHouse, tcell.ColorBlue, tcell.ColorBlack},
	"church":    {mapHouse, tcell.ColorBlue, tcell.ColorBlack},
	"inn":       {mapHouse, tcell.ColorBeige, tcell.ColorBlack},
	"shrub":     {mapPlant, tcell.ColorGreen, tcell.ColorBlack},
	"brush":     {mapPlant, tcell.ColorGreen, tcell.ColorBlack},
	"tree":      {mapTree, tcell.ColorGreen, tcell.ColorBlack},
	"lamp":      {mapLight, tcell.ColorYellow, tcell.ColorBlack},
	"pond":      {mapPond, tcell.ColorBlue, tcell.ColorBlack},
	"lake":      {mapPond, tcell.ColorBlue, tcell.ColorBlack},
	"grasspond": {mapPond, tcell.ColorBlue, tcell.ColorGreen},
	"fountain":  {mapFountain, tcell.ColorBlue, tcell.ColorBlack},
	"sign":      {mapSign, tcell.ColorWhite, tcell.ColorBlack},
	"crossroad": {mapSign, tcell.ColorWhite, tcell.ColorBlack},
	"statue":    {mapStatue, tcell.ColorWhite, tcell.ColorBlack},
	"well":      {mapWell, tcell.ColorBlue, tcell.ColorBlack},
	"woods":     {mapTree, tcell.ColorGreen, tcell.ColorBlack},
	"empty":     {mapEmpty, tcell.ColorBlack, tcell.ColorBlack},
	"clock":     {mapClock, tcell.ColorWhite, tcell.ColorBlack},
}

func nameToTile(name string) mapTile {
	for k, v := range nameToMapTile {
		if strings.Contains(name, k) {
			return v
		}
	}
	debug.Debug("missing image: ", name)
	return mapTile{mapRune(name[0]), tcell.ColorWhite, tcell.ColorBlack}
}

var faceToRuneMap = map[uint16]mapTile{}

type Map struct {
	view  *tview.Box
	cells [11][11]mapTile // TODO: Make resizeable.
}

func (m *Map) Clear() {
	for y := 0; y < len(m.cells); y++ {
		for x := 0; x < len(m.cells[y]); x++ {
			m.cells[y][x] = mapTile{' ', tcell.ColorWhite, tcell.ColorBlack}
		}
	}
}

func (m *Map) Shift(dx, dy int) {
	newCells := [11][11]mapTile{}
	for y := 0; y < 11; y++ {
		for x := 0; x < 11; x++ {
			newX := x + dx
			newY := y + dy
			if newX < 0 || newX >= 11 || newY < 0 || newY >= 11 {
				continue
			}
			newCells[newY][newX] = m.cells[y][x]
		}
	}
	m.cells = newCells
}

func (m *Map) SetCell(x, y int, r mapRune, fg tcell.Color, bg tcell.Color) {
	if x < 0 || y < 0 {
		return
	}
	if x >= 11 || y >= 11 {
		return
	}
	m.cells[y][x] = mapTile{r, fg, bg}
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
	mapp      Map
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
		game.App().SetFocus(p.mapp.view)
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
	right.SetBorder(true)
	right.SetTitle("right")
	right.SetDirection(tview.FlexRow)

	flex.AddItem(left, 0, 1, false)
	flex.AddItem(middle, 0, 2, false)
	flex.AddItem(right, 0, 1, false)

	p.messages.view = tview.NewTextView()
	p.messages.view.SetScrollable(true)
	p.messages.view.SetDynamicColors(true)
	p.messages.view.SetWrap(true)
	p.messages.view.SetWordWrap(true)
	right.AddItem(p.messages.view, 0, 1, false)

	p.mapp.view = tview.NewBox()
	p.mapp.view.SetBorder(true)
	p.mapp.view.SetTitle("Map")
	p.mapp.view.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// offset so we render in the center.
		x += (width - 11) / 2
		y += (height - 11) / 2

		for my := 0; my < 11; my++ {
			for mx := 0; mx < 11; mx++ {
				r, fg, bg := p.mapp.cells[my][mx].r, p.mapp.cells[my][mx].f, p.mapp.cells[my][mx].b
				screen.SetContent(x+mx+1, y+my+1, rune(r), nil, tcell.StyleDefault.Foreground(fg).Background(bg))
			}
		}
		return x, y, width, height
	})
	p.mapp.view.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
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

	middle.AddItem(p.mapp.view, 0, 1, true)

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
		if event.Key() == tcell.KeyRune {
			switch event.Rune() {
			case 'i':
				if pg, _ := game.Pages().GetFrontPage(); pg != "inventory" {
					game.Pages().SwitchToPage("inventory")
					game.App().SetFocus(p.inventory.ListView)
				}
			}
		} else if event.Key() == tcell.KeyEsc {
			game.Pages().SwitchToPage("play")
		}
		return event
	})
	game.Pages().AddAndSwitchToPage("play", flex, true)
	game.Redraw()

	p.On(&messages.MessageNewMap{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		p.mapp.Clear()
	})

	p.On(&messages.MessageFace2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageFace2)
		r := nameToTile(m.Name)
		faceToRuneMap[uint16(m.Num)] = r
		debug.Debug("face2!", msg.Value())
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
			t     mapTile
			layer int
		}

		for _, m := range m.Coords {
			if m.Type == messages.MessageMap2CoordTypeScrollInformation {
				p.mapp.Shift(-int(m.X), -int(m.Y))
			}

			if len(m.Data) == 0 {
				// I think this is a "you are here" type message???
				if m.X == 5 && m.Y == 5 {
					p.mapp.SetCell(m.X, m.Y, mapRunePlayer, tcell.ColorWhite, tcell.ColorBlack)
				}
				continue
			}
			for _, c := range m.Data {
				switch d := c.(type) {
				case *messages.MessageMap2CoordDataClear:
					p.mapp.SetCell(m.X, m.Y, ' ', tcell.ColorBlack, tcell.ColorBlack)
				case *messages.MessageMap2CoordDataImage:
					t, ok := faceToRuneMap[d.FaceNum]
					if !ok {
						t = mapTile{'?', tcell.ColorWhite, tcell.ColorBlack}
					}
					found := false
					for i, change := range setChanges {
						if change.x == m.X && change.y == m.Y {
							found = true
							if change.layer <= int(d.Layer) {
								setChanges[i].t = t
								break
							}
						}
					}
					if !found {
						setChanges = append(setChanges, struct {
							x, y  int
							t     mapTile
							layer int
						}{m.X, m.Y, t, int(d.Layer)})
					}
				}
			}
		}
		for _, change := range setChanges {
			p.mapp.SetCell(change.x, change.y, change.t.r, change.t.f, change.t.b)
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
