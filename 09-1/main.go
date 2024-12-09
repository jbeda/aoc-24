package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type FileID int

const FreeSpace FileID = -1

type FileSystem []FileID

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	binput, err := os.ReadFile("input.txt")
	if err != nil {
		log.Fatal(err)
	}

	input := strings.TrimSpace(string(binput))

	fs := FileSystem{}

	// Initialize the filesystem
	var nextFileID FileID = 0
	var isFileNext = true
	for _, d := range input {
		numBlocks, err := strconv.Atoi(string(d))
		if err != nil {
			log.Fatal(err)
		}

		for i := 0; i < numBlocks; i++ {
			if isFileNext {
				fs = append(fs, nextFileID)
			} else {
				fs = append(fs, FreeSpace)
			}
		}

		if isFileNext {
			nextFileID++
		}
		isFileNext = !isFileNext
	}

	// Optimize the filesystem
	front := 0
	back := len(fs) - 1
	for {
		for fs[front] != FreeSpace {
			front++
		}
		for fs[back] == FreeSpace {
			back--
		}
		if front >= back {
			break
		}
		fs[front], fs[back] = fs[back], fs[front]
		front++
		back--
	}

	// Create the "checksum"
	checksum := 0
	for i, id := range fs {
		if id != FreeSpace {
			checksum += i * int(id)
		}
	}
	fmt.Println(checksum)
}
