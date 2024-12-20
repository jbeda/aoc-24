package main

import (
	"fmt"
	"log"
	"maps"
	"math"
	"slices"
	"time"
)

type Cell struct {
	Pos         Vector
	Wall        bool
	DistToEnd   int
	DistToStart int
}

func NewCell(pos Vector, wall bool) Cell {
	return Cell{pos, wall, math.MaxInt, math.MaxInt}
}

// --------------------------------------------------------------------
type Maze struct {
	Cells [][]Cell
	Size  Vector
	Start Vector
	End   Vector
}

func NewMaze(lines []string) *Maze {
	m := Maze{Cells: make([][]Cell, len(lines))}
	for y, line := range lines {
		m.Cells[y] = make([]Cell, len(line))
		for x, r := range line {
			m.Cells[y][x] = NewCell(Vector{x, y}, r == '#')
			if r == 'S' {
				m.Start = Vector{x, y}
			} else if r == 'E' {
				m.End = Vector{x, y}
			}
		}
	}

	m.Size = Vector{len(m.Cells[0]), len(m.Cells)}
	return &m
}

func (m *Maze) At(pos Vector) *Cell {
	if pos.IsOOB(m.Size) {
		return nil
	}
	return &m.Cells[pos.Y][pos.X]
}

func (m *Maze) String() string {
	var s string
	for y, row := range m.Cells {
		for x, cell := range row {
			if m.Start.X == x && m.Start.Y == y {
				s += "S"
			} else if m.End.X == x && m.End.Y == y {
				s += "E"
			} else if cell.Wall {
				s += "#"
			} else {
				s += "."
			}
		}
		s += "\n"
	}
	return s
}

// Compute all the distances from the end to each cell
func (m *Maze) BackwardsSolve() {
	m.At(m.End).DistToEnd = 0
	frontiers := []Vector{m.End}
	for len(frontiers) > 0 {
		frontier := frontiers[0]
		frontiers = frontiers[1:]
		for _, neighbor := range frontier.Neighbors4() {
			if !m.At(neighbor).Wall && m.At(neighbor).DistToEnd == math.MaxInt {
				m.At(neighbor).DistToEnd = m.At(frontier).DistToEnd + 1
				frontiers = append(frontiers, neighbor)
			}
		}
	}
}

type Shortcut struct {
	Pos1 Vector
	Pos2 Vector
	Dist int
}

// List all the neighbors within X steps
func (v Vector) NeighborsX(dist int) []Vector {
	neighbors := []Vector{}

	for dy := -20; dy <= 20; dy++ {
		stepsLeft := 20 - AbsInt(dy)
		for dx := -stepsLeft; dx <= stepsLeft; dx++ {
			if dx == 0 && dy == 0 {
				continue
			}
			neighbors = append(neighbors, Vector{v.X + dx, v.Y + dy})
		}
	}

	return neighbors
}

func (m *Maze) SolveShortcuts(fastest int) []Shortcut {
	shortcuts := []Shortcut{}

	m.At(m.Start).DistToStart = 0
	frontiers := []Vector{m.Start}
	for len(frontiers) > 0 {
		frontier := frontiers[0]
		frontiers = frontiers[1:]
		for _, neighbor := range frontier.Neighbors4() {
			if !m.At(neighbor).Wall && m.At(neighbor).DistToStart == math.MaxInt {
				m.At(neighbor).DistToStart = m.At(frontier).DistToStart + 1
				frontiers = append(frontiers, neighbor)
			}
		}

		// Check if we can make a shortcut. The first move is always a wall and the
		// second must not be a wall.
		for _, n := range frontier.NeighborsX(20) {
			if n.IsOOB(m.Size) || m.At(n).Wall {
				continue
			}

			dist := m.At(frontier).DistToStart + frontier.ManhattanDist(n) + m.At(n).DistToEnd
			if dist < fastest {
				shortcuts = append(shortcuts, Shortcut{frontier, n, dist})
			}
		}
	}

	return shortcuts
}

// --------------------------------------------------------------------
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	lines := ReadFileLines("input.txt")

	m := NewMaze(lines)
	// fmt.Println(m)

	m.BackwardsSolve()
	fastest := m.At(m.Start).DistToEnd
	fmt.Println("Fastest: ", fastest)

	shortcuts := m.SolveShortcuts(fastest)
	fmt.Println("len(shortcuts):", len(shortcuts))

	// Build/Print histogram of savings
	if false {
		shortcutHistogram := make(map[int]int)
		for _, s := range shortcuts {
			savings := fastest - s.Dist
			shortcutHistogram[savings]++
		}

		histKeys := slices.Collect(maps.Keys(shortcutHistogram))
		slices.Sort(histKeys)
		for _, savings := range histKeys {
			fmt.Printf("Savings: %d, Count: %d\n", savings, shortcutHistogram[savings])
		}
	}

	if false {
		for _, s := range shortcuts {
			savings := fastest - s.Dist
			if savings == 76 {
				fmt.Printf("Shortcut: %v -> %v, Dist: %d, Savings: %d\n", s.Pos1, s.Pos2, s.Dist, savings)
			}
		}
	}

	// Count the number of shortcuts that save over 100 steps
	over100 := 0
	for _, s := range shortcuts {
		if fastest-s.Dist >= 100 {
			over100++
		}
	}
	fmt.Println(">= 100:", over100)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
