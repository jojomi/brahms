package brahms

import (
	_ "fmt"
	"strings"
)

type CircleOfFifths struct {
	BaseNote  *Note
	ScaleType ScaleType
}

func NewCircleOfFifths() *CircleOfFifths {
	cof := &CircleOfFifths{}
	return cof
}

func (c CircleOfFifths) GetScale() *Scale {
	return NewScaleByNote(c.BaseNote, c.ScaleType)
}

var circleSteps = map[ScaleType]map[string]int8{
	Major: {
		"Ces": -7,
		"Ges": -6,
		"Des": -5,
		"As":  -4,
		"Es":  -3,
		"B":   -2,
		"F":   -1,
		"C":   0,
		"G":   1,
		"D":   2,
		"A":   3,
		"E":   4,
		"H":   5,
		"Fis": 6,
		"Cis": 7,
	},
	Minor: {
		"as":  -7,
		"es":  -6,
		"b":   -5,
		"f":   -4,
		"c":   -3,
		"g":   -2,
		"d":   -1,
		"a":   0,
		"e":   1,
		"h":   2,
		"fis": 3,
		"cis": 4,
		"gis": 5,
		"dis": 6,
		"ais": 7,
	},
}

func (c *CircleOfFifths) getPosition() (position int8) {
	scalePositions := circleSteps[c.ScaleType]
	for _, n := range []NoteNamingType{NamingFlat, NamingSharp} {
		if val, ok := scalePositions[c.BaseNote.Name(n, c.ScaleType)]; ok {
			position = val
			return
		}
	}
	return
}

func (c *CircleOfFifths) Flats() (flats uint8) {
	position := c.getPosition()
	if position >= 0 {
		return 0
	}
	return uint8(-position)
}

func (c *CircleOfFifths) Sharps() (sharps uint8) {
	position := c.getPosition()
	if position <= 0 {
		return 0
	}
	return uint8(position)
}

func (c *CircleOfFifths) SelectByName(name string) {
	firstLetter := string(name[0])
	major := strings.ToUpper(firstLetter) == firstLetter

	c.BaseNote = NewNoteByName(name)
	if major {
		c.ScaleType = Major
	} else {
		c.ScaleType = Minor
	}
}

func (c *CircleOfFifths) SelectByKey(flats, sharps uint8, scaleType ScaleType) {
	subcircle := circleSteps[c.ScaleType]

	for key, val := range subcircle {
		if (val > 0 && val == int8(flats)) || (val < 0 && -val == int8(sharps)) {
			c.BaseNote = NewNoteByName(key)
			c.ScaleType = scaleType
			return
		}
		if val == 0 && flats == 0 && val == 0 && sharps == 0 {
			c.BaseNote = NewNoteByName(key)
			c.ScaleType = scaleType
			return
		}
	}

}

func (c *CircleOfFifths) SelectByIndex(index HalftoneIndex, scaleType ScaleType) {
	c.BaseNote = NewNoteByIndex(index)
	c.ScaleType = scaleType
}
