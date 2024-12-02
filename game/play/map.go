package play

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/debug"
	"github.com/rivo/tview"
)

type MapRune rune

const (
	MapRuneWall   MapRune = '█'
	MapRuneWindow         = '▓'
	MapRuneStones         = '·'
	MapRuneDirt           = '░'
	MapRuneDoor           = '+'
	MapRuneGate           = '‡'
	MapRuneWater          = '~'
	MapRunePlayer         = '@'
	MapRuneCoin           = '¢'
	MapRuneBed            = '&'
	MapRuneTable          = 'T'
	MapRuneChair          = 'h'
	MapRuneScroll         = '!'
	MapRuneLever          = '/'
	MapHouse              = '#'
	MapShop               = '#'
	MapTower              = '#'
	MapPlant              = '♣'
	MapTree               = '♠'
	MapLight              = '☼'
	MapPond               = '≈'
	MapFountain           = '⌂'
	MapSign               = '☺'
	MapStatue             = '☻'
	MapWell               = 'O'
	MapEmpty              = ' '
	MapClock              = '♦'
)

type MapTile struct {
	R MapRune
	F tcell.Color
	B tcell.Color
}

var NameToMapTile = map[string]MapTile{
	"wall":      {MapRuneWall, tcell.ColorWhite, tcell.ColorBlack},
	"window":    {MapRuneWindow, tcell.ColorWhite, tcell.ColorBlack},
	"floor":     {MapRuneDirt, tcell.ColorWhite, tcell.ColorBlack},
	"stones":    {MapRuneStones, tcell.ColorGray, tcell.ColorBlack},
	"dirt":      {MapRuneDirt, tcell.ColorBrown, tcell.ColorBlack},
	"grass":     {MapRuneDirt, tcell.ColorGreen, tcell.ColorBlack},
	"ground":    {MapRuneDirt, tcell.ColorWhite, tcell.ColorBlack},
	"cobble":    {MapRuneDirt, tcell.ColorGray, tcell.ColorBlack},
	"door":      {MapRuneDoor, tcell.ColorWhite, tcell.ColorBlack},
	"gate":      {MapRuneGate, tcell.ColorGray, tcell.ColorBlack},
	"water":     {MapRuneWater, tcell.ColorBlue, tcell.ColorBlack},
	"player":    {MapRunePlayer, tcell.ColorWhite, tcell.ColorBlack},
	"coin":      {MapRuneCoin, tcell.ColorYellow, tcell.ColorBlack},
	"bed":       {MapRuneBed, tcell.ColorRed, tcell.ColorBlack},
	"table":     {MapRuneTable, tcell.ColorBeige, tcell.ColorBlack},
	"chair":     {MapRuneChair, tcell.ColorBeige, tcell.ColorBlack},
	"scroll":    {MapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"card":      {MapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"book":      {MapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"lever":     {MapRuneLever, tcell.ColorGray, tcell.ColorBlack},
	"house":     {MapHouse, tcell.ColorBlack, tcell.ColorWhite},
	"barrack":   {MapHouse, tcell.ColorBrown, tcell.ColorWhite},
	"tavern":    {MapHouse, tcell.ColorBeige, tcell.ColorWhite},
	"guild":     {MapHouse, tcell.ColorDarkGray, tcell.ColorWhite},
	"fort":      {MapHouse, tcell.ColorBlack, tcell.ColorWhite},
	"tower":     {MapTower, tcell.ColorBlack, tcell.ColorWhite},
	"shop":      {MapShop, tcell.ColorBlack, tcell.ColorYellow},
	"store":     {MapShop, tcell.ColorBlack, tcell.ColorYellow},
	"market":    {MapShop, tcell.ColorBlack, tcell.ColorYellow},
	"bank":      {MapShop, tcell.ColorBlack, tcell.ColorYellow},
	"shrine":    {MapHouse, tcell.ColorBlue, tcell.ColorBlack},
	"church":    {MapHouse, tcell.ColorBlue, tcell.ColorBlack},
	"inn":       {MapHouse, tcell.ColorBeige, tcell.ColorBlack},
	"shrub":     {MapPlant, tcell.ColorGreen, tcell.ColorBlack},
	"brush":     {MapPlant, tcell.ColorGreen, tcell.ColorBlack},
	"tree":      {MapTree, tcell.ColorGreen, tcell.ColorBlack},
	"lamp":      {MapLight, tcell.ColorYellow, tcell.ColorBlack},
	"pond":      {MapPond, tcell.ColorBlue, tcell.ColorBlack},
	"lake":      {MapPond, tcell.ColorBlue, tcell.ColorBlack},
	"grasspond": {MapPond, tcell.ColorBlue, tcell.ColorGreen},
	"fountain":  {MapFountain, tcell.ColorBlue, tcell.ColorBlack},
	"sign":      {MapSign, tcell.ColorWhite, tcell.ColorBlack},
	"crossroad": {MapSign, tcell.ColorWhite, tcell.ColorBlack},
	"statue":    {MapStatue, tcell.ColorWhite, tcell.ColorBlack},
	"well":      {MapWell, tcell.ColorBlue, tcell.ColorBlack},
	"woods":     {MapTree, tcell.ColorGreen, tcell.ColorBlack},
	"empty":     {MapEmpty, tcell.ColorBlack, tcell.ColorBlack},
	"clock":     {MapClock, tcell.ColorWhite, tcell.ColorBlack},
}

func NameToTile(name string) MapTile {
	for k, v := range NameToMapTile {
		if strings.Contains(name, k) {
			return v
		}
	}
	debug.Debug("missing image: ", name)
	return MapTile{MapRune(name[0]), tcell.ColorWhite, tcell.ColorBlack}
}

var FaceToRuneMap = map[uint16]MapTile{}

type Map struct {
	View       *tview.Box
	cells      [][]MapTile // TODO: Make resizeable.
	width      int
	height     int
	viewWidth  int
	viewHeight int
	onResize   func(width, height int)
}

func (m *Map) SetOnResize(onResize func(width, height int)) {
	m.onResize = onResize
}

func (m *Map) Init() {
	m.View = tview.NewBox()
	m.View.SetBorder(true)
	m.View.SetTitle("Map")
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
			debug.Debug("resize map: ", m.viewWidth, m.viewHeight)
		}
		// offset so we render in the center.
		x += (width - m.width) / 2
		y += (height - m.height) / 2

		for my := 0; my < m.height; my++ {
			for mx := 0; mx < m.width; mx++ {
				r, fg, bg := m.cells[my][mx].R, m.cells[my][mx].F, m.cells[my][mx].B
				screen.SetContent(x+mx+1, y+my+1, rune(r), nil, tcell.StyleDefault.Foreground(fg).Background(bg))
			}
		}
		return x, y, width, height
	})
}

func (m *Map) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.cells = make([][]MapTile, m.height)
	for y := 0; y < m.height; y++ {
		m.cells[y] = make([]MapTile, m.width)
	}
}

func (m *Map) Clear() {
	for y := 0; y < len(m.cells); y++ {
		for x := 0; x < len(m.cells[y]); x++ {
			m.cells[y][x] = MapTile{' ', tcell.ColorWhite, tcell.ColorBlack}
		}
	}
}

func (m *Map) Shift(dx, dy int) {
	newCells := make([][]MapTile, m.height)
	for y := 0; y < m.height; y++ {
		newCells[y] = make([]MapTile, m.width)
	}

	for y := 0; y < m.height; y++ {
		for x := 0; x < m.width; x++ {
			newX := x + dx
			newY := y + dy
			if newX < 0 || newX >= m.width || newY < 0 || newY >= m.height {
				continue
			}
			newCells[newY][newX] = m.cells[y][x]
		}
	}
	m.cells = newCells
}

func (m *Map) SetCell(x, y int, r MapRune, fg tcell.Color, bg tcell.Color) {
	if x < 0 || y < 0 {
		return
	}
	if x >= m.width || y >= m.height {
		return
	}
	m.cells[y][x] = MapTile{r, fg, bg}
}
