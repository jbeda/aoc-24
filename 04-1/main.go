package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var grid [][]rune

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		grid = append(grid, []rune(line))
	}

	var target = []rune("XMAS")
	var offsets = []struct{ x, y int }{
		{1, 0}, {-1, 0}, // Horizontal
		{0, 1}, {0, -1}, // Vertical
		{1, 1}, {-1, -1}, {1, -1}, {-1, 1}, // Diagonal
	}

	var tot int
	for y, row := range grid {
		for x, ch := range row {
			if ch != target[0] {
				continue
			}

			for _, offset := range offsets {
				var found = true
				for i, ch := range target {
					if x+i*offset.x < 0 || x+i*offset.x >= len(row) {
						found = false
						break
					}
					if y+i*offset.y < 0 || y+i*offset.y >= len(grid) {
						found = false
						break
					}
					if grid[y+i*offset.y][x+i*offset.x] != ch {
						found = false
						break
					}
				}
				if found {
					log.Printf("Found %s at (%d, %d)\n", string(target), x, y)
					tot++
				}
			}
		}
	}

	fmt.Println(tot)
}
