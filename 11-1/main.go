package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// input := "125 17"
	input := "5 89749 6061 43 867 1965860 0 206250"

	ss := strings.Fields(input)

	var stones []int

	for _, s := range ss {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		stones = append(stones, i)
	}

	for i := 0; i < 25; i++ {
		var stones2 []int

		for _, s := range stones {
			if s == 0 {
				stones2 = append(stones2, 1)
				continue
			}

			sText := strconv.Itoa(s)
			if len(sText)%2 == 0 {
				s1, err := strconv.Atoi(sText[:len(sText)/2])
				if err != nil {
					log.Fatal(err)
				}
				s2, err := strconv.Atoi(sText[len(sText)/2:])
				if err != nil {
					log.Fatal(err)
				}
				stones2 = append(stones2, s1, s2)
				continue
			}

			stones2 = append(stones2, s*2024)
		}

		stones = stones2
	}

	fmt.Println(len(stones))
}
