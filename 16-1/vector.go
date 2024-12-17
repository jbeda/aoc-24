package main

import "fmt"

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

func (v Vector) Neg() Vector {
	return Vector{-v.X, -v.Y}
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
