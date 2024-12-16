package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
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

// -------------------------------------
type Vector struct {
	X, Y int
}

func (v1 Vector) Add(v2 Vector) Vector {
	return Vector{v1.X + v2.X, v1.Y + v2.Y}
}

func (v1 Vector) AddInt(i int) Vector {
	return Vector{v1.X + i, v1.Y + i}
}

func (v1 Vector) Sub(v2 Vector) Vector {
	return Vector{v1.X - v2.X, v1.Y - v2.Y}
}

func (v1 Vector) SubInt(i int) Vector {
	return Vector{v1.X - i, v1.Y - i}
}

func (v Vector) Mul(i int) Vector {
	return Vector{v.X * i, v.Y * i}
}

func (v Vector) Div(d int) (res Vector, rem Vector) {
	res = Vector{v.X / d, v.Y / d}
	rem = Vector{v.X % d, v.Y % d}
	return
}

func (v Vector) Abs() Vector {
	return Vector{AbsInt(v.X), AbsInt(v.Y)}
}

func (v Vector) Wrap(w Vector) Vector {
	if v.X < 0 {
		v.X += w.X
	}
	if v.Y < 0 {
		v.Y += w.Y
	}
	return Vector{v.X % w.X, v.Y % w.Y}
}

func (v *Vector) String() string {
	return fmt.Sprintf("(%d, %d)", v.X, v.Y)
}

// -------------------------------------
type Robot struct {
	Pos Vector
	Vel Vector
}

func (r Robot) String() string {
	return fmt.Sprintf("Pos: %v, Vel: %v", r.Pos, r.Vel)
}

// -------------------------------------
type Board struct {
	Robots []Robot
	Size   Vector
}

func (b *Board) AddRobot(r Robot) {
	b.Robots = append(b.Robots, r)
}

func (b *Board) Step() {
	for i := range b.Robots {
		r := &b.Robots[i]
		r.Pos = r.Pos.Add(r.Vel)
		r.Pos = r.Pos.Wrap(b.Size)
	}
}

func (b *Board) QuadCounts() (q [4]int) {
	mid, _ := b.Size.Div(2)
	for _, r := range b.Robots {
		if r.Pos.X < mid.X {
			if r.Pos.Y < mid.Y {
				q[0]++
			} else if r.Pos.Y > mid.Y {
				q[1]++
			}
		} else if r.Pos.X > mid.X {
			if r.Pos.Y < mid.Y {
				q[2]++
			} else if r.Pos.Y > mid.Y {
				q[3]++
			}
		}
	}
	return
}

func (b *Board) Score() int {
	q := b.QuadCounts()
	return q[0] * q[1] * q[2] * q[3]
}

func (b Board) String() string {
	cells := b.Grid()

	var res string
	for y := 0; y < b.Size.Y; y++ {
		for x := 0; x < b.Size.X; x++ {
			if cells[y][x] == 0 {
				res += "."
			} else {
				res += strconv.Itoa(cells[y][x])
			}
		}
		res += "\n"
	}
	return res
}

func (b Board) Grid() [][]int {
	cells := make([][]int, b.Size.Y)
	for i := range cells {
		cells[i] = make([]int, b.Size.X)
	}

	for _, r := range b.Robots {
		cells[r.Pos.Y][r.Pos.X]++
	}
	return cells
}

// -------------------------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	infile := "input.txt"
	size := Vector{101, 103}
	// infile := "test.txt"
	// size := Vector{11, 7}

	fmt.Println("Size:", size)

	f, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	board := Board{Size: size}

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		r := Robot{}

		// p=0,4 v=3,-3
		line := scan.Text()

		re := regexp.MustCompile(`p=(-?\d+),(-?\d+) v=(-?\d+),(-?\d+)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) != 5 {
			log.Fatalf("Couldn't parse line: %q", line)
		}

		r.Pos.X = MustAtoi(matches[1])
		r.Pos.Y = MustAtoi(matches[2])
		r.Vel.X = MustAtoi(matches[3])
		r.Vel.Y = MustAtoi(matches[4])

		fmt.Printf("Robot: %v\n", r)
		board.AddRobot(r)
	}

	// fmt.Println("Initial:")
	// fmt.Println(board)
	// fmt.Println()

	for i := 0; i < 1000000; i++ {
		if i%1000 == 0 {
			fmt.Println("Step:", i+1)
		}
		board.Step()

		grid := board.Grid()

		longestRun := 0
		for y := 0; y < size.Y; y++ {
			run := 0
			for x := 0; x < size.X; x++ {
				if grid[y][x] > 0 {
					run++
					if run > longestRun {
						longestRun = run
					}
				} else {
					run = 0
				}
			}
		}

		if longestRun > 10 {
			fmt.Println("Step:", i+1)
			fmt.Println(board)
			fmt.Println()
			//break outer
		}
	}
}
