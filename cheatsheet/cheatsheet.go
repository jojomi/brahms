package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jojomi/brahms"
	"github.com/jojomi/go-latex"
)

var verbose bool

func main() {
	var rootCmd = &cobra.Command{
		Use: "app",
		Run: runCheatsheet,
	}
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.Execute()
}

type ChordData struct {
	Name  string
	Notes string
}

type ScaleData struct {
	Name                string
	Flats               string
	Sharps              string
	Notes               string
	NoteStyle           string
	ParallelScale       string
	Tonika              ChordData
	Subdominant         ChordData
	Dominant            ChordData
	Dominantsept        ChordData
	Tonikaparallel      ChordData
	Subdominantparallel ChordData
	Dominantparallel    ChordData
	Min                 ChordData
	ScaleNotes          map[string]string
}

func getScaleData(scale *brahms.Scale) *ScaleData {
	result := &ScaleData{}
	result.Name = scale.LongName()
	result.Flats = strings.Repeat("$\\flat$", int(scale.Flats))
	result.Sharps = strings.Repeat("$\\sharp$", int(scale.Sharps))
	result.Notes = getNoteHint(scale.Notes)
	result.NoteStyle = "size=1.0em,height=2.20,ratio=0.60"
	result.ParallelScale = scale.GetParallel().LongName()

	chords := scale.GetChords()
	result.Tonika.Name = chords["tonica"].LongName(scale.NamingType())
	result.Tonika.Notes = getNoteHint(chords["tonica"].GetNotes(0))
	result.Subdominant.Name = chords["subdominant"].LongName(scale.NamingType())
	result.Subdominant.Notes = getNoteHint(chords["subdominant"].GetNotes(0))
	result.Dominant.Name = chords["dominant"].LongName(scale.NamingType())
	result.Dominant.Notes = getNoteHint(chords["dominant"].GetNotes(0))

	result.Dominantsept.Name = chords["dominantsept"].LongName(scale.NamingType())
	result.Dominantsept.Notes = getNoteHint(chords["dominantsept"].GetNotes(0))
	result.Tonikaparallel.Name = chords["tonicaparallel"].LongName(scale.NamingType())
	result.Tonikaparallel.Notes = getNoteHint(chords["tonicaparallel"].GetNotes(0))
	result.Subdominantparallel.Name = chords["subdominantparallel"].LongName(scale.NamingType())
	result.Subdominantparallel.Notes = getNoteHint(chords["subdominantparallel"].GetNotes(0))
	result.Dominantparallel.Name = chords["dominantparallel"].LongName(scale.NamingType())
	result.Dominantparallel.Notes = getNoteHint(chords["dominantparallel"].GetNotes(0))
	result.Min.Name = strings.Replace(chords["min"].LongName(scale.NamingType()), "°", "${}^{\\circ}$", -1)
	result.Min.Notes = getNoteHint(chords["min"].GetNotes(0))

	result.ScaleNotes = make(map[string]string, len(scale.Notes))
	for i, note := range scale.Notes {
		result.ScaleNotes["no"+strconv.Itoa(i+1)] = note.Name(scale.NamingType(), scale.Type)
	}
	return result
}

func getNoteHint(notes []*brahms.Note) (result string) {
	noteOffsets := make([]string, len(notes))

	minOffset := 10000
	for _, note := range notes {
		if int(note.HalftoneIndex) < minOffset {
			minOffset = int(note.HalftoneIndex)
		}
	}
	baseOffset := minOffset - (minOffset % 12)

	for i, note := range notes {
		intVal := int(note.HalftoneIndex) - baseOffset
		if intVal < 12 {
			noteOffsets[i] = strconv.Itoa(intVal)
		} else {
			noteOffsets[i] = strconv.Itoa(intVal-12) + "'"
		}
	}
	result = strings.Join(noteOffsets, ",")
	return
}

func runCheatsheet(cmd *cobra.Command, args []string) {
	data := []*ScaleData{}

	circle := brahms.NewCircleOfFifths()
	for _, scaleType := range []brahms.ScaleType{brahms.Major, brahms.Minor} {
		for _, index := range []uint8{0, 1, 2, 3, 4, 5, 6, 7} {
			circle.SelectByKey(index, 0, scaleType)
			scale := circle.GetScale()
			data = append(data, getScaleData(scale))

			if index == 0 {
				continue
			}
			circle.SelectByKey(0, index, scaleType)
			scale = circle.GetScale()
			data = append(data, getScaleData(scale))
		}
	}

	// generate the document
	basePath := "./"
	l := latex.NewCompileTask()
	l.SetSourceDir(basePath)
	l.SetCompileFilename("cheatsheet-template")
	l.SetResolveSymlinks(false)
	l.CopyToCompileDir("")
	defer os.RemoveAll(l.CompileDir())
	fmt.Printf("Building Cheatsheets...\n")
	t, mainFile := l.Template("")
	t.ParseFiles(mainFile)
	l.ExecuteTemplate(t, data, "", "")
	for i := 0; i < 1; i++ {
		l.Lualatex("")
	}
	l.Optimize("", "printer")
	outPath, _ := filepath.Abs(path.Join(basePath, "cheatsheet.pdf"))
	err := l.MoveToDest("", outPath)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Copied to %s.\n", outPath)
}
