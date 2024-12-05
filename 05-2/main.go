package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Map of after rules. Key is a number.  Value is another map of numbers
	// that must come after the first.
	var rules map[int]map[int]bool = make(map[int]map[int]bool)

	// Read the rules.  In the form of <int>|<int>
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()

		if len(line) == 0 {
			break
		}

		// Parse the rule
		var before, after int
		ss := strings.Split(line, "|")
		if len(ss) != 2 {
			log.Fatalf("invalid rule: %s", line)
		}

		before, err = strconv.Atoi(ss[0])
		if err != nil {
			log.Fatal(err)
		}
		after, err = strconv.Atoi(ss[1])
		if err != nil {
			log.Fatal(err)
		}

		rule, found := rules[before]
		if !found {
			rule = make(map[int]bool)
			rules[before] = rule
		}
		rule[after] = true
	}

	var tot int
	for scan.Scan() {
		line := scan.Text()

		var input []int

		ss := strings.Split(line, ",")
		for _, s := range ss {
			n, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			input = append(input, n)
		}

		// Sort the list according to the rules
		sorted := make([]int, len(input))
		copy(sorted, input)
		slices.SortStableFunc(sorted, func(i, j int) int {
			rule, found := rules[i]
			if found {
				_, found = rule[j]
				if found {
					// We know that j must come after i so return -1.
					return -1
				}
			}

			rule, found = rules[j]
			if found {
				_, found = rule[i]
				if found {
					// We know that i must come after j so return 1.
					return 1
				}
			}

			// We have no opinion on ordering
			return 0
		})

		fixed := false
		for i := range input {
			if input[i] != sorted[i] {
				fixed = true
				break
			}
		}

		// Find the middle number and add it to tot
		if fixed {
			tot += sorted[len(sorted)/2]
		}
	}

	fmt.Println(tot)
}
