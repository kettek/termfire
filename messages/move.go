package messages

import "fmt"

type MessageMove struct {
	To   int32
	Tag  int32
	Nrof int32
}

func (m *MessageMove) UnmarshalBinary(data []byte) error {
	// TODO
	return nil
}

func (m MessageMove) Kind() string {
	return "move"
}

func (m MessageMove) Value() string {
	return ""
}

func (m MessageMove) Bytes() []byte {
	var results []byte
	results = append(results, []byte(m.Kind())...)
	results = append(results, ' ')
	results = append(results, []byte(fmt.Sprintf("%d", m.To))...)
	results = append(results, ' ')
	results = append(results, []byte(fmt.Sprintf("%d", m.Tag))...)
	results = append(results, ' ')
	results = append(results, []byte(fmt.Sprintf("%d", m.Nrof))...)
	return results
}

func init() {
	gMessages = append(gMessages, &MessageMove{})
}
