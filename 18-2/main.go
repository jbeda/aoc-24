package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"time"
)

type Cell struct {
	Blocked bool
	Loc     Vector
	Score   int
}

func NewCell(loc Vector) Cell {
	return Cell{false, loc, math.MaxInt}
}

// --------------------------------------------------------------------
type Board struct {
	Size  Vector
	Cells [][]Cell
	Start Vector
	End   Vector
}

func NewBoard(size Vector) *Board {
	cells := make([][]Cell, size.Y)
	for y := range cells {
		cells[y] = make([]Cell, size.X)
		for x := range cells[y] {
			cells[y][x] = NewCell(Vector{x, y})
		}
	}
	return &Board{size, cells, Vector{0, 0}, size.SubInt(1)}
}

func (b *Board) At(v Vector) *Cell {
	if v.IsOOB(b.Size) {
		return nil
	}
	return &b.Cells[v.Y][v.X]
}

func (b *Board) String() string {
	var sb strings.Builder
	for y := range b.Cells {
		for x := range b.Cells[y] {
			if b.Cells[y][x].Blocked {
				sb.WriteString("#")
			} else {
				sb.WriteString(".")
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func (b *Board) Clone() *Board {
	clone := NewBoard(b.Size)
	for y := range b.Cells {
		copy(clone.Cells[y], b.Cells[y])
	}
	return clone
}

func (b *Board) Reset() {
	for y := range b.Cells {
		for x := range b.Cells[y] {
			b.Cells[y][x].Score = math.MaxInt
		}
	}
}

type PriorityQueue []*Cell

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Score < pq[j].Score }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*Cell)) }
func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	x := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return x
}

func (b *Board) Solve() int {
	pq := PriorityQueue{}
	start := b.At(b.Start)
	start.Score = 0
	heap.Push(&pq, start)

	for len(pq) > 0 {
		cell := heap.Pop(&pq).(*Cell)

		if cell.Loc == b.End {
			break
		}

		for _, neighbor := range cell.Loc.Neighbors4() {
			if next := b.At(neighbor); next != nil && !next.Blocked {
				nextScore := cell.Score + 1
				if nextScore < next.Score {
					next.Score = nextScore
					heap.Push(&pq, next)
				}
			}
		}
	}

	return b.At(b.End).Score
}

// --------------------------------------------------------------------
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	var board *Board
	var fn string

	if false {
		board = NewBoard(Vector{7, 7})
		fn = "test.txt"
	} else {
		board = NewBoard(Vector{71, 71})
		fn = "input.txt"
	}

	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var events []Vector
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		ss := strings.Split(line, ",")
		event := Vector{MustAtoi(ss[0]), MustAtoi(ss[1])}
		events = append(events, event)
	}

	// Find the first solution that blocks the path.
	// Only optimization here is to only resolve if we've visited the event
	// already.  Optimal might be doing a binary search across events.  This naive
	// solution took me 2.1s to run.
	board.Solve()
	for i := 0; i < len(events); i++ {
		event := events[i]
		eventCell := board.At(event)
		eventCell.Blocked = true
		if eventCell.Score != math.MaxInt {
			board.Reset()
			score := board.Solve()
			if score == math.MaxInt {
				fmt.Printf("No solution at event: %d,%d\n", event.X, event.Y)
				break
			}
		}
	}

	fmt.Printf("Elapsed: %v\n", time.Since(timeStart))
}
