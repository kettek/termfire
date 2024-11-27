package messages

import (
	"strconv"
	"strings"
)

func LengthPrefixedString(s string) []byte {
	// 8-bit length followed by bytes.
	result := []byte{byte(len(s))}
	result = append(result, []byte(s)...)
	return result
}

type MessageVersion struct {
	CLVersion string
	SVVersion string
	SVName    string
	value     string
}

func (m *MessageVersion) UnmarshalBinary(data []byte) error {
	parts := strings.SplitN(string(data), " ", 3)

	if len(parts) > 0 {
		m.CLVersion = parts[0]
	}
	if len(parts) > 1 {
		m.SVVersion = parts[1]
	}
	if len(parts) > 2 {
		m.SVName = parts[2]
	}

	return nil
}

func (m MessageVersion) Value() string {
	return m.CLVersion + " " + m.SVVersion + " " + m.SVName
}

func (m MessageVersion) Kind() string {
	return "version"
}

func (m MessageVersion) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(m.CLVersion)...)
	result = append(result, ' ')
	result = append(result, []byte(m.SVVersion)...)
	result = append(result, ' ')
	result = append(result, []byte(m.SVName)...)
	return result
}

type MessageFailure struct {
	Command string
	Reason  string
}

func (m *MessageFailure) UnmarshalBinary(data []byte) error {
	parts := strings.SplitN(string(data), " ", 2)

	if len(parts) > 0 {
		m.Command = parts[0]
	}
	if len(parts) > 1 {
		m.Reason = parts[1]
	}

	return nil
}

func (m MessageFailure) Value() string {
	return m.Command + " " + m.Reason
}

func (m MessageFailure) Kind() string {
	return "failure"
}

func (m MessageFailure) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(m.Command)...)
	result = append(result, ' ')
	result = append(result, []byte(m.Reason)...)
	return result
}

type MessageSetup struct {
	FaceCache   bool
	LoginMethod string
}

func (m *MessageSetup) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageSetup) Kind() string {
	return "setup"
}

func (m MessageSetup) Value() string {
	return strconv.FormatBool(m.FaceCache)
}

func (m MessageSetup) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte("facecache")...)
	result = append(result, ' ')
	result = append(result, '1')
	result = append(result, ' ')
	result = append(result, []byte("loginmethod")...)
	result = append(result, ' ')
	result = append(result, []byte(m.LoginMethod)...)
	return result
}

func init() {
	gMessages = append(gMessages, &MessageFailure{})
	gMessages = append(gMessages, &MessageVersion{})
	gMessages = append(gMessages, &MessageSetup{})
}
