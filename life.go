package main

import (
	"encoding/csv"
	"log"
	"math/rand"
	"os"
	"time"

	tm "github.com/buger/goterm"
	kb "github.com/eiannone/keyboard"
)

const (
	FIELD_WIDTH  = 10
	FIELD_HEIGHT = FIELD_WIDTH
	FRAME_RATE   = 1
)

type Field [][]bool

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
	kb.Close()
	os.Exit(0)
}

func main() {
	tm.Clear()
	go waitForEsc()

	fieldA := generateField(FIELD_WIDTH, FIELD_HEIGHT)
	fieldA.populate()
	fieldB := generateField(len(fieldA), len(fieldA[0]))

	isAturn := true
	gen := 0
	for {
		if isAturn {
			fieldA.print()
			fieldA.update(fieldB)
		} else {
			fieldB.print()
			fieldB.update(fieldA)
		}
		isAturn = !isAturn
		gen++
		tm.Printf("Generation %v\n", gen)
		time.Sleep(time.Second / FRAME_RATE)
	}
}

func (field Field) print() {
	tm.MoveCursor(1, 1)
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

func (field Field) populate() {
	for _, row := range field {
		for i, _ := range row {
			row[i] = rand.Int()%2 == 0
		}
	}
}

func generateField(w, h int) (field Field) {
	field = make([][]bool, w)
	rows := make([]bool, w*h)
	for i := 0; i < w; i++ {
		field[i] = rows[i*h : (i+1)*h]
	}
	return
}

func readFromCsv(filePath string) (field Field) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+filePath, err)
	}

	rows := len(records)
	cols := len(records[0])
	field = generateField(rows, cols)

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if records[i][j] != "" {
				field[i][j] = true
			}
		}
	}
	tm.Println("readed field:")
	tm.Println(field)
	time.Sleep(time.Second + 10)
	return
}

func (fieldA Field) update(fieldB Field) {
	for i := 0; i < len(fieldA); i++ {
		for j := 0; j < len(fieldA[0]); j++ {
			nb := fieldA.calculateNeibours(i, j)
			// Cell is lonely and dies or dies due to over population
			if fieldA[i][j] && (nb < 2 || nb > 3) {
				fieldB[i][j] = false
				// A new cell is born
			} else if !fieldA[i][j] && nb == 3 {
				fieldB[i][j] = true
			} else {
				// Remains the same
				fieldB[i][j] = fieldA[i][j]
			}
		}
	}
}

func (field Field) calculateNeibours(row, col int) (nb int) {
	w := len(field)
	h := len(field[0])
	var rn, cn int
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			rn = row + i
			cn = col + j
			if rn < 0 {
				rn += w
			}
			if rn >= w {
				rn -= w
			}
			if cn < 0 {
				cn += h
			}
			if cn >= h {
				cn -= h
			}
			if field[rn][cn] {
				nb++
			}
		}
	}
	return
}
