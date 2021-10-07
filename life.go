package main

import (
	"log"
	"os"
	"time"

	tm "github.com/buger/goterm"
	kb "github.com/eiannone/keyboard"
)

const (
	FRAME_RATE = 1
)

func waitForEsc() {
	keysEvents, err := kb.GetKeys(10)
	if err != nil {
		log.Fatal("Unable to get keyboard input", err)
	}
	for {
		event := <-keysEvents
		if event.Err != nil {
			log.Fatal("KeyEvent error", event.Err)
		}
		if event.Key == kb.KeyEsc {
			break
		}
	}
	defer os.Exit(0)
	defer kb.Close()
}

func main() {
	tm.Clear()
	go waitForEsc()

	//fieldA := generateField(FIELD_WIDTH, FIELD_HEIGHT)
	//fieldA.populate()
	fieldA := readFromCsv("field2.csv")
	fieldB := generateField(len(fieldA), len(fieldA[0]))

	isAturn := true
	gen := 0
	var err error
	for {
		if isAturn {
			fieldA.print(&gen)
			err = fieldA.update(fieldB)
		} else {
			fieldB.print(&gen)
			err = fieldB.update(fieldA)
		}
		if err != nil {
			break
		}
		isAturn = !isAturn
		time.Sleep(time.Second / FRAME_RATE)
		gen++
	}
}

func (field Field) print(gen *int) {
	tm.MoveCursor(1, 1)
	tm.Printf("Generation %v\n", *gen)
	for i := 0; i < len(field); i++ {
		for j := 0; j < len(field[0]); j++ {
			if field[i][j] {
				tm.Print("\u25A0 ")
			} else {
				tm.Print("\u25A1 ")
			}
		}
		tm.Println()
	}
	tm.Flush()
}
