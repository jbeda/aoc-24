package main

import (
	"fmt"
	"log"
	"time"
)

func Generate(i int) int {
	i2 := i << 6        // i * 64
	i = i ^ i2          // i XOR i2 (mix)
	i = i & (1<<24 - 1) // i MOD 2^24 (prune)

	i2 = i >> 5         // i2 / 32
	i = i ^ i2          // i XOR i2 (mix)
	i = i & (1<<24 - 1) // i MOD 2^24 (prune)

	i2 = i << 11        // i * 2048
	i = i ^ i2          // i XOR i2 (mix)
	i = i & (1<<24 - 1) // i MOD 2^24 (prune)

	return i
}

func GenerateN(i int, n int) int {
	for j := 0; j < n; j++ {
		i = Generate(i)
	}
	return i
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	lines := ReadFileLines("input.txt")
	var inputs []int
	for _, line := range lines {
		inputs = append(inputs, MustAtoi(line))
	}
	_ = inputs

	sum := 0
	for _, i := range inputs {
		DebugLogf("%d: ", i)
		i = GenerateN(i, 2000)
		DebugLogf("%d\n", i)
		sum += i
	}
	fmt.Println("Sum:", sum)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
