package messages

import "fmt"

type MessageSound struct {
	X         int8
	Y         int8
	Direction uint8 // 0-8
	Volume    uint8 // 0-100
	Type      uint8 //???
	Action    string
	Name      string
}

func (m *MessageSound) UnmarshalBinary(data []byte) error {
	m.X = int8(data[0])
	m.Y = int8(data[1])
	m.Direction = data[2]
	m.Volume = data[3]
	m.Type = data[4]
	offset := 5
	m.Action, offset = readLengthPrefixedString(data, offset)
	m.Name, _ = readLengthPrefixedString(data, offset)
	return nil
}

func (m MessageSound) Kind() string {
	return "sound2"
}

func (m MessageSound) Value() string {
	return fmt.Sprintf("%d %d %d %d %d %s %s", m.X, m.Y, m.Direction, m.Volume, m.Type, m.Action, m.Name)
}

func (m MessageSound) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageSound{})
}
