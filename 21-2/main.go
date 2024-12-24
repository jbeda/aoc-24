package main

import (
	"container/heap"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Action int
type State rune

const (
	Up Action = iota
	Right
	Down
	Left
	PressOld
	Invalid = -1
)

var Dirs = []Action{Up, Right, Down, Left}
var Actions = []Action{Up, Right, Down, Left, PressOld}

var ActionReverse = map[Action]Action{
	Up:       Down,
	Right:    Left,
	Down:     Up,
	Left:     Right,
	PressOld: PressOld,
}

var ActionMap = map[Action]State{
	Up:       '^',
	Right:    '>',
	Down:     'v',
	Left:     '<',
	PressOld: 'A',
}

func (a Action) String() string {
	return string(ActionMap[a])
}

func (a Action) State() State {
	return ActionMap[a]
}

func (a Action) Reverse() Action {
	return ActionReverse[a]
}

var ActionMapInv = map[State]Action{
	'^': Up,
	'>': Right,
	'v': Down,
	'<': Left,
	'A': PressOld,
}

func (s State) Action() Action {
	if a, ok := ActionMapInv[s]; ok {
		return a
	}

	return Invalid
}

func (s State) String() string {
	return string(s)
}

// Maps from the position, direction and resulting
const NULL State = 0

type MoveMap map[State][4]State

var keypadMoves = MoveMap{
	'A': {'3', NULL, NULL, '0'},
	'0': {'2', 'A', NULL, NULL},
	'1': {'4', '2', NULL, NULL},
	'2': {'5', '3', '0', '1'},
	'3': {'6', NULL, 'A', '2'},
	'4': {'7', '5', '1', NULL},
	'5': {'8', '6', '2', '4'},
	'6': {'9', NULL, '3', '5'},
	'7': {NULL, '8', '4', NULL},
	'8': {NULL, '9', '5', '7'},
	'9': {NULL, NULL, '6', '8'},
}
var robotMoves = MoveMap{
	'A': {NULL, NULL, '>', '^'},
	'<': {NULL, 'v', NULL, NULL},
	'v': {'^', '>', NULL, '<'},
	'>': {'A', NULL, NULL, 'v'},
	'^': {NULL, 'A', 'v', NULL},
}

// --------------------------------------------------------------------
type CacheKey struct {
	From State
	To   State
}

type Machine struct {
	Name   string
	Depth  int
	Pos    State
	Moves  MoveMap
	Parent *Machine

	Cache map[CacheKey]int
}

func NewMachine(name string, depth int, moves MoveMap) *Machine {
	return &Machine{
		Name:  name,
		Depth: depth,
		Pos:   'A',
		Moves: moves,
		Cache: make(map[CacheKey]int),
	}
}

func (m *Machine) DebugLogf(format string, v ...interface{}) {
	if Debug {
		prefix := fmt.Sprintf("%s%s ", strings.Repeat(" ", m.Depth), m.Name)
		s := fmt.Sprintf(format, v...)
		log.Output(2, prefix+s)
	}
}

type QueueItem struct {
	S        State
	ParentS  State
	D        int
	Pressing bool
}

func (qi *QueueItem) Score() int {
	return qi.D
}

func (qi *QueueItem) String() string {
	return fmt.Sprintf("[S: %s, pS: %s D: %d P: %v]", qi.S, qi.ParentS, qi.D, qi.Pressing)
}

type PriorityQueue []*QueueItem

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].Score() < pq[j].Score() }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(*QueueItem)) }
func (pq *PriorityQueue) Pop() interface{} {
	n := len(*pq)
	x := (*pq)[n-1]
	*pq = (*pq)[:n-1]
	return x
}

func (m *Machine) Press(from, to State) int {
	m.DebugLogf("%s -> %s Press called", from, to)
	if from == to {
		m.DebugLogf("%s -> %s Same state. Dist: 1\n", from, to)
		return 1
	}

	if dist, ok := m.Cache[CacheKey{from, to}]; ok {
		m.DebugLogf("%s -> %s Cache hit: %d\n", from, to, dist)
		return dist
	}

	frontiers := PriorityQueue{}
	heap.Init(&frontiers)
	// We assume the parent is starting on 'A'
	qi := &QueueItem{from, 'A', 0, false}
	heap.Push(&frontiers, qi)
	m.DebugLogf("%s -> %s Pushing: %s\n", from, to, qi)

	visited := make(map[State]bool)
	visited[from] = true

	for len(frontiers) > 0 {
		frontier := heap.Pop(&frontiers).(*QueueItem)
		m.DebugLogf("%s -> %s Popping: %s\n", from, to, frontier)

		visited[frontier.S] = true

		if frontier.S == to {
			if frontier.Pressing {
				m.DebugLogf("%s -> %s Found path: %d\n", from, to, frontier.D)
				m.Cache[CacheKey{from, to}] = frontier.D
				m.DebugLogf("%s -> %s Cache set: %d\n", from, to, frontier.D)
				return frontier.D
			}

			m.DebugLogf("%s -> %s Found path (w/o press): %d\n", from, to, frontier.D)

			if m.Parent == nil {
				m.DebugLogf("%s -> %s No parent. User moves to %s\n", from, to, to)
				m.Cache[CacheKey{from, to}] = frontier.D + 1
				m.DebugLogf("%s -> %s Cache set: %d\n", from, to, frontier.D+1)
				return frontier.D + 1
			}

			nextDist := m.Parent.Press(frontier.ParentS, 'A')
			dist := frontier.D + nextDist

			qi := &QueueItem{frontier.S, frontier.ParentS, dist, true}
			heap.Push(&frontiers, qi)
			m.DebugLogf("%s -> %s Pushing: %s\n", from, to, qi)
			continue
		}

		nextStates := m.Moves[frontier.S]
		for i, next := range nextStates {
			if next == NULL || visited[next] {
				continue
			}

			m.DebugLogf("%s -> %s Curr: %s Next: %s\n", from, to, frontier.S, next)

			var nextDist int
			var parentS State = NULL

			// If we don't have a parent then there is a human at the controls and the
			// cost of each move is 1
			if m.Parent == nil {
				m.DebugLogf("%s -> %s No parent. User moves to %s\n", from, to, Action(i).State())
				nextDist = 1
			} else {
				parentKey := Action(i).State()
				nextDist = m.Parent.Press(frontier.ParentS, parentKey)
				parentS = parentKey
			}

			dist := frontier.D + nextDist

			qi := &QueueItem{next, parentS, dist, false}
			heap.Push(&frontiers, qi)
			m.DebugLogf("%s -> %s Pushing: %s\n", from, to, qi)
		}
	}

	Assert(false, "No path found")
	return 0
}

// --------------------------------------------------------------------

const NumLayers = 26

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	Debug = false
	timeStart := time.Now()

	var codes []string

	// Input date
	test := false
	if test {
		codes = []string{
			"029A",
			"980A",
			"179A",
			"456A",
			"379A",
		}
	} else {
		codes = []string{
			"964A",
			"246A",
			"973A",
			"682A",
			"180A",
		}
	}

	// Build machine chain
	keypad := NewMachine("kp", 0, keypadMoves)
	{
		c := keypad
		for i := 0; i < NumLayers-1; i++ {
			robot := NewMachine("m"+strconv.Itoa(i), i+1, robotMoves)
			c.Parent = robot
			c = robot
		}
	}

	complexity := 0

	for _, code := range codes {
		totalDist := 0

		currDigit := State('A')
		for _, nextRune := range code {
			nextDigit := State(nextRune)
			dist := keypad.Press(currDigit, nextDigit)
			fmt.Printf("Pressing %s -> %s: %d\n", currDigit, nextDigit, dist)
			totalDist += dist // Account for pressing A
			currDigit = nextDigit
		}
		fmt.Printf("Code: %s, Dist: %d\n", code, totalDist)
		currComplexity := totalDist * MustAtoi(code[0:3])
		complexity += currComplexity
		fmt.Printf("  Complexity: %d\n", currComplexity)
	}

	fmt.Println("Complexity:", complexity)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
