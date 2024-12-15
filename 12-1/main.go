package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Vector struct {
	X, Y int
}

type Cell struct {
	Crop     rune
	Visited  bool
	Location Vector
}

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

// Returns the neighboring cells along with the number of edges that are out of bounds
func (b *Board) GetNeighbors(v Vector) ([]*Cell, int) {
	var neighbors []*Cell
	var outOfBounds int
	for _, d := range []Vector{{0, 1}, {0, -1}, {1, 0}, {-1, 0}} {
		v2 := Vector{v.X + d.X, v.Y + d.Y}
		if b.InBounds(v2) {
			neighbors = append(neighbors, b.GetCell(v2))
		} else {
			outOfBounds++
		}
	}
	return neighbors, outOfBounds
}

func (b *Board) GetCost(v Vector) (area, fence int) {
	me := b.GetCell(v)

	me.Visited = true
	area = 1

	neighbors, outOfBounds := b.GetNeighbors(v)
	fence += outOfBounds
	for _, c := range neighbors {
		if c.Crop != me.Crop {
			fence++
			continue
		}
		if !c.Visited {
			carea, cfence := b.GetCost(Vector{c.Location.X, c.Location.Y})
			area += carea
			fence += cfence
		}
	}

	return
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	b := LoadBoard("input.txt")

	var cost int
	for y, row := range b.Cells {
		for x, cell := range row {
			if !cell.Visited {
				area, fence := b.GetCost(Vector{x, y})
				fmt.Println(area, fence)
				cost += area * fence
			}
		}
	}

	fmt.Println(cost)
}
