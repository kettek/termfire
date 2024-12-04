package messages

import "fmt"

type MessageLookAt struct {
	DX, DY int
}

func (m *MessageLookAt) UnmarshalBinary(data []byte) error {
	// TODO
	return nil
}

func (m *MessageLookAt) Kind() string {
	return "lookat"
}

func (m MessageLookAt) Value() string {
	return ""
}

func (m MessageLookAt) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(fmt.Sprintf("%d %d", m.DX, m.DY))...)
	return result
}

func init() {
	gMessages = append(gMessages, &MessageLookAt{})
}
