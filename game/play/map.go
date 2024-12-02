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
	View  *tview.Box
	cells [11][11]MapTile // TODO: Make resizeable.
}

func (m *Map) Init() {
	m.View = tview.NewBox()
	m.View.SetBorder(true)
	m.View.SetTitle("Map")
	m.View.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		// offset so we render in the center.
		x += (width - 11) / 2
		y += (height - 11) / 2

		for my := 0; my < 11; my++ {
			for mx := 0; mx < 11; mx++ {
				r, fg, bg := m.cells[my][mx].R, m.cells[my][mx].F, m.cells[my][mx].B
				screen.SetContent(x+mx+1, y+my+1, rune(r), nil, tcell.StyleDefault.Foreground(fg).Background(bg))
			}
		}
		return x, y, width, height
	})
}

func (m *Map) Clear() {
	for y := 0; y < len(m.cells); y++ {
		for x := 0; x < len(m.cells[y]); x++ {
			m.cells[y][x] = MapTile{' ', tcell.ColorWhite, tcell.ColorBlack}
		}
	}
}

func (m *Map) Shift(dx, dy int) {
	newCells := [11][11]MapTile{}
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

func (m *Map) SetCell(x, y int, r MapRune, fg tcell.Color, bg tcell.Color) {
	if x < 0 || y < 0 {
		return
	}
	if x >= 11 || y >= 11 {
		return
	}
	m.cells[y][x] = MapTile{r, fg, bg}
}
