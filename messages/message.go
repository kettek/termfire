package messages

import (
	"errors"

	"github.com/qdm12/reprint"
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
			// This is bad, but my prior foolish decisions have led to this. FIXME: Replace the above unmarshaling into with a function that returns a new instance of the message.
			return reprint.This(m).(Message), nil
		}
	}
	return nil, errors.New("unknown message type" + msgType)
}

var gMessages []Message
