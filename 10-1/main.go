package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Vector struct {
	X, Y int
}

var board [][]int
var size Vector

func IsOutOfBounds(v Vector) bool {
	return v.X < 0 || v.X >= size.X || v.Y < 0 || v.Y >= size.Y
}

func SearchBoard(v Vector, hPrev int) int {
	if IsOutOfBounds(v) {
		return 0
	}

	h := board[v.Y][v.X]
	if h != hPrev+1 {
		return 0
	}

	if h == 9 {
		return 1
	}

	ret := 0
	ret += SearchBoard(Vector{v.X - 1, v.Y}, h)
	ret += SearchBoard(Vector{v.X + 1, v.Y}, h)
	ret += SearchBoard(Vector{v.X, v.Y - 1}, h)
	ret += SearchBoard(Vector{v.X, v.Y + 1}, h)

	return ret
}

func ScoreTrailhead(v Vector) int {
	res := SearchBoard(v, -1)
	return res
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Load the board
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		var row []int
		for _, c := range line {
			h, err := strconv.Atoi(string(c))
			if err != nil {
				log.Fatal(err)
			}
			row = append(row, h)
		}
		board = append(board, row)
	}
	size.X = len(board[0])
	size.Y = len(board)

	// Find all trailheads
	tot := 0
	for y, row := range board {
		for x, h := range row {
			if h == 0 {
				score := ScoreTrailhead(Vector{x, y})

				fmt.Printf("Trailhead at (%d, %d) has score %d\n", x, y, score)

				tot += score
			}
		}
	}

	fmt.Println(tot)
}
