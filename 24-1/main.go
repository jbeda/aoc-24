package main

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

type NodeType int

const (
	Constant NodeType = iota
	AND
	OR
	XOR
	Unknown
)

type LogicNode struct {
	Name     string
	Type     NodeType
	Val      bool
	Computed bool
	Inputs   [2]*LogicNode
	Outputs  []*LogicNode
}

func (ln *LogicNode) Compute() bool {
	if ln.Computed {
		return ln.Val
	}

	switch ln.Type {
	case Constant:
		Assert(false, "Constant node should have been computed already")
	case AND:
		ln.Val = ln.Inputs[0].Compute() && ln.Inputs[1].Compute()
	case OR:
		ln.Val = ln.Inputs[0].Compute() || ln.Inputs[1].Compute()
	case XOR:
		ln.Val = ln.Inputs[0].Compute() != ln.Inputs[1].Compute()
	case Unknown:
		Assert(false, "Unknown node type")
	}
	ln.Computed = true
	return ln.Val
}

// ----------------------------------------
type LogicGraph struct {
	Nodes   map[string]*LogicNode
	Outputs []*LogicNode
}

func NewLogicGraph() *LogicGraph {
	return &LogicGraph{
		Nodes: make(map[string]*LogicNode),
	}
}

func (lg *LogicGraph) GetNode(name string) *LogicNode {
	node, ok := lg.Nodes[name]
	if !ok {
		node = &LogicNode{
			Name: name,
			Type: Unknown,
		}
		lg.Nodes[name] = node

		if name[0] == 'z' {
			lg.Outputs = append(lg.Outputs, node)
		}
	}
	return node
}

func (lg *LogicGraph) AddConstant(name string, val bool) {
	node := lg.GetNode(name)
	node.Type = Constant
	node.Val = val
	node.Computed = true
}

func (lg *LogicGraph) AddRule(in1, in2 string, op NodeType, out string) {
	node := lg.GetNode(out)
	node.Type = op
	node.Inputs[0] = lg.GetNode(in1)
	node.Inputs[1] = lg.GetNode(in2)
	node.Inputs[0].Outputs = append(node.Inputs[0].Outputs, node)
	node.Inputs[1].Outputs = append(node.Inputs[1].Outputs, node)
}

func (lg *LogicGraph) Compute() {
	for _, node := range lg.Outputs {
		node.Compute()
	}
}

func (lg *LogicGraph) Reset() {
	for _, node := range lg.Nodes {
		if node.Type != Constant {
			node.Computed = false
		}
	}
}

func (lg *LogicGraph) Load(lines []string) {
	var i int

	re := regexp.MustCompile(`^(\w+): (0|1)$`)
	for i = 0; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			break
		}

		matches := re.FindStringSubmatch(line)
		Assert(matches != nil, "Invalid constant line: %s", line)
		name := matches[1]
		val := matches[2] == "1"
		lg.AddConstant(name, val)
	}

	re = regexp.MustCompile(`^(\w+) (AND|OR|XOR) (\w+) -> (\w+)$`)
	for i++; i < len(lines); i++ {
		line := lines[i]
		matches := re.FindStringSubmatch(line)
		Assert(matches != nil, "Invalid rule line: %s", line)
		in1 := matches[1]
		op := matches[2]
		in2 := matches[3]
		out := matches[4]

		var opType NodeType
		switch op {
		case "AND":
			opType = AND
		case "OR":
			opType = OR
		case "XOR":
			opType = XOR
		default:
			Assert(false, "Unknown operation: %s", op)
		}

		lg.AddRule(in1, in2, opType, out)
	}
}

func (lg *LogicGraph) GetOutput() int {
	lg.Compute()
	result := 0
	for _, node := range lg.Outputs {
		bitPos := MustAtoi(node.Name[1:])
		if node.Val {
			result |= 1 << bitPos
		}
	}
	return result
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	lines := ReadFileLines("input.txt")
	lg := NewLogicGraph()
	lg.Load(lines)
	fmt.Println("Output:", lg.GetOutput())

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
