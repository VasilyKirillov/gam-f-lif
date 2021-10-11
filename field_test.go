package main

import "testing"

//Test check that field size is correct and all elements has zero value.
func TestGenerateField(t *testing.T) {
	w := 10
	h := 5
	f := generateField(w, h)
	if len(f) != w || len(f[0]) != h {
		t.Fatalf("Generated Field' size not equal to %vx%x", w, h)
	}
	var bi interface{}
	for i, r := range f {
		for j, c := range r {
			bi = c
			b, ok := bi.(bool)
			if !ok || b {
				t.Fatalf("value f[%v][%v] isn't false", i, j)
			}
		}
	}
}

//Test checks that counting naighbours works correctly.
//If one of the neighbour indeces is out of bounds, value should be taken from the other side of the Feild.
func TestCountNeighbours(t *testing.T) {
	f := Field{
		{false, false, true},
		{false, false, false},
		{false, false, true},
		{false, false, true},
	}
	nb := f.countNeighbours(3, 1)
	if nb != 3 {
		t.Fatalf("Expected 3 neighbours but got %v", nb)
	}
}
