package messages

import (
	"fmt"
	"strconv"
	"strings"
)

type ItemFlags int32

func (f ItemFlags) Applied() bool {
	return int32(f)&0x000f != 0
}

func (f ItemFlags) Unidentified() bool {
	return int32(f)&0x0010 != 0
}

func (f ItemFlags) Unpaid() bool {
	return int32(f)&0x0200 != 0
}

func (f ItemFlags) Magic() bool {
	return int32(f)&0x0400 != 0
}

func (f ItemFlags) Cursed() bool {
	return int32(f)&0x0800 != 0
}

func (f ItemFlags) Damned() bool {
	return int32(f)&0x1000 != 0
}

func (f ItemFlags) Open() bool {
	return int32(f)&0x2000 != 0
}

func (f ItemFlags) NoPick() bool {
	return int32(f)&0x4000 != 0
}

func (f ItemFlags) Locked() bool {
	return int32(f)&0x8000 != 0
}

func (f ItemFlags) Blessed() bool {
	return int32(f)&0x0100 != 0
}

func (f ItemFlags) Read() bool {
	return int32(f)&0x0020 != 0
}

type ItemType int16

func (t ItemType) IsSpecial() bool {
	return t >= 1 && t <= 49
}

func (t ItemType) IsMeleeWeapon() bool {
	return t >= 100 && t <= 149
}

func (t ItemType) IsRangedWeapon() bool {
	return t >= 150 && t <= 199
}

// IsAmmo returns if the type is a ammo. Note that IsRangedWeapon contains the ammo type.
func (t ItemType) IsAmmo() bool {
	return t == 159 || t == 165 || t == 170 // arrows, bolts, and bombs.
}

func (t ItemType) IsArmor() bool {
	return t >= 250 && t <= 399
}

func (t ItemType) IsBodyArmor() bool {
	return t >= 250 && t <= 257
}

func (t ItemType) IsShield() bool {
	return t >= 260 && t <= 269
}

func (t ItemType) IsHeadwear() bool {
	return t >= 270 && t <= 279
}

func (t ItemType) IsCloak() bool {
	return t == 280 || t == 281
}

func (t ItemType) IsBoots() bool {
	return t == 290 || t == 291
}

func (t ItemType) IsGloves() bool {
	return t == 300 || t == 301 || t == 305
}

func (t ItemType) IsBracers() bool {
	return t == 310 || t == 311
}

func (t ItemType) IsGirdle() bool {
	return t == 321
}

func (t ItemType) IsAmulet() bool {
	return t == 381
}

func (t ItemType) IsRing() bool {
	return t == 390 || t == 391
}

func (t ItemType) IsSkillObject() bool {
	return t >= 450 && t <= 459
}

func (t ItemType) IsFoodOrAlchemy() bool {
	return t >= 600 && t <= 649
}

func (t ItemType) IsFood() bool {
	return t == 601
}

func (t ItemType) IsDrink() bool {
	return t == 611
}

func (t ItemType) IsFlesh() bool {
	return t >= 620 && t <= 627
}

func (t ItemType) IsAlchemical() bool {
	return t >= 628 && t <= 645
}

func (t ItemType) IsSpellCastingConsumable() bool {
	return t >= 650 && t <= 699
}

func (t ItemType) IsPotion() bool {
	return t == 651
}

func (t ItemType) IsBalmOrDust() bool {
	return t == 652
}

func (t ItemType) IsFigurine() bool {
	return t == 653
}

func (t ItemType) IsScroll() bool {
	return t == 661
}

func (t ItemType) IsSpellCastingItem() bool {
	return t >= 700 && t <= 749
}

func (t ItemType) IsRod() bool {
	return t == 701
}

func (t ItemType) IsWand() bool {
	return t == 711
}

func (t ItemType) IsStaff() bool {
	return t == 712
}

func (t ItemType) IsHorn() bool {
	return t == 721
}

func (t ItemType) IsKey() bool {
	return t >= 800 && t <= 849
}

func (t ItemType) IsReadable() bool {
	return t >= 1000 && t <= 1049
}

func (t ItemType) IsLightSource() bool {
	return t >= 1100 && t <= 1149
}

func (t ItemType) IsValuables() bool {
	return t >= 2000 && t <= 2049
}

func (t ItemType) IsMisc() bool {
	return t >= 8000 && t <= 8999
}

type ItemObject struct {
	Tag         int32
	Flags       ItemFlags
	Weight      int32
	TotalWeight int32
	Face        int32
	Name        string
	PluralName  string
	Anim        int16
	AnimSpeed   int8
	Nrof        int32
	Type        ItemType
}

func (o ItemObject) GetName() string {
	if o.Nrof > 1 {
		return fmt.Sprintf("%d %s", o.Nrof, o.PluralName)
	}
	return o.Name
}

type MessageItem2 struct {
	Location int32
	Objects  []ItemObject
}

func (m *MessageItem2) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Location = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
	offset += 4
	m.Objects = make([]ItemObject, 0)
	for offset < len(data) {
		var obj ItemObject
		obj.Tag = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		obj.Flags = ItemFlags(int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3]))
		offset += 4
		obj.Weight = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		obj.Face = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		obj.Name, offset = readLengthPrefixedString(data, offset)
		{ // SC 1024 support
			parts := strings.Split(obj.Name, "\x00")
			if len(parts) > 1 {
				obj.Name = parts[0]
				obj.PluralName = parts[1]
			}
		}
		obj.Anim = int16(data[offset])<<8 | int16(data[offset+1])
		offset += 2
		obj.AnimSpeed = int8(data[offset])
		offset++
		obj.Nrof = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		// I don't know why, but Nrof for some single-item objects is 0, so we set it to 1 if 0 here...
		if obj.Nrof == 0 {
			obj.Nrof = 1
		}
		{
			obj.TotalWeight = obj.Weight * obj.Nrof
		}
		offset += 4
		obj.Type = ItemType(int16(data[offset])<<8 | int16(data[offset+1]))
		offset += 2
		m.Objects = append(m.Objects, obj)
	}
	return nil
}

func (m MessageItem2) Kind() string {
	return "item2"
}

func (m MessageItem2) Value() string {
	var result string
	result += "location: " + strconv.Itoa(int(m.Location)) + "\n"
	for _, o := range m.Objects {
		result += strconv.Itoa(int(o.Nrof)) + " " + o.Name + "/" + o.PluralName + " " + strconv.Itoa(int(o.Flags)) + " " + strconv.Itoa(int(o.Weight)) + " " + strconv.Itoa(int(o.TotalWeight)) + " " + strconv.Itoa(int(o.Face)) + " " + strconv.Itoa(int(o.Anim)) + " " + strconv.Itoa(int(o.AnimSpeed)) + " " + strconv.Itoa(int(o.Type)) + "\n"
	}
	return result
}

func (m MessageItem2) Bytes() []byte {
	return nil
}

type MessageUpdateItem struct {
	Tag    int32
	Flags  int8
	Values []any // TODO
}

func (m *MessageUpdateItem) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageUpdateItem) Kind() string {
	return "upditem"
}

func (m MessageUpdateItem) Value() string {
	return ""
}

func (m MessageUpdateItem) Bytes() []byte {
	return nil
}

type MessageDeleteItem struct {
	Tags []int32
}

func (m *MessageDeleteItem) UnmarshalBinary(data []byte) error {
	m.Tags = nil
	for i := 0; i < len(data); i += 4 {
		m.Tags = append(m.Tags, int32(data[i])<<24|int32(data[i+1])<<16|int32(data[i+2])<<8|int32(data[i+3]))
	}
	return nil
}

func (m MessageDeleteItem) Kind() string {
	return "delitem"
}

func (m MessageDeleteItem) Value() string {
	return fmt.Sprintf("%+v", m.Tags)
}

func (m MessageDeleteItem) Bytes() []byte {
	return nil
}

type MessageDeleteInventory struct {
	Tag int32
}

func (m *MessageDeleteInventory) UnmarshalBinary(data []byte) error {
	v, _ := strconv.Atoi(string(data))
	m.Tag = int32(v)
	return nil
}

func (m MessageDeleteInventory) Kind() string {
	return "delinv"
}

func (m MessageDeleteInventory) Value() string {
	return fmt.Sprintf("%d", m.Tag)
}

func (m MessageDeleteInventory) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageItem2{})
	gMessages = append(gMessages, &MessageUpdateItem{})
	gMessages = append(gMessages, &MessageDeleteItem{})
	gMessages = append(gMessages, &MessageDeleteInventory{})
}
