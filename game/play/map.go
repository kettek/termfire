package play

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/debug"
	"github.com/rivo/tview"
)

type MapRune rune

const (
	MapRuneWall     MapRune = '█'
	MapRuneWindow           = '▓'
	MapRuneStones           = '·'
	MapRuneDirt             = '░'
	MapRuneDoor             = '+'
	MapRuneGate             = '‡'
	MapRuneWater            = '~'
	MapRunePlayer           = '@'
	MapRuneCoin             = '¢'
	MapRuneBed              = '&'
	MapRuneTable            = 'T'
	MapRuneChair            = 'h'
	MapRuneScroll           = '!'
	MapRuneAmulet           = '¤'
	MapRuneTalisman         = '¥'
	MapRuneRing             = '°'
	MapRuneLever            = '/'
	MapHouse                = '▲'
	MapShop                 = '#'
	MapTower                = '#'
	MapPlant                = '♣'
	MapTree                 = '♠'
	MapJungle               = '♣'
	MapLight                = '¡'
	MapPond                 = '≈'
	MapFountain             = '⌂'
	MapSign                 = '¶'
	MapStatue               = '☻'
	MapWell                 = 'O'
	MapEmpty                = ' '
	MapClock                = '♦'
	MapHills                = '∆'
	MapMountain             = '^'
	MapBridge               = '='
	MapAltar                = '†'
	MapCorpse               = ','
	MapHole                 = '>'
	MapStairs               = '<'
)

type MapTiles []MapTile

func (t *MapTiles) SetCell(layer int, r MapRune, fg tcell.Color, bg tcell.Color) {
	if layer >= len(*t) {
		*t = append(*t, MapTile{})
	}
	for len(*t) <= layer {
		*t = append(*t, MapTile{})
	}
	(*t)[layer] = MapTile{r, fg, bg}
}

func (t *MapTiles) RemoveLayer(layer int) {
	if layer >= len(*t) {
		return
	}
	*t = append((*t)[:layer], (*t)[layer+1:]...)
}

func (t *MapTiles) GetCell(layer int) MapTile {
	if layer >= len(*t) || layer < 0 {
		return MapTile{}
	}
	return (*t)[layer]
}

func (t *MapTiles) GetTopCell() MapTile {
	// Get it if the rune is non-zero
	for i := len(*t) - 1; i >= 0; i-- {
		if (*t)[i].R != 0 {
			return (*t)[i]
		}
	}
	return t.GetCell(len(*t) - 1)
}

func (t *MapTiles) GetBottomCell() MapTile {
	return t.GetCell(0)
}

type MapTile struct {
	R MapRune
	F tcell.Color
	B tcell.Color
}

var NameToMapTile = map[string]MapTile{
	"wall":      {MapRuneWall, tcell.ColorWhite, tcell.ColorBlack},
	"window":    {MapRuneWindow, tcell.ColorWhite, tcell.ColorBlack},
	"mine":      {MapRuneWall, tcell.ColorDarkGray, tcell.ColorBlack},
	"floor":     {MapRuneDirt, tcell.ColorBlack, tcell.ColorWhite},
	"stones":    {MapRuneStones, tcell.ColorWhite, tcell.ColorGray},
	"stone":     {MapRuneStones, tcell.ColorWhite, tcell.ColorGray},
	"medston":   {MapRuneStones, tcell.ColorWhite, tcell.ColorGray},
	"pier":      {MapRuneDirt, tcell.ColorTan, tcell.ColorBrown},
	"dirt":      {MapRuneDirt, tcell.ColorBrown, tcell.ColorBlack},
	"farm":      {MapRuneDirt, tcell.ColorGold, tcell.ColorBrown},
	"grass":     {MapRuneDirt, tcell.ColorBlack, tcell.ColorGreen},
	"beach":     {MapRuneDirt, tcell.ColorTan, tcell.ColorYellow},
	"desert":    {MapRuneDirt, tcell.ColorYellow, tcell.ColorGold},
	"ground":    {MapRuneDirt, tcell.ColorBlack, tcell.ColorWhite},
	"cobble":    {MapRuneDirt, tcell.ColorBlack, tcell.ColorGray},
	"steppe":    {MapRuneDirt, tcell.ColorBlack, tcell.ColorTan},
	"door":      {MapRuneDoor, tcell.ColorWhite, tcell.ColorBlack},
	"gate":      {MapRuneGate, tcell.ColorGray, tcell.ColorBlack},
	"grate":     {MapRuneGate, tcell.ColorGray, tcell.ColorBlack},
	"water":     {MapRuneWater, tcell.ColorBlue, tcell.ColorBlack},
	"river":     {MapRuneWater, tcell.ColorBlue, tcell.ColorLightBlue},
	"sea":       {MapRuneDirt, tcell.ColorBlue, tcell.ColorLightBlue},
	"branch_":   {MapRuneWater, tcell.ColorBlue, tcell.ColorLightBlue},
	"swamp":     {MapRuneDirt, tcell.ColorLightGreen, tcell.ColorDarkGreen},
	"player":    {MapRunePlayer, tcell.ColorWhite, tcell.ColorBlack},
	"coin":      {MapRuneCoin, tcell.ColorYellow, tcell.ColorBlack},
	"amulet":    {MapRuneAmulet, tcell.ColorWhite, tcell.ColorBlack},
	"talisman":  {MapRuneTalisman, tcell.ColorWhite, tcell.ColorBlack},
	"ring":      {MapRuneRing, tcell.ColorWhite, tcell.ColorBlack},
	"bed":       {MapRuneBed, tcell.ColorRed, tcell.ColorBlack},
	"table":     {MapRuneTable, tcell.ColorBeige, tcell.ColorBlack},
	"chair":     {MapRuneChair, tcell.ColorBeige, tcell.ColorBlack},
	"scroll":    {MapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"card":      {MapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"book":      {MapRuneScroll, tcell.ColorWhite, tcell.ColorBlack},
	"lever":     {MapRuneLever, tcell.ColorGray, tcell.ColorBlack},
	"handle":    {MapRuneLever, tcell.ColorGray, tcell.ColorBlack},
	"house":     {MapHouse, tcell.ColorBlack, tcell.ColorWhite},
	"barrack":   {MapHouse, tcell.ColorBrown, tcell.ColorWhite},
	"barn":      {MapHouse, tcell.ColorBrown, tcell.ColorWhite},
	"hut":       {MapHouse, tcell.ColorBrown, tcell.ColorWhite},
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
	"zoo":       {MapHouse, tcell.ColorWhite, tcell.ColorBlack},
	"shrub":     {MapPlant, tcell.ColorGreen, tcell.ColorBlack},
	"brush":     {MapPlant, tcell.ColorDarkGreen, tcell.ColorGreen},
	"tree":      {MapTree, tcell.ColorGreen, tcell.ColorBlack},
	"jungle":    {MapJungle, tcell.ColorLightGreen, tcell.ColorDarkGreen},
	"lamp":      {MapLight, tcell.ColorYellow, tcell.ColorBlack},
	"brazier":   {MapLight, tcell.ColorRed, tcell.ColorBlack},
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
	"hills":     {MapHills, tcell.ColorBlack, tcell.ColorGreen},
	"mountain":  {MapMountain, tcell.ColorWhite, tcell.ColorGray},
	"bridge":    {MapBridge, tcell.ColorBrown, tcell.ColorBlack},
	"altar":     {MapAltar, tcell.ColorWhite, tcell.ColorBlack},
	"corpse":    {MapCorpse, tcell.ColorTan, tcell.ColorBlack},
	"hole":      {MapHole, tcell.ColorDarkGray, tcell.ColorBlack},
	"stair":     {MapStairs, tcell.ColorWhite, tcell.ColorBlack},
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
	cells      [][]MapTiles // TODO: Make resizeable.
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
			debug.Debug("resize map: ", m.viewWidth, m.viewHeight)
		}
		// offset so we render in the center.
		x += (width - m.width) / 2
		y += (height - m.height) / 2

		for my := 0; my < m.height; my++ {
			for mx := 0; mx < m.width; mx++ {
				top := m.cells[my][mx].GetTopCell()
				bot := m.cells[my][mx].GetBottomCell()

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
	m.cells = make([][]MapTiles, m.height)
	for y := 0; y < m.height; y++ {
		m.cells[y] = make([]MapTiles, m.width)
		for x := 0; x < m.width; x++ {
			m.cells[y][x] = MapTiles{}
		}
	}
}

func (m *Map) Clear() {
	for y := 0; y < len(m.cells); y++ {
		for x := 0; x < len(m.cells[y]); x++ {
			m.cells[y][x] = MapTiles{}
		}
	}
}

func (m *Map) Shift(dx, dy int) {
	newCells := make([][]MapTiles, m.height)
	for y := 0; y < m.height; y++ {
		newCells[y] = make([]MapTiles, m.width)
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

func (m *Map) SetCell(x, y int, layer int, r MapRune, fg tcell.Color, bg tcell.Color) {
	if x < 0 || y < 0 {
		return
	}
	if x >= m.width || y >= m.height {
		return
	}
	m.cells[y][x].SetCell(layer, r, fg, bg)
}

func (m *Map) ClearCell(x, y int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= m.width || y >= m.height {
		return
	}
	m.cells[y][x] = MapTiles{}
}

func (m *Map) RemoveCellLayer(x, y, layer int) {
	if x < 0 || y < 0 {
		return
	}
	if x >= m.width || y >= m.height {
		return
	}
	m.cells[y][x].RemoveLayer(layer)
}

func (m *Map) GetCell(x, y int) MapTiles {
	if x < 0 || y < 0 {
		return MapTiles{}
	}
	if x >= m.width || y >= m.height {
		return MapTiles{}
	}
	return m.cells[y][x]
}
