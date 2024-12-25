package main

import (
	"fmt"
	"log"
	"maps"
	"regexp"
	"slices"
	"time"
)

type NodeID int
type Graph struct {
	Edges     [][]bool
	Degrees   []int
	NodeNames map[NodeID]string
	NodeIDs   map[string]NodeID
}

type Clique []NodeID
type Cliques struct {
	D []Clique
}

func (c Clique) String(g *Graph) string {
	s := ""
	for _, node := range c {
		s += g.GetNodeName(node) + " "
	}
	return s[:len(s)-1]
}

func (c Clique) Clone() Clique {
	clone := make(Clique, len(c))
	copy(clone, c)
	return clone
}

func (c *Cliques) Add(ids ...NodeID) {
	c.D = append(c.D, Clique(ids).Clone())
}

func (c *Cliques) Merge(other Cliques) {
	c.D = append(c.D, other.D...)
}

func NewGraph() *Graph {
	g := &Graph{
		NodeNames: make(map[NodeID]string),
		NodeIDs:   make(map[string]NodeID),
	}
	g.RegisterNode("INVALID")
	return g
}

func (g *Graph) InitEdges() {
	Assert(len(g.NodeNames) > 0, "No nodes registered")
	numNodes := g.GetNumNodes()
	g.Edges = make([][]bool, numNodes)
	for i := range g.Edges {
		g.Edges[i] = make([]bool, numNodes)
	}
	g.Degrees = make([]int, numNodes)
}

func (g *Graph) RegisterNode(name string) {
	Assert(g.Edges == nil, "Node registered after graph init")
	_ = g.GetNodeID(name)
}

func (g *Graph) GetNodeID(name string) NodeID {
	if id, ok := g.NodeIDs[name]; ok {
		return id
	}
	Assert(g.Edges == nil, "Node registered after graph init")
	id := NodeID(len(g.NodeNames))
	g.NodeNames[id] = name
	g.NodeIDs[name] = id
	return id
}

func (g *Graph) GetNodeName(id NodeID) string {
	return g.NodeNames[id]
}

func (g *Graph) GetNumNodes() int {
	return len(g.NodeNames)
}

func (g *Graph) GetNodeDegree(id NodeID) int {
	Assert(id >= 0 && id < NodeID(len(g.Degrees)), "Invalid node")
	return g.Degrees[id]
}

func (g *Graph) AddEdge(from, to NodeID) {
	Assert(g.Edges != nil, "Graph not initialized")
	Assert(from >= 0 && from < NodeID(len(g.Edges)), "Invalid from node")
	Assert(to >= 0 && to < NodeID(len(g.Edges)), "Invalid to node")
	g.Edges[from][to] = true
	g.Edges[to][from] = true
	g.Degrees[from]++
	g.Degrees[to]++

	DebugLogf("Edge: %s-%s\n", g.GetNodeName(from), g.GetNodeName(to))
}

func (g *Graph) LoadGraph(lines []string) {
	re := regexp.MustCompile(`(.{2})-(.{2})`)

	// Collect all the nodes
	nodes := make(map[string]bool)
	for _, line := range lines {
		ss := re.FindStringSubmatch(line)
		nodes[ss[1]] = true
		nodes[ss[2]] = true
	}

	// Sort the nodes
	nodeList := slices.Collect(maps.Keys(nodes))
	slices.Sort(nodeList)

	// Register all nodes
	for _, node := range nodeList {
		g.RegisterNode(node)
	}

	g.InitEdges()

	// Populate the edge matrix
	for _, line := range lines {
		ss := re.FindStringSubmatch(line)
		from := g.GetNodeID(ss[1])
		to := g.GetNodeID(ss[2])
		g.AddEdge(from, to)
	}
}

func (g *Graph) FindTriangleCliques() Cliques {
	working := make([]NodeID, 3)

	return g.FindCliquesInner(1, 1, 3, working)
}

// Is the clique in the working array a triangle of size size?
func (g *Graph) IsClique(size int, working []NodeID) bool {
	for i := 0; i < size; i++ {
		for j := i + 1; j < size; j++ {
			if !g.Edges[working[i]][working[j]] {
				return false
			}
		}
	}
	return true
}

func (g *Graph) FindCliquesInner(
	startID NodeID, currSize int, size int, working []NodeID) Cliques {

	numNodes := g.GetNumNodes()
	cliques := Cliques{}

	for i := startID; i < NodeID(numNodes); i++ {
		// Only look at nodes that have enough edges to possibly be part of a clique
		if g.GetNodeDegree(i) >= size-1 {
			working[currSize-1] = i

			// Check to see if we have a clique of the correct size
			if g.IsClique(currSize, working) {
				if currSize == size {
					cliques.Add(working...)
				} else {
					// Recurse to find the next node in the clique
					nextCliques := g.FindCliquesInner(i+1, currSize+1, size, working)
					cliques.Merge(nextCliques)
				}
			}
		}
	}

	return cliques
}

// --------------------------------------------------------------------
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	Debug = false

	lines := ReadFileLines("input.txt")

	g := NewGraph()
	g.LoadGraph(lines)

	cliques := g.FindTriangleCliques()
	fmt.Println("Number of triangle cliques:", len(cliques.D))

	count := 0
	for _, c := range cliques.D {
		if g.GetNodeName(c[0])[0] == 't' ||
			g.GetNodeName(c[1])[0] == 't' ||
			g.GetNodeName(c[2])[0] == 't' {
			count++
		}
	}

	fmt.Println("Number of cliques with 't': ", count)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
