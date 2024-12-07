package play

import (
	"slices"
	"strings"
)

var GlobalObjectMapper ObjectMapper

type ObjectMapper struct {
	FaceToName map[uint16]string
	FaceToRune map[uint16]MapTile
	FaceToSize map[uint16]RuneSize
	Runes      []RuneDefinition
	Objects    []ObjectMap
}

func (m *ObjectMapper) Reset() {
	m.FaceToName = make(map[uint16]string)
	m.FaceToRune = make(map[uint16]MapTile)
	m.FaceToSize = make(map[uint16]RuneSize)
}

func (m *ObjectMapper) UnmarshalBinary(data []byte) error {
	// split by lines
	lines := strings.Split(string(data), "\n")
	var mode int
	for _, line := range lines {
		if strings.TrimSpace(string(line)) == "" {
			mode = 1
			continue
		} else if len(line) > 1 && line[0] == '/' && line[1] == '/' {
			continue
		}
		if mode == 0 {
			var r RuneDefinition
			r.UnmarshalBinary(line)
			m.Runes = append(m.Runes, r)
		} else if mode == 1 {
			var o ObjectMap
			o.UnmarshalBinary(line)
			m.Objects = append(m.Objects, o)
		}
	}
	return nil
}

func (m *ObjectMapper) GetRuneAndColors(name string) (rune, string, string) {
	// Get our best possible rune.
	var results []RuneDefinition
	exact := false
	for _, d := range m.Runes {
		for _, s := range d.Strings {
			if s == name {
				exact = true
				results = []RuneDefinition{d}
				break
			} else if strings.Contains(name, s) {
				results = append(results, d)
			}
		}
		if exact {
			break
		}
	}

	// Get our best possible name match.
	if len(results) > 1 {
		slices.SortFunc(results, func(a RuneDefinition, b RuneDefinition) int {
			// Get the closest string for a.
			var aClosest string
			var aClosestDistance int
			for _, s := range a.Strings {
				if strings.Contains(name, s) {
					if aClosest == "" {
						aClosest = s
						aClosestDistance = len(s) - len(name)
					} else {
						if len(s)-len(name) < aClosestDistance {
							aClosest = s
							aClosestDistance = len(s) - len(name)
						}
					}
				}
			}
			// Same for b.
			var bClosest string
			var bClosestDistance int
			for _, s := range b.Strings {
				if strings.Contains(name, s) {
					if bClosest == "" {
						bClosest = s
						bClosestDistance = len(s) - len(name)
					} else {
						if len(s)-len(name) < bClosestDistance {
							bClosest = s
							bClosestDistance = len(s) - len(name)
						}
					}
				}
			}
			return aClosestDistance - bClosestDistance
		})
	}

	var bestRuneWeight float64
	var bestRune rune
	var bestForegroundWeight float64
	var bestForeground string
	var bestBackgroundWeight float64
	var bestBackground string

	for _, d := range m.Objects {
		for _, s := range d.Strings {
			if strings.Contains(name, s) {
				if d.Rune.Weight >= bestRuneWeight {
					bestRuneWeight = d.Rune.Weight
					bestRune = d.Rune.Value
				}
				if d.Foreground.Weight >= bestForegroundWeight {
					bestForegroundWeight = d.Foreground.Weight
					bestForeground = d.Foreground.Color
				}
				if d.Background.Weight >= bestBackgroundWeight {
					bestBackgroundWeight = d.Background.Weight
					bestBackground = d.Background.Color
				}
			}
		}
	}

	if bestRune == 0 && len(results) > 0 {
		bestRune = results[0].Rune
	}

	return bestRune, bestForeground, bestBackground
}

func (m *ObjectMapper) GetRune(r rune) (RuneDefinition, bool) {
	for _, d := range m.Runes {
		if d.Rune == r {
			return d, true
		}
	}
	return RuneDefinition{}, false
}
