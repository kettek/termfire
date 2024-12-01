package messages

var FloatMulti float32 = 10000

var gMessageStats []MessageStat

type MessageStat interface {
	UnmarshalBinary([]byte) (int, error)
	Matches(byte) bool
}

type MessageStatHP int16

func (m *MessageStatHP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatHP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatHP) Matches(id byte) bool {
	return id == 1
}

type MessageStatMaxHP int16

func (m *MessageStatMaxHP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatMaxHP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatMaxHP) Matches(id byte) bool {
	return id == 2
}

type MessageStatMaxSP int16

func (m *MessageStatMaxSP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatMaxSP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatMaxSP) Matches(id byte) bool {
	return id == 3
}

type MessageStatSP int16

func (m *MessageStatSP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatSP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatSP) Matches(id byte) bool {
	return id == 4
}

type MessageStatStr int16

func (m *MessageStatStr) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatStr)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatStr) Matches(id byte) bool {
	return id == 5
}

type MessageStatInt int16

func (m *MessageStatInt) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatInt)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatInt) Matches(id byte) bool {
	return id == 6
}

type MessageStatWis int16

func (m *MessageStatWis) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatWis)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatWis) Matches(id byte) bool {
	return id == 7
}

type MessageStatDex int16

func (m *MessageStatDex) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatDex)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatDex) Matches(id byte) bool {
	return id == 8
}

type MessageStatCon int16

func (m *MessageStatCon) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatCon)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatCon) Matches(id byte) bool {
	return id == 9
}

type MessageStatCha int16

func (m *MessageStatCha) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatCha)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatCha) Matches(id byte) bool {
	return id == 10
}

type MessageStatLevel int16

func (m *MessageStatLevel) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatLevel)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatLevel) Matches(id byte) bool {
	return id == 12
}

type MessageStatWC int16

func (m *MessageStatWC) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatWC)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatWC) Matches(id byte) bool {
	return id == 13
}

type MessageStatAC int16

func (m *MessageStatAC) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAC)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAC) Matches(id byte) bool {
	return id == 14
}

type MessageStatDam int16

func (m *MessageStatDam) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatDam)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatDam) Matches(id byte) bool {
	return id == 15
}

type MessageStatArmour int16

func (m *MessageStatArmour) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatArmour)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatArmour) Matches(id byte) bool {
	return id == 16
}

type MessageStatSpeed float32

func (m *MessageStatSpeed) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatSpeed)(float32(int32(data[0])<<24|int32(data[1])<<16|int32(data[2])<<8|int32(data[1])) / FloatMulti)
	return 4, nil
}

func (m MessageStatSpeed) Matches(id byte) bool {
	return id == 17
}

type MessageStatFood int16

func (m *MessageStatFood) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatFood)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatFood) Matches(id byte) bool {
	return id == 18
}

type MessageStatWeaponSpeed float32

func (m *MessageStatWeaponSpeed) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatWeaponSpeed)(float32(int32(data[0])<<24|int32(data[1])<<16|int32(data[2])<<8|int32(data[1])) / FloatMulti)
	return 4, nil
}

func (m MessageStatWeaponSpeed) Matches(id byte) bool {
	return id == 19
}

type MessageStatRange string

func (m *MessageStatRange) UnmarshalBinary(data []byte) (int, error) {
	str, len := readLengthPrefixedString(data, 0)
	*m = MessageStatRange(str)
	return len, nil
}

func (m MessageStatRange) Matches(id byte) bool {
	return id == 20
}

type MessageStatTitle string

func (m *MessageStatTitle) UnmarshalBinary(data []byte) (int, error) {
	str, len := readLengthPrefixedString(data, 0)
	*m = MessageStatTitle(str)
	return len, nil
}

func (m MessageStatTitle) Matches(id byte) bool {
	return id == 21
}

type MessageStatPow int16

func (m *MessageStatPow) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatPow)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatPow) Matches(id byte) bool {
	return id == 22
}

type MessageStatGrace int16

func (m *MessageStatGrace) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatGrace)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatGrace) Matches(id byte) bool {
	return id == 23
}

type MessageStatMaxGrace int16

func (m *MessageStatMaxGrace) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatMaxGrace)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatMaxGrace) Matches(id byte) bool {
	return id == 24
}

type MessageStatFlags int16

func (m *MessageStatFlags) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatFlags)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatFlags) Matches(id byte) bool {
	return id == 25
}

type MessageStatWeightLimit int32

func (m *MessageStatWeightLimit) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatWeightLimit)(int32(data[0])<<24 | int32(data[1])<<16 | int32(data[2])<<8 | int32(data[3]))
	return 4, nil
}

func (m MessageStatWeightLimit) Matches(id byte) bool {
	return id == 26
}

type MessageStatExp64 int64

func (m *MessageStatExp64) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatExp64)(int64(data[0])<<54 | int64(data[1])<<48 | int64(data[2])<<40 | int64(data[3])<<32 | int64(data[4])<<24 | int64(data[5])<<16 | int64(data[6])<<8 | int64(data[7]))

	return 8, nil
}

func (m MessageStatExp64) Matches(id byte) bool {
	return id == 28
}

type MessageStatSpellAttune int16

func (m *MessageStatSpellAttune) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatSpellAttune)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatSpellAttune) Matches(id byte) bool {
	return id == 29
}

type MessageStatSpellRepel int16

func (m *MessageStatSpellRepel) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatSpellRepel)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatSpellRepel) Matches(id byte) bool {
	return id == 30
}

type MessageStatSpellDeny int16

func (m *MessageStatSpellDeny) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatSpellDeny)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatSpellDeny) Matches(id byte) bool {
	return id == 31
}

func init() {
	statHP := MessageStatHP(0)
	gMessageStats = append(gMessageStats, &statHP)
	statMaxHP := MessageStatMaxHP(0)
	gMessageStats = append(gMessageStats, &statMaxHP)
	statSP := MessageStatSP(0)
	gMessageStats = append(gMessageStats, &statSP)
	statMaxSP := MessageStatMaxSP(0)
	gMessageStats = append(gMessageStats, &statMaxSP)
	statStr := MessageStatStr(0)
	gMessageStats = append(gMessageStats, &statStr)
	statInt := MessageStatInt(0)
	gMessageStats = append(gMessageStats, &statInt)
	statWis := MessageStatWis(0)
	gMessageStats = append(gMessageStats, &statWis)
	statDex := MessageStatDex(0)
	gMessageStats = append(gMessageStats, &statDex)
	statCon := MessageStatCon(0)
	gMessageStats = append(gMessageStats, &statCon)
	statCha := MessageStatCha(0)
	gMessageStats = append(gMessageStats, &statCha)
	statLevel := MessageStatLevel(0)
	gMessageStats = append(gMessageStats, &statLevel)
	statWC := MessageStatWC(0)
	gMessageStats = append(gMessageStats, &statWC)
	statAC := MessageStatAC(0)
	gMessageStats = append(gMessageStats, &statAC)
	statDam := MessageStatDam(0)
	gMessageStats = append(gMessageStats, &statDam)
	statArmour := MessageStatArmour(0)
	gMessageStats = append(gMessageStats, &statArmour)
	statSpeed := MessageStatSpeed(0)
	gMessageStats = append(gMessageStats, &statSpeed)
	statFood := MessageStatFood(0)
	gMessageStats = append(gMessageStats, &statFood)
	statWeaponSpeed := MessageStatWeaponSpeed(0)
	gMessageStats = append(gMessageStats, &statWeaponSpeed)
	statRange := MessageStatRange("")
	gMessageStats = append(gMessageStats, &statRange)
	statTitle := MessageStatTitle("")
	gMessageStats = append(gMessageStats, &statTitle)
	statPow := MessageStatPow(0)
	gMessageStats = append(gMessageStats, &statPow)
	statGrace := MessageStatGrace(0)
	gMessageStats = append(gMessageStats, &statGrace)
	statMaxGrace := MessageStatMaxGrace(0)
	gMessageStats = append(gMessageStats, &statMaxGrace)
	statFlags := MessageStatFlags(0)
	gMessageStats = append(gMessageStats, &statFlags)
	statWeightLimit := MessageStatWeightLimit(0)
	gMessageStats = append(gMessageStats, &statWeightLimit)
	statExp64 := MessageStatExp64(0)
	gMessageStats = append(gMessageStats, &statExp64)
	statSpellAttune := MessageStatSpellAttune(0)
	gMessageStats = append(gMessageStats, &statSpellAttune)
	statSpellRepel := MessageStatSpellRepel(0)
	gMessageStats = append(gMessageStats, &statSpellRepel)
	statSpellDeny := MessageStatSpellDeny(0)
	gMessageStats = append(gMessageStats, &statSpellDeny)
}

type MessageStats struct {
	Stats []MessageStat
}

func (m *MessageStats) UnmarshalBinary(data []byte) error {
	for i := 0; i < len(data); {
		kind := data[i]
		for _, s := range gMessageStats {
			if s.Matches(kind) {
				if count, err := s.UnmarshalBinary(data[i+1:]); err != nil {
					return err
				} else {
					i += count
				}
				m.Stats = append(m.Stats, s)
				break
			}
		}
	}

	return nil
}

func (m MessageStats) Kind() string {
	return "stats"
}

func (m MessageStats) Value() string {
	return ""
}

func (m MessageStats) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageStats{})
}
