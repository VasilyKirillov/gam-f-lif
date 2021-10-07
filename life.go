package main

import (
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
		panic(err)
	}
	for {
		event := <-keysEvents
		if event.Err != nil {
			panic(event.Err)
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
	fieldB := generateField(FIELD_WIDTH, FIELD_HEIGHT)
	isAturn := true
	fieldA.populate()
	for {
		if isAturn {
			fieldA.print()
			fieldA.update(fieldB)
		} else {
			fieldB.print()
			fieldB.update(fieldA)
		}
		isAturn = !isAturn
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

func readFromCsv() (field Field) {
	panic("not implemented")
}

func (fieldA Field) update(fieldB Field) Field {
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
	return fieldB
}

func (field Field) calculateNeibours(row, col int) (nb int) {
	nb = 0
	w := len(field)
	h := len(field[0])
	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			if i == 0 && j == 0 {
				continue
			}
			rn := row + i
			cn := col + j
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
