package main

import (
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
func NextState(s state, action rune) *state {
	var nextState *state

	// This would only happen if it was A all the way down.
	if len(s) == 0 {
		return nil
	}

	currPos := rune(s[0])
	if action == 'A' {
		nextSubState := NextState(s[1:], currPos)
		if nextSubState != nil {
			s = state(currPos) + *nextSubState
			nextState = &s
		}
	} else {
		moveMap := robotMoves
		if len(s) == 1 {
			moveMap = keypadMoves
		}

		new0State := moveMap[currPos][ActionMapInv[action]]
		if new0State == NULL {
			return nil
		}

		s = state(new0State) + s[1:]
		nextState = &s
	}

	return nextState
}

type CacheKey struct {
	From, To state
}

// Cache the path to get from the start state to any specific state
var PathCache = map[CacheKey]string{}

type Frontier struct {
	P string
	S state
}

type SearchState struct {
	Frontiers []Frontier
}

// Map of the start and current frontiers
var SearchCache = map[state][]Frontier{}

// Do a BFS to find the minimum number of moves to get to the target
func FindShortestPath(from state, to state) (string, bool) {
	var frontiers []Frontier

	DebugLogf("Finding path from %s to %s\n", from, to)

	// See if this is a search that we can pick up again
	if cachedFrontiers, ok := SearchCache[from]; ok {
		DebugLogf("Reloading search from cache\n")
		if path, ok := PathCache[CacheKey{from, to}]; ok {
			return path, true
		}
		frontiers = cachedFrontiers
	} else {
		// We know that distance from 'from' to 'from' is 0
		PathCache[CacheKey{from, from}] = ""

		frontiers = []Frontier{{"", from}}
	}

	for len(frontiers) > 0 {
		frontier := frontiers[0]
		frontiers = frontiers[1:]

		DebugLogf("Popping frontier: %s\n", frontier.P)

		if frontier.S == to {
			// Save the state of this search so we can continue with another target
			SearchCache[from] = frontiers
			return frontier.P, true
		}

		for _, action := range Actions {
			nextState := NextState(frontier.S, ActionMap[action])
			if nextState == nil {
				continue
			}

			if _, ok := PathCache[CacheKey{from, *nextState}]; !ok {
				newPath := frontier.P + string(ActionMap[action])
				frontiers = append(frontiers, Frontier{newPath, *nextState})
				PathCache[CacheKey{from, *nextState}] = newPath
				DebugLogf("Adding frontier: %s\n", *nextState)
			}
		}
	}

	return "", false
}

// --------------------------------------------------------------------
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
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

	const NumLayers = 3

	complexity := 0

	for _, code := range codes {
		path := ""

		startState := state(strings.Repeat("A", NumLayers))
		for _, digit := range code {
			destState := state(strings.Repeat("A", NumLayers-1) + string(digit))
			subPath, ok := FindShortestPath(state(startState), destState)
			if !ok {
				log.Fatalf("No path found for %s", code)
			}

			fmt.Printf("Pressing %s -> %s: %d\n", string(startState[NumLayers-1]), string(destState[NumLayers-1]), len(subPath)+1)

			path += subPath + "A"
			startState = destState
		}
		currComplexity := len(path) * MustAtoi(code[0:3])
		fmt.Printf("Code: %s, Path: %s\n", code, path)
		fmt.Printf("  Complexity: %d\n", currComplexity)
		complexity += len(path) * MustAtoi(code[0:3])
	}

	fmt.Println("Complexity:", complexity)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
