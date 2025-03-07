package messages

import (
	"strconv"
	"strings"

	"github.com/kettek/termfire/debug"
)

func LengthPrefixedString(s string) []byte {
	// 8-bit length followed by bytes.
	result := []byte{byte(len(s))}
	result = append(result, []byte(s)...)
	return result
}

func readLengthPrefixedString(data []byte, length int) (string, int) {
	strlen := int(data[length])
	length++
	return string(data[length : length+strlen]), length + strlen
}

func readLengthPrefixedString2(data []byte, length int) (string, int) {
	strlen := int(data[length])<<8 | int(data[length+1])
	length += 2
	return string(data[length : length+strlen]), length + strlen
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
	FaceCache struct {
		Use   bool
		Value bool
	}
	FaceSet struct {
		Use   bool
		Value uint8
	}
	LoginMethod struct {
		Use   bool
		Value string
	}
	ExtendedStats struct {
		Use   bool
		Value bool
	}
	MapSize struct {
		Use   bool
		Value string
	}
	Sound2 struct {
		Use   bool
		Value uint8
	}
	SpellMon struct {
		Use   bool
		Value uint8 // 0, 1, 2
	}
	Tick struct {
		Use   bool
		Value uint8 // 0, 1
	}
}

func (m *MessageSetup) UnmarshalBinary(data []byte) error {
	parts := strings.Split(string(data), " ")
	for i := 0; i < len(parts); i += 2 {
		switch parts[i] {
		case "facecache":
			m.FaceCache.Use = true
			m.FaceCache.Value, _ = strconv.ParseBool(parts[i+1])
		case "faceset":
			m.FaceSet.Use = true
			v, _ := strconv.ParseUint(parts[i+1], 10, 8)
			m.FaceSet.Value = uint8(v)
		case "loginmethod":
			m.LoginMethod.Use = true
			m.LoginMethod.Value = parts[i+1]
		case "extendedstats":
			m.ExtendedStats.Use = true
			m.ExtendedStats.Value, _ = strconv.ParseBool(parts[i+1])
		case "mapsize":
			m.MapSize.Use = true
			m.MapSize.Value = parts[i+1]
		case "sound2":
			m.Sound2.Use = true
			v, _ := strconv.ParseUint(parts[i+1], 10, 8)
			m.Sound2.Value = uint8(v)
		case "spellmon":
			m.SpellMon.Use = true
			v, _ := strconv.ParseUint(parts[i+1], 10, 8)
			m.SpellMon.Value = uint8(v)
		case "tick":
			m.Tick.Use = true
			v, _ := strconv.ParseUint(parts[i+1], 10, 8)
			m.Tick.Value = uint8(v)
		}
	}
	return nil
}

func (m MessageSetup) Kind() string {
	return "setup"
}

func (m MessageSetup) Value() string {
	return "TODO"
}

func (m MessageSetup) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	if m.ExtendedStats.Use {
		result = append(result, ' ')
		result = append(result, []byte("extended_stats")...)
		result = append(result, ' ')
		result = append(result, '1')
	}
	if m.FaceCache.Use {
		result = append(result, ' ')
		result = append(result, []byte("facecache")...)
		result = append(result, ' ')
		result = append(result, '1')
	}
	if m.FaceSet.Use {
		result = append(result, ' ')
		result = append(result, []byte("faceset")...)
		result = append(result, ' ')
		result = append(result, []byte(strconv.Itoa(int(m.FaceSet.Value)))...)
	}
	if m.LoginMethod.Use {
		result = append(result, ' ')
		result = append(result, []byte("loginmethod")...)
		result = append(result, ' ')
		result = append(result, []byte(m.LoginMethod.Value)...)
	}
	if m.MapSize.Use {
		result = append(result, ' ')
		result = append(result, []byte("mapsize")...)
		result = append(result, ' ')
		result = append(result, []byte(m.MapSize.Value)...)
	}
	if m.Sound2.Use {
		result = append(result, ' ')
		result = append(result, []byte("sound2")...)
		result = append(result, ' ')
		result = append(result, []byte(strconv.Itoa(int(m.Sound2.Value)))...)
	}
	if m.SpellMon.Use {
		result = append(result, ' ')
		result = append(result, []byte("spellmon")...)
		result = append(result, ' ')
		result = append(result, []byte(strconv.Itoa(int(m.SpellMon.Value)))...)
	}
	if m.Tick.Use {
		result = append(result, ' ')
		result = append(result, []byte("tick")...)
		result = append(result, ' ')
		result = append(result, []byte(strconv.Itoa(int(m.Tick.Value)))...)
	}
	debug.Debug("Bytes:", string(result))
	return result
}

func init() {
	gMessages = append(gMessages, &MessageFailure{})
	gMessages = append(gMessages, &MessageVersion{})
	gMessages = append(gMessages, &MessageSetup{})
}
