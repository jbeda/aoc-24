package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Trie struct {
	Children map[rune]*Trie

	Token bool
}

func NewTrie() *Trie {
	return &Trie{
		Children: make(map[rune]*Trie),
	}
}

func (t *Trie) InsertToken(word string) {
	node := t
	for _, r := range word {
		child, ok := node.Children[r]
		if !ok {
			child = NewTrie()
			node.Children[r] = child
		}
		node = child
	}
	node.Token = true
}

func (t *Trie) Solve(root *Trie, word string) bool {
	node := t
	for i, r := range word {
		child, ok := node.Children[r]
		if !ok {
			if node.Token {
				return root.Solve(root, word[i:])
			}
			return false
		}

		if node.Token && root.Solve(root, word[i:]) {
			return true
		}

		node = child
	}
	return node.Token
}

//---------------------------------------------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	t := NewTrie()

	lines := ReadFileLines("input.txt")

	tokenLines := lines[0]
	tokens := strings.Split(tokenLines, ", ")

	for _, token := range tokens {
		t.InsertToken(token)
	}

	// Skip a blank line before we get all the inputs
	inputs := lines[2:]

	tot := 0
	for _, input := range inputs {
		if t.Solve(t, input) {
			tot++
		}
	}

	fmt.Println(tot)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
