package messages

import (
	"strconv"
)

type MessageAccountLogin struct {
	Account  string
	Password string
}

func (m *MessageAccountLogin) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageAccountLogin) Kind() string {
	return "accountlogin"
}

func (m MessageAccountLogin) Value() string {
	return m.Account + " " + m.Password
}

func (m MessageAccountLogin) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(LengthPrefixedString(m.Account))...)
	result = append(result, []byte(LengthPrefixedString(m.Password))...)
	return result
}

type Character struct {
	Name    string
	Level   int
	Class   string
	Race    string
	Face    string
	Party   string
	Map     string
	FaceNum int
}

type MessageAccountPlayers struct {
	Characters []Character
}

const (
	ACL_BLANK int = iota
	ACL_NAME
	ACL_CLASS
	ACL_RACE
	ACL_LEVEL
	ACL_FACE
	ACL_PARTY
	ACL_MAP
	ACL_FACE_NUM
)

func (m *MessageAccountPlayers) UnmarshalBinary(data []byte) error {
	offset := 0
	count := int(data[offset])
	if count == 0 {
		m.Characters = make([]Character, 0)
		return nil
	}
	offset++
	m.Characters = make([]Character, count)

	for char := 0; char < count-1; {
		fieldLen := int(data[offset])
		offset++
		if fieldLen == 0 {
			char++
			continue
		}

		switch int(data[offset]) {
		case ACL_NAME:
			m.Characters[char].Name = string(data[offset+1 : offset+fieldLen])
		case ACL_LEVEL:
			m.Characters[char].Level = int(data[offset+1])<<8 + int(data[offset+2])
		case ACL_CLASS:
			m.Characters[char].Class = string(data[offset+1 : offset+fieldLen])
		case ACL_FACE:
			m.Characters[char].Face = string(data[offset+1 : offset+fieldLen])
		case ACL_PARTY:
			m.Characters[char].Party = string(data[offset+1 : offset+fieldLen])
		case ACL_RACE:
			m.Characters[char].Race = string(data[offset+1 : offset+fieldLen])
		case ACL_MAP:
			m.Characters[char].Map = string(data[offset+1 : offset+fieldLen])
		case ACL_FACE_NUM:
			m.Characters[char].FaceNum = int(data[offset+1])<<8 + int(data[offset+2])
		}

		offset += fieldLen
	}

	return nil
}

func (m MessageAccountPlayers) Kind() string {
	return "accountplayers"
}

func (m MessageAccountPlayers) Value() string {
	var result string
	result += "count: " + strconv.Itoa(len(m.Characters)) + " "
	for _, c := range m.Characters {
		if c.Name == "" {
			continue
		}
		result += strconv.Itoa(c.Level) + " " + c.Name + " " + c.Class + " @ " + c.Map + " | " + c.Face + "(" + strconv.Itoa(c.FaceNum) + ") " + c.Party + " "
	}
	return result
}

func (m MessageAccountPlayers) Bytes() []byte {
	var result []byte

	result = append(result, []byte(m.Kind())...)

	result = append(result, byte(len(m.Characters)))

	for _, c := range m.Characters {
		result = append(result, []byte(c.Name)...)
		result = append(result, byte(c.Level))
	}

	return result
}

type MessageAccountPlay struct {
	Character string
}

func (m *MessageAccountPlay) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageAccountPlay) Kind() string {
	return "accountplay"
}

func (m MessageAccountPlay) Value() string {
	return m.Character
}

func (m MessageAccountPlay) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(m.Character)...)
	return result
}

type MessageAccountAddmeSuccess struct{}

func (m *MessageAccountAddmeSuccess) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageAccountAddmeSuccess) Kind() string {
	return "addme_success"
}

func (m MessageAccountAddmeSuccess) Value() string {
	return ""
}

func (m MessageAccountAddmeSuccess) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageAccountLogin{})
	gMessages = append(gMessages, &MessageAccountPlayers{})
	gMessages = append(gMessages, &MessageAccountPlay{})
	gMessages = append(gMessages, &MessageAccountAddmeSuccess{})
}
