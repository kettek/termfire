package messages

type MessageInscribe struct {
	SpellTag  int32
	ScrollTag int32
}

func (m *MessageInscribe) UnmarshalBinary(data []byte) error {
	// TODO
	return nil
}

func (m MessageInscribe) Kind() string {
	return "inscribe"
}

func (m MessageInscribe) Value() string {
	return ""
}

func (m MessageInscribe) Bytes() []byte {
	var results []byte
	results = append(results, []byte(m.Kind())...)
	results = append(results, ' ')
	results = append(results, byte(0))
	results = append(results, []byte{byte(m.SpellTag >> 24), byte(m.SpellTag >> 16), byte(m.SpellTag >> 8), byte(m.SpellTag)}...)
	results = append(results, []byte{byte(m.ScrollTag >> 24), byte(m.ScrollTag >> 16), byte(m.ScrollTag >> 8), byte(m.ScrollTag)}...)
	return results
}

func init() {
	gMessages = append(gMessages, &MessageInscribe{})
}
