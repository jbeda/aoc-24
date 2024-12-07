package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Search possible combination of operators to produce total based on the
// inputs. Return true if a combination is found, false otherwise.
func Search(inputs []int, total int) bool {
	// Base case: last input stands alone
	if len(inputs) == 1 {
		return inputs[0] == total
	}

	// Pop off the last elemeent
	a := inputs[:len(inputs)-1]
	b := inputs[len(inputs)-1]

	// try addition
	if Search(a, total-b) {
		return true
	}

	// try multiplication
	q := total / b
	r := total % b
	if r == 0 && Search(a, q) {
		return true
	}

	return false
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var tot int

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()

		// Split the line at the colon
		ss := strings.Split(line, ":")
		if len(ss) != 2 {
			log.Fatalf("invalid line: %s", line)
		}

		total, err := strconv.Atoi(ss[0])
		if err != nil {
			log.Fatal(err)
		}

		// Split the input values
		ss = strings.Fields(ss[1])
		inputs := make([]int, len(ss))
		for i, s := range ss {
			inputs[i], err = strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
		}

		if Search(inputs, total) {
			tot += total
		}
	}

	fmt.Println(tot)
}
