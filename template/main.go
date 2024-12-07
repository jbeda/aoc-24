package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		fmt.Println(line)
	}
}
