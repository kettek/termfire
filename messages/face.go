package messages

import "strconv"

type MessageFace2 struct {
	Num      int16
	SetNum   int8
	Checksum int32
	Name     string
}

func (m *MessageFace2) UnmarshalBinary(data []byte) error {
	m.Num = int16(data[0])<<8 | int16(data[1])
	m.SetNum = int8(data[2])
	m.Checksum = int32(data[3])<<24 | int32(data[4])<<16 | int32(data[5])<<8 | int32(data[6])
	m.Name = string(data[7:])
	return nil
}

func (m MessageFace2) Kind() string {
	return "face2"
}

func (m MessageFace2) Value() string {
	return strconv.Itoa(int(m.Num)) + " " + strconv.Itoa(int(m.SetNum)) + " " + strconv.Itoa(int(m.Checksum)) + " " + m.Name
}

func (m MessageFace2) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageFace2{})
}
