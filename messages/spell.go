package messages

type SpellUsage uint8

const (
	SpellNoArgument             SpellUsage = 0
	SpellNeedsOtherSpell        SpellUsage = 1
	SpellUsesFreeformString     SpellUsage = 2
	SpellRequiresFreeformString SpellUsage = 3
)

type MessageAddSpell struct {
	Tag         int32
	Level       int16
	CastingTime int16
	Mana        int16
	Grace       int16
	Damage      int16
	Skill       uint8
	Path        uint32
	Face        int32
	Name        string
	Description string
	// Only set if spellmon 2
	Usage        SpellUsage
	Requirements string
}

func (m MessageAddSpell) Bytes() []byte {
	return nil
}

func (m *MessageAddSpell) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Tag = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
	offset += 4
	m.Level = int16(data[offset])<<8 | int16(data[offset+1])
	offset += 2
	m.CastingTime = int16(data[offset])<<8 | int16(data[offset+1])
	offset += 2
	m.Mana = int16(data[offset])<<8 | int16(data[offset+1])
	offset += 2
	m.Grace = int16(data[offset])<<8 | int16(data[offset+1])
	offset += 2
	m.Damage = int16(data[offset])<<8 | int16(data[offset+1])
	offset += 2
	m.Skill = uint8(data[offset])
	offset++
	m.Path = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
	offset += 4
	m.Face = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
	offset += 4
	name, length := readLengthPrefixedString(data, offset)
	m.Name = name
	offset = length
	description, length := readLengthPrefixedString2(data, offset)
	m.Description = description
	offset = length
	if offset < len(data) {
		m.Usage = SpellUsage(data[offset])
		offset++
		requirements, length := readLengthPrefixedString(data, offset)
		m.Requirements = requirements
		offset = length
	}
	return nil
}

func (m MessageAddSpell) Kind() string {
	return "addspell"
}

func (m MessageAddSpell) Value() string {
	return ""
}

type SpellFlags uint8

func (f SpellFlags) Mana() bool {
	return f&0x01 != 0
}

func (f SpellFlags) Grace() bool {
	return f&0x02 != 0
}

func (f SpellFlags) Damage() bool {
	return f&0x04 != 0
}

type MessageUpdateSpellMana int16
type MessageUpdateSpellGrace int16
type MessageUpdateSpellDamage int16

type MessageUpdateSpell struct {
	Flags  SpellFlags
	Tag    int32
	Fields []any
}

func (m MessageUpdateSpell) Bytes() []byte {
	return nil
}

func (m *MessageUpdateSpell) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Flags = SpellFlags(data[offset])
	offset++

	m.Tag = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
	offset += 4

	m.Fields = make([]any, 0)
	for offset < len(data) {
		switch m.Flags {
		case 0x01:
			m.Fields = append(m.Fields, MessageUpdateSpellMana(int16(data[offset])<<8|int16(data[offset+1])))
			offset += 2
		case 0x02:
			m.Fields = append(m.Fields, MessageUpdateSpellGrace(int16(data[offset])<<8|int16(data[offset+1])))
			offset += 2
		case 0x04:
			m.Fields = append(m.Fields, MessageUpdateSpellDamage(int16(data[offset])<<8|int16(data[offset+1])))
			offset += 2
		}
	}
	return nil
}

func (m MessageUpdateSpell) Kind() string {
	return "updspell"
}

func (m MessageUpdateSpell) Value() string {
	return ""
}

type MessageDeleteSpell struct {
	Tag int32
}

func (m MessageDeleteSpell) Bytes() []byte {
	return nil
}

func (m *MessageDeleteSpell) UnmarshalBinary(data []byte) error {
	m.Tag = int32(data[0])<<24 | int32(data[1])<<16 | int32(data[2])<<8 | int32(data[3])
	return nil
}

func (m MessageDeleteSpell) Kind() string {
	return "delspell"
}

func (m MessageDeleteSpell) Value() string {
	return ""
}

func init() {
	gMessages = append(gMessages, &MessageAddSpell{})
	gMessages = append(gMessages, &MessageUpdateSpell{})
	gMessages = append(gMessages, &MessageDeleteSpell{})
}
