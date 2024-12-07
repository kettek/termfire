package play

import (
	"strconv"
)

type RuneDefinition struct {
	identifier int // interally used
	Rune       rune
	Strings    []string
}

func (t *RuneDefinition) UnmarshalBinary(data string) error {
	var start int
	var mode int
	var lastRune rune
	for i, r := range data {
		switch mode {
		case 0: // rune (first space)
			if r == ' ' && i > 0 {
				t.Rune = lastRune
				mode = 1
				start = i + 1
			}
		case 1: // strings
			if r == ' ' {
				t.Strings = append(t.Strings, data[start:i])
				start = i + 1
			} else if i == len(data)-1 {
				t.Strings = append(t.Strings, data[start:i+1])
			}
		}
		lastRune = r
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

func (t *ObjectMap) UnmarshalBinary(data string) error {
	var start int
	var mode int
	var submode int
	for i := 0; i < len(data); i++ {
		switch mode {
		case 0: // strings
			if data[i] == ' ' {
				t.Strings = append(t.Strings, data[start:i])
				start = i + 1
				mode = 1
			} else if data[i] == ',' {
				t.Strings = append(t.Strings, data[start:i])
				start = i + 1
			}
		case 1: // foreground,weight
			switch submode {
			case 0: // foreground
				if data[i] == ',' {
					t.Foreground.Color = data[start:i]
					start = i + 1
					submode = 1
				} else if data[i] == ' ' {
					t.Foreground.Color = data[start:i]
					start = i + 1
					submode = 0
					mode = 2
				} else if i == len(data)-1 {
					t.Foreground.Color = data[start : i+1]
				}
			case 1: // weight
				if data[i] == ' ' {
					w, err := strconv.ParseFloat(data[start:i], 64)
					if err != nil {
						return err
					}
					t.Foreground.Weight = w
					start = i + 1
					submode = 0
					mode = 2
				} else if i == len(data)-1 {
					w, err := strconv.ParseFloat(data[start:i+1], 64)
					if err != nil {
						return err
					}
					t.Foreground.Weight = w
				}
			}
		case 2: // background,weight
			switch submode {
			case 0: // background
				if data[i] == ',' {
					t.Background.Color = data[start:i]
					start = i + 1
					submode = 1
				} else if data[i] == ' ' {
					t.Background.Color = data[start:i]
					start = i + 1
					submode = 0
					mode = 3
				} else if i == len(data)-1 {
					t.Background.Color = data[start : i+1]
				}
			case 1: // weight
				if data[i] == ' ' {
					w, err := strconv.ParseFloat(data[start:i], 64)
					if err != nil {
						return err
					}
					t.Background.Weight = w
					start = i + 1
					submode = 0
					mode = 3
				} else if i == len(data)-1 {
					w, err := strconv.ParseFloat(data[start:i+1], 64)
					if err != nil {
						return err
					}
					t.Background.Weight = w
				}
			}
		case 3: // rune,weight
			switch submode {
			case 0: // rune
				if data[i] == ',' {
					t.Rune.Value = rune(data[start:i][0])
					start = i + 1
					submode = 1
				} else if data[i] == ' ' {
					t.Rune.Value = rune(data[start:i][0])
					start = i + 1
					submode = 0
				} else if i == len(data)-1 {
					t.Rune.Value = rune(data[start : i+1][0])
				}
			case 1: // weight
				if data[i] == ' ' {
					w, err := strconv.ParseFloat(data[start:i], 64)
					if err != nil {
						return err
					}
					t.Rune.Weight = w
					start = i + 1
					submode = 0
				} else if i == len(data)-1 {
					w, err := strconv.ParseFloat(data[start:i+1], 64)
					if err != nil {
						return err
					}
					t.Rune.Weight = w
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
