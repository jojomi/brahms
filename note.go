package brahms

import (
	"strings"
)

type HalftoneIndex uint8

type Note struct {
	HalftoneIndex HalftoneIndex
}

func NewNoteByName(name string) *Note {
	note := &Note{}
	note.HalftoneIndex = NoteToIndex(name)
	return note
}

func NewNoteByIndex(index HalftoneIndex) *Note {
	note := &Note{}
	note.HalftoneIndex = index
	return note
}

type NoteNamingType uint8

const (
	NamingFlat NoteNamingType = iota
	NamingSharp
)

func (n *Note) Name(namingType NoteNamingType, scaleType ScaleType) string {
	name := IndexToNote(n.HalftoneIndex, namingType)
	if scaleType == Major {
		name = strings.ToUpper(string(name[0])) + string(name[1:])
	}
	return name
}

var noteTypeOffsets = map[string]HalftoneIndex{
	"c":   0,
	"cis": 1,
	"des": 1,
	"d":   2,
	"dis": 3,
	"es":  3,
	"e":   4,
	"f":   5,
	"fis": 6,
	"ges": 6,
	"g":   7,
	"gis": 8,
	"as":  8,
	"a":   9,
	"b":   10,
	"ais": 10,
	"hes": 10,
	"h":   11,
	"ces": 11,
}

var noteOffsetTypes = map[NoteNamingType]map[HalftoneIndex]string{
	NamingFlat: {
		0:  "c",
		1:  "des",
		2:  "d",
		3:  "es",
		4:  "e",
		5:  "f",
		6:  "ges",
		7:  "g",
		8:  "as",
		9:  "a",
		10: "b",
		11: "h",
	},
	NamingSharp: {
		0:  "c",
		1:  "cis",
		2:  "d",
		3:  "dis",
		4:  "e",
		5:  "f",
		6:  "fis",
		7:  "g",
		8:  "gis",
		9:  "a",
		10: "ais",
		11: "h",
	},
}

// c -> 0, c0 -> 0, d0 -> 2, cis1 -> 12
func NoteToIndex(note string) HalftoneIndex {
	return noteTypeOffsets[strings.ToLower(note)]
}

func IndexToNote(index HalftoneIndex, namingType NoteNamingType) string {
	return noteOffsetTypes[namingType][index%HalftoneIndex(len(noteOffsetTypes[namingType]))]
}
