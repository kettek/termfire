package messages

import (
	"strconv"
)

type MessageColor int

const (
	MessageColorBlack      MessageColor = 0
	MessageColorWhite                   = 1
	MessageColorNavy                    = 2
	MessageColorRed                     = 3
	MessageColorOrange                  = 4
	MessageColorBlue                    = 5
	MessageColorDarkOrange              = 6
	MessageColorGreen                   = 7
	MessageColorLightGreen              = 8
	MessageColorGrey                    = 9
	MessageColorBrown                   = 10
	MessageColorGold                    = 11
	MessageColorTan                     = 12
	MessageColorMax                     = 12
	MessageColorMask                    = 0xff
	// Control flags
	MessageUnique      = 0x100
	MessageAll         = 0x200
	MessageAllDMs      = 0x400
	MessageNoTranslate = 0x800
	MessageDelayed     = 0x1000
)

type MessageType int

const (
	_ MessageType = iota
	MessageTypeBook
	MessageTypeCard
	MessageTypePaper
	MessageTypeSign
	MessageTypeMonument
	MessageTypeDialog
	MessageTypeMOTD
	MessageTypeAdmin
	MessageTypeShop
	MessageTypeCommand
	MessageTypeAttribute
	MessageTypeSkill
	MessageTypeApply
	MessageTypeAttack
	MessageTypeCommunication
	MessageTypeSpell
	MessageTypeItem
	MessageTypeMisc
	MessageTypeVictim
	MessageTypeClient
	MessageTypeLast = MessageTypeClient
)

type SubMessageType int

var SubMessageTypeNone SubMessageType = 0

const (
	_ SubMessageType = iota
	SubMessageTypeBookClasp1
	SubMessageTypeBookClasp2
	SubMessageTypeBookElegant1
	SubMessageTypeBookElegant2
	SubMessageTypeBookQuarto1
	SubMessageTypeBookQuarto2
	SubMessageTypeBookSpellEvoker
	SubMessageTypeBookSpellPrayer
	SubMessageTypeBookSpellPyro
	SubMessageTypeBookSpellSorcerer
	SubMessageTypeBookSpellSummoner
)

const (
	_ SubMessageType = iota
	SubMessageTypeCardSimple1
	SubMessageTypeCardSimple2
	SubMessageTypeCardSimple3
	SubMessageTypeCardElegant1
	SubMessageTypeCardElegant2
	SubMessageTypeCardElegant3
	SubMessageTypeCardStrange1
	SubMessageTypeCardStrange2
	SubMessageTypeCardStrange3
	SubMessageTypeCardMoney1
	SubMessageTypeCardMoney2
	SubMessageTypeCardMoney3
)

const (
	_ SubMessageType = iota
	SubMessageTypePaperNote1
	SubMessageTypePaperNote2
	SubMessageTypePaperNote3
	SubMessageTypePaperLetterOld1
	SubMessageTypePaperLetterOld2
	SubMessageTypePaperLetterNew1
	SubMessageTypePaperLetterNew2
	SubMessageTypePaperEnvelope1
	SubMessageTypePaperEnvelope2
	SubMessageTypePaperScrollOld1
	SubMessageTypePaperScrollOld2
	SubMessageTypePaperScrollNew1
	SubMessageTypePaperScrollNew2
	SubMessageTypePaperMagic
)

const (
	_ SubMessageType = iota
	SubMessageTypeSignBasic
	SubMessageTypeSignDirLeft
	SubMessageTypeSignDirRight
	SubMessageTypeSignDirBoth
	SubMessageTypeSignMagicMouth
)

const (
	_ SubMessageType = iota
	SubMessageTypeMonumentStone1
	SubMessageTypeMonumentStone2
	SubMessageTypeMonumentStone3
	SubMessageTypeMonumentStatue1
	SubMessageTypeMonumentStatue2
	SubMessageTypeMonumentStatue3
	SubMessageTypeMonumentGravestone1
	SubMessageTypeMonumentGravestone2
	SubMessageTypeMonumentGravestone3
	SubMessageTypeMonumentWall1
	SubMessageTypeMonumentWall2
	SubMessageTypeMonumentWall3
)

const (
	_ SubMessageType = iota
	SubMessageTypeDialogNpc
	SubMessageTypeDialogAltar
	SubMessageTypeDialogMagicEar
)

const (
	_ SubMessageType = iota
	SubMessageTypeAdminRules
	SubMessageTypeAdminNews
	SubMessageTypeAdminPlayer
	SubMessageTypeAdminDM
	SubMessageTypeAdminHiscore
	SubMessageTypeAdminLoadSave
	SubMessageTypeAdminLogin
	SubMessageTypeAdminVersion
	SubMessageTypeAdminError
)

const (
	_ SubMessageType = iota
	SubMessageTypeShopListing
	SubMessageTypeShopPayment
	SubMessageTypeShopSell
	SubMessageTypeShopMisc
)

const (
	_ SubMessageType = iota
	SubMessageTypeCommandWho
	SubMessageTypeCommandMaps
	SubMessageTypeCommandBody
	SubMessageTypeCommandMalloc
	SubMessageTypeCommandWeather
	SubMessageTypeCommandStatistics
	SubMessageTypeCommandConfig
	SubMessageTypeCommandInfo
	SubMessageTypeCommandQuests
	SubMessageTypeCommandDebug
	SubMessageTypeCommandError
	SubMessageTypeCommandSuccess
	SubMessageTypeCommandFailure
	SubMessageTypeCommandExamine
	SubMessageTypeCommandInventory
	SubMessageTypeCommandHelp
	SubMessageTypeCommandDM
	SubMessageTypeCommandNewPlayer
)

const (
	_ SubMessageType = iota
	SubMessageTypeAttributeAttacktypeGain
	SubMessageTypeAttributeAttacktypeLoss
	SubMessageTypeAttributeProtectionGain
	SubMessageTypeAttributeProtectionLoss
	SubMessageTypeAttributeMove
	SubMessageTypeAttributeRace
	SubMessageTypeAttributeBadEffectStart
	SubMessageTypeAttributeBadEffectEnd
	SubMessageTypeAttributeStatGain
	SubMessageTypeAttributeStatLoss
	SubMessageTypeAttributeLevelGain
	SubMessageTypeAttributeLevelLoss
	SubMessageTypeAttributeGoodEffectStart
	SubMessageTypeAttributeGoodEffectEnd
	SubMessageTypeAttributeGod
)

const (
	_ SubMessageType = iota
	SubMessageTypeSkillMissing
	SubMessageTypeSkillError
	SubMessageTypeSkillSuccess
	SubMessageTypeSkillFailure
	SubMessageTypeSkillPray
	SubMessageTypeSkillList
)

const (
	_ SubMessageType = iota
	SubMessageTypeApplyError
	SubMessageTypeApplyUnapply
	SubMessageTypeApplySuccess
	SubMessageTypeApplyFailure
	SubMessageTypeApplyCursed
	SubMessageTypeApplyTrap
	SubMessageTypeApplyBadBody
	SubMessageTypeApplyProhibition
	SubMessageTypeApplyBuild
)

const (
	_ SubMessageType = iota
	SubMessageTypeAttackDidHit
	SubMessageTypeAttackPetHit
	SubMessageTypeAttackFumble
	SubMessageTypeAttackDidKill
	SubMessageTypeAttackPetDied
	SubMessageTypeAttackNoKey
	SubMessageTypeAttackNoAttack
	SubMessageTypeAttackPushed
	SubMessageTypeAttackMissed
)

const (
	_ SubMessageType = iota
	SubMessageTypeCommunicationRandom
	SubMessageTypeCommunicationSay
	SubMessageTypeCommunicationMe
	SubMessageTypeCommunicationTell
	SubMessageTypeCommunicationEmote
	SubMessageTypeCommunicationParty
	SubMessageTypeCommunicationShout
	SubMessageTypeCommunicationChat
)

const (
	_ SubMessageType = iota
	SubMessageTypeSpellHeal
	SubMessageTypeSpellPet
	SubMessageTypeSpellFailure
	SubMessageTypeSpellEnd
	SubMessageTypeSpellSuccess
	SubMessageTypeSpellError
	SubMessageTypeSpellPerceiveSelf
	SubMessageTypeSpellTarget
	SubMessageTypeSpellInfo
)

const (
	_ SubMessageType = iota
	SubMessageTypeItemRemove
	SubMessageTypeItemAdd
	SubMessageTypeItemChange
	SubMessageTypeItemInfo
)

const (
	_ SubMessageType = iota
	SubMessageTypeVictimSwamp
	SubMessageTypeVictimWasHit
	SubMessageTypeVictimSteal
	SubMessageTypeVictimDied
	SubMessageTypeVictimWasPushed
)

const (
	_ SubMessageType = iota
	SubMessageTypeClientConfig
	SubMessageTypeClientServer
	SubMessageTypeClientCommand
	SubMessageTypeClientQuery
	SubMessageTypeClientDebug
	SubMessageTypeClientNotice
	SubMessageTypeClientMetaserver
	SubMessageTypeClientScript
	SubMessageTypeClientError
)

type MessageDrawExtInfo struct {
	Color   MessageColor
	Type    MessageType
	Subtype SubMessageType
	Message string
}

func (m *MessageDrawExtInfo) UnmarshalBinary(data []byte) error {
	step := 0
	start := 0
	for i := 0; i < len(data); i++ {
		if data[i] == ' ' {
			switch step {
			case 0: // MessageColor
				c, _ := strconv.Atoi(string(data[start:i]))
				m.Color = MessageColor(c)
			case 1: // Type
				t, _ := strconv.Atoi(string(data[start:i]))
				m.Type = MessageType(t)
			case 2: // SubType
				s, _ := strconv.Atoi(string(data[start:i]))
				m.Subtype = SubMessageType(s)
			default: // Message (rest of data)
				m.Message = string(data[start:])
				i = len(data)
				break
			}

			step++
			start = i + 1
			continue
		}
	}
	return nil
}

func (m MessageDrawExtInfo) Kind() string {
	return "drawextinfo"
}

func (m MessageDrawExtInfo) Value() string {
	return strconv.Itoa(int(m.Color)) + " " + strconv.Itoa(int(m.Type)) + " " + strconv.Itoa(int(m.Subtype)) + " " + m.Message
}

func (m MessageDrawExtInfo) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageDrawExtInfo{})
}
