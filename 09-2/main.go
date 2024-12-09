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

type FileEntry struct {
	ID       FileID
	Length   int
	Location int
}

type FileEntries []FileEntry

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

	// Build map of file size and location
	files := FileEntries{}
	freeSpans := FileEntries{}

	for i := 0; i < len(fs); {
		id := fs[i]

		fe := FileEntry{}
		fe.ID = id
		fe.Location = i
		fe.Length = 0
		for i < len(fs) && id == fs[i] {
			i++
			fe.Length++
		}

		if id == FreeSpace {
			freeSpans = append(freeSpans, fe)
		} else {
			files = append(files, fe)
		}
	}

	// Optimize the filesystem
Files:
	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]

		// Find the first free span that can fit the file
		for j := 0; j < len(freeSpans); j++ {
			free := freeSpans[j]
			if free.Location > file.Location {
				// We didn't find a free span to use
				continue Files
			}

			if free.Length >= file.Length {
				// Move the file
				for k := 0; k < file.Length; k++ {
					fs[free.Location+k] = file.ID
					fs[file.Location+k] = FreeSpace
				}

				// Update the free span
				freeSpans[j].Location += file.Length
				freeSpans[j].Length -= file.Length

				break
			}
		}
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
