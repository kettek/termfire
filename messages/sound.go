package messages

import "fmt"

type MessageSound struct {
	X         int8
	Y         int8
	Direction int8 // 0-8
	Volume    int8 // 0-100
	Type      int8 //???
	Action    string
	Name      string
}

func (m *MessageSound) UnmarshalBinary(data []byte) error {
	m.X = int8(data[0])
	m.Y = int8(data[1])
	m.Direction = int8(data[2])
	m.Volume = int8(data[3])
	m.Type = int8(data[4])
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
