package messages

import "fmt"

type MessageAnim struct {
	AnimID int16
	Flags  int16
	Faces  []int16
}

func (m MessageAnim) Bytes() []byte {
	return nil
}

func (m *MessageAnim) UnmarshalBinary(data []byte) error {
	offset := 0
	m.AnimID = int16(data[offset])<<8 | int16(data[offset+1])
	offset += 2
	m.Flags = int16(data[offset])<<8 | int16(data[offset+1])
	offset += 2
	m.Faces = make([]int16, 0)
	for offset < len(data) {
		m.Faces = append(m.Faces, int16(data[offset])<<8|int16(data[offset+1]))
		offset += 2
	}
	return nil
}

func (m MessageAnim) Kind() string {
	return "anim"
}

func (m MessageAnim) Value() string {
	return ""
}

func (m MessageAnim) String() string {
	return fmt.Sprintf("AnimID: %d, Flags: %d, Faces: %v", m.AnimID, m.Flags, m.Faces)
}

func init() {
	gMessages = append(gMessages, &MessageAnim{})
}
