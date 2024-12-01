package messages

var gMessageStats []MessageStat

type MessageStat interface {
	UnmarshalBinary([]byte) (int, error)
	Matches(byte) bool
}

type MessageStatHP int16

func (m *MessageStatHP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatHP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatHP) Matches(id byte) bool {
	return id == 1
}

type MessageStatMaxHP int16

func (m *MessageStatMaxHP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatMaxHP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatMaxHP) Matches(id byte) bool {
	return id == 2
}

type MessageStatMaxSP int16

func (m *MessageStatMaxSP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatMaxSP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatMaxSP) Matches(id byte) bool {
	return id == 3
}

type MessageStatSP int16

func (m *MessageStatSP) UnmarshalBinary(data []byte) (int, error) {
	*m = (MessageStatSP)(int16(data[0])<<8 | int16(data[1]))
	return 2, nil
}

func (m MessageStatSP) Matches(id byte) bool {
	return id == 4
}

func init() {
	statHP := MessageStatHP(0)
	gMessageStats = append(gMessageStats, &statHP)
	statMaxHP := MessageStatMaxHP(0)
	gMessageStats = append(gMessageStats, &statMaxHP)
	statSP := MessageStatSP(0)
	gMessageStats = append(gMessageStats, &statSP)
	statMaxSP := MessageStatMaxSP(0)
	gMessageStats = append(gMessageStats, &statMaxSP)
}

type MessageStats struct {
	Stats []MessageStat
}

func (m *MessageStats) UnmarshalBinary(data []byte) error {
	for i := 0; i < len(data); {
		kind := data[i]
		for _, s := range gMessageStats {
			if s.Matches(kind) {
				if count, err := s.UnmarshalBinary(data[i+1:]); err != nil {
					return err
				} else {
					i += count
				}
				m.Stats = append(m.Stats, s)
				break
			}
		}
	}

	return nil
}

func (m MessageStats) Kind() string {
	return "stats"
}

func (m MessageStats) Value() string {
	return ""
}

func (m MessageStats) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageStats{})
}
