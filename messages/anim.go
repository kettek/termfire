package messages

type MessageAnim struct {
	AnimID uint16
	Flags  uint16
	Faces  []uint16
}

func (m MessageAnim) Bytes() []byte {
	return nil
}

func (m *MessageAnim) UnmarshalBinary(data []byte) error {
	offset := 0
	m.AnimID = uint16(data[offset])<<8 | uint16(data[offset+1])
	offset += 2
	m.Flags = uint16(data[offset])<<8 | uint16(data[offset+1])
	offset += 2
	m.Faces = make([]uint16, 0)
	for offset < len(data) {
		m.Faces = append(m.Faces, uint16(data[offset])<<8|uint16(data[offset+1]))
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

func init() {
	gMessages = append(gMessages, &MessageAnim{})
}
