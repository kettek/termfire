package messages

import (
	"fmt"
	"strconv"
	"strings"
)

// TODO: image_sums, exp_table, knowledge_info, skill_info, skill_extra, spell_paths, race_list, race_info, class_list, class_info, startingmap, newcharinfo
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

type MessageRequestInfoNews struct{}

func (m MessageRequestInfoNews) Kind() string {
	return "news"
}

func (m MessageRequestInfoNews) Bytes() []byte {
	return nil
}

type MessageRequestInfoRules struct{}

func (m MessageRequestInfoRules) Kind() string {
	return "rules"
}

func (m MessageRequestInfoRules) Bytes() []byte {
	return nil
}

type MessageRequestInfoMotd struct{}

func (m MessageRequestInfoMotd) Kind() string {
	return "motd"
}

func (m MessageRequestInfoMotd) Bytes() []byte {
	return nil
}

type MessageRequestInfoRaceList struct{}

func (m MessageRequestInfoRaceList) Kind() string {
	return "race_list"
}

func (m MessageRequestInfoRaceList) Bytes() []byte {
	return nil
}

type MessageRequestInfoRaceInfo string

func (m MessageRequestInfoRaceInfo) Kind() string {
	return "race_info"
}

func (m MessageRequestInfoRaceInfo) Bytes() []byte {
	return []byte(m)
}

type MessageRequestInfoClassList struct{}

func (m MessageRequestInfoClassList) Kind() string {
	return "class_list"
}

func (m MessageRequestInfoClassList) Bytes() []byte {
	return nil
}

type MessageRequestInfoClassInfo string

func (m MessageRequestInfoClassInfo) Kind() string {
	return "class_info"
}

func (m MessageRequestInfoClassInfo) Bytes() []byte {
	return []byte(m)
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
		bytes := m.Data.Bytes()
		if len(bytes) > 0 {
			result = append(result, ' ')
			result = append(result, bytes...)
		}
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

type MessageReplyInfoDataNews string

type MessageReplyInfoDataRules string

type MessageReplyInfoDataMotd string

type MessageReplyInfoDataRaceList []string

type MessageReplyInfoDataClassList []string

type MessageReplyInfoDataRaceOrClassInfo struct {
	Arch        string
	Name        string
	Description string
	Stats       []MessageStat
	Choices     []Choice
}

type MessageReplyInfoDataRaceInfo MessageReplyInfoDataRaceOrClassInfo
type MessageReplyInfoDataClassInfo MessageReplyInfoDataRaceOrClassInfo

func (msg *MessageReplyInfoDataRaceOrClassInfo) UnmarshalBinary(data []byte) (int, error) {
	adjust := 0
	// Race is until newline.
	for i, b := range data {
		if b == '\n' {
			msg.Arch = string(data[:i])
			data = data[i+1:]
			adjust += i + 1
			break
		}
	}
	done := false
	for !done {
		var kind string
		for i, b := range data {
			if b == ' ' {
				kind = string(data[:i])
				data = data[i+1:]
				adjust += i + 1
				break
			}
		}
		offset := 0
		switch kind {
		case "name":
			msg.Name, offset = readLengthPrefixedString(data, offset)
			data = data[offset:]
			adjust += offset
		case "msg":
			msg.Description, offset = readLengthPrefixedString2(data, offset)
			data = data[offset:]
			adjust += offset
		case "stats":
			for i := 0; i < len(data); i++ {
				kind := data[i]
				if kind == 0 { // 0 signifies done.
					data = data[i+1:]
					adjust += i + 1
					break
				}
				// Re-use stat message processing, I suppose.
				for _, s := range gMessageStats {
					if s.Matches(kind) {
						count, err := s.UnmarshalBinary(data)
						if err != nil {
							return adjust, err
						}
						i += count
						adjust += count
						msg.Stats = append(msg.Stats, s)
						break
					}
				}
			}
		case "choice":
			choice := Choice{}
			// Documentation for this in protocols.txt is a rather bad to parse out (I know a certain someone would disagree), but it seems to be this.
			choice.Name, offset = readLengthPrefixedString(data, offset)
			choice.Description, offset = readLengthPrefixedString(data, offset)
			// Loop reading each "arch", then check next byte for 0.
			for i := 0; i < len(data); i++ {
				option := [2]string{}
				option[0], offset = readLengthPrefixedString(data, offset)
				option[1], offset = readLengthPrefixedString(data, offset)
				choice.Options = append(choice.Options, option)
				if data[offset] == 0 {
					offset++
					adjust += offset
					break
				}
			}
			data = data[offset:]
			adjust += offset
			msg.Choices = append(msg.Choices, choice)
		default:
			done = true
			break
		}
	}
	return adjust, nil
}

type Choice struct {
	Name        string
	Description string
	Options     [][2]string // Name and Description pair
}

type MessageReplyInfo struct {
	Data MessageReplyInfoData
}

func (m MessageReplyInfo) Kind() string {
	return "replyinfo"
}

func (m *MessageReplyInfo) UnmarshalBinary(data []byte) error {
	// Read info type.
	var infoType string
	for i, b := range data {
		if b == ' ' || b == '\n' {
			infoType = string(data[:i])
			data = data[i+1:]
			break
		}
	}
	switch infoType {
	case "image_info":
		parts := strings.Split(string(data), "\n")
		if len(parts) < 3 {
			return fmt.Errorf("Not enough parts for image_info")
		}
		lastImageNumber, _ := strconv.Atoi(parts[0])
		checksum, _ := strconv.Atoi(parts[1])

		data := MessageReplyInfoDataImageInfo{
			LastImageNumber: lastImageNumber,
			Checksum:        checksum,
		}
		// Get our image sets.
		for i := 2; i < len(parts); i++ {
			if parts[i] == "" {
				continue
			}
			imageParts := strings.Split(parts[i], ":")
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
	case "news":
		m.Data = MessageReplyInfoDataNews(data)
	case "rules":
		m.Data = MessageReplyInfoDataRules(data)
	case "motd":
		m.Data = MessageReplyInfoDataMotd(data)
	case "race_list":
		// NOTE: delimiter is "|", not ":", unlike what the CF protocol file says!
		races := strings.Split(string(data), "|")
		// Skip the first, because for whatever reason we start with "|"
		races = races[1:]
		m.Data = MessageReplyInfoDataRaceList(races)
	case "race_info":
		msg := MessageReplyInfoDataRaceOrClassInfo{}
		_, err := msg.UnmarshalBinary(data)
		if err != nil {
			return err
		}
		m.Data = MessageReplyInfoDataRaceInfo(msg)
	case "class_list":
		classes := strings.Split(string(data), "|")
		classes = classes[1:]
		m.Data = MessageReplyInfoDataClassList(classes)
	case "class_info":
		msg := MessageReplyInfoDataRaceOrClassInfo{}
		_, err := msg.UnmarshalBinary(data)
		if err != nil {
			return err
		}
		m.Data = MessageReplyInfoDataClassInfo(msg)
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
