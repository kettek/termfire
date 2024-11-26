package main

import (
	"errors"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Game struct {
	logPanel *LogPanel
	panels   []Panel
	conn     net.Conn
}

func (g *Game) SendMessage(msg Message) error {
	bytes := msg.Bytes()
	if len(bytes) > 0 {
		g.logPanel.Add("C->S" + string(bytes))
		g.conn.Write([]byte{byte(len(bytes) >> 8), byte(len(bytes))})
		g.conn.Write(bytes)
		return nil
	}
	return errors.New("empty message")
}

var gMessages []Message

type Message interface {
	Bytes() []byte
	Kind() string
	Value() string
	UnmarshalBinary([]byte) error
}

type MessageVersion struct {
	CLVersion string
	SVVersion string
	SVName    string
	value     string
}

func (m *MessageVersion) UnmarshalBinary(data []byte) error {
	parts := strings.SplitN(string(data), " ", 3)

	if len(parts) > 0 {
		m.CLVersion = parts[0]
	}
	if len(parts) > 1 {
		m.SVVersion = parts[1]
	}
	if len(parts) > 2 {
		m.SVName = parts[2]
	}

	return nil
}

func (m MessageVersion) Value() string {
	return m.CLVersion + " " + m.SVVersion + " " + m.SVName
}

func (m MessageVersion) Kind() string {
	return "version"
}

func (m MessageVersion) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(m.CLVersion)...)
	result = append(result, ' ')
	result = append(result, []byte(m.SVVersion)...)
	result = append(result, ' ')
	result = append(result, []byte(m.SVName)...)
	return result
}

type MessageFailure struct {
	Command string
	Reason  string
}

func (m *MessageFailure) UnmarshalBinary(data []byte) error {
	parts := strings.SplitN(string(data), " ", 2)

	if len(parts) > 0 {
		m.Command = parts[0]
	}
	if len(parts) > 1 {
		m.Reason = parts[1]
	}

	return nil
}

func (m MessageFailure) Value() string {
	return m.Command + " " + m.Reason
}

func (m MessageFailure) Kind() string {
	return "failure"
}

func (m MessageFailure) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(m.Command)...)
	result = append(result, ' ')
	result = append(result, []byte(m.Reason)...)
	return result
}

type MessageSetup struct {
	FaceCache bool
}

func (m *MessageSetup) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageSetup) Kind() string {
	return "setup"
}

func (m MessageSetup) Value() string {
	return strconv.FormatBool(m.FaceCache)
}

func (m MessageSetup) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte("facecache")...)
	result = append(result, ' ')
	result = append(result, '1')
	return result
}

type MessageAccountLogin struct {
	Account  string
	Password string
}

func (m *MessageAccountLogin) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageAccountLogin) Kind() string {
	return "accountlogin"
}

func (m MessageAccountLogin) Value() string {
	return m.Account + " " + m.Password
}

func (m MessageAccountLogin) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(LengthPrefixedString(m.Account))...)
	result = append(result, []byte(LengthPrefixedString(m.Password))...)
	return result
}

type Character struct {
	Name  string
	Level int
}

type MessageAccountPlayers struct {
	Characters []Character
}

const (
	ACL_NAME int = 1 << iota
	ACL_CLASS
	ACL_RACE
	ACL_LEVEL
	ACL_FACE
	ACL_PARTY
	ACL_MAP
	ACL_FACE_NUM
)

func (m *MessageAccountPlayers) UnmarshalBinary(data []byte) error {
	count := int(data[0])
	m.Characters = make([]Character, count)

	offset := 1
	for i := 0; i < count; i++ {
		length := int(data[offset])
		kind := int(data[offset+1])
		switch kind {
		case ACL_NAME:
			m.Characters[i].Name = string(data[offset+2 : offset+2+length])
		case ACL_LEVEL:
			m.Characters[i].Level = int(data[offset+2])<<8 | int(data[offset+3])
		}
		offset += length
	}

	return nil
}

func (m MessageAccountPlayers) Kind() string {
	return "accountplayers"
}

func (m MessageAccountPlayers) Value() string {
	var result string
	for _, c := range m.Characters {
		result += c.Name + " " + strconv.Itoa(c.Level) + "\n"
	}
	return result
}

func (m MessageAccountPlayers) Bytes() []byte {
	var result []byte

	result = append(result, []byte(m.Kind())...)

	result = append(result, byte(len(m.Characters)))

	for _, c := range m.Characters {
		result = append(result, []byte(c.Name)...)
		result = append(result, byte(c.Level))
	}

	return result
}

func init() {
	gMessages = []Message{
		&MessageFailure{},
		&MessageVersion{},
		&MessageAccountPlayers{},
	}
}

func UnmarshalMessage(data []byte) (Message, error) {
	msgLength := int(data[0])<<8 | int(data[1])
	msgType := ""
	for i := 2; i < msgLength; i++ {
		if data[i] == ' ' {
			msgType = string(data[2:i])
			break
		}
	}
	for _, m := range gMessages {
		if m.Kind() == msgType {
			if err := m.UnmarshalBinary(data[2+len(msgType)+1:]); err != nil {
				return nil, err
			} else {
				return m, nil
			}
		}
	}
	return nil, errors.New("unknown message type")
}

func LengthPrefixedString(s string) []byte {
	// 8-bit length followed by bytes.
	result := []byte{byte(len(s))}
	result = append(result, []byte(s)...)
	return result
}

func main() {
	if len(os.Args) < 4 {
		println("Usage: termfire <server> <account> <password>")
		os.Exit(1)
	}
	targetServer := os.Args[1]
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
			contents:  "ah",
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
				g.logPanel.Add("S->C" + string(buf[:n]))
				msg, err := UnmarshalMessage(buf[:n])
				if err != nil {
					s.PostEvent(tcell.NewEventInterrupt(err.Error()))
				} else {
					s.PostEvent(tcell.NewEventInterrupt(msg))
				}
			}
		}
	}()

	g.SendMessage(&MessageVersion{
		CLVersion: "1023",
		SVVersion: "1030",
		SVName:    "termfire",
	})

	g.SendMessage(&MessageSetup{
		FaceCache: true,
	})

	g.SendMessage(&MessageAccountLogin{
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
				g.logPanel.Add("RAW " + string(t))
			case string:
				g.logPanel.Add("STR " + t)
			case Message:
				g.logPanel.Add("MSG " + t.Kind() + " " + t.Value())
			}
		}

		for _, panel := range g.panels {
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
		p.CorePanel.contents += l
		if i != p.h {
			p.CorePanel.contents += "\n"
		}
	}
	p.CorePanel.Draw(s)
}
