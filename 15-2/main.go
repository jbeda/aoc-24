package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// -------------------------------------

type Cell rune

const InBox Cell = 'O'
const LBox Cell = '['
const RBox Cell = ']'
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
			c1 := Cell(r)
			c2 := c1
			if c1 == Player {
				b.Pos = Vector{len(row), len(*cells)}
				c2 = Empty
			}
			if c1 == InBox {
				c1 = LBox
				c2 = RBox
			}

			row = append(row, c1)
			row = append(row, c2)
		}
		*cells = append(*cells, row)
	}
	return b
}

func (b *Board) CanMove(pos Vector, move MoveType) bool {
	newPos := MoveVector(pos, move)
	newCell := b.Get(newPos)

	if newCell == Wall {
		return false
	}

	if newCell == Empty {
		return true
	}

	if move == Up || move == Down {
		if newCell == LBox {
			newPos2 := Vector{newPos.X + 1, newPos.Y}
			return b.CanMove(newPos, move) && b.CanMove(newPos2, move)
		}
		if newCell == RBox {
			newPos2 := Vector{newPos.X - 1, newPos.Y}
			return b.CanMove(newPos, move) && b.CanMove(newPos2, move)
		}
	}

	return b.CanMove(newPos, move)
}

func (b *Board) Move(pos Vector, move MoveType) {
	cell := b.Get(pos)
	newPos := MoveVector(pos, move)
	newCell := b.Get(newPos)

	if newCell == Empty {
		b.Set(pos, Empty)
		b.Set(newPos, cell)
		return
	}

	if newCell == Wall {
		return
	}

	var newPos2 *Vector = nil
	if move == Up || move == Down {
		if newCell == LBox {
			newPos2 = &Vector{newPos.X + 1, newPos.Y}
		}
		if newCell == RBox {
			newPos2 = &Vector{newPos.X - 1, newPos.Y}
		}
	}

	b.Move(newPos, move)

	if newPos2 != nil {
		b.Move(*newPos2, move)
	}

	b.Set(pos, Empty)
	b.Set(newPos, cell)
}

func (b *Board) MoveRobot(move MoveType) {
	if b.CanMove(b.Pos, move) {
		b.Move(b.Pos, move)
		b.Pos = MoveVector(b.Pos, move)
	}
}

func (b *Board) Score() int {
	var score int
	for y, row := range b.Cells {
		for x, cell := range row {
			if cell == LBox {
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
	_ = moves

	fmt.Println(board)

	for _, move := range moves {
		board.MoveRobot(move)

		// fmt.Println("Move:", string(move))
		// fmt.Println(board)
	}

	fmt.Println("Score:", board.Score())
}
