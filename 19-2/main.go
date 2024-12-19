package main

import (
	"fmt"
	"log"
	"strings"
	"time"
)

type Trie struct {
	root  *TrieNode
	cache map[string]int
}

type TrieNode struct {
	Children map[rune]*TrieNode

	Token bool
}

func NewTrieNode() *TrieNode {
	return &TrieNode{
		Children: make(map[rune]*TrieNode),
	}
}

func NewTrie() *Trie {
	return &Trie{
		root:  NewTrieNode(),
		cache: make(map[string]int),
	}
}

func (t *Trie) InsertToken(word string) {
	node := t.root
	for _, r := range word {
		child, ok := node.Children[r]
		if !ok {
			child = NewTrieNode()
			node.Children[r] = child
		}
		node = child
	}
	node.Token = true
}

func (t *Trie) Tokens() []string {
	return t.root.Tokens()
}

func (t *TrieNode) Tokens() []string {
	tokens := make([]string, 0)

	for r, child := range t.Children {
		if child.Token {
			tokens = append(tokens, string(r))
		}

		for _, token := range child.Tokens() {
			tokens = append(tokens, string(r)+token)
		}
	}

	return tokens
}

func (t *Trie) Solve(word string) int {
	tot := 0
	node := t.root

	if total, ok := t.cache[word]; ok {
		return total
	}

	for i, r := range word {
		if node.Token {
			remtotal := t.Solve(word[i:])
			tot += remtotal
		}

		child, ok := node.Children[r]
		if !ok {
			t.cache[word] = tot
			return tot
		}

		node = child
	}

	if node.Token {
		tot++
	}

	t.cache[word] = tot

	return tot
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
		fmt.Print(input, " ")
		inputtot := t.Solve(input)

		fmt.Println(inputtot)

		tot += inputtot
	}

	fmt.Println(tot)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
