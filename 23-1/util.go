package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"runtime"
	"runtime/debug"
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

func ReadFileLines(path string) []string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	lines := make([]string, 0)
	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		lines = append(lines, line)
	}

	return lines
}

var Debug bool

func DebugLogf(format string, v ...interface{}) {
	if Debug {
		fmt.Printf(format, v...)
	}
}

func Assert(cond bool, msg string) {
	if !cond {
		debug.PrintStack()
		log.Print(msg)
		runtime.Breakpoint()
	}
}
