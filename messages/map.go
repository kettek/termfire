package messages

import (
	"fmt"
	"strconv"

	"github.com/kettek/termfire/debug"
)

type MessageMap2CoordType int16

const (
	MessageMap2CoordTypeNormal MessageMap2CoordType = iota
	MessageMap2CoordTypeScrollInformation
)

type MessageMap2CoordData interface {
}

type MessageMap2CoordDataClear struct {
}

type MessageMap2CoordDataClearLayer struct {
	Layer int8
}

type MessageMap2CoordDataDarkness struct {
	Darkness int8
}

type MessageMap2CoordDataImage struct {
	Layer        int8
	FaceNum      int16
	HasAnimSpeed bool
	AnimSpeed    int8
	HasSmooth    bool
	Smooth       int8
}

type MessageMap2CoordDataAnim struct {
	Layer  int8
	Anim   int16
	Flags  int8
	Speed  int8
	Smooth int8
}

func (m MessageMap2CoordDataAnim) String() string {
	return fmt.Sprintf("Layer: %d, Anim: %d, Flags: %d, Speed: %d, Smooth: %d", m.Layer, m.Anim, m.Flags, m.Speed, m.Smooth)
}

func (m MessageMap2CoordDataImage) String() string {
	r := fmt.Sprintf("Layer: %d, FaceNum: %d", m.Layer, m.FaceNum)
	if m.HasAnimSpeed {
		r += fmt.Sprintf(", AnimSpeed: %d", m.AnimSpeed)
	}
	if m.HasSmooth {
		r += fmt.Sprintf(", Smooth: %d", m.Smooth)
	}
	return r
}

type MessageMapCoord struct {
	X, Y int
	Type MessageMap2CoordType
	Data []MessageMap2CoordData
}

func (m MessageMapCoord) String() string {
	r := fmt.Sprintf("X: %d, Y: %d, ", m.X, m.Y)
	if m.Type == MessageMap2CoordTypeNormal {
		r += "Normal"
	} else {
		r += "ScrollInformation"
	}
	for _, d := range m.Data {
		switch d.(type) {
		case *MessageMap2CoordDataClear:
			r += ", Clear"
		case *MessageMap2CoordDataDarkness:
			r += fmt.Sprintf(", Darkness: %d", d.(*MessageMap2CoordDataDarkness).Darkness)
		case *MessageMap2CoordDataImage:
			r += fmt.Sprintf(", Image: [%s]", d.(*MessageMap2CoordDataImage))
		}
	}
	return r
}

func (m *MessageMapCoord) UnmarshalBinary(data []byte) (int, error) {
	var offset int
	var coord int16
	coord = int16(data[offset])<<8 | int16(data[offset+1])
	// X is the first 6 bits.
	m.X = (int((coord) >> 10 & 0x3F)) - 15
	// Y is the next 6 bits after X.
	m.Y = (int((coord) >> 4 & 0x3F)) - 15
	// Type is LSB 0-3
	m.Type = MessageMap2CoordType(coord & 0x3)
	//m.Type = MessageMap2CoordType(coord & 0x0f)
	offset += 2

	if m.Type&0x01 != 0 {
		// Just scroll information
		return offset, nil
	}

	for offset < len(data) {
		var lenType uint8
		lenType = data[offset]
		offset++

		if lenType == 255 {
			break
		}

		// len is the top 3 bits.
		var dataLen uint8
		//		dataLen = (lenType >> 5) & 0x07
		dataLen = (lenType >> 5)
		// type is the bottom 5 bits.
		var dataType uint8
		dataType = lenType & 0x1F

		// FIXME: This is only true if server protocol is 1030+. Probably bail if <1030 or something.
		if dataLen == 0x07 {
			dataLen = data[offset]
			offset++
		}

		switch dataType {
		case 0x0:
			m.Data = append(m.Data, &MessageMap2CoordDataClear{})
		case 0x1:
			if dataLen != 1 {
				return 0, fmt.Errorf("dataLen for darkness is not 1, got %d", dataLen)
			}
			var darkness MessageMap2CoordDataDarkness
			darkness.Darkness = int8(data[offset])
			m.Data = append(m.Data, &darkness)
		case 0x2: // label SC 1030
			labelType := data[offset]
			labelLen := data[offset+1]
			label := string(data[offset+2 : offset+2+int(labelLen)])
			debug.Debug("label", labelType, label)
		default:
			if dataType >= 0x10 && dataType <= 0x19 {
				var image MessageMap2CoordDataImage
				image.Layer = int8(dataType - 0x10)

				faceOrAnim := int16(data[offset])<<8 | int16(data[offset+1])
				var isAnim bool
				if (faceOrAnim >> 15) != 0 {
					isAnim = true
				}

				if dataLen == 2 {
					// No smooth
				} else if dataLen == 3 {
					if isAnim {
						image.AnimSpeed = int8(data[offset+2])
						image.HasAnimSpeed = true
					} else {
						image.Smooth = int8(data[offset+2])
						image.HasSmooth = true
					}
				} else if dataLen == 4 {
					image.AnimSpeed = int8(data[offset+2])
					image.Smooth = int8(data[offset+3])
					image.HasAnimSpeed = true
					image.HasSmooth = true
				} else {
					return 0, fmt.Errorf("dataLen is not 2, 3, or 4, got %d", dataLen)
				}

				if faceOrAnim == 0 {
					// Clear layer!
					m.Data = append(m.Data, MessageMap2CoordDataClearLayer{Layer: image.Layer})
				} else if isAnim {
					animFlags := (faceOrAnim >> 6) & 0x03
					animation := int16(faceOrAnim) & 0x1fff
					anim := MessageMap2CoordDataAnim{
						Layer:  image.Layer,
						Anim:   animation,
						Flags:  int8(animFlags),
						Speed:  image.AnimSpeed,
						Smooth: image.Smooth,
					}
					m.Data = append(m.Data, &anim)
				} else {
					image.FaceNum = faceOrAnim
					m.Data = append(m.Data, &image)
				}
			} else {
				return 0, fmt.Errorf("Unknown data type %d", dataType)
			}
		}
		offset += int(dataLen)
	}

	return offset, nil
}

type MessageMap2 struct {
	Coords []MessageMapCoord
}

func (m *MessageMap2) UnmarshalBinary(data []byte) error {
	m.Coords = make([]MessageMapCoord, 0)
	for i := 0; i < len(data); {
		var mapCoord MessageMapCoord

		if count, err := mapCoord.UnmarshalBinary(data[i:]); err != nil {
			return err
		} else {
			m.Coords = append(m.Coords, mapCoord)
			i += count
		}
	}
	return nil
}

func (m MessageMap2) Kind() string {
	return "map2"
}

func (m MessageMap2) Value() string {
	r := ""
	for _, c := range m.Coords {
		r += c.String() + "\n"
	}
	return r
}

func (m MessageMap2) Bytes() []byte {
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
	Face      int16
	SmoothPic int16
}

func (m *MessageSmooth) UnmarshalBinary(data []byte) error {
	m.Face = int16(data[0])<<8 | int16(data[1])
	m.SmoothPic = int16(data[2])<<8 | int16(data[3])
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
	Tag    int32
	Weight int32
	Face   int32
	Name   string
}

func (m *MessagePlayer) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Tag = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
	offset += 4
	m.Weight = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
	offset += 4
	m.Face = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
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
	gMessages = append(gMessages, &MessageNewMap{})
	gMessages = append(gMessages, &MessageSmooth{})
	gMessages = append(gMessages, &MessagePlayer{})
}
