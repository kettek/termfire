package play

import (
	"bytes"
	"slices"
	"strconv"
	"strings"
)

var GlobalObjectMapper ObjectMapper

type ObjectMapper struct {
	Runes   []RuneDefinition
	Objects []ObjectMap
}

func (m *ObjectMapper) UnmarshalBinary(data []byte) error {
	// split by lines
	lines := bytes.Split(data, []byte("\n"))
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

type RuneDefinition struct {
	Rune    rune
	Strings []string
}

func (t *RuneDefinition) UnmarshalBinary(data []byte) error {
	d := string(data)
	var start int
	var mode int
	for i := 0; i < len(d); i++ {
		switch mode {
		case 0: // rune (first space)
			if d[i] == ' ' && i > 0 {
				t.Rune = rune(d[i-1])
				mode = 1
				start = i + 1
			}
		case 1: // strings
			if d[i] == ' ' || i == len(d)-1 {
				t.Strings = append(t.Strings, d[start:i])
				start = i + 1
			}
		}
	}
	return nil
}

type ObjectMap struct {
	Strings []string
	Rune    struct {
		Weight float64
		Value  rune
	}
	Foreground struct {
		Weight float64
		Color  string
	}
	Background struct {
		Weight float64
		Color  string
	}
}

func (t *ObjectMap) UnmarshalBinary(data []byte) error {
	d := string(data)
	var start int
	var mode int
	var submode int
	for i := 0; i < len(d); i++ {
		switch mode {
		case 0: // strings
			if d[i] == ' ' {
				t.Strings = append(t.Strings, d[start:i])
				start = i + 1
				mode = 1
			} else if d[i] == ',' {
				t.Strings = append(t.Strings, d[start:i])
				start = i + 1
			}
		case 1: // foreground,weight
			switch submode {
			case 0: // foreground
				if d[i] == ',' || i == len(d)-1 {
					t.Foreground.Color = d[start:i]
					start = i + 1
					submode = 1
				} else if d[i] == ' ' || i == len(d)-1 {
					t.Foreground.Color = d[start:i]
					start = i + 1
					submode = 0
					mode = 2
				}
			case 1: // weight
				if d[i] == ' ' || i == len(d)-1 {
					w, err := strconv.ParseFloat(d[start:i], 64)
					if err != nil {
						return err
					}
					t.Foreground.Weight = w
					start = i + 1
					submode = 0
					mode = 2
				}
			}
		case 2: // background,weight
			switch submode {
			case 0: // background
				if d[i] == ',' {
					t.Background.Color = d[start:i]
					start = i + 1
					submode = 1
				} else if d[i] == ' ' || i == len(d)-1 {
					t.Background.Color = d[start:i]
					start = i + 1
					submode = 0
					mode = 3
				}
			case 1: // weight
				if d[i] == ' ' || i == len(d)-1 {
					w, err := strconv.ParseFloat(d[start:i], 64)
					if err != nil {
						return err
					}
					t.Background.Weight = w
					start = i + 1
					submode = 0
					mode = 3
				}
			}
		case 3: // rune,weight
			switch submode {
			case 0: // rune
				if d[i] == ',' || i == len(d)-1 {
					t.Rune.Value = rune(d[start:i][0])
					start = i + 1
					submode = 1
				} else if d[i] == ' ' || i == len(d)-1 {
					t.Rune.Value = rune(d[start:i][0])
					start = i + 1
					submode = 0
				}
			case 1: // weight
				if d[i] == ' ' || i == len(d)-1 {
					w, err := strconv.ParseFloat(d[start:i], 64)
					if err != nil {
						return err
					}
					t.Rune.Weight = w
					start = i + 1
					submode = 0
				}
			}
		}
	}
	return nil
}

/**
"wall":
	rune:
		weight: 0.5
		value: @
	fg:
		weight: 0.5
		color: white
	bg:
		weight: 0.5
		color: transparent

if our own parser...

$:coin,gold,money
#:wall,store,whatever
/:switch,lever

wall,1;white,1;transparent,1
window,1;white,1;transparent,1
floor;white;transparent
# colors
wood,0;brown,2;dark brown,2
stone,0;gray,2;dark gray,2

<rune>:<string>,<string>,<string>

----

<string>:<rune>,<weight>;<foreground>,<weight>;<background>,<weight>

// result => wooden wall = # brown fg, dark brown bg
// result => stone wall = # gray fg, dark gray fg

**/
