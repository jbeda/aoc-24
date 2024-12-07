package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type CellStatus struct {
	visited VisitStatus
	dir     Dir
}

type VisitStatus int

const (
	CellEmpty VisitStatus = iota
	CellObstacle
	CellVisited
)

type Dir int

const (
	DirUp Dir = 1 << iota
	DirRight
	DirDown
	DirLeft
)

type Pos struct {
	X, Y int
}

type Board struct {
	grid      [][]CellStatus
	playerPos Pos
	playerDir Dir
}

func LoadBoard(fn string) *Board {
	b := new(Board)
	// Load the board
	f, err := os.Open(fn)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Load the board
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		lineLength := len(line)
		gridLine := make([]CellStatus, lineLength)
		b.grid = append(b.grid, gridLine)
		var status VisitStatus
		for i, c := range line {
			switch c {
			case '.':
				status = CellEmpty
			case '#':
				status = CellObstacle
			case '^':
				status = CellEmpty
				b.playerPos = Pos{}
				b.playerPos.Y = len(b.grid) - 1
				b.playerPos.X = i
				b.playerDir = DirUp
			}
			gridLine[i].visited = status
		}
	}
	return b
}

func (b *Board) Clone() *Board {
	ret := new(Board)
	ret.grid = make([][]CellStatus, len(b.grid))
	for y, line := range b.grid {
		ret.grid[y] = make([]CellStatus, len(line))
		copy(ret.grid[y], line)
	}
	ret.playerPos = b.playerPos
	ret.playerDir = b.playerDir
	return ret
}

func (b *Board) nextPos() (Pos, bool) {
	switch b.playerDir {
	case DirUp:
		if b.playerPos.Y <= 0 {
			return Pos{}, false
		}
		return Pos{b.playerPos.X, b.playerPos.Y - 1}, true
	case DirRight:
		if b.playerPos.X >= len(b.grid[b.playerPos.Y])-1 {
			return Pos{}, false
		}
		return Pos{b.playerPos.X + 1, b.playerPos.Y}, true
	case DirDown:
		if b.playerPos.Y >= len(b.grid)-1 {
			return Pos{}, false
		}
		return Pos{b.playerPos.X, b.playerPos.Y + 1}, true
	case DirLeft:
		if b.playerPos.X <= 0 {
			return Pos{}, false
		}
		return Pos{b.playerPos.X - 1, b.playerPos.Y}, true
	}
	log.Fatalf("invalid direction: %d", b.playerDir)
	return Pos{}, false
}

func (b *Board) TurnRight() {
	switch b.playerDir {
	case DirUp:
		b.playerDir = DirRight
	case DirRight:
		b.playerDir = DirDown
	case DirDown:
		b.playerDir = DirLeft
	case DirLeft:
		b.playerDir = DirUp
	}
}

// Look for a loop or an exit. Return 1 for a loop, 0 for an exit.
func (b *Board) LoopOrExit() int {
	for {
		//b.Print()

		// Check if we are about to go off the board
		nextPos, ok := b.nextPos()
		if !ok {
			return 0
		}

		// Check if we have visited this cell before in this direction
		cell := b.grid[b.playerPos.Y][b.playerPos.X]
		if cell.visited == CellVisited && (cell.dir&b.playerDir) != 0 {
			//b.Print()
			return 1
		}

		// Mark this cell as visited
		b.grid[b.playerPos.Y][b.playerPos.X].dir |= b.playerDir
		b.grid[b.playerPos.Y][b.playerPos.X].visited = CellVisited

		// Move the player
		if b.grid[nextPos.Y][nextPos.X].visited == CellObstacle {
			b.TurnRight()
		} else {
			b.playerPos = nextPos
		}
	}
}

func (b *Board) Print() {
	fmt.Println()
	for y, line := range b.grid {
		for x, cell := range line {
			if x == b.playerPos.X && y == b.playerPos.Y {
				switch b.playerDir {
				case DirUp:
					fmt.Print("^")
				case DirRight:
					fmt.Print(">")
				case DirDown:
					fmt.Print("v")
				case DirLeft:
					fmt.Print("<")
				}
			} else {
				switch cell.visited {
				case CellEmpty:
					fmt.Print(".")
				case CellObstacle:
					fmt.Print("#")
				case CellVisited:
					fmt.Print(cell.dir)
				}
			}
		}
		fmt.Println()
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	b := LoadBoard("input.txt")

	var loops int

	// Move the player
	for {
		//b.Print()

		// Check if we are about to go off the board
		nextPos, ok := b.nextPos()
		if !ok {
			break
		}

		cell := &b.grid[b.playerPos.Y][b.playerPos.X]
		cellNext := &b.grid[nextPos.Y][nextPos.X]

		// Fork as if there were an obsticle ahead and we haven't visited that
		// cell as then we've already tested it.
		if cellNext.visited == CellEmpty {
			b1 := b.Clone()
			cell1 := &(b1.grid[nextPos.Y][nextPos.X])
			cell1.visited = CellObstacle
			loops += b1.LoopOrExit()
		}

		// Mark this cell as visited
		cell.dir |= b.playerDir
		cell.visited = CellVisited

		if cellNext.visited == CellObstacle {
			b.TurnRight()
		} else {
			// Continue with main loop and move the player
			b.playerPos = nextPos
		}
	}

	fmt.Println(loops)

}
