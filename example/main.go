package main

import (
	"fmt"
	"github.com/jojomi/brahms"
)

func main() {
	//brahms.NewScaleByName("C").Print()
	//brahms.NewScaleByName("a").Print()

	// Alle Tonleitern
	circle := brahms.NewCircleOfFifths()
	for _, index := range []uint8{0, 1, 2, 3, 4, 5, 6, 7} {
		name, err := circle.FindByKey(index, 0, true)
		if err != nil {
			panic(err)
		}
		scale := brahms.NewScaleByName(name)
		chords := scale.GetChords()
		scale.Print()

		for i, c := range chords {
			fmt.Println(i)
			for _, n := range c.GetNotes(0) {
				fmt.Println(n.Name)
			}
		}

		if index == 0 {
			continue
		}

		name, err = circle.FindByKey(0, index, true)
		if err != nil {
			panic(err)
		}
		brahms.NewScaleByName(name).Print()
	}
}
