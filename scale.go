package brahms

import (
	"fmt"
	"strings"
)

type ScaleType uint8

const (
	Major ScaleType = iota
	Minor
)

type Scale struct {
	BaseNote *Note
	Notes    []*Note
	Flats    uint8
	Sharps   uint8
	Type     ScaleType
}

func NewScaleByNote(note *Note, scaleType ScaleType) *Scale {
	scale := &Scale{}

	scale.BaseNote = note
	scale.Type = scaleType

	cof := NewCircleOfFifths()
	cof.SelectByIndex(note.HalftoneIndex, scaleType)
	scale.Flats, scale.Sharps = cof.Flats(), cof.Sharps()

	scale.setNotes()

	return scale
}

func NewScaleByName(definition string) *Scale {
	scale := &Scale{}

	scale.BaseNote = NewNoteByName(strings.ToLower(definition))

	// Major or Minor key?
	firstLetter := string(definition[0])
	capitalized := strings.ToUpper(firstLetter) == firstLetter
	if capitalized {
		scale.Type = Major
	} else {
		scale.Type = Minor
	}

	// Flats and Sharps
	cof := NewCircleOfFifths()
	cof.SelectByName(definition)
	scale.Flats, scale.Sharps = cof.Flats(), cof.Sharps()

	scale.setNotes()

	return scale
}

func (s *Scale) setNotes() {
	// Notes
	steps := scaleSteps[s.Type]
	s.Notes = make([]*Note, len(steps))
	note := s.BaseNote
	s.Notes[0] = note
	for i, step := range steps {
		if i == len(steps)-1 {
			break
		}
		offset := note.HalftoneIndex + HalftoneIndex(step)
		note = NewNoteByIndex(offset)
		s.Notes[i+1] = note
	}
}

func (s *Scale) GetParallelScaleType() (scaleType ScaleType) {
	switch s.Type {
	case Major:
		scaleType = Minor
	case Minor:
		scaleType = Major
	}
	return
}

func (s *Scale) GetParallelOffset() (offset int8) {
	switch s.Type {
	case Major:
		offset = -flatThird
	case Minor:
		offset = flatThird
	}
	return
}

func (s *Scale) GetParallel() *Scale {
	offsetParallel := s.GetParallelOffset()
	newLevel := HalftoneIndex((int(s.BaseNote.HalftoneIndex) + 12 + int(offsetParallel)) % 12)
	return NewScaleByNote(NewNoteByIndex(newLevel), s.GetParallelScaleType())
}

func (s *Scale) GetChords() map[string]*Chord {
	// Tonika, Dominante, Subdominante
	var ctBase ChordType
	var ctOpposite ChordType
	var ctSept ChordType
	var ctMin ChordType
	var offsetParallel int8
	switch s.Type {
	case Major:
		ctBase = ChordMajor
		ctOpposite = ChordMinor
		ctSept = ChordMajorSept
		ctMin = ChordMin
		offsetParallel = -flatThird
	case Minor:
		ctBase = ChordMinor
		ctOpposite = ChordMajor
		ctSept = ChordMinorSept
		ctMin = ChordMin
		offsetParallel = flatThird
	}

	if s.Notes == nil {
		s.setNotes()
	}
	tonikaParallelIndex := (int(s.Notes[0].HalftoneIndex) + 12 + int(offsetParallel)) % 12
	subdominantParallelIndex := (int(s.Notes[3].HalftoneIndex) + 12 + int(offsetParallel)) % 12
	dominantParallelIndex := (int(s.Notes[4].HalftoneIndex) + 12 + int(offsetParallel)) % 12

	result := map[string]*Chord{
		"tonica":              NewChordByIndex(s.Notes[0].HalftoneIndex, ctBase),
		"subdominant":         NewChordByIndex(s.Notes[3].HalftoneIndex, ctBase),
		"dominant":            NewChordByIndex(s.Notes[4].HalftoneIndex, ctBase),
		"dominantsept":        NewChordByIndex(s.Notes[4].HalftoneIndex, ctSept),
		"tonicaparallel":      NewChordByIndex(HalftoneIndex(tonikaParallelIndex), ctOpposite),
		"subdominantparallel": NewChordByIndex(HalftoneIndex(subdominantParallelIndex), ctOpposite),
		"dominantparallel":    NewChordByIndex(HalftoneIndex(dominantParallelIndex), ctOpposite),
		"min":                 NewChordByIndex(s.Notes[6].HalftoneIndex, ctMin),
	}
	return result
}

func (s *Scale) NamingType() (naming NoteNamingType) {
	if s.Flats > 0 {
		naming = NamingFlat
	} else if s.Sharps > 0 {
		naming = NamingSharp
	}
	return
}

func (s *Scale) LongName() (longName string) {
	naming := s.NamingType()

	switch s.Type {
	case Major:
		longName = s.BaseNote.Name(naming, s.Type) + "-Dur"
	case Minor:
		longName = s.BaseNote.Name(naming, s.Type) + "-Moll"
	}
	return
}

func (s *Scale) Print() {
	noteStrings := make([]string, len(s.Notes))
	for i, note := range s.Notes {
		noteStrings[i] = note.Name(s.NamingType(), s.Type)
	}
	notes := strings.Join(noteStrings, ", ")
	keys := []string{}
	if s.Flats > 0 {
		keys = append(keys, fmt.Sprintf("%d flats", s.Flats))
	}
	if s.Sharps > 0 {
		keys = append(keys, fmt.Sprintf("%d sharps", s.Sharps))
	}
	keyString := strings.Join(keys, ", ")
	if len(keyString) > 0 {
		keyString = "(" + keyString + ") "
	}
	fmt.Printf("%s %s– %s\n", s.LongName(), keyString, notes)
}

const Half = 1
const Full = 2

var scaleSteps = map[ScaleType][]uint8{
	Major: {Full, Full, Half, Full, Full, Full, Half},
	Minor: {Full, Half, Full, Full, Half, Full, Full},
}
