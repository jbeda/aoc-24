package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"
)

// --------------------------------------------------------------------
type Cell struct {
	Type CellType
	Dirs Dir
}

type CellType rune

const (
	Empty      CellType = '.'
	Wall       CellType = '#'
	Start      CellType = 'S'
	End        CellType = 'E'
	Breadcrumb CellType = '*'
)

func (c *Cell) Breadcrumb(dir Dir) {
	c.Type = Breadcrumb
	if dir == Up || dir == Down {
		c.Dirs |= Up | Down
	} else {
		c.Dirs |= Left | Right
	}
}

// --------------------------------------------------------------------
type Dir int

const (
	Up Dir = 1 << iota
	Right
	Down
	Left
)

func (d Dir) Vector() Vector {
	switch d {
	case Up:
		return Vector{0, -1}
	case Right:
		return Vector{1, 0}
	case Down:
		return Vector{0, 1}
	case Left:
		return Vector{-1, 0}
	}
	return Vector{}
}

// --------------------------------------------------------------------
type Maze struct {
	Cells [][]Cell
	Start Vector
	End   Vector
	Pos   Vector
}

func ReadMaze(fn string) Maze {
	m := Maze{}
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		row := []Cell{}
		for _, r := range line {
			ct := CellType(r)
			if ct == Start {
				m.Start = Vector{len(row), len(m.Cells)}
				ct = Empty
			} else if ct == End {
				m.End = Vector{len(row), len(m.Cells)}
				ct = Empty
			}
			row = append(row, Cell{Type: ct})
		}
		m.Cells = append(m.Cells, row)
	}

	return m
}

func (m Maze) String() string {
	var s string
	for y, row := range m.Cells {
		for x, cell := range row {
			v := Vector{x, y}
			if v == m.Start {
				s += string(Start)
			} else if v == m.End {
				s += string(End)
			} else {
				s += string(cell.Type)
			}
		}
		s += "\n"
	}
	return s
}

func (m Maze) SavePNG(curr *Vector, frontiers []*Frontier) {
	sizeX := len(m.Cells[0])
	sizeY := len(m.Cells)

	img := image.NewRGBA(image.Rect(0, 0, sizeX, sizeY))

	for y, row := range m.Cells {
		for x, cell := range row {
			if x == m.Start.X && y == m.Start.Y {
				cell.Type = Start
			} else if x == m.End.X && y == m.End.Y {
				cell.Type = End
			}
			var c color.RGBA
			switch cell.Type {
			case Empty:
				c = color.RGBA{255, 255, 255, 255}
			case Wall:
				c = color.RGBA{0, 0, 0, 255}
			case Start:
				c = color.RGBA{0, 255, 0, 255}
			case End:
				c = color.RGBA{255, 0, 0, 255}
			case Breadcrumb:
				c = color.RGBA{128, 128, 128, 255}
			}

			img.Set(x, y, c)
		}
	}

	for _, f := range frontiers {
		img.Set(f.Pos.X, f.Pos.Y, color.RGBA{0, 0, 255, 255})
	}

	if curr != nil {
		img.Set(curr.X, curr.Y, color.RGBA{0, 255, 255, 255})
	}

	f, err := os.Create("debug.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		log.Fatal(err)
	}
}

func (m Maze) At(v Vector) *Cell {
	return &m.Cells[v.Y][v.X]
}

// DFSSolve the maze returning the best score and if a solution was found
func (m *Maze) DFSSolve(pos Vector, dir Dir, bestScore int) (int, bool) {
	if pos == m.End {
		fmt.Println(m)
		return 0, true
	}

	curr := m.At(pos)
	curr.Type = Breadcrumb
	curr.Dirs |= dir

	solved := false
	for _, newDir := range []Dir{Up, Right, Down, Left} {
		newPos := pos.Add(newDir.Vector())
		newCell := m.At(newPos)
		if newCell.Type == Empty ||
			(newCell.Type == Breadcrumb && newCell.Dirs&newDir == 0) {
			var cost int
			if dir == newDir {
				cost = 1
			} else {
				cost = 1001
			}
			// If we are already over the best score, don't bother
			if cost >= bestScore {
				fmt.Println(m)
				continue
			}
			subscore, ok := m.DFSSolve(newPos, newDir, bestScore-cost)
			if ok {
				solved = true
				if subscore+cost < bestScore {
					bestScore = subscore + cost
				}
			}
		}
	}

	m.At(pos).Type = Empty

	return bestScore, solved
}

type Frontier struct {
	Pos   Vector
	Dir   Dir
	Score int
}

type PriorityQueue []*Frontier

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Score < pq[j].Score }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Frontier)) }
func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	x := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return x
}

func (m Maze) BFSSolve() int {
	// Use a heap to keep track of the frontier with the lowest score
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	heap.Push(&pq, &Frontier{m.Start, Right, 0})
	cell := m.At(m.Start)
	cell.Type = Breadcrumb
	cell.Breadcrumb(Right)

	for len(pq) > 0 {
		// Pop the first frontier
		f := heap.Pop(&pq).(*Frontier)
		cell := m.At(f.Pos)
		cell.Type = Breadcrumb
		cell.Breadcrumb(f.Dir)

		// m.SavePNG(&f.Pos, pq)

		if f.Pos == m.End {
			return f.Score
		}

		for _, newDir := range []Dir{Up, Right, Down, Left} {
			newPos := f.Pos.Add(newDir.Vector())
			newCell := m.At(newPos)
			if newCell.Type == Empty ||
				(newCell.Type == Breadcrumb && newCell.Dirs&newDir == 0) {
				if f.Dir == newDir {
					heap.Push(&pq, &Frontier{newPos, newDir, f.Score + 1})
				} else {
					heap.Push(&pq, &Frontier{f.Pos, newDir, f.Score + 1000})
				}
			}
		}
	}

	return math.MaxInt
}

// --------------------------------------------------------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m := ReadMaze("input.txt")
	fmt.Println(m)
	fmt.Println(m.BFSSolve())
	m.SavePNG(nil, nil)
}
