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

	re := regexp.MustCompile(`mul\((\d*),(\d*)\)`)

	var tot int

	for _, v := range re.FindAllSubmatch(dat, -1) {
		fmt.Printf("%q\n", v)

		as := v[1]
		a, err := strconv.Atoi(string(as))
		if err != nil {
			log.Fatal(err)
		}

		bs := v[2]
		b, err := strconv.Atoi(string(bs))
		if err != nil {
			log.Fatal(err)
		}

		tot += a * b
	}

	fmt.Println(tot)
}
