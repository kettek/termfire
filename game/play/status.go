package play

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/kettek/termfire/messages"
	"github.com/rivo/tview"
)

type Status struct {
	View *tview.Box
	//
	title string
	//
	hp       int
	maxHP    int
	sp       int
	maxSP    int
	grace    int
	maxGrace int
	//
	strength     int
	dexterity    int
	constitution int
	intelligence int
	wisdom       int
	power        int
	charisma     int
	//
	baseStrength     int
	baseDexterity    int
	baseConstitution int
	baseIntelligence int
	baseWisdom       int
	basePower        int
	baseCharisma     int
}

func (s *Status) Init() {
	s.View = tview.NewBox()

	s.View.SetDrawFunc(func(screen tcell.Screen, x, y, width, height int) (int, int, int, int) {
		tview.Print(screen, s.title, x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorTan])
		y++
		// Status
		tview.Print(screen, fmt.Sprintf("HP: %d/%d", s.hp, s.maxHP), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("SP: %d/%d", s.sp, s.maxSP), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("Grace: %d/%d", s.grace, s.maxGrace), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		// Stats
		tview.Print(screen, fmt.Sprintf("Str: %d (%d)", s.strength, s.baseStrength), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("Dex: %d (%d)", s.dexterity, s.baseDexterity), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("Con: %d (%d)", s.constitution, s.baseConstitution), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("Int: %d (%d)", s.intelligence, s.baseIntelligence), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("Wis: %d (%d)", s.wisdom, s.baseWisdom), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("Pow: %d (%d)", s.power, s.basePower), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++
		tview.Print(screen, fmt.Sprintf("Cha: %d (%d)", s.charisma, s.baseCharisma), x, y, width, tview.AlignLeft, CF2W3CColor[messages.MessageColorGreen])
		y++

		return x, y, width, height
	})
}

func (s *Status) Update(msg *messages.MessageStats) {
	for _, stat := range msg.Stats {
		switch stat := stat.(type) {
		case *messages.MessageStatHP:
			s.hp = int(*stat)
		case *messages.MessageStatMaxHP:
			s.maxHP = int(*stat)
		case *messages.MessageStatSP:
			s.sp = int(*stat)
		case *messages.MessageStatMaxSP:
			s.maxSP = int(*stat)
		case *messages.MessageStatGrace:
			s.grace = int(*stat)
		case *messages.MessageStatMaxGrace:
			s.maxGrace = int(*stat)
		case *messages.MessageStatStr:
			s.strength = int(*stat)
		case *messages.MessageStatCon:
			s.constitution = int(*stat)
		case *messages.MessageStatDex:
			s.dexterity = int(*stat)
		case *messages.MessageStatInt:
			s.intelligence = int(*stat)
		case *messages.MessageStatWis:
			s.wisdom = int(*stat)
		case *messages.MessageStatPow:
			s.power = int(*stat)
		case *messages.MessageStatCha:
			s.charisma = int(*stat)
		case *messages.MessageStatBaseStr:
			s.baseStrength = int(*stat)
		case *messages.MessageStatBaseCon:
			s.baseConstitution = int(*stat)
		case *messages.MessageStatBaseDex:
			s.baseDexterity = int(*stat)
		case *messages.MessageStatBaseInt:
			s.baseIntelligence = int(*stat)
		case *messages.MessageStatBaseWis:
			s.baseWisdom = int(*stat)
		case *messages.MessageStatBasePow:
			s.basePower = int(*stat)
		case *messages.MessageStatBaseCha:
			s.baseCharisma = int(*stat)
		}
	}
}

func (s *Status) SetTitle(t string) {
	s.title = t
}
