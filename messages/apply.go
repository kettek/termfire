package messages

import "fmt"

type MessageApply struct {
	Tag int32
}

func (m *MessageApply) UnmarshalBinary(data []byte) error {
	// TODO
	return nil
}

func (m *MessageApply) Kind() string {
	return "apply"
}

func (m MessageApply) Value() string {
	return ""
}

func (m MessageApply) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(fmt.Sprintf("%d", m.Tag))...)
	return result
}

func init() {
	gMessages = append(gMessages, &MessageApply{})
}
