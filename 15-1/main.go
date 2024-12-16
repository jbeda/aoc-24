package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// -------------------------------------

type Cell rune

const Box Cell = 'O'
const Empty Cell = '.'
const Wall Cell = '#'
const Player Cell = '@'

type MoveType rune

const Up MoveType = '^'
const Down MoveType = 'v'
const Left MoveType = '<'
const Right MoveType = '>'

// -------------------------------------
func MoveVector(v Vector, m MoveType) Vector {
	switch m {
	case Up:
		v.Y--
	case Down:
		v.Y++
	case Left:
		v.X--
	case Right:
		v.X++
	}
	return v
}

// -------------------------------------
type Board struct {
	Cells [][]Cell
	Pos   Vector
}

func (b Board) Get(v Vector) Cell {
	return b.Cells[v.Y][v.X]
}

func (b *Board) Set(v Vector, c Cell) {
	b.Cells[v.Y][v.X] = c
}

func (b Board) String() string {
	var s string
	for _, row := range b.Cells {
		for _, cell := range row {
			s += string(cell)
		}
		s += "\n"
	}
	return s
}

// Read board from scanner up until the first blank line
func ReadBoard(scan *bufio.Scanner) Board {
	var b Board
	cells := &b.Cells
	for scan.Scan() {
		line := scan.Text()
		if len(line) == 0 {
			break
		}
		var row []Cell
		for _, r := range line {
			c := Cell(r)
			if c == Player {
				b.Pos = Vector{len(row), len(*cells)}
			}

			row = append(row, c)
		}
		*cells = append(*cells, row)
	}
	return b
}

func (b *Board) MoveObject(pos Vector, move MoveType) bool {
	cell := b.Get(pos)
	newPos := MoveVector(pos, move)
	newCell := b.Get(newPos)

	if newCell == Wall {
		return false
	}

	if b.Get(newPos) == Empty {
		b.Set(pos, Empty)
		b.Set(newPos, cell)
		return true
	}

	if newCell == Box {
		if b.MoveObject(newPos, move) {
			b.Set(pos, Empty)
			b.Set(newPos, cell)
			return true
		}
	}
	return false
}

func (b *Board) MoveRobot(move MoveType) {
	if b.MoveObject(b.Pos, move) {
		b.Pos = MoveVector(b.Pos, move)
	}
}

func (b *Board) Score() int {
	var score int
	for y, row := range b.Cells {
		for x, cell := range row {
			if cell == Box {
				score += x + y*100
			}
		}
	}
	return score
}

// -------------------------------------

func ReadMoves(scan *bufio.Scanner) []MoveType {
	var moves []MoveType
	for scan.Scan() {
		line := scan.Text()
		for _, r := range line {
			moves = append(moves, MoveType(r))
		}
	}
	return moves
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	board := ReadBoard(scan)
	moves := ReadMoves(scan)

	fmt.Println(board)

	for _, move := range moves {
		board.MoveRobot(move)

		// fmt.Println("Move:", string(move))
		// fmt.Println(board)
	}

	fmt.Println("Score:", board.Score())
}
