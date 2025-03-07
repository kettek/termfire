package game

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/assets"
	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/game/play"
	"github.com/kettek/termfire/messages"
	"github.com/kettek/termfire/startup"
	"github.com/rivo/tview"
)

type Messages struct {
	view *tview.TextView
}

func (m *Messages) Init() {
	m.view = tview.NewTextView()
	m.view.SetScrollable(true)
	m.view.SetDynamicColors(true)
	m.view.SetWrap(true)
	m.view.SetWordWrap(true)
}

func (m *Messages) Add(msg string, color messages.MessageColor) {
	// Most terminals are black background, so I guess we can swap black to something else...
	if color == messages.MessageColorBlack {
		color = messages.MessageColorAltBlack
	}

	colorizedText := fmt.Sprintf("[%s]%s[%s]", play.CF2W3CColor[color], msg, play.CF2W3CColor[messages.MessageColorWhite])

	txt := m.view.GetText(false)
	m.view.SetText(txt + "\n" + colorizedText)
	m.view.ScrollToEnd()
}

type KeyBind struct {
	Key     tcell.Key
	Rune    rune
	ModMask tcell.ModMask
	Command string
}

type Play struct {
	messages.MessageHandler
	game            Game
	character       string
	playerTag       int32
	objectDebugView play.ObjectDebugView
	inventory       play.Container
	ground          play.Container
	status          play.Status
	input           *tview.InputField
	mapp            play.Map
	messages        Messages
	sounds          Messages
	topPacket       uint16
	lastDir         string // Not sure if we can query this instead...
	pendingBind     string
	binds           []KeyBind
}

func (p *Play) Init(game Game) (tidy func()) {
	p.game = game

	game.SendMessage(&messages.MessageAccountPlay{Character: p.character})

	p.Once(&messages.MessageAccountPlay{}, &messages.MessageAccountPlay{}, func(msg messages.Message, failure *messages.MessageFailure) {
		debug.Debug("bad character!")
		// TODO: Boot back to Login, but with a preserved login state...
	})

	// Load in our tilemap to objectMapper
	tbytes, err := assets.FS.ReadFile("tilemap.txt")
	if err != nil {
		panic(err)
	}
	play.GlobalObjectMapper.UnmarshalBinary(tbytes)

	// Setup our UI.
	flex := tview.NewFlex()
	flex.SetDirection(tview.FlexColumn)
	flex.SetFocusFunc(func() {
		game.App().SetFocus(p.mapp.View)
	})

	p.objectDebugView.Init()

	left := tview.NewFlex()
	left.SetDirection(tview.FlexRow)

	p.status.Init()
	left.AddItem(p.status.View, 0, 2, false)

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

	p.ground.Init("Ground", []string{"Take", "Examine", "Apply"})
	p.ground.SetOnTrigger(func(button string, object messages.ItemObject, index int) {
		debug.Debug("triggered: ", button, object.Tag)

		debug.Debug("object ", fmt.Sprintf("%d", object.Tag))
		if button == "Take" {
			game.SendMessage(&messages.MessageMove{To: p.playerTag, Tag: object.Tag})
		} else if button == "Examine" {
			game.SendMessage(&messages.MessageExamine{Tag: object.Tag})
		} else if button == "Apply" {
			game.SendMessage(&messages.MessageApply{Tag: object.Tag})
		}
	})
	left.AddItem(p.ground.GetContainer(), 0, 1, false)

	p.sounds.Init()
	p.messages.Init()

	p.input = tview.NewInputField()
	p.input.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			txt := p.input.GetText()
			if strings.HasPrefix(txt, "bind") {
				if txt == "bind" {
					p.messages.Add("bind requires one or more commands", messages.MessageColorRed)
				} else {
					p.pendingBind = txt[5:]
					p.messages.Add(fmt.Sprintf("Press a key to bind %s to: ", p.pendingBind), messages.MessageColorGrey)
				}
			} else if strings.HasPrefix(txt, "unbind") {
				if txt == "unbind" {
					for i := 0; i < len(p.binds); i++ {
						p.messages.Add(fmt.Sprintf("%d: %s", i, p.binds[i].Command), messages.MessageColorGrey)
					}
				} else {
					remaining := txt[7:]
					index, err := strconv.Atoi(remaining)
					if err != nil {
						p.messages.Add("Invalid index", messages.MessageColorRed)
					} else {
						p.Unbind(index)
						p.messages.Add(fmt.Sprintf("Unbound %d", index), messages.MessageColorGrey)
					}
				}
			} else {
				msg := &messages.MessageCommand{Command: txt}
				msg.Packet = p.topPacket
				msg.Repeat = 1
				p.topPacket++
				game.SendMessage(msg)
			}
			p.input.SetText("")
		} else if key == tcell.KeyEsc {
			p.input.SetText("")
		}
		game.Pages().SwitchToPage("play")
	})

	right.AddItem(p.sounds.view, 0, 1, false)
	right.AddItem(p.messages.view, 0, 4, false)
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
		// debug.Debug("event: ", event.Key(), event.Rune(), event.Modifiers(), event.Name())
		if p.pendingBind != "" {
			if event.Key() == tcell.KeyEsc {
				p.pendingBind = ""
				p.messages.Add("Bind canceled", messages.MessageColorGrey)
			} else {
				p.Bind(event.Key(), event.Rune(), event.Modifiers(), p.pendingBind)
				p.messages.Add(fmt.Sprintf("Bound %s to %s", p.pendingBind, event.Name()), messages.MessageColorGrey)
				p.pendingBind = ""
			}
			return event
		}

		if p.CheckBinds(event.Key(), event.Rune(), event.Modifiers()) {
			return event
		}

		firing := false

		// Add run for directional ctrl.
		if event.Key() >= tcell.KeyUp && event.Key() <= tcell.KeyLeft {
			if event.Modifiers() == tcell.ModShift {
				firing = true
			} else if event.Modifiers() == tcell.ModCtrl {
				p.SendCommand(&messages.MessageCommand{Command: "run"})
			} else {
				p.SendCommand(&messages.MessageCommand{Command: "run_stop"})
			}
		}

		cmd := ""

		if firing {
			cmd = "fire "
		}

		if event.Key() == tcell.KeyUp {
			cmd += "north"
			p.lastDir = "north"
		} else if event.Key() == tcell.KeyDown {
			cmd += "south"
			p.lastDir = "south"
		} else if event.Key() == tcell.KeyLeft {
			cmd += "west"
			p.lastDir = "west"
		} else if event.Key() == tcell.KeyRight {
			cmd += "east"
			p.lastDir = "east"
		} else if event.Key() == tcell.KeyRune {
			if event.Rune() == 'a' {
				cmd = "apply"
			} else if event.Rune() == '\'' {
				game.App().SetFocus(p.input)
			} else if event.Rune() == 'i' {
				if pg, _ := game.Pages().GetFrontPage(); pg != "inventory" {
					game.Pages().SwitchToPage("inventory")
					game.App().SetFocus(p.inventory.GetList())
				}
			}
		}

		if cmd != "" {
			p.SendCommand(&messages.MessageCommand{Command: cmd})
		}

		if firing {
			p.SendCommand(&messages.MessageCommand{Command: "fire_stop"})
		}

		return event
	})
	p.mapp.SetOnPostDraw(func(screen tcell.Screen, x, y, width, height int) {
		// TODO: Make this optional/configurable.
		cx := p.mapp.CenterX()
		cy := p.mapp.CenterY()
		switch p.lastDir {
		case "north":
			screen.SetContent(x+cx+1, y+cy+1, '↑', nil, tcell.StyleDefault)
		case "south":
			screen.SetContent(x+cx+1, y+cy+1, '↓', nil, tcell.StyleDefault)
		case "west":
			screen.SetContent(x+cx+1, y+cy+1, '←', nil, tcell.StyleDefault)
		case "east":
			screen.SetContent(x+cx+1, y+cy+1, '→', nil, tcell.StyleDefault)
		}
	})
	p.mapp.SetOnClick(func(button int, x, y int) {
		dx := x - p.mapp.CenterX() - 1
		dy := y - p.mapp.CenterY() - 1
		game.SendMessage(&messages.MessageLookAt{DX: dx, DY: dy})
	})

	middle.AddItem(p.mapp.View, 0, 1, true)

	p.inventory.Init("Inventory", []string{"Apply", "Drop", "Examine", "Lock", "Mark"})
	p.inventory.SetOnTrigger(func(button string, object messages.ItemObject, index int) {
		if button == "Apply" {
			game.SendMessage(&messages.MessageApply{Tag: object.Tag})
		} else if button == "Drop" {
			game.SendMessage(&messages.MessageMove{To: 0, Tag: object.Tag})
		} else if button == "Examine" {
			game.SendMessage(&messages.MessageExamine{Tag: object.Tag})
		} else if button == "Lock" {
			// TODO: Check if object is locked or not and toggle based on that.
			game.SendMessage(&messages.MessageLock{Tag: object.Tag})
		} else if button == "Mark" {
			game.SendMessage(&messages.MessageMark{Tag: object.Tag})
		}
	})

	game.Pages().AddPage("inventory", p.inventory.GetContainer(), true, true)

	game.Pages().AddPage("objectDebug", p.objectDebugView.GetContainer(), true, true)
	p.objectDebugView.GetContainer().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			game.Pages().SwitchToPage("play")
		}
		return event
	})

	game.Pages().SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			game.Pages().SwitchToPage("play")
		} else if event.Key() == tcell.KeyF1 {
			game.Pages().SwitchToPage("objectDebug")
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

	p.On(&messages.MessageAnim{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageAnim)
		play.GlobalObjectMapper.AnimToFaces[m.AnimID] = m.Faces
	})

	p.On(&messages.MessageImage2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageImage2)
		play.GlobalObjectMapper.FaceToSize[int16(m.Face)] = play.RuneSize{Width: uint8(m.Width / 32), Height: uint8(m.Height / 32)}
	})

	p.On(&messages.MessageFace2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageFace2)

		if regexp.MustCompile(`(.*?)\.x\d\d`).MatchString(m.Name) {
			game.SendMessage(&messages.MessageAskFace{Face: int32(m.Num)})
		}

		r, fg, bg := play.GlobalObjectMapper.GetRuneAndColors(m.Name)
		if r == 0 {
			r = rune(m.Name[0])
		}
		play.GlobalObjectMapper.FaceToName[int16(m.Num)] = m.Name
		play.GlobalObjectMapper.FaceToRune[int16(m.Num)] = play.MapTile{R: play.MapRune(r), F: tcell.GetColor(fg), B: tcell.GetColor(bg)}
		p.objectDebugView.Refresh()
	})

	p.On(&messages.MessageSmooth{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
	})

	p.On(&messages.MessageDrawExtInfo{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageDrawExtInfo)
		p.messages.Add(m.Message, m.Color)
		game.Redraw()
	})

	p.On(&messages.MessageDeleteInventory{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		p.inventory.Clear()
		game.Redraw()
	})

	p.On(&messages.MessageItem2{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageItem2)
		if m.Location == 0 {
			p.ground.Clear()
			for _, item := range m.Objects {
				p.ground.AddItem(item)
			}
		} else {
			if m.Location == p.playerTag {
				for _, item := range m.Objects {
					p.inventory.AddItem(item)
				}
			} else {
				debug.Debug("some items @ ", m.Location, ": ", m.Objects)
			}
		}
		game.Redraw()
	})

	p.On(&messages.MessageUpdateItem{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageUpdateItem)
		debug.Debug("update item: ", m)
	})

	p.On(&messages.MessageStats{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageStats)
		p.status.Update(m)
	})

	p.On(&messages.MessageSound{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageSound)

		var dirstring string
		if m.Y < 0 {
			dirstring = "north"
		} else if m.Y > 0 {
			dirstring = "south"
		}
		if m.X < 0 {
			dirstring += "west"
		} else if m.X > 0 {
			dirstring += "east"
		}
		if dirstring == "" {
			dirstring = "under you"
		} else {
			dirstring = "to the " + dirstring
		}

		prefix := "a"
		if strings.ContainsAny(m.Name[:1], "aeiou") {
			prefix = "an"
		}

		p.sounds.Add(fmt.Sprintf("You hear %s %s %s", prefix, m.Action, dirstring), messages.MessageColorGrey)
		game.Redraw()
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
					// p.mapp.SetCell(m.X, m.Y, play.MapRunePlayer, tcell.ColorWhite, tcell.ColorBlack)
				}
				continue
			}
			for _, c := range m.Data {
				switch d := c.(type) {
				case messages.MessageMap2CoordDataClear:
					p.mapp.ClearCell(m.X, m.Y)
				case messages.MessageMap2CoordDataClearLayer:
					p.mapp.RemoveCellLayer(m.X, m.Y, int(d.Layer))
				case messages.MessageMap2CoordDataAnim:
					if faces, ok := play.GlobalObjectMapper.AnimToFaces[d.Anim]; ok && len(faces) > 0 {
						setChanges = append(setChanges, struct {
							x, y  int
							t     play.MapTile
							layer int
						}{m.X, m.Y, play.GlobalObjectMapper.FaceToRune[faces[0]], int(d.Layer)})
					}
				case messages.MessageMap2CoordDataDarkness:
					// debug.Debug("darkness: ", d)
				case messages.MessageMap2CoordDataImage:
					// debug.Debug("image: ", d)
					t, ok := play.GlobalObjectMapper.FaceToRune[d.FaceNum]
					if !ok {
						t = play.MapTile{'?', tcell.ColorWhite, tcell.ColorBlack}
					}
					if d.FaceNum == 0 {
						p.mapp.RemoveCellLayer(m.X, m.Y, int(d.Layer))
						if size, ok := play.GlobalObjectMapper.FaceToSize[d.FaceNum]; ok {
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
						if size, ok := play.GlobalObjectMapper.FaceToSize[d.FaceNum]; ok {
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
		p.status.SetTitle(m.Name)
		game.Redraw()
		debug.Debug("player!", msg.Value())
	})

	p.On(&messages.MessageAccountPlayers{}, nil, func(msg messages.Message, failure *messages.MessageFailure) {
		m := msg.(*messages.MessageAccountPlayers)
		// Set startup character to blank so we don't relogin.
		startup.Character = ""
		game.SetState(&Characters{Characters: m.Characters})
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

func (p *Play) SendCommand(msg *messages.MessageCommand) {
	// debug.Debug("command: ", msg.Command)
	msg.Packet = p.topPacket
	// msg.Repeat = 1
	p.topPacket++
	p.game.SendMessage(msg)
}

func (p *Play) Bind(key tcell.Key, r rune, modMask tcell.ModMask, command string) {
	p.binds = append(p.binds, KeyBind{Key: key, Rune: r, ModMask: modMask, Command: command})
}

func (p *Play) Unbind(index int) {
	if index < 0 || index >= len(p.binds) {
		return
	}
	p.binds = append(p.binds[:index], p.binds[index+1:]...)
}

func (p *Play) CheckBinds(key tcell.Key, r rune, modMask tcell.ModMask) bool {
	for _, bind := range p.binds {
		if bind.Key == key && bind.Rune == r && bind.ModMask == modMask {
			p.SendCommand(&messages.MessageCommand{Command: bind.Command})
			return true
		}
	}
	return false
}
