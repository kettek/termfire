package messages

import (
	"strconv"
	"strings"

	"github.com/kettek/termfire/debug"
)

type MessageMap2CoordType uint16

const (
	MessageMap2CoordTypeNormal MessageMap2CoordType = iota
	MessageMap2CoordTypeScrollInformation
)

type MessageMap2CoordData interface {
}

type MessageMap2CoordDataClear struct {
}

type MessageMap2CoordDataDarkness struct {
	Darkness uint8
}

type MessageMap2CoordDataImage struct {
	Layer        uint8
	FaceNum      uint16
	HasAnimSpeed bool
	AnimSpeed    uint8
	HasSmooth    bool
	Smooth       uint8
}

type MessageMapCoord struct {
	X, Y uint16
	Type MessageMap2CoordType
	Data []MessageMap2CoordData
}

func (m *MessageMapCoord) UnmarshalBinary(data []byte) (int, error) {
	var offset int
	var coord uint16
	coord = uint16(data[offset])<<8 | uint16(data[offset+1])
	// X is the first 6 bits.
	m.X = uint16(coord >> 10)
	// Y is the next 6 bits after X.
	m.Y = uint16(coord & 0x3F)
	// Type is LSB 0-3
	m.Type = MessageMap2CoordType(coord & 0x3)
	offset += 2

	for offset < len(data) {
		var lenType uint8
		lenType = data[offset]

		// len is the top 3 bits.
		dataLen := int(lenType >> 5)
		// type is the bottom 5 bits.
		dataType := lenType & 0x1F

		switch dataType {
		case 0x0:
			m.Data = append(m.Data, &MessageMap2CoordDataClear{})
		case 0x1:
			var darkness MessageMap2CoordDataDarkness
			darkness.Darkness = data[offset+1]
			m.Data = append(m.Data, &darkness)
		default:
			// FIXME: 99% this is wrong.
			if dataType >= 0x10 && dataType <= 0x19 {
				var image MessageMap2CoordDataImage
				if dataLen == 2 {
					image.FaceNum = uint16(data[offset+1])<<8 | uint16(data[offset+2])
				} else if dataLen == 3 {
					image.FaceNum = uint16(data[offset+1])<<8 | uint16(data[offset+2])
					// If facenum's high bit is set, it has an animation.
					if image.FaceNum&0x8000 != 0 {
						image.AnimSpeed = uint8(data[offset+3])
						image.HasAnimSpeed = true
					} else {
						image.Smooth = uint8(data[offset+3])
						image.HasSmooth = true
					}
				} else if dataLen == 4 {
					image.FaceNum = uint16(data[offset+1])<<8 | uint16(data[offset+2])
					image.AnimSpeed = uint8(data[offset+3])
					image.Smooth = uint8(data[offset+4])
					image.HasAnimSpeed = true
					image.HasSmooth = true
				}
				m.Data = append(m.Data, &image)
			}
			offset += dataLen
		}
	}

	return offset, nil
}

type MessageMap2 struct {
	Coords []MessageMapCoord
}

func (m *MessageMap2) UnmarshalBinary(data []byte) error {
	for i := 0; i < len(data); i++ {
		var mapCoord MessageMapCoord

		if count, err := mapCoord.UnmarshalBinary(data[i:]); err != nil {
			return err
		} else {
			i += count
		}

	}
	debug.Debug("FUCK", m.Coords)
	return nil
}

func (m MessageMap2) Kind() string {
	return "map2"
}

func (m MessageMap2) Value() string {
	return "TODO"
}

func (m MessageMap2) Bytes() []byte {
	return nil
}

type ItemObject struct {
	Tag         uint32
	Flags       uint32
	Weight      uint32
	TotalWeight uint32
	Face        uint32
	Name        string
	PluralName  string
	Anim        uint16
	AnimSpeed   uint8
	Nrof        uint32
	Type        uint16
}

type MessageItem2 struct {
	Location uint32
	Objects  []ItemObject
}

func (m *MessageItem2) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Location = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
	offset += 4
	for offset < len(data) {
		var obj ItemObject
		obj.Tag = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
		offset += 4
		obj.Flags = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
		offset += 4
		obj.Weight = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
		offset += 4
		obj.Face = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
		offset += 4
		obj.Name, offset = readLengthPrefixedString(data, offset)
		{ // SC 1024 support
			parts := strings.Split(obj.Name, "\x00")
			if len(parts) > 1 {
				obj.Name = parts[0]
				obj.PluralName = parts[1]
			}
		}
		obj.Anim = uint16(data[offset])<<8 | uint16(data[offset+1])
		offset += 2
		obj.AnimSpeed = uint8(data[offset])
		offset++
		obj.Nrof = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
		{
			obj.TotalWeight = obj.Weight * obj.Nrof
		}
		offset += 4
		obj.Type = uint16(data[offset])<<8 | uint16(data[offset+1])
		offset += 2
		m.Objects = append(m.Objects, obj)
	}
	return nil
}

func (m MessageItem2) Kind() string {
	return "item2"
}

func (m MessageItem2) Value() string {
	var result string
	result += "location: " + strconv.Itoa(int(m.Location)) + "\n"
	for _, o := range m.Objects {
		result += strconv.Itoa(int(o.Nrof)) + " " + o.Name + "/" + o.PluralName + " " + strconv.Itoa(int(o.Flags)) + " " + strconv.Itoa(int(o.Weight)) + " " + strconv.Itoa(int(o.TotalWeight)) + " " + strconv.Itoa(int(o.Face)) + " " + strconv.Itoa(int(o.Anim)) + " " + strconv.Itoa(int(o.AnimSpeed)) + " " + strconv.Itoa(int(o.Type)) + "\n"
	}
	return result
}

func (m MessageItem2) Bytes() []byte {
	return nil
}

type MessageNewMap struct{}

func (m *MessageNewMap) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageNewMap) Kind() string {
	return "newmap"
}

func (m MessageNewMap) Value() string {
	return ""
}

func (m MessageNewMap) Bytes() []byte {
	return nil
}

type MessageSmooth struct {
	Face      uint16
	SmoothPic uint16
}

func (m *MessageSmooth) UnmarshalBinary(data []byte) error {
	m.Face = uint16(data[0])<<8 | uint16(data[1])
	m.SmoothPic = uint16(data[2])<<8 | uint16(data[3])
	return nil
}

func (m MessageSmooth) Kind() string {
	return "smooth"
}

func (m MessageSmooth) Value() string {
	return strconv.Itoa(int(m.Face)) + " " + strconv.Itoa(int(m.SmoothPic))
}

func (m MessageSmooth) Bytes() []byte {
	return nil
}

type MessagePlayer struct {
	Tag    uint32
	Weight uint32
	Face   uint32
	Name   string
}

func (m *MessagePlayer) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Tag = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
	offset += 4
	m.Weight = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
	offset += 4
	m.Face = uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
	offset += 4
	m.Name, _ = readLengthPrefixedString(data, offset)
	return nil
}

func (m MessagePlayer) Kind() string {
	return "player"
}

func (m MessagePlayer) Value() string {
	return "@ " + m.Name + "(" + strconv.Itoa(int(m.Tag)) + ") " + strconv.Itoa(int(m.Weight)) + " " + strconv.Itoa(int(m.Face))
}

func (m MessagePlayer) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageMap2{})
	gMessages = append(gMessages, &MessageItem2{})
	gMessages = append(gMessages, &MessageNewMap{})
	gMessages = append(gMessages, &MessageSmooth{})
	gMessages = append(gMessages, &MessagePlayer{})
}
