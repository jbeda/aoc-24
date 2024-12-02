package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
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

	var in1, in2 []int

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()

		ss := strings.Fields(line)
		s1, s2 := ss[0], ss[1]
		n1, err := strconv.Atoi(s1)
		if err != nil {
			log.Fatal(err)
		}
		in1 = append(in1, n1)

		n2, err := strconv.Atoi(s2)
		if err != nil {
			log.Fatal(err)
		}
		in2 = append(in2, n2)
	}

	slices.Sort(in1)
	slices.Sort(in2)

	var tot int = 0
	for i := range in1 {
		tot += Abs(in1[i] - in2[i])
	}

	fmt.Println(tot)
}
