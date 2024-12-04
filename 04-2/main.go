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

	var targets = [][][]rune{
		{
			[]rune("M.M"),
			[]rune(".A."),
			[]rune("S.S"),
		},
		{
			[]rune("S.S"),
			[]rune(".A."),
			[]rune("M.M"),
		},
		{
			[]rune("M.S"),
			[]rune(".A."),
			[]rune("M.S"),
		},
		{
			[]rune("S.M"),
			[]rune(".A."),
			[]rune("S.M"),
		},
	}

	var tot int
	for y, row := range grid {
		for x := range row {
			for _, target := range targets {
				var found = true
			target:
				for j, trow := range target {
					for i, tch := range trow {
						if x+i >= len(row) || y+j >= len(grid) {
							found = false
							break target
						}
						if tch == '.' {
							continue
						}
						if grid[y+j][x+i] != tch {
							found = false
							break target
						}
					}
				}
				if found {
					tot++
				}
			}
		}
	}

	fmt.Println(tot)
}
