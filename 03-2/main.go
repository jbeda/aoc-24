package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	dat, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	re := regexp.MustCompile(`mul\((\d*),(\d*)\)|do\(\)|don't\(\)`)

	var tot int
	enabled := true

	for _, v := range re.FindAllSubmatch(dat, -1) {
		fmt.Printf("enabled: %t - %q\n", enabled, v)

		switch string(v[0]) {
		case "do()":
			enabled = true
		case "don't()":
			enabled = false
		default:
			a, err := strconv.Atoi(string(v[1]))
			if err != nil {
				log.Fatal(err)
			}

			b, err := strconv.Atoi(string(v[2]))
			if err != nil {
				log.Fatal(err)
			}
			if enabled {
				tot += a * b
			}
		}
	}

	fmt.Println(tot)
}
