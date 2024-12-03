package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/exp/constraints"
)

func Abs[T constraints.Integer](x T) T {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	var tot int

report:
	for scan.Scan() {
		line := scan.Text()

		var report []int
		ss := strings.Fields(line)
		for _, s := range ss {
			n, err := strconv.Atoi(s)
			if err != nil {
				log.Fatal(err)
			}
			report = append(report, n)
		}

		last := report[0]
		increasing := false
		if report[1] > last {
			increasing = true
		}

		var bad int
		for i := 1; i < len(report); i++ {
			curr := report[i]
			if (increasing && curr < last) ||
				(!increasing && curr > last) ||
				curr == last {
				bad++
			}

			diff := Abs(curr - last)
			if diff > 3 {
				bad++
			}
			last = curr
		}

		if bad > 1 {
			continue report
		}

		tot++
	}

	fmt.Println(tot)
}
