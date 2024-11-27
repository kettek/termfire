package messages

import "strconv"

type MessageFace2 struct {
	num      int16
	setnum   int8
	checksum int32
	name     string
}

func (m *MessageFace2) UnmarshalBinary(data []byte) error {
	m.num = int16(data[0])<<8 | int16(data[1])
	m.setnum = int8(data[2])
	m.checksum = int32(data[3])<<24 | int32(data[4])<<16 | int32(data[5])<<8 | int32(data[6])
	m.name = string(data[7:])
	return nil
}

func (m MessageFace2) Kind() string {
	return "face2"
}

func (m MessageFace2) Value() string {
	return strconv.Itoa(int(m.num)) + " " + strconv.Itoa(int(m.setnum)) + " " + strconv.Itoa(int(m.checksum)) + " " + m.name
}

func (m MessageFace2) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageFace2{})
}
