package brahms

import (
	_ "fmt"
)

type Chord struct {
	BaseNote *Note
	Type     ChordType
	Notes    []*Note
}

type ChordType uint8

const (
	ChordMajor ChordType = iota
	ChordMinor
	ChordMajorSept
	ChordMinorSept
	ChordMin
)

const flatThird = 3
const majorThird = 4

var indices = map[ChordType][]HalftoneIndex{
	ChordMajor:     []HalftoneIndex{majorThird, flatThird},
	ChordMinor:     []HalftoneIndex{flatThird, majorThird},
	ChordMajorSept: []HalftoneIndex{majorThird, flatThird, flatThird},
	ChordMinorSept: []HalftoneIndex{flatThird, majorThird, flatThird},
	ChordMin:       []HalftoneIndex{flatThird, flatThird},
}

var indicesMinor = map[ChordType][]HalftoneIndex{}

func NewChordByName(baseNoteName string, chordType ChordType) *Chord {
	chord := &Chord{
		BaseNote: NewNoteByName(baseNoteName),
		Type:     chordType,
	}

	return chord
}

func NewChordByIndex(baseNoteIndex HalftoneIndex, chordType ChordType) *Chord {
	chord := &Chord{
		BaseNote: NewNoteByIndex(baseNoteIndex),
		Type:     chordType,
	}

	return chord
}

func (c *Chord) GetNotes(invertLevel int) []*Note {
	steps := indices[c.Type]
	notes := make([]*Note, len(steps)+1)
	note := NewNoteByIndex(c.BaseNote.HalftoneIndex)
	notes[0] = note
	for i, step := range steps {
		offset := note.HalftoneIndex + HalftoneIndex(step)
		note = NewNoteByIndex(offset)
		notes[i+1] = note
	}

	return notes
}

func (c *Chord) LongName(direction NoteNamingType) (longName string) {
	switch c.Type {
	case ChordMajor:
		longName = c.BaseNote.Name(direction, Major) + "-Dur"
	case ChordMinor:
		longName = c.BaseNote.Name(direction, Minor) + "-Moll"
	case ChordMajorSept:
		longName = c.BaseNote.Name(direction, Major) + "7"
	case ChordMinorSept:
		longName = c.BaseNote.Name(direction, Minor) + "7"
	case ChordMin:
		longName = c.BaseNote.Name(direction, Major) + "°"
	}
	return
}
