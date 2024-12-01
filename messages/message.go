package messages

import (
	"errors"
)

type Message interface {
	Bytes() []byte
	Kind() string
	Value() string
	UnmarshalBinary([]byte) error
}

func UnmarshalMessage(data []byte) (Message, error) {
	msgType := ""
	var msgData []byte
	noEnd := true
	for i := 0; i < len(data); i++ {
		if data[i] == ' ' {
			msgType = string(data[:i])
			noEnd = false
			break
		}
	}
	if noEnd {
		msgType = string(data)
	} else {
		msgData = data[len(msgType)+1:]
	}
	for _, m := range gMessages {
		if m.Kind() == msgType {
			if err := m.UnmarshalBinary(msgData); err != nil {
				return nil, err
			}
			return m, nil
		}
	}
	return nil, errors.New("unknown message type" + msgType)
}

var gMessages []Message
