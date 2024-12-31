package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type LockType int

const (
	LTKey LockType = iota
	LTLock
)

const NumTumblers = 5
const NumLevels = 5

type LockPart struct {
	Type    LockType
	Profile [NumTumblers]int
}

// Returns a lock part (key or lock) and the number of lines consumed.
func LoadLockPart(lines []string) (LockPart, int) {
	var part LockPart
	part.Profile = [NumTumblers]int{}

	// First line (. or #) tells us type of part.
	if lines[0][0] == '.' {
		Assert(lines[0] == strings.Repeat(".", NumTumblers), "Invalid key part")
		part.Type = LTKey
	} else {
		Assert(lines[0] == strings.Repeat("#", NumTumblers), "Invalid key part")
		part.Type = LTLock
	}

	lines = lines[1:]
	for i := 0; i < NumLevels; i++ {
		line := lines[i]
		Assert(len(line) == NumTumblers, "Invalid line length")
		for j, c := range line {
			if c == '#' {
				part.Profile[j]++
			}
		}
	}

	// Lines used: 2 for bookends, NumLevels for profile.
	return part, NumLevels + 2
}

type Locks []LockPart
type Keys []LockPart

func LoadLocksAndKeys(lines []string) (Locks, Keys) {
	var locks Locks
	var keys Keys

	for len(lines) > 0 {
		for len(lines) > 0 && lines[0] == "" {
			lines = lines[1:]
		}

		part, consumed := LoadLockPart(lines)
		lines = lines[consumed:]

		if part.Type == LTKey {
			keys = append(keys, part)
		} else {
			locks = append(locks, part)
		}
	}

	return locks, keys
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	lines := ReadFileLines("input.txt")
	locks, keys := LoadLocksAndKeys(lines)
	fmt.Printf("Loaded %d locks and %d keys\n", len(locks), len(keys))

	// Simply compare each key with each lock.
	totalFits := 0
	for _, lock := range locks {
	NextKey:
		for _, key := range keys {
			for i := 0; i < NumTumblers; i++ {
				if key.Profile[i]+lock.Profile[i] > NumLevels {
					continue NextKey
				}
			}
			totalFits++
		}
	}

	fmt.Println("Total fits:", totalFits)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
