//go:build exclude

package main

import (
	"container/heap"
	"fmt"
	"log"
	"strings"
	"time"
)

type Action int

const (
	Up Action = iota
	Right
	Down
	Left
	Press
)

var Dirs = []Action{Up, Right, Down, Left}
var Actions = []Action{Up, Right, Down, Left, Press}

var ActionReverse = map[Action]Action{
	Up:    Down,
	Right: Left,
	Down:  Up,
	Left:  Right,
	Press: Press,
}

var ActionMap = map[Action]rune{
	Up:    '^',
	Right: '>',
	Down:  'v',
	Left:  '<',
	Press: 'A',
}

func (a Action) String() string {
	return string(ActionMap[a])
}

var ActionMapInv = map[rune]Action{
	'^': Up,
	'>': Right,
	'v': Down,
	'<': Left,
	'A': Press,
}

// Maps from the position, direction and resulting
const NULL rune = 0

type MoveMap map[rune][4]rune

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

// The state is a string that is 3 digits long that represents the position of
// the two intermediate robots along with the final robot (in front of the
// keypad), in that order.
type state string

// Returns possible next states from a state, return is the results from
// pressing DIR along with, the final state: A.
func NextState(s state, action rune, backward bool) *state {
	var nextState *state

	// This would only happen if it was A all the way down.
	if len(s) == 0 {
		return nil
	}

	currPos := rune(s[0])
	if action == 'A' {
		nextSubState := NextState(s[1:], currPos, backward)
		if nextSubState != nil {
			s = state(currPos) + *nextSubState
			nextState = &s
		}
	} else {
		moveMap := robotMoves
		if len(s) == 1 {
			moveMap = keypadMoves
		}

		actionEnum := ActionMapInv[action]
		if backward {
			actionEnum = ActionReverse[actionEnum]
		}

		new0State := moveMap[currPos][actionEnum]
		if new0State == NULL {
			return nil
		}

		s = state(new0State) + s[1:]
		nextState = &s
	}

	return nextState
}

func NextStates(s state, backward bool) []state {
	var nextStates []state

	for _, action := range Actions {
		if nextState := NextState(s, ActionMap[action], backward); nextState != nil {
			nextStates = append(nextStates, *nextState)
		}
	}

	return nextStates
}

// --------------------------------------------------------------------

type FromTo struct {
	From, To state
}

var keypadDistances = map[FromTo]int{}

func InitKeypadDistances() {
	pos := map[state]Vector{
		"7": {0, 0},
		"8": {1, 0},
		"9": {2, 0},
		"4": {0, 1},
		"5": {1, 1},
		"6": {2, 1},
		"1": {0, 2},
		"2": {1, 2},
		"3": {2, 2},
		"0": {1, 3},
		"A": {2, 3},
	}

	for from, fromPos := range pos {
		for to, toPos := range pos {
			keypadDistances[FromTo{from, to}] = fromPos.ManhattanDist(toPos)
		}
	}
}

func DistEstimate(from, to state) int {
	keypadFrom := state(from[len(from)-1])
	keypadTo := state(to[len(to)-1])

	estDist := keypadDistances[FromTo{keypadFrom, keypadTo}]

	if len(from) > 1 {
		intFrom := state(from[0 : len(from)-1])
		intTo := state(to[0 : len(to)-1])

		for i, fromDigit := range intFrom {
			if fromDigit != rune(intTo[i]) {
				estDist += 1
			}
		}
	}

	return estDist
}

type QueueItem struct {
	State     state
	DistFrom  int
	EstDistTo int
}

func (qi *QueueItem) Score() int {
	return qi.DistFrom + qi.EstDistTo
}

func (qi *QueueItem) String() string {
	return fmt.Sprintf("[State: %s, DistFrom: %d, EstDistTo: %d]", qi.State, qi.DistFrom, qi.EstDistTo)
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

var DistCache = map[FromTo]int{}

func AStarFindShortestDist(from, to state) (int, bool) {
	pq := PriorityQueue{}
	start := &QueueItem{from, 0, DistEstimate(from, to)}
	heap.Push(&pq, start)

	DistCache[FromTo{from, from}] = 0
	visited := map[state]bool{}
	visited[from] = true

	DebugLogf("Finding shortest path from %s to %s\n", from, to)
	DebugLogf("Pushing start: %s Queue Len: %d\n", start, len(pq))

	for len(pq) > 0 {
		item := heap.Pop(&pq).(*QueueItem)

		DebugLogf("Popping item: %s\n", item)

		if item.State == to {
			DebugLogf("Found path. Dist: %d\n", item.DistFrom)
			DistCache[FromTo{from, to}] = item.DistFrom
			return item.DistFrom, true
		}

		for _, nextState := range NextStates(item.State, false) {
			// Check if we have already visited this state
			if _, ok := visited[nextState]; ok {
				continue
			}
			visited[nextState] = true

			distFrom := item.DistFrom + 1

			// Put this partial distance in the cache
			DistCache[FromTo{from, nextState}] = distFrom

			// Check if we have a cached distance from this state to the target
			if dist, ok := DistCache[FromTo{nextState, to}]; ok {
				totalDist := distFrom + dist
				DebugLogf("Cache hit. from: %s nextState: %s to: %s dist: %d totalDist: %d\n",
					from, nextState, to, dist, totalDist)
				return totalDist, true
			}

			estDistTo := DistEstimate(nextState, to)
			nextItem := &QueueItem{nextState, distFrom, estDistTo}
			heap.Push(&pq, nextItem)
			DebugLogf("Pushing item: %s Queue Len: %d\n", nextItem, len(pq))
		}
	}

	return 0, false
}

// --------------------------------------------------------------------

const NumLayers = 10

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	InitKeypadDistances()
	Debug = false
	timeStart := time.Now()

	var codes []string

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

	complexity := 0

	for _, code := range codes {
		totalDist := 0

		startState := state(strings.Repeat("A", NumLayers))
		for _, digit := range code {
			destState := state(strings.Repeat("A", NumLayers-1) + string(digit))
			dist, ok := AStarFindShortestDist(startState, destState)
			if !ok {
				log.Fatalf("No path found for %s", code)
			}
			totalDist += dist + 1 // Account for pressing A
			startState = destState
		}
		fmt.Printf("Code: %s, Dist: %d\n", code, totalDist)
		fmt.Printf("  DistCache Size: %d\n", len(DistCache))
		complexity += totalDist * MustAtoi(code[0:3])
	}

	fmt.Println("Complexity:", complexity)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
