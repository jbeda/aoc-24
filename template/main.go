package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	lines := ReadFileLines("input.txt")
	fmt.Print(lines)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
