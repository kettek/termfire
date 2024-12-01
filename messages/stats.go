package messages

import (
	"fmt"
)

var FloatMulti float32 = 10000

// FIXME: Uh... should probably just use enums for all this stuff...

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

type MessageStatRaceStr int16

func (m *MessageStatRaceStr) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatRaceStr)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatRaceStr) Matches(id byte) bool {
	return id == 32
}

type MessageStatRaceInt int16

func (m *MessageStatRaceInt) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatRaceInt)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatRaceInt) Matches(id byte) bool {
	return id == 33
}

type MessageStatRaceWis int16

func (m *MessageStatRaceWis) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatRaceWis)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatRaceWis) Matches(id byte) bool {
	return id == 34
}

type MessageStatRaceDex int16

func (m *MessageStatRaceDex) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatRaceDex)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatRaceDex) Matches(id byte) bool {
	return id == 35
}

type MessageStatRaceCon int16

func (m *MessageStatRaceCon) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatRaceCon)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatRaceCon) Matches(id byte) bool {
	return id == 36
}

type MessageStatRaceCha int16

func (m *MessageStatRaceCha) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatRaceCha)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatRaceCha) Matches(id byte) bool {
	return id == 37
}

type MessageStatRacePow int16

func (m *MessageStatRacePow) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatRacePow)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatRacePow) Matches(id byte) bool {
	return id == 38
}

type MessageStatBaseStr int16

func (m *MessageStatBaseStr) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatBaseStr)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatBaseStr) Matches(id byte) bool {
	return id == 39
}

type MessageStatBaseInt int16

func (m *MessageStatBaseInt) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatBaseInt)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatBaseInt) Matches(id byte) bool {
	return id == 40
}

type MessageStatBaseWis int16

func (m *MessageStatBaseWis) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatBaseWis)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatBaseWis) Matches(id byte) bool {
	return id == 41
}

type MessageStatBaseDex int16

func (m *MessageStatBaseDex) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatBaseDex)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatBaseDex) Matches(id byte) bool {
	return id == 42
}

type MessageStatBaseCon int16

func (m *MessageStatBaseCon) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatBaseCon)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatBaseCon) Matches(id byte) bool {
	return id == 43
}

type MessageStatBaseCha int16

func (m *MessageStatBaseCha) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatBaseCha)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatBaseCha) Matches(id byte) bool {
	return id == 44
}

type MessageStatBasePow int16

func (m *MessageStatBasePow) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatBasePow)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatBasePow) Matches(id byte) bool {
	return id == 45
}

type MessageStatAppliedStr int16

func (m *MessageStatAppliedStr) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAppliedStr)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAppliedStr) Matches(id byte) bool {
	return id == 46
}

type MessageStatAppliedInt int16

func (m *MessageStatAppliedInt) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAppliedInt)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAppliedInt) Matches(id byte) bool {
	return id == 47
}

type MessageStatAppliedWis int16

func (m *MessageStatAppliedWis) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAppliedWis)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAppliedWis) Matches(id byte) bool {
	return id == 48
}

type MessageStatAppliedDex int16

func (m *MessageStatAppliedDex) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAppliedDex)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAppliedDex) Matches(id byte) bool {
	return id == 49
}

type MessageStatAppliedCon int16

func (m *MessageStatAppliedCon) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAppliedCon)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAppliedCon) Matches(id byte) bool {
	return id == 50
}

type MessageStatAppliedCha int16

func (m *MessageStatAppliedCha) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAppliedCha)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAppliedCha) Matches(id byte) bool {
	return id == 51
}

type MessageStatAppliedPow int16

func (m *MessageStatAppliedPow) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatAppliedPow)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatAppliedPow) Matches(id byte) bool {
	return id == 52
}

type MessageStatGolemHP int16

func (m *MessageStatGolemHP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatGolemHP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatGolemHP) Matches(id byte) bool {
	return id == 53
}

type MessageStatGolemMaxHP int16

func (m *MessageStatGolemMaxHP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatGolemMaxHP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatGolemMaxHP) Matches(id byte) bool {
	return id == 54
}

type MessageStatCharacterFlags int16

func (m *MessageStatCharacterFlags) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatCharacterFlags)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatCharacterFlags) Matches(id byte) bool {
	return id == 55
}

type MessageStatGodName string

func (m *MessageStatGodName) UnmarshalBinary(data []byte) (int, error) {
	str, len := readLengthPrefixedString(data, 0)
	*m = MessageStatGodName(str)
	return len, nil
}

func (m MessageStatGodName) Matches(id byte) bool {
	return id == 56
}

type MessageStatOverload float32

func (m *MessageStatOverload) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatOverload)(float32(int32(data[0])<<24|int32(data[1])<<16|int32(data[2])<<8|int32(data[1])) / FloatMulti)
	return 4, nil
}

func (m MessageStatOverload) Matches(id byte) bool {
	return id == 57
}

type MessageStatItemPower int16

func (m *MessageStatItemPower) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatItemPower)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatItemPower) Matches(id byte) bool {
	return id == 58
}

type MessageStatResistStart int16

func (m *MessageStatResistStart) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistStart)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistStart) Matches(id byte) bool {
	return id == 100
}

type MessageStatResistEnd int16

func (m *MessageStatResistEnd) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistEnd)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistEnd) Matches(id byte) bool {
	return id == 117
}

type MessageStatResistPhys int16

func (m *MessageStatResistPhys) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistPhys)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistPhys) Matches(id byte) bool {
	return id == 100
}

type MessageStatResistMag int16

func (m *MessageStatResistMag) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistMag)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistMag) Matches(id byte) bool {
	return id == 101
}

type MessageStatResistFire int16

func (m *MessageStatResistFire) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistFire)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistFire) Matches(id byte) bool {
	return id == 102
}

type MessageStatResistElec int16

func (m *MessageStatResistElec) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistElec)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistElec) Matches(id byte) bool {
	return id == 103
}

type MessageStatResistCold int16

func (m *MessageStatResistCold) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistCold)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistCold) Matches(id byte) bool {
	return id == 104
}

type MessageStatResistConf int16

func (m *MessageStatResistConf) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistConf)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistConf) Matches(id byte) bool {
	return id == 105
}

type MessageStatResistAcid int16

func (m *MessageStatResistAcid) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistAcid)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistAcid) Matches(id byte) bool {
	return id == 106
}

type MessageStatResistDrain int16

func (m *MessageStatResistDrain) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistDrain)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistDrain) Matches(id byte) bool {
	return id == 107
}

type MessageStatResistGhostHit int16

func (m *MessageStatResistGhostHit) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistGhostHit)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistGhostHit) Matches(id byte) bool {
	return id == 108
}

type MessageStatResistPoison int16

func (m *MessageStatResistPoison) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistPoison)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistPoison) Matches(id byte) bool {
	return id == 109
}

type MessageStatResistSlow int16

func (m *MessageStatResistSlow) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistSlow)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistSlow) Matches(id byte) bool {
	return id == 110
}

type MessageStatResistPara int16

func (m *MessageStatResistPara) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistPara)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistPara) Matches(id byte) bool {
	return id == 111
}

type MessageStatTurnUndead int16

func (m *MessageStatTurnUndead) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatTurnUndead)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatTurnUndead) Matches(id byte) bool {
	return id == 112
}

type MessageStatResistFear int16

func (m *MessageStatResistFear) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistFear)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistFear) Matches(id byte) bool {
	return id == 113
}

type MessageStatResistDeplete int16

func (m *MessageStatResistDeplete) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistDeplete)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistDeplete) Matches(id byte) bool {
	return id == 114
}

type MessageStatResistDeath int16

func (m *MessageStatResistDeath) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistDeath)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistDeath) Matches(id byte) bool {
	return id == 115
}

type MessageStatResistHolyWord int16

func (m *MessageStatResistHolyWord) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistHolyWord)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistHolyWord) Matches(id byte) bool {
	return id == 116
}

type MessageStatResistBlind int16

func (m *MessageStatResistBlind) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatResistBlind)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatResistBlind) Matches(id byte) bool {
	return id == 117
}

type MessageStatSkill struct {
	Level int8
	Exp   int64
}

func (m *MessageStatSkill) UnmarshalBinary(data []byte) (int, error) {
	m.Level = int8(data[0])
	m.Exp = (int64(data[1])<<54 | int64(data[2])<<48 | int64(data[3])<<40 | int64(data[4])<<32 | int64(data[5])<<24 | int64(data[6])<<16 | int64(data[7])<<8 | int64(data[8]))
	return 9, nil
}

func (m MessageStatSkill) Matches(id byte) bool {
	return id > 140
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
	statRaceStr := MessageStatRaceStr(0)
	gMessageStats = append(gMessageStats, &statRaceStr)
	statRaceInt := MessageStatRaceInt(0)
	gMessageStats = append(gMessageStats, &statRaceInt)
	statRaceWis := MessageStatRaceWis(0)
	gMessageStats = append(gMessageStats, &statRaceWis)
	statRaceDex := MessageStatRaceDex(0)
	gMessageStats = append(gMessageStats, &statRaceDex)
	statRaceCon := MessageStatRaceCon(0)
	gMessageStats = append(gMessageStats, &statRaceCon)
	statRaceCha := MessageStatRaceCha(0)
	gMessageStats = append(gMessageStats, &statRaceCha)
	statRacePow := MessageStatRacePow(0)
	gMessageStats = append(gMessageStats, &statRacePow)
	statBaseStr := MessageStatBaseStr(0)
	gMessageStats = append(gMessageStats, &statBaseStr)
	statBaseInt := MessageStatBaseInt(0)
	gMessageStats = append(gMessageStats, &statBaseInt)
	statBaseWis := MessageStatBaseWis(0)
	gMessageStats = append(gMessageStats, &statBaseWis)
	statBaseDex := MessageStatBaseDex(0)
	gMessageStats = append(gMessageStats, &statBaseDex)
	statBaseCon := MessageStatBaseCon(0)
	gMessageStats = append(gMessageStats, &statBaseCon)
	statBaseCha := MessageStatBaseCha(0)
	gMessageStats = append(gMessageStats, &statBaseCha)
	statBasePow := MessageStatBasePow(0)
	gMessageStats = append(gMessageStats, &statBasePow)
	statAppliedStr := MessageStatAppliedStr(0)
	gMessageStats = append(gMessageStats, &statAppliedStr)
	statAppliedInt := MessageStatAppliedInt(0)
	gMessageStats = append(gMessageStats, &statAppliedInt)
	statAppliedWis := MessageStatAppliedWis(0)
	gMessageStats = append(gMessageStats, &statAppliedWis)
	statAppliedDex := MessageStatAppliedDex(0)
	gMessageStats = append(gMessageStats, &statAppliedDex)
	statAppliedCon := MessageStatAppliedCon(0)
	gMessageStats = append(gMessageStats, &statAppliedCon)
	statAppliedCha := MessageStatAppliedCha(0)
	gMessageStats = append(gMessageStats, &statAppliedCha)
	statAppliedPow := MessageStatAppliedPow(0)
	gMessageStats = append(gMessageStats, &statAppliedPow)
	statGolemHP := MessageStatGolemHP(0)
	gMessageStats = append(gMessageStats, &statGolemHP)
	statGolemMaxHP := MessageStatGolemMaxHP(0)
	gMessageStats = append(gMessageStats, &statGolemMaxHP)
	statCharacterFlags := MessageStatCharacterFlags(0)
	gMessageStats = append(gMessageStats, &statCharacterFlags)
	statGodName := MessageStatGodName("")
	gMessageStats = append(gMessageStats, &statGodName)
	statOverload := MessageStatOverload(0)
	gMessageStats = append(gMessageStats, &statOverload)
	statItemPower := MessageStatItemPower(0)
	gMessageStats = append(gMessageStats, &statItemPower)
	statResistStart := MessageStatResistStart(0)
	gMessageStats = append(gMessageStats, &statResistStart)
	statResistEnd := MessageStatResistEnd(0)
	gMessageStats = append(gMessageStats, &statResistEnd)
	statResistPhys := MessageStatResistPhys(0)
	gMessageStats = append(gMessageStats, &statResistPhys)
	statResistMag := MessageStatResistMag(0)
	gMessageStats = append(gMessageStats, &statResistMag)
	statResistFire := MessageStatResistFire(0)
	gMessageStats = append(gMessageStats, &statResistFire)
	statResistElec := MessageStatResistElec(0)
	gMessageStats = append(gMessageStats, &statResistElec)
	statResistCold := MessageStatResistCold(0)
	gMessageStats = append(gMessageStats, &statResistCold)
	statResistConf := MessageStatResistConf(0)
	gMessageStats = append(gMessageStats, &statResistConf)
	statResistAcid := MessageStatResistAcid(0)
	gMessageStats = append(gMessageStats, &statResistAcid)
	statResistDrain := MessageStatResistDrain(0)
	gMessageStats = append(gMessageStats, &statResistDrain)
	statResistGhostHit := MessageStatResistGhostHit(0)
	gMessageStats = append(gMessageStats, &statResistGhostHit)
	statResistPoison := MessageStatResistPoison(0)
	gMessageStats = append(gMessageStats, &statResistPoison)
	statResistSlow := MessageStatResistSlow(0)
	gMessageStats = append(gMessageStats, &statResistSlow)
	statResistPara := MessageStatResistPara(0)
	gMessageStats = append(gMessageStats, &statResistPara)
	statTurnUndead := MessageStatTurnUndead(0)
	gMessageStats = append(gMessageStats, &statTurnUndead)
	statResistFear := MessageStatResistFear(0)
	gMessageStats = append(gMessageStats, &statResistFear)
	statResistDeplete := MessageStatResistDeplete(0)
	gMessageStats = append(gMessageStats, &statResistDeplete)
	statResistDeath := MessageStatResistDeath(0)
	gMessageStats = append(gMessageStats, &statResistDeath)
	statResistHolyWord := MessageStatResistHolyWord(0)
	gMessageStats = append(gMessageStats, &statResistHolyWord)
	statResistBlind := MessageStatResistBlind(0)
	gMessageStats = append(gMessageStats, &statResistBlind)
	statSkill := MessageStatSkill{}
	gMessageStats = append(gMessageStats, &statSkill)
}

type MessageStats struct {
	Stats []MessageStat
}

func (m *MessageStats) UnmarshalBinary(data []byte) error {
	m.Stats = make([]MessageStat, 0)
	for i := 0; i < len(data); {
		kind := data[i]
		match := false
		for _, s := range gMessageStats {
			if s.Matches(kind) {
				if count, err := s.UnmarshalBinary(data[i+1:]); err != nil {
					return err
				} else {
					i += count
				}
				match = true
				m.Stats = append(m.Stats, s)
				break
			}
		}
		if !match {
			return fmt.Errorf("Unknown stat %d", kind)
		}
		i++
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
