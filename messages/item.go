package messages

import (
	"fmt"
	"strconv"
	"strings"
)

type ItemObject struct {
	Tag         int32
	Flags       int32
	Weight      int32
	TotalWeight int32
	Face        int32
	Name        string
	PluralName  string
	Anim        int16
	AnimSpeed   int8
	Nrof        int32
	Type        int16
}

func (o ItemObject) GetName() string {
	if o.Nrof > 1 {
		return fmt.Sprintf("%d %s", o.Nrof, o.PluralName)
	}
	return o.Name
}

type MessageItem2 struct {
	Location int32
	Objects  []ItemObject
}

func (m *MessageItem2) UnmarshalBinary(data []byte) error {
	offset := 0
	m.Location = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
	offset += 4
	m.Objects = make([]ItemObject, 0)
	for offset < len(data) {
		var obj ItemObject
		obj.Tag = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		obj.Flags = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		obj.Weight = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		obj.Face = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		offset += 4
		obj.Name, offset = readLengthPrefixedString(data, offset)
		{ // SC 1024 support
			parts := strings.Split(obj.Name, "\x00")
			if len(parts) > 1 {
				obj.Name = parts[0]
				obj.PluralName = parts[1]
			}
		}
		obj.Anim = int16(data[offset])<<8 | int16(data[offset+1])
		offset += 2
		obj.AnimSpeed = int8(data[offset])
		offset++
		obj.Nrof = int32(data[offset])<<24 | int32(data[offset+1])<<16 | int32(data[offset+2])<<8 | int32(data[offset+3])
		{
			obj.TotalWeight = obj.Weight * obj.Nrof
		}
		offset += 4
		obj.Type = int16(data[offset])<<8 | int16(data[offset+1])
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

type MessageUpdateItem struct {
	Tag    int32
	Flags  int8
	Values []any // TODO
}

func (m *MessageUpdateItem) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageUpdateItem) Kind() string {
	return "upditem"
}

func (m MessageUpdateItem) Value() string {
	return ""
}

func (m MessageUpdateItem) Bytes() []byte {
	return nil
}

type MessageDeleteItem struct {
	Tags []int32
}

func (m *MessageDeleteItem) UnmarshalBinary(data []byte) error {
	for i := 0; i < len(data); i += 4 {
		m.Tags = append(m.Tags, int32(data[i])<<24|int32(data[i+1])<<16|int32(data[i+2])<<8|int32(data[i+3]))
	}
	return nil
}

func (m MessageDeleteItem) Kind() string {
	return "delitem"
}

func (m MessageDeleteItem) Value() string {
	return fmt.Sprintf("%+v", m.Tags)
}

func (m MessageDeleteItem) Bytes() []byte {
	return nil
}

type MessageDeleteInventory struct {
	Tag int32
}

func (m *MessageDeleteInventory) UnmarshalBinary(data []byte) error {
	v, _ := strconv.Atoi(string(data))
	m.Tag = int32(v)
	return nil
}

func (m MessageDeleteInventory) Kind() string {
	return "delinv"
}

func (m MessageDeleteInventory) Value() string {
	return fmt.Sprintf("%d", m.Tag)
}

func (m MessageDeleteInventory) Bytes() []byte {
	return nil
}

func init() {
	gMessages = append(gMessages, &MessageItem2{})
	gMessages = append(gMessages, &MessageUpdateItem{})
	gMessages = append(gMessages, &MessageDeleteItem{})
	gMessages = append(gMessages, &MessageDeleteInventory{})
}
