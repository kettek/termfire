package messages

type SpellUsage uint8

const (
	SpellNoArgument             SpellUsage = 0
	SpellNeedsOtherSpell        SpellUsage = 1
	SpellUsesFreeformString     SpellUsage = 2
	SpellRequiresFreeformString SpellUsage = 3
)

type Spell struct {
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

type MessageAddSpell struct {
	Spells []Spell
}

func (m MessageAddSpell) Bytes() []byte {
	return nil
}

func (m *MessageAddSpell) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Spells = make([]Spell, 0)
	for offset < len(data) {
		spell := Spell{}
		spell.Tag = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		spell.Level = int16(data[offset])<<8 | int16(data[offset+1])
		offset += 2
		spell.CastingTime = int16(data[offset])<<8 | int16(data[offset+1])
		offset += 2
		spell.Mana = int16(data[offset])<<8 | int16(data[offset+1])
		offset += 2
		spell.Grace = int16(data[offset])<<8 | int16(data[offset+1])
		offset += 2
		spell.Damage = int16(data[offset])<<8 | int16(data[offset+1])
		offset += 2
		spell.Skill = uint8(data[offset])
		offset++
		spell.Path = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
		offset += 4
		spell.Face = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		name, length := readLengthPrefixedString(data, offset)
		spell.Name = name
		offset = length
		description, length := readLengthPrefixedString2(data, offset)
		spell.Description = description
		offset = length
		// FIXME: We need unmarshalling to have some setup! For now we're assuming spellmon 2
		spell.Usage = SpellUsage(data[offset])
		offset++
		requirements, length := readLengthPrefixedString(data, offset)
		spell.Requirements = requirements
		offset = length
		m.Spells = append(m.Spells, spell)
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
	if m.Flags.Mana() {
		m.Fields = append(m.Fields, MessageUpdateSpellMana(int16(data[offset])<<8|int16(data[offset+1])))
		offset += 2
	}
	if m.Flags.Grace() {
		m.Fields = append(m.Fields, MessageUpdateSpellGrace(int16(data[offset])<<8|int16(data[offset+1])))
		offset += 2
	}
	if m.Flags.Damage() {
		m.Fields = append(m.Fields, MessageUpdateSpellDamage(int16(data[offset])<<8|int16(data[offset+1])))
		offset += 2
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
