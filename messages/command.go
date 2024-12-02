package messages

type MessageCommand struct {
	Packet  uint16
	Repeat  uint32
	Command string
}

func (m *MessageCommand) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageCommand) Kind() string {
	return "ncom"
}

func (m MessageCommand) Value() string {
	return m.Command
}

func (m MessageCommand) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, []byte{byte(m.Packet >> 8), byte(m.Packet)}...)
	result = append(result, []byte{byte(m.Repeat >> 24), byte(m.Repeat >> 16), byte(m.Repeat >> 8), byte(m.Repeat)}...)
	result = append(result, []byte(m.Command)...)
	return result
}

type MessageCommandCompleted struct {
	Packet uint16
	Time   uint32
}

func (m *MessageCommandCompleted) UnmarshalBinary(data []byte) error {
	m.Packet = uint16(data[0])<<8 | uint16(data[1])
	m.Time = uint32(data[2])<<24 | uint32(data[3])<<16 | uint32(data[4])<<8 | uint32(data[5])
	return nil
}

func (m MessageCommandCompleted) Kind() string {
	return "comc"
}

func (m MessageCommandCompleted) Value() string {
	return ""
}

func (m MessageCommandCompleted) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageCommand{})
	gMessages = append(gMessages, &MessageCommandCompleted{})
}
