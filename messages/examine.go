package messages

import "fmt"

type MessageExamine struct {
	Tag int32
}

func (m *MessageExamine) UnmarshalBinary(data []byte) error {
	// TODO
	return nil
}

func (m *MessageExamine) Kind() string {
	return "examine"
}

func (m MessageExamine) Value() string {
	return ""
}

func (m MessageExamine) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(fmt.Sprintf("%d", m.Tag))...)
	return result
}

func init() {
	gMessages = append(gMessages, &MessageExamine{})
}
