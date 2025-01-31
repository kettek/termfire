package messages

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kettek/termfire/debug"
)

// TODO: image_sums, exp_table, knowledge_info, skill_info, skill_extra, spell_paths, race_list, race_info, class_list, class_info, startingmap, newcharinfo, news, rules, motd
type MessageRequestInfoData interface {
	Kind() string
	Bytes() []byte
}

type MessageRequestInfoDataImageInfo struct{}

func (m MessageRequestInfoDataImageInfo) Kind() string {
	return "image_info"
}

func (m MessageRequestInfoDataImageInfo) Bytes() []byte {
	return nil
}

type MessageRequestInfo struct {
	Data MessageRequestInfoData
}

func (m *MessageRequestInfo) UnmarshalBinary(data []byte) error {
	return nil
}

func (m MessageRequestInfo) Kind() string {
	return "requestinfo"
}

func (m *MessageRequestInfo) Value() string {
	return ""
}

func (m MessageRequestInfo) Bytes() []byte {
	var result []byte
	result = append(result, []byte(m.Kind())...)
	result = append(result, ' ')
	if m.Data != nil {
		result = append(result, []byte(m.Data.Kind())...)
		result = append(result, m.Data.Bytes()...)
	}
	return result
}

type MessageReplyInfoData interface {
}

type MessageReplyInfoDataImageInfoSet struct {
	Index          int
	Extension      string
	Name           string
	Fallback       string
	Width          int
	Height         int
	OtherExtension string
	Description    string
}

type MessageReplyInfoDataImageInfo struct {
	LastImageNumber int
	Checksum        int
	Sets            []MessageReplyInfoDataImageInfoSet
}

type MessageReplyInfo struct {
	Data MessageReplyInfoData
}

func (m MessageReplyInfo) Kind() string {
	return "replyinfo"
}

func (m *MessageReplyInfo) UnmarshalBinary(data []byte) error {
	parts := strings.Split(string(data), "\n")
	if len(parts) == 0 {
		return nil
	}
	switch parts[0] {
	case "image_info":
		if len(parts) < 3 {
			return fmt.Errorf("Not enough parts for image_info")
		}
		lastImageNumber, _ := strconv.Atoi(parts[1])
		checksum, _ := strconv.Atoi(parts[2])

		data := MessageReplyInfoDataImageInfo{
			LastImageNumber: lastImageNumber,
			Checksum:        checksum,
		}
		// Get our image sets.
		for i := 3; i < len(parts); i++ {
			if parts[i] == "" {
				continue
			}
			imageParts := strings.Split(parts[i], ":")
			debug.Debug(imageParts, len(imageParts))
			if len(imageParts) != 7 {
				return fmt.Errorf("Not enough parts for image_info image")
			}
			index, _ := strconv.Atoi(imageParts[0])
			extension := imageParts[1]
			name := imageParts[2]
			fallback := imageParts[3]
			geom := strings.Split(imageParts[4], "x")
			width, _ := strconv.Atoi(geom[0])
			height, _ := strconv.Atoi(geom[1])
			otherExtension := imageParts[5]
			description := imageParts[6]
			data.Sets = append(data.Sets, MessageReplyInfoDataImageInfoSet{
				Index:          index,
				Extension:      extension,
				Name:           name,
				Fallback:       fallback,
				Width:          width,
				Height:         height,
				OtherExtension: otherExtension,
				Description:    description,
			})
		}
		m.Data = data
		// TODO: image_sums, exp_table, knowledge_info, skill_info, skill_extra, spell_paths, race_list, race_info, class_list, class_info, startingmap, newcharinfo, news, rules, motd
	}
	return nil
}

func (m *MessageReplyInfo) Bytes() []byte {
	return nil
}

func (m MessageReplyInfo) Value() string {
	return ""
}

func init() {
	gMessages = append(gMessages, &MessageReplyInfo{})
}
