package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

// -------------------------------------
type Matrix2x2 struct {
	A, B, C, D int
}

func (m Matrix2x2) Determinant() int {
	return m.A*m.D - m.B*m.C
}

func (m Matrix2x2) String() string {
	return fmt.Sprintf("[[%d, %d], [%d, %d]]", m.A, m.B, m.C, m.D)
}

// -------------------------------------
type Vector struct {
	X, Y int
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func (v1 Vector) Sub(v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}

func (v Vector) Div(d int) (res Vector, rem Vector) {
	res = Vector{v.X / d, v.Y / d}
	rem = Vector{v.X % d, v.Y % d}
	return
}

func (v *Vector) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

// -------------------------------------

func LinearSolve(a Matrix2x2, c Vector) (Vector, bool) {
	det := a.Determinant()
	if det == 0 {
		return Vector{}, false
	}

	v := Vector{a.D*c.X - a.B*c.Y, a.A*c.Y - a.C*c.X}
	res, rem := v.Div(det)
	if rem.X != 0 || rem.Y != 0 {
		return Vector{}, false
	}
	return res, true
}

// -------------------------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	f, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var totCost int

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		a := Matrix2x2{}
		c := Vector{}

		// Line 1 - Button A: X+94, Y+34
		line1 := scan.Text()
		if len(line1) == 0 {
			continue
		}

		re := regexp.MustCompile(`Button A: X\+(\d*), Y\+(\d*)`)
		matches := re.FindStringSubmatch(line1)
		if len(matches) != 3 {
			log.Fatalf("Couldn't parse line 1: %q", line1)
		}

		i, err := strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		a.A = i

		i, err = strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		a.C = i

		// Line 2 - Button B: X+22, Y+67
		ok := scan.Scan()
		if !ok {
			log.Fatal("Couldn't read line 2")
		}
		line2 := scan.Text()

		re = regexp.MustCompile(`Button B: X\+(\d*), Y\+(\d*)`)
		matches = re.FindStringSubmatch(line2)
		if len(matches) != 3 {
			log.Fatalf("Couldn't parse line 2: %q", line2)
		}

		i, err = strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		a.B = i

		i, err = strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		a.D = i

		// Line 3 - Prize: X=8400, Y=5400
		ok = scan.Scan()
		if !ok {
			log.Fatal("Couldn't read line 3")
		}
		line3 := scan.Text()

		re = regexp.MustCompile(`Prize: X=(\d*), Y=(\d*)`)
		matches = re.FindStringSubmatch(line3)
		if len(matches) != 3 {
			log.Fatalf("Couldn't parse line 3: %q", line3)
		}

		i, err = strconv.Atoi(matches[1])
		if err != nil {
			log.Fatal(err)
		}
		c.X = i

		i, err = strconv.Atoi(matches[2])
		if err != nil {
			log.Fatal(err)
		}
		c.Y = i

		// Now solve
		fmt.Println("A: ", a)
		fmt.Println("C: ", c)
		v, ok := LinearSolve(a, c)
		if !ok {
			fmt.Println("No solution")
		} else {
			fmt.Println("Solution: ", v)
			cost := 3*v.X + v.Y
			fmt.Println("Cost: ", cost)
			totCost += cost
		}
		fmt.Println()
	}

	fmt.Println("Total cost: ", totCost)
}
