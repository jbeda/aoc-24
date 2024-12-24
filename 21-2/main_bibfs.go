//go:build exclude

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

type CacheKey struct {
	From, To state
}

// Cache the distance to get between any two states
var DistCache = map[CacheKey]int{}

type Frontier struct {
	D int // Distance
	S state
}

type SearchState struct {
	Frontiers []Frontier
}

// Map of the start and current frontiers
var ForwardSearchCache = map[state][]Frontier{}
var BackwardSearchCache = map[state][]Frontier{}

// Do a BFS to find the minimum number of moves to get to the target
func FindShortestDist(from state, to state) (int, bool) {
	var forwardFrontiers, backwardFrontiers []Frontier
	_ = backwardFrontiers

	DebugLogf("Finding path from %s to %s\n", from, to)

	if dist, ok := DistCache[CacheKey{from, to}]; ok {
		DebugLogf("Found cached path from %s to %s, length: %d\n", from, to, dist)
		return dist, true
	}

	// See if this is a search that we can pick up again
	if cachedFrontiers, ok := ForwardSearchCache[from]; ok {
		DebugLogf("Reloading forward search from cache\n")
		forwardFrontiers = cachedFrontiers
	} else {
		// We know that distance from 'from' to 'from' is 0
		DistCache[CacheKey{from, from}] = 0

		forwardFrontiers = []Frontier{{0, from}}
	}
	if cachedFrontiers, ok := BackwardSearchCache[to]; ok {
		DebugLogf("Reloading backward search from cache\n")
		backwardFrontiers = cachedFrontiers
	} else {
		// We know that distance from 'to' to 'to' is 0
		DistCache[CacheKey{to, to}] = 0

		backwardFrontiers = []Frontier{{0, to}}
	}

	backwards := false
	for len(forwardFrontiers) > 0 || len(backwardFrontiers) > 0 {
		if !backwards {
			if len(forwardFrontiers) > 0 {
				frontier := forwardFrontiers[0]
				forwardFrontiers = forwardFrontiers[1:]

				DebugLogf("F Popping frontier: %s\n", frontier.S)

				// Search cache for the rest of the path
				if dist, ok := DistCache[CacheKey{frontier.S, to}]; ok {
					// Found the rest of the path in the cache
					totalDist := frontier.D + dist
					DebugLogf("F Found cached path from %s to %s, length: %d\n"+
						"  Total path from %s to %s is %d\n",
						frontier.S, to, dist,
						from, to, totalDist)
					DistCache[CacheKey{from, to}] = totalDist

					// Save the state of this search so we can continue with another target
					ForwardSearchCache[from] = forwardFrontiers
					BackwardSearchCache[to] = backwardFrontiers

					return totalDist, true
				}

				nextStates := NextStates(frontier.S, false)
				for _, nextState := range nextStates {
					cacheKey := CacheKey{from, nextState}
					if _, ok := DistCache[cacheKey]; !ok {
						newDist := frontier.D + 1
						forwardFrontiers = append(forwardFrontiers, Frontier{newDist, nextState})
						DistCache[cacheKey] = newDist
						DebugLogf("F Adding frontier: %s, Dist: %d\n", nextState, newDist)
					} else {
						DebugLogf("F Skipping frontier: %s\n", nextState)
					}
				}
			}
		} else {
			if len(backwardFrontiers) > 0 {
				frontier := backwardFrontiers[0]
				backwardFrontiers = backwardFrontiers[1:]

				DebugLogf("B Popping frontier: %s\n", frontier.S)

				// Search cache for the rest of the path
				if dist, ok := DistCache[CacheKey{from, frontier.S}]; ok {
					// Found the rest of the path in the cache
					totalDist := frontier.D + dist
					DebugLogf("B Found cached path from %s to %s, length: %d\n"+
						"  Total path from %s to %s is %d\n",
						from, frontier.S, dist,
						from, to, totalDist)
					DistCache[CacheKey{from, to}] = totalDist

					// Save the state of this search so we can continue with another target
					ForwardSearchCache[from] = forwardFrontiers
					BackwardSearchCache[to] = backwardFrontiers

					return totalDist, true
				}

				nextStates := NextStates(frontier.S, true)
				for _, nextState := range nextStates {
					cacheKey := CacheKey{nextState, to}
					if _, ok := DistCache[cacheKey]; !ok {
						newDist := frontier.D + 1
						backwardFrontiers = append(backwardFrontiers, Frontier{newDist, nextState})
						DistCache[cacheKey] = newDist
						DebugLogf("B Adding frontier: %s, Dist: %d\n", nextState, newDist)
					} else {
						DebugLogf("B Skipping frontier: %s\n", nextState)
					}
				}
			}
		}

		backwards = !backwards
	}

	return 0, false
}

// --------------------------------------------------------------------

const NumLayers = 5

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

	complexity := 0

	for _, code := range codes {
		totalDist := 0

		startState := state(strings.Repeat("A", NumLayers))
		for _, digit := range code {
			destState := state(strings.Repeat("A", NumLayers-1) + string(digit))
			ForwardSearchCache = map[state][]Frontier{}
			BackwardSearchCache = map[state][]Frontier{}
			DistCache = map[CacheKey]int{}
			dist, ok := FindShortestDist(state(startState), destState)
			if !ok {
				log.Fatalf("No path found for %s", code)
			}
			totalDist += dist + 1 // Account for pressing A
			startState = destState
		}
		fmt.Printf("Code: %s, Dist: %d\n", code, totalDist)
		complexity += totalDist * MustAtoi(code[0:3])
	}

	fmt.Println("Complexity:", complexity)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
