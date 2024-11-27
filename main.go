package main

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/debug"
	"github.com/kettek/termfire/messages"
)

type Game struct {
	logPanel *LogPanel
	panels   []Panel
	conn     net.Conn
}

func (g *Game) SendMessage(msg messages.Message) error {
	bytes := msg.Bytes()
	if len(bytes) > 0 {
		g.logPanel.Add("C->S" + string(bytes))
		g.conn.Write([]byte{byte(len(bytes) >> 8), byte(len(bytes))})
		g.conn.Write(bytes)
		return nil
	}
	return errors.New("empty message")
}

func bytesToStringAndHex(b []byte) string {
	result := ""
	for _, c := range b {
		result += string(c) + "  "
	}
	result += "\n"
	for _, c := range b {
		result += strconv.FormatInt(int64(c), 16) + " "
	}
	return result
}

func main() {
	if len(os.Args) < 4 {
		println("Usage: termfire <server> <account> <password>")
		os.Exit(1)
	}
	targetServer := os.Args[1]

	if strings.Contains(targetServer, ":") == false {
		targetServer += ":13327"
	}

	account := os.Args[2]
	password := os.Args[3]

	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := s.Init(); err != nil {
		panic(err)
	}

	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	s.Clear()

	if err := debug.Start(); err != nil {
		panic(err)
	}

	var g Game

	w, h := s.Size()

	g.logPanel = &LogPanel{
		CorePanel: CorePanel{
			x:         0,
			y:         0,
			w:         w - 1,
			h:         h - 1,
			style:     defStyle,
			fillWidth: true,
			contents:  "",
		},
		lines: []string{},
	}
	g.panels = append(g.panels, g.logPanel)

	quit := func() {
		s.Fini()
		os.Exit(0)
	}

	conn, err := net.Dial("tcp", targetServer)
	if err != nil {
		g.logPanel.contents = "Failed to connect to server"
	} else {
		g.logPanel.contents = "Connected to server"
	}
	g.conn = conn

	go func() {
		for {
			var buf [4096]byte
			if n, err := g.conn.Read(buf[:]); err != nil {
				s.PostEvent(tcell.NewEventInterrupt("fail"))
			} else {
				var offset int
				for {
					msgLength := int(buf[offset])<<8 | int(buf[offset+1])
					offset += 2
					if msgLength == 0 || offset+msgLength > n {
						break
					}
					bytes := buf[offset : offset+msgLength]
					//g.logPanel.Add("S->C" + strconv.Itoa(msgLength) + "\n" + bytesToStringAndHex(bytes))
					msg, err := messages.UnmarshalMessage(bytes)
					if err != nil {
						s.PostEvent(tcell.NewEventInterrupt(err))
					} else {
						s.PostEvent(tcell.NewEventInterrupt(msg))
					}
					offset += msgLength
					if offset >= n {
						break
					}
				}
				//g.logPanel.Add("S->C" + string(buf[:n]))
				//g.logPanel.Add("S->C\n" + bytesToStringAndHex(buf[:n]))
				/*msg, err := UnmarshalMessage(buf[:n])
				if err != nil {
					s.PostEvent(tcell.NewEventInterrupt(err.Error()))
				} else {
					s.PostEvent(tcell.NewEventInterrupt(msg))
				}*/
			}
		}
	}()

	g.SendMessage(&messages.MessageVersion{
		CLVersion: "1023",
		SVVersion: "1030",
		SVName:    "termfire",
	})

	g.SendMessage(&messages.MessageSetup{
		FaceCache:   true,
		LoginMethod: "2",
	})

	g.SendMessage(&messages.MessageAccountLogin{
		Account:  account,
		Password: password,
	})

	for {
		s.Show()

		ev := s.PollEvent()
		w, h := s.Size()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
			for i := range g.panels {
				g.panels[i].Layout(w, h)
			}
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				quit()
			}
		case *tcell.EventInterrupt:
			switch t := ev.Data().(type) {
			case []byte:
				g.logPanel.Add("RAW " + string(t) + "\n")
			case string:
				g.logPanel.Add("STR " + t + "\n")
			case messages.Message:
				g.logPanel.Add("MSG " + t.Kind() + " " + t.Value() + "\n")
				g.logPanel.Add(t.Kind() + " " + t.Value() + "\n")
				for i, c := range t.Kind() + t.Value() {
					s.SetContent(i, 2, c, nil, defStyle)
				}
			case error:
				g.logPanel.Add("ERR " + t.Error() + "\n")
			default:
				g.logPanel.Add("UNKNOWN\n")
			}
		}

		for _, panel := range g.panels {
			panel.Layout(w, h)
			panel.Draw(s)
		}
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

type Panel interface {
	Draw(tcell.Screen)
	Layout(int, int)
}

type CorePanel struct {
	x, y, w, h int
	style      tcell.Style
	fillWidth  bool
	fillHeight bool
	contents   string
}

func (p *CorePanel) Draw(s tcell.Screen) {
	parts := strings.Split(p.contents, "\n")
	for i, part := range parts {
		drawText(s, p.x, p.y+i, p.x+p.w, p.y+i, p.style, part)
	}
}

func (p *CorePanel) Layout(w, h int) {
	if p.fillWidth {
		p.w = w - 1
	} else {
		if p.w > w {
			p.w = w - 1
		}
	}
	if p.fillHeight {
		p.h = h - 1
	} else {
		if p.h > h {
			p.h = h - 1
		}
	}
}

type LogPanel struct {
	CorePanel
	lines []string
}

func (p *LogPanel) Add(line string) {
	p.lines = append(p.lines, strconv.FormatInt(time.Now().UnixMilli(), 10)+": "+line)
}

func (p *LogPanel) Draw(s tcell.Screen) {
	p.CorePanel.contents = ""
	for i, l := range p.lines {
		// Split string if longer than p.w
		for len(l) > p.w {
			p.CorePanel.contents += l[:p.w]
			p.CorePanel.contents += "\n"
			l = l[p.w:]
		}
		if i != p.h {
			p.CorePanel.contents += "\n"
		}
	}
	p.CorePanel.Draw(s)
}
