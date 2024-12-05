package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	// Map of rules. Key is the first number.  Value is another map of numbers
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
input_line:
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

		for i, n := range input {
			// Check if the number is in the rules
			rule, found := rules[n]
			if !found {
				continue
			}

			// Check each number before this one to make sure it isn't in the rules to
			// come after this one.  If it is found this is invalid and we skip to the
			// next input line.
			for j := 0; j < i; j++ {
				b := input[j]
				if _, found := rule[b]; found {
					continue input_line
				}
			}

		}
		// Find the middle number and add it to tot
		tot += input[len(input)/2]
	}

	fmt.Println(tot)
}
