package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
)

// -------------------------------------
type Vector struct {
	X, Y int
}

func VectorCmpLRTB(v1, v2 Vector) int {
	if v1.Y < v2.Y {
		return -1
	}
	if v1.Y > v2.Y {
		return 1
	}
	if v1.X < v2.X {
		return -1
	}
	if v1.X > v2.X {
		return 1
	}
	return 0
}

func VectorCmpTBLR(v1, v2 Vector) int {
	if v1.X < v2.X {
		return -1
	}
	if v1.X > v2.X {
		return 1
	}
	if v1.Y < v2.Y {
		return -1
	}
	if v1.Y > v2.Y {
		return 1
	}
	return 0
}

// -------------------------------------
type Cell struct {
	Crop     rune
	Visited  bool
	Location Vector
}

// -------------------------------------
type Region struct {
	Area int
	Crop rune

	// There is a fence above the cell at these Vector locations
	HFences map[Vector]bool

	// There is a fence below the cell at these Vector locations
	VFences map[Vector]bool
}

func NewRegion(crop rune) *Region {
	return &Region{Crop: crop, HFences: make(map[Vector]bool), VFences: make(map[Vector]bool)}
}

func (r *Region) NumSides() int {
	// sort the hfences and vfences
	hfences := slices.Collect(maps.Keys(r.HFences))
	slices.SortFunc(hfences, VectorCmpLRTB)

	vfences := slices.Collect(maps.Keys(r.VFences))
	slices.SortFunc(vfences, VectorCmpTBLR)

	// count the number of sides
	var sides int
	prev := Vector{-1, -1}
	for _, v := range hfences {
		if _, ok := r.VFences[v]; ok {
			sides++
		} else if prev.Y != v.Y || prev.X+1 != v.X {
			sides++
		}
		prev = v
	}

	prev = Vector{-1, -1}
	for _, v := range vfences {
		if _, ok := r.HFences[v]; ok {
			sides++
		} else if prev.X != v.X || prev.Y+1 != v.Y {
			sides++
		}
		prev = v
	}

	return sides
}

// -------------------------------------
type Board struct {
	Cells [][]Cell
	Size  Vector
}

func LoadBoard(filename string) *Board {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	b := new(Board)
	b.Cells = make([][]Cell, 0)

	scan := bufio.NewScanner(f)
	var y int
	for scan.Scan() {
		line := scan.Text()
		var x int
		var row []Cell
		for _, r := range line {
			cell := Cell{r, false, Vector{x, y}}
			row = append(row, cell)
			x++
		}
		b.Cells = append(b.Cells, row)
		y++
	}

	b.Size = Vector{len(b.Cells[0]), len(b.Cells)}
	return b
}

func (b *Board) GetCell(v Vector) *Cell {
	return &b.Cells[v.Y][v.X]
}

func (b *Board) GetCellXY(x, y int) *Cell {
	return &b.Cells[y][x]
}

func (b *Board) InBounds(v Vector) bool {
	return v.X >= 0 && v.X < b.Size.X && v.Y >= 0 && v.Y < b.Size.Y
}

func (b *Board) GetCrop(v Vector) rune {
	if !b.InBounds(v) {
		return '.'
	}
	return b.GetCell(v).Crop
}

// Returns the neighboring cells along with the neighbors that are out of bounds
func (b *Board) GetNeighbors(v Vector) ([]*Cell, []Vector) {
	var neighbors []*Cell
	var outOfBounds []Vector
	for _, d := range []Vector{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		v2 := Vector{v.X + d.X, v.Y + d.Y}
		if b.InBounds(v2) {
			neighbors = append(neighbors, b.GetCell(v2))
		} else {
			outOfBounds = append(outOfBounds, v2)
		}
	}
	return neighbors, outOfBounds
}

func (b *Board) Walk(v Vector, r *Region) {
	me := b.GetCell(v)

	me.Visited = true
	r.Area++

	// Up
	v2 := Vector{v.X, v.Y - 1}
	if b.InBounds(v2) {
		c := b.GetCell(v2)
		if me.Crop != c.Crop {
			r.HFences[v] = true
		} else if !c.Visited {
			b.Walk(v2, r)
		}
	} else {
		r.HFences[v] = true
	}

	// Down
	v2 = Vector{v.X, v.Y + 1}
	if b.InBounds(v2) {
		c := b.GetCell(v2)
		if me.Crop != c.Crop {
			r.HFences[v2] = true
		} else if !c.Visited {
			b.Walk(v2, r)
		}
	} else {
		r.HFences[v2] = true
	}

	// Left
	v2 = Vector{v.X - 1, v.Y}
	if b.InBounds(v2) {
		c := b.GetCell(v2)
		if me.Crop != c.Crop {
			r.VFences[v] = true
		} else if !c.Visited {
			b.Walk(v2, r)
		}
	} else {
		r.VFences[v] = true
	}

	// Right
	v2 = Vector{v.X + 1, v.Y}
	if b.InBounds(v2) {
		c := b.GetCell(v2)
		if me.Crop != c.Crop {
			r.VFences[v2] = true
		} else if !c.Visited {
			b.Walk(v2, r)
		}
	} else {
		r.VFences[v2] = true
	}
}

// -------------------------------------
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	b := LoadBoard("input.txt")

	var cost int
	for y, row := range b.Cells {
		for x, cell := range row {
			if !cell.Visited {
				r := NewRegion(cell.Crop)
				b.Walk(Vector{x, y}, r)
				nSides := r.NumSides()
				fmt.Println(string(r.Crop), r.Area, nSides)
				cost += r.Area * nSides
			}
		}
	}

	fmt.Println(cost)
}
