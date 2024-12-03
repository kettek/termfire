package messages

type MessageLock struct {
	Lock bool
	Tag  int32
}

func (m *MessageLock) UnmarshalBinary(data []byte) error {
	// TODO
	return nil
}

func (m MessageLock) Kind() string {
	return "lock"
}

func (m MessageLock) Value() string {
	return ""
}

func (m MessageLock) Bytes() []byte {
	var results []byte
	results = append(results, []byte(m.Kind())...)
	results = append(results, ' ')
	if m.Lock {
		results = append(results, byte(1))
	} else {
		results = append(results, byte(0))
	}
	results = append(results, []byte{byte(m.Tag >> 24), byte(m.Tag >> 16), byte(m.Tag >> 8), byte(m.Tag)}...)
	return results
}

func init() {
	gMessages = append(gMessages, &MessageLock{})
}
