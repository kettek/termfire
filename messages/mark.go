package messages

type MessageMark struct {
	Tag int32
}

func (m *MessageMark) UnmarshalBinary(data []byte) error {
	// TODO
	return nil
}

func (m MessageMark) Kind() string {
	return "mark"
}

func (m MessageMark) Value() string {
	return ""
}

func (m MessageMark) Bytes() []byte {
	var results []byte
	results = append(results, []byte(m.Kind())...)
	results = append(results, ' ')
	results = append(results, []byte{byte(m.Tag >> 24), byte(m.Tag >> 16), byte(m.Tag >> 8), byte(m.Tag)}...)
	return results
}

func init() {
	gMessages = append(gMessages, &MessageMark{})
}
