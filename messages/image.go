package messages

import (
	"fmt"
)

type MessageImage2 struct {
	Face uint32
	Set  uint8
	// We extract the width and height manually for the TUI client, so it doesn't have to actually manage PNGs.
	Width  int
	Height int
	Data   []byte // PNG payload.
}

func (m *MessageImage2) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Face = uint32(data[0])<<24 | uint32(data[1])<<16 | uint32(data[2])<<8 | uint32(data[3])
	offset += 4
	m.Set = uint8(data[4])
	offset++
	dataLen := uint32(data[5])<<24 | uint32(data[6])<<16 | uint32(data[7])<<8 | uint32(data[8])
	// NOTE: We could make this optional, in the event the client doesn't actually want the data...
	m.Data = data[offset : offset+int(dataLen)]
	offset += 4
	// For now we just want the width and height so we can determine columns and rows the image contains. If this library is used for graphics, we'd want the actual data.
	for i := 0; i < len(data); i++ {
		offset += 8 // skip header
		offset += 4 // skip chunk length
		offset += 4 // skip chunk type
		// Get width and height
		m.Width = int(data[offset])<<24 | int(data[offset+1])<<16 | int(data[offset+2])<<8 | int(data[offset+3])
		offset += 4
		m.Height = int(data[offset])<<24 | int(data[offset+1])<<16 | int(data[offset+2])<<8 | int(data[offset+3])
		break
	}
	return nil
}

func (m MessageImage2) Kind() string {
	return "image2"
}

func (m MessageImage2) Value() string {
	return fmt.Sprintf("%d %d %d %d", m.Face, m.Set, m.Width, m.Height)
}

func (m MessageImage2) Bytes() []byte {
	return nil
}

type MessageAskFace struct {
	Face uint32
}

func (m *MessageAskFace) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageAskFace) Kind() string {
	return "askface"
}

func (m MessageAskFace) Value() string {
	return ""
}

func (m MessageAskFace) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	result = append(result, []byte(fmt.Sprintf("%d", m.Face))...)
	return result
}

func init() {
	gMessages = append(gMessages, &MessageImage2{})
}
