package main

import (
	"log"
	"strconv"
)

func MustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	return i
}

func AbsInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
