package main

import (
	"bufio"
	"fmt"
	"log"
	"maps"
	"os"
	"slices"
)

type Vector struct {
	X, Y int
}

func Add(v1, v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func Sub(v1, v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}

func (v *Vector) GoString() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

func (v *Vector) IsOutOfBounds() bool {
	return v.X < 0 || v.X >= boardSize.X || v.Y < 0 || v.Y >= boardSize.Y
}

type Board map[Vector]bool

var boardSize Vector

var antinodes Board = make(Board)
var freqs = make(map[rune]Board)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	y := 0
	for scan.Scan() {
		line := scan.Text()

		if boardSize.X == 0 {
			boardSize.X = len(line)
		}

		for x, r := range line {
			if r == '.' {
				continue
			}

			freqBoard, ok := freqs[r]
			if !ok {
				freqBoard = make(Board)
				freqs[r] = freqBoard
			}

			freqBoard[Vector{x, y}] = true
		}
		y++
	}
	boardSize.Y = y

	for freq, board := range freqs {
		fmt.Println("Freq: ", string(freq))

		antennas := slices.Collect(maps.Keys(board))

		// Go through each pair and figure out antinodes
		for i, a1 := range antennas {
			for j, a2 := range antennas {
				if i >= j {
					continue
				}

				diff := Sub(a2, a1)
				antinode := Sub(a1, diff)
				if !antinode.IsOutOfBounds() {
					antinodes[antinode] = true
				}

				antinode = Add(a2, diff)
				if !antinode.IsOutOfBounds() {
					antinodes[antinode] = true
				}
			}
		}
	}

	antinodesVector := slices.Collect(maps.Keys(antinodes))
	fmt.Println("Antinodes: ", antinodesVector)
	fmt.Println("Antinodes count: ", len(antinodesVector))
}
