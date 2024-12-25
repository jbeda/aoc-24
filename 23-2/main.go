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

func NewClique(ids ...NodeID) Clique {
	ret := make(Clique, len(ids))
	copy(ret, ids)
	return ret
}

func (c Clique) String(g *Graph) string {
	s := ""
	for _, node := range c {
		s += g.GetNodeName(node) + ","
	}
	return s[:len(s)-1]
}

func (c Clique) Clone() Clique {
	clone := make(Clique, len(c))
	copy(clone, c)
	return clone
}

func (c Clique) Expand(n NodeID) Clique {
	return append(c, n)
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

func (g *Graph) FindMaxClique() Clique {
	numNodes := g.GetNumNodes()
	working := make([]NodeID, numNodes)

	return g.FindMaxCliqueInner(1, 1, working)
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

func (g *Graph) FindMaxCliqueInner(
	startID NodeID, currSize int, working []NodeID) Clique {

	numNodes := g.GetNumNodes()
	maxClique := Clique{}

	for i := startID; i < NodeID(numNodes); i++ {
		working[currSize-1] = i

		// Check to see if we have a clique of the correct size
		if g.IsClique(currSize, working) {
			if currSize > len(maxClique) {
				maxClique = NewClique(working[:currSize]...)
			}

			// Recurse to find the next node in the clique
			nextClique := g.FindMaxCliqueInner(i+1, currSize+1, working)
			if len(nextClique) > len(maxClique) {
				maxClique = nextClique
			}
		}
	}

	return maxClique
}

// --------------------------------------------------------------------
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	Debug = false

	lines := ReadFileLines("input.txt")

	g := NewGraph()
	g.LoadGraph(lines)

	clique := g.FindMaxClique()
	fmt.Println("Size: ", len(clique))
	fmt.Println("Clique: ", clique.String(g))

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
