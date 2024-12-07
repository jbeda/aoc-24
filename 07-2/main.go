package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// Do a textual split of two integers.  Remove b from teh end of a.
func split(a, b int64) (int64, bool) {
	sa := strconv.FormatInt(a, 10)
	sb := strconv.FormatInt(b, 10)

	if strings.HasSuffix(sa, sb) {
		trimmed := strings.TrimSuffix(sa, sb)
		if trimmed == "" || trimmed[0] == '-' {
			return 0, false
		}
		r, err := strconv.ParseInt(trimmed, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		return r, true
	}
	return 0, false
}

// Search possible combination of operators to produce total based on the
// inputs. Return true if a combination is found, false otherwise.
func Search(inputs []int64, total int64) bool {
	// Base case: last input stands alone
	if len(inputs) == 1 {
		return inputs[0] == total
	}

	// Pop off the last elemeent
	a := inputs[:len(inputs)-1]
	b := inputs[len(inputs)-1]

	// try addition
	if b < total && Search(a, total-b) {
		return true
	}

	// try multiplication
	q := total / b
	r := total % b
	if r == 0 && Search(a, q) {
		return true
	}

	// try concatenation
	s, ok := split(total, b)
	if ok && Search(a, s) {
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

	var tot int64

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()

		// Split the line at the colon
		ss := strings.Split(line, ":")
		if len(ss) != 2 {
			log.Fatalf("invalid line: %s", line)
		}

		total, err := strconv.ParseInt(ss[0], 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		// Split the input values
		ss = strings.Fields(ss[1])
		inputs := make([]int64, len(ss))
		for i, s := range ss {
			inputs[i], err = strconv.ParseInt(s, 10, 64)
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
