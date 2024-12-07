package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type CellStatus int

const (
	CellEmpty CellStatus = iota
	CellObstacle
	CellVisited
)

type Dir int

const (
	DirUp Dir = iota
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

func (b *Board) Print() {
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
				switch cell {
				case CellEmpty:
					fmt.Print(".")
				case CellObstacle:
					fmt.Print("#")
				case CellVisited:
					fmt.Print("X")
				}
			}
		}
		fmt.Println()
	}
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var b Board

	// Load the board
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		lineLength := len(line)
		gridLine := make([]CellStatus, lineLength)
		b.grid = append(b.grid, gridLine)
		var status CellStatus
		for i, c := range line {
			switch c {
			case '.':
				status = CellEmpty
			case '#':
				status = CellObstacle
			case '^':
				status = CellVisited
				b.playerPos = Pos{}
				b.playerPos.Y = len(b.grid) - 1
				b.playerPos.X = i
				b.playerDir = DirUp
			}
			gridLine[i] = status
		}
	}

	// Move the player
	for {
		nextPos, ok := b.nextPos()
		if !ok {
			break
		}

		if b.grid[nextPos.Y][nextPos.X] == CellObstacle {
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
		} else {
			b.playerPos = nextPos
			b.grid[b.playerPos.Y][b.playerPos.X] = CellVisited
		}
	}

	// Count the visited cells
	var tot int
	for _, line := range b.grid {
		for _, cell := range line {
			if cell == CellVisited {
				tot++
			}
		}
	}

	fmt.Println(tot)

}
