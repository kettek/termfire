package messages

type MessageTick uint32

func (m *MessageTick) UnmarshalBinary(data []byte) error {
	*m = MessageTick(uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3]))
	return nil
}

func (m MessageTick) Kind() string {
	return "tick"
}

func (m MessageTick) Value() string {
	return ""
}

func (m MessageTick) Bytes() []byte {
	return nil
}

func init() {
	tick := MessageTick(0)
	gMessages = append(gMessages, &tick)
}
