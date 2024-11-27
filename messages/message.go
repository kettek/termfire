package messages

import "errors"

type Message interface {
	Bytes() []byte
	Kind() string
	Value() string
	UnmarshalBinary([]byte) error
}

func UnmarshalMessage(data []byte) (Message, error) {
	msgType := ""
	for i := 0; i < len(data); i++ {
		if data[i] == ' ' {
			msgType = string(data[:i])
			break
		}
	}
	for _, m := range gMessages {
		if m.Kind() == msgType {
			if err := m.UnmarshalBinary(data[len(msgType)+1:]); err != nil {
				return nil, err
			}
			return m, nil
		}
	}
	return nil, errors.New("unknown message type")
}

var gMessages []Message
