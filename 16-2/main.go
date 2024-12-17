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
type Dir int

const (
	Up Dir = iota
	Right
	Down
	Left
)

var Dirs = []Dir{Up, Right, Down, Left}

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

func (d Dir) Turns() []Dir {
	switch d {
	case Up, Down:
		return []Dir{Left, Right}
	case Right, Left:
		return []Dir{Up, Down}
	}
	return []Dir{}
}

// --------------------------------------------------------------------
type Node struct {
	Pos  Vector
	Dir  Dir
	Dist int
	Cell *Cell
}

// --------------------------------------------------------------------
type Edge struct {
	From, To *Node
	Cost     int
}

// --------------------------------------------------------------------
type Cell struct {
	Type   CellType
	Nodes  [4]*Node
	OnPath bool
}

type CellType rune

const (
	Empty CellType = '.'
	Wall  CellType = '#'
	Start CellType = 'S'
	End   CellType = 'E'
)

func NewCell(pos Vector, t CellType) *Cell {
	c := Cell{Type: t}
	for _, d := range Dirs {
		c.Nodes[d] = &Node{Pos: pos, Dir: d, Dist: math.MaxInt, Cell: &c}
	}

	return &c
}

func (c Cell) IsAnyVisited() bool {
	for _, n := range c.Nodes {
		if n.Dist != math.MaxInt {
			return true
		}
	}
	return false
}

func (c Cell) GetNode(d Dir) *Node {
	return c.Nodes[d]
}

// --------------------------------------------------------------------
type Maze struct {
	Cells [][]*Cell
	Start Vector
	End   Vector
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
		row := []*Cell{}
		for _, r := range line {
			ct := CellType(r)
			pos := Vector{len(row), len(m.Cells)}
			if ct == Start {
				m.Start = pos
				ct = Empty
			} else if ct == End {
				m.End = pos
				ct = Empty
			}
			row = append(row, NewCell(pos, ct))
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

func (m Maze) SavePNG(curr *Vector, frontiers []*Node) {
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
				if cell.IsAnyVisited() {
					c = color.RGBA{128, 128, 128, 255}
				}
			case Wall:
				c = color.RGBA{0, 0, 0, 255}
			case Start:
				c = color.RGBA{0, 255, 0, 255}
			case End:
				c = color.RGBA{255, 0, 0, 255}
			}

			if cell.OnPath {
				c = color.RGBA{255, 255, 0, 255}
			}

			img.Set(x, y, c)
		}
	}

	for _, n := range frontiers {
		img.Set(n.Pos.X, n.Pos.Y, color.RGBA{0, 0, 255, 255})
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
	return m.Cells[v.Y][v.X]
}

func (m Maze) GetOutEdges(n *Node) []*Edge {
	edges := []*Edge{}

	// Create an edge for forward movement
	{
		nextPos := n.Pos.Add(n.Dir.Vector())
		nextCell := m.At(nextPos)
		if nextCell.Type == Empty {
			nextNode := nextCell.GetNode(n.Dir)
			edges = append(edges, &Edge{n, nextNode, 1})
		}
	}

	// Now create edges for each turn
	for _, d := range n.Dir.Turns() {
		nextNode := n.Cell.GetNode(d)
		edges = append(edges, &Edge{n, nextNode, 1000})
	}
	return edges
}

func (m Maze) GetInEdges(n *Node) []*Edge {
	edges := []*Edge{}

	// Create an edge for backward movement
	{
		nextPos := n.Pos.Sub(n.Dir.Vector())
		nextCell := m.At(nextPos)
		if nextCell.Type == Empty {
			nextNode := nextCell.GetNode(n.Dir)
			edges = append(edges, &Edge{nextNode, n, 1})
		}
	}

	// Now create edges for each turn
	for _, d := range n.Dir.Turns() {
		nextNode := n.Cell.GetNode(d)
		edges = append(edges, &Edge{nextNode, n, 1000})
	}
	return edges
}

type PriorityQueue []*Node

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Dist < pq[j].Dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Node)) }
func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	x := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return x
}

func (m Maze) Solve() int {
	// Use a heap to keep track of the frontier with the lowest score
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)

	// Add the start node to the queue
	{
		cell := m.At(m.Start)
		node := cell.GetNode(Right)
		node.Dist = 0
		heap.Push(&pq, node)
	}

	for len(pq) > 0 {
		// Pop the first frontier
		n := heap.Pop(&pq).(*Node)

		// m.SavePNG(&f.Pos, pq)

		if n.Pos == m.End {
			return n.Dist
		}

		for _, e := range m.GetOutEdges(n) {
			if e.To.Dist > n.Dist+e.Cost {
				e.To.Dist = n.Dist + e.Cost
				heap.Push(&pq, e.To)
			}
		}
	}

	return math.MaxInt
}

func (m Maze) MarkPath() {
	var q []*Node

	// Add all the end nodes
	{
		endCell := m.At(m.End)
		endCell.OnPath = true
		for _, n := range endCell.Nodes {
			if n.Pos == m.End {
				q = append(q, n)
			}
		}
	}

	for len(q) > 0 {
		n := q[0]
		q = q[1:]

		if n.Pos == m.Start {
			continue
		}

		edges := m.GetInEdges(n)

		for _, e := range edges {
			if e.From.Dist != math.MaxInt && e.From.Dist+e.Cost == e.To.Dist {
				e.From.Cell.OnPath = true
				q = append(q, e.From)
			}
		}
	}
}

// --------------------------------------------------------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m := ReadMaze("input.txt")
	fmt.Println(m.Solve())
	m.MarkPath()
	m.SavePNG(nil, nil)

	nPath := 0
	for _, row := range m.Cells {
		for _, cell := range row {
			if cell.OnPath {
				nPath++
			}
		}
	}

	fmt.Println(nPath)
}
