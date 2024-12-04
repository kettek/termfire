package play

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MapRune rune

type MapTile struct {
	R MapRune
	F tcell.Color
	B tcell.Color
}

var FaceToRuneMap = map[uint16]MapTile{}

func ResetFaceToRuneMap() {
	FaceToRuneMap = map[uint16]MapTile{}
}

type RuneSize struct {
	Width  uint8
	Height uint8
}

var FaceToSizeMap = map[uint16]RuneSize{}

func ResetFaceToSizeMap() {
	FaceToSizeMap = map[uint16]RuneSize{}
}

type Map struct {
	View       *tview.Box
	layers     [10][][]MapTile
	width      int
	height     int
	viewWidth  int
	viewHeight int
	onResize   func(width, height int)
	onPostDraw func(screen tcell.Screen, x, y, width, height int)
}

func (m *Map) SetOnResize(onResize func(width, height int)) {
	m.onResize = onResize
}

func (m *Map) SetOnPostDraw(onPostDraw func(screen tcell.Screen, x, y, width, height int)) {
	m.onPostDraw = onPostDraw
}

func (m *Map) CenterX() int {
	return m.width / 2
}

func (m *Map) CenterY() int {
	return m.height / 2
}

func (m *Map) Init() {
	m.View = tview.NewBox()
	m.SetSize(11, 11)
	m.viewWidth = 11
	m.viewHeight = 11
	m.View.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		if width-2 != m.viewWidth || height-2 != m.viewHeight {
			m.viewWidth = width - 2
			m.viewHeight = height - 2
			if m.onResize != nil {
				m.onResize(m.viewWidth, m.viewHeight)
			}
		}
		// offset so we render in the center.
		x += (width - m.width) / 2
		y += (height - m.height) / 2

		// Draw from bottom-right to top-left
		for my := m.height - 1; my >= 0; my-- {
			for mx := m.width - 1; mx >= 0; mx-- {
				top := m.GetTopCell(mx, my)
				bot := m.GetBottomCell(mx, my)

				r := top.R
				fg := top.F
				bg := bot.B
				if top.F == bot.B {
					bg = top.B
				}

				style := tcell.StyleDefault

				if mx == m.CenterX() && my == m.CenterY() {
					style = style.Blink(true).Bold(true).Underline(true)
				}

				screen.SetContent(x+mx+1, y+my+1, rune(r), nil, style.Foreground(fg).Background(bg))
			}
		}
		if m.onPostDraw != nil {
			m.onPostDraw(screen, x, y, width, height)
		}
		return x, y, width, height
	})
}

func (m *Map) SetSize(width, height int) {
	m.width = width
	m.height = height

	m.layers = [10][][]MapTile{}
	for i := 0; i < 10; i++ {
		m.layers[i] = make([][]MapTile, m.height)
		for y := 0; y < m.height; y++ {
			m.layers[i][y] = make([]MapTile, m.width)
		}
	}
}

func (m *Map) Clear() {
	for i := 0; i < 10; i++ {
		for y := 0; y < len(m.layers[i]); y++ {
			for x := 0; x < len(m.layers[i][y]); x++ {
				m.layers[i][y][x] = MapTile{}
			}
		}
	}
}

func (m *Map) Shift(dx, dy int) {
	newLayers := [10][][]MapTile{}
	for i := 0; i < 10; i++ {
		newLayers[i] = make([][]MapTile, m.height)
		for y := 0; y < m.height; y++ {
			newLayers[i][y] = make([]MapTile, m.width)
		}
	}

	for i := 0; i < 10; i++ {
		for y := 0; y < m.height; y++ {
			for x := 0; x < m.width; x++ {
				newX := x + dx
				newY := y + dy
				if newX < 0 || newX >= m.width || newY < 0 || newY >= m.height {
					continue
				}
				newLayers[i][newY][newX] = m.layers[i][y][x]
			}
		}
	}
	m.layers = newLayers
}

func (m *Map) SetCell(x, y int, layer int, r MapRune, fg tcell.Color, bg tcell.Color) {
	if layer < 0 || layer >= 10 {
		return
	}
	if x < 0 || y < 0 {
		return
	}
	if x >= m.width || y >= m.height {
		return
	}
	m.layers[layer][y][x] = MapTile{r, fg, bg}
}

func (m *Map) ClearCell(x, y int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= m.width || y >= m.height {
		return
	}

	for i := 0; i < 10; i++ {
		m.layers[i][y][x] = MapTile{}
	}
}

func (m *Map) RemoveCellLayer(x, y, layer int) {
	if layer < 0 || layer >= 10 {
		return
	}
	if x < 0 || y < 0 {
		return
	}
	if x >= m.width || y >= m.height {
		return
	}
	m.layers[layer][y][x] = MapTile{}
}

func (m *Map) GetTopCell(x, y int) MapTile {
	for i := 9; i >= 0; i-- {
		if m.layers[i][y][x].R != 0 {
			return m.layers[i][y][x]
		}
	}

	return MapTile{}
}

func (m *Map) GetBottomCell(x, y int) MapTile {
	for i := 0; i < 10; i++ {
		if m.layers[i][y][x].R != 0 {
			return m.layers[i][y][x]
		}
	}
	return MapTile{}
}
