package play

import (
	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/messages"
)

// CF2W3CColor is das colours from cf to w3c
var CF2W3CColor = map[messages.MessageColor]tcell.Color{
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
	messages.MessageColorAltBlack:   tcell.ColorDarkGray,
}
