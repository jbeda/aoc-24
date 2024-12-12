package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

// Input to the problem: the value of a stone and the number of iterations to
// perform
type question struct {
	value int
	iters int
}

// A memoization table for answers.  The key is the question, the value is the
// number of stones after the iterations.
var answers map[question]int

// Given a number on a stone and the number of blinks, how many stones are left
// at the end?
func Blink(value int, iters int) int {
	if iters == 0 {
		return 1
	}

	// Look it up in the memoization table
	if v, ok := answers[question{value, iters}]; ok {
		return v
	}

	if value == 0 {
		ret := Blink(1, iters-1)
		answers[question{value, iters}] = ret
		return ret
	}

	vText := strconv.Itoa(value)
	if len(vText)%2 == 0 {
		s1, err := strconv.Atoi(vText[:len(vText)/2])
		if err != nil {
			log.Fatal(err)
		}
		s2, err := strconv.Atoi(vText[len(vText)/2:])
		if err != nil {
			log.Fatal(err)
		}

		ret := Blink(s1, iters-1) + Blink(s2, iters-1)
		answers[question{value, iters}] = ret
		return ret
	}

	ret := Blink(value*2024, iters-1)
	answers[question{value, iters}] = ret
	return ret
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// input := "125 17"
	input := "5 89749 6061 43 867 1965860 0 206250"

	ss := strings.Fields(input)

	var stones []int

	// Load stones
	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		stones = append(stones, i)
	}

	answers = make(map[question]int)

	tot := 0
	for _, s := range stones {
		tot += Blink(s, 75)
	}

	fmt.Println(tot)
}
