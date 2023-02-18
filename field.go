package main

import (
	"encoding/csv"
	"errors"
	"log"
	"math/rand"
	"os"
)

const (
	FIELD_WIDTH  = 35
	FIELD_HEIGHT = FIELD_WIDTH
)

type Field [][]bool

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
	return
}

func (fieldA Field) update(fieldB Field) error {
	for i := 0; i < len(fieldA); i++ {
		for j := 0; j < len(fieldA[0]); j++ {
			nb := fieldA.countNeighbours(i, j)
			if fieldA[i][j] && (nb < 2 || nb > 3) {
				// Cell is lonely and dies or dies due to over population
				fieldB[i][j] = false
			} else if !fieldA[i][j] && nb == 3 {
				// A new cell is born
				fieldB[i][j] = true
			} else {
				// Remains the same
				fieldB[i][j] = fieldA[i][j]
			}
		}
	}
	if fieldA.equals(fieldB) {
		return errors.New("update fails, there are no changes")
	}
	return nil
}

func (a Field) equals(b Field) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		for j := 0; j < len(a[0]); j++ {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func (field Field) countNeighbours(row, col int) (nb int) {
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
