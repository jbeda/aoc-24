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

func (nt NodeType) String() string {
	switch nt {
	case Constant:
		return "Constant"
	case AND:
		return "AND"
	case OR:
		return "OR"
	case XOR:
		return "XOR"
	case Unknown:
		return "Unknown"
	}
	return "Invalid"
}

type LogicNode struct {
	Name     string
	Alias    string
	Type     NodeType
	Val      bool
	Computed bool
	Inputs   [2]*LogicNode
	Outputs  []*LogicNode
}

func (ln *LogicNode) NameString() string {
	if ln.Alias != "" {
		return fmt.Sprintf("%s/%s", ln.Alias, ln.Name)
	}
	return ln.Name
}

func (ln *LogicNode) ValString() string {
	if ln.Computed {
		if ln.Val {
			return "1"
		}
		return "0"
	}
	return "?"
}

func (ln *LogicNode) Compute() bool {
	if ln.Computed {
		return ln.Val
	}

	for _, in := range ln.Inputs {
		if in != nil && !in.Computed {
			in.Compute()
		}
	}

	switch ln.Type {
	case Constant:
		Assert(false, "Constant node should have been computed already")
	case AND:
		ln.Val = ln.Inputs[0].Val && ln.Inputs[1].Val
	case OR:
		ln.Val = ln.Inputs[0].Val || ln.Inputs[1].Val
	case XOR:
		ln.Val = ln.Inputs[0].Val != ln.Inputs[1].Val
	case Unknown:
		Assert(false, "Unknown node type")
	}
	ln.Computed = true
	return ln.Val
}

// ----------------------------------------
type LogicGraph struct {
	Nodes   map[string]*LogicNode
	Aliases map[string]*LogicNode
	Outputs []*LogicNode
	Swaps   map[string]string
}

func NewLogicGraph() *LogicGraph {
	return &LogicGraph{
		Nodes:   make(map[string]*LogicNode),
		Aliases: make(map[string]*LogicNode),
		Swaps:   make(map[string]string),
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

func (lg *LogicGraph) AddSwap(n1, n2 string) {
	lg.Swaps[n1] = n2
	lg.Swaps[n2] = n1
}

func (lg *LogicGraph) GetByAliasOrName(name string) *LogicNode {
	node, ok := lg.Aliases[name]
	if !ok {
		node = lg.Nodes[name]
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

	if swap, ok := lg.Swaps[out]; ok {
		out = swap
	}

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

func (lg *LogicGraph) SetAlias(name, alias string) {
	node := lg.Nodes[name]
	node.Alias = alias
	lg.Aliases[alias] = node
}

func (lg *LogicGraph) PrintLogic(name string, depth int) {
	lg.PrintLogicInner(name, "", depth)
}

func (lg *LogicGraph) PrintLogicInner(name string, indent string, depth int) {
	if depth == 0 {
		return
	}

	node := lg.GetByAliasOrName(name)
	if node == nil {
		fmt.Printf("!! Could not find node for %s\n", name)
		return
	}

	fmt.Printf("%s%s: ", indent, node.NameString())
	switch node.Type {
	case Constant:
		fmt.Printf("Constant")
	case AND:
		fmt.Printf("%s AND %s", node.Inputs[0].NameString(), node.Inputs[1].NameString())
	case OR:
		fmt.Printf("%s OR %s", node.Inputs[0].NameString(), node.Inputs[1].NameString())
	case XOR:
		fmt.Printf("%s XOR %s", node.Inputs[0].NameString(), node.Inputs[1].NameString())
	case Unknown:
		fmt.Printf("Unknown")
	}

	if len(node.Outputs) > 0 {
		fmt.Printf(" ->")
		for _, out := range node.Outputs {
			fmt.Printf(" %s", out.NameString())
		}
	}
	fmt.Println()

	for _, in := range node.Inputs {
		if in != nil {
			lg.PrintLogicInner(in.Name, indent+"  ", depth-1)
		}
	}
}

// Find a node that fits the in/op pattern and give it an alias
func (lg *LogicGraph) CreateAlias(in1 string, op NodeType, in2 string, alias string) {
	n1 := lg.GetByAliasOrName(in1)
	if n1 == nil {
		fmt.Printf("!! CreateAlias: Could not find node for %s\n", in1)
		return
	}
	n2 := lg.GetByAliasOrName(in2)
	if n2 == nil {
		fmt.Printf("!! CreateAlias: Could not find node for %s\n", in2)
		return
	}

	for _, node := range lg.Nodes {
		if node.Type == op && ((node.Inputs[0] == n1 && node.Inputs[1] == n2) ||
			(node.Inputs[0] == n2 && node.Inputs[1] == n1)) {
			lg.SetAlias(node.Name, alias)
			return
		}
	}
	fmt.Printf("!! CreateAlias: Could not find node for %s %s %s to set alias %s\n", in1, op, in2, alias)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	lines := ReadFileLines("input.txt")
	lg := NewLogicGraph()

	lg.AddSwap("kmb", "z10")
	lg.AddSwap("tvp", "z15")
	lg.AddSwap("dpg", "z25")
	lg.AddSwap("mmf", "vdk")

	lg.Load(lines)
	fmt.Println("Output:", lg.GetOutput())

	fmt.Printf("Digit 0\n")
	lg.PrintLogic("z00", 1)
	lg.CreateAlias("x00", AND, "y00", "c00")
	lg.PrintLogic("c00", 1)
	fmt.Println()

	for digit := 1; digit < 45; digit++ {
		fmt.Printf("Digit %d\n", digit)

		prevNum := fmt.Sprintf("%02d", digit-1)
		currNum := fmt.Sprintf("%02d", digit)

		x := "x" + currNum
		y := "y" + currNum
		z := "z" + currNum
		o := "o" + currNum
		i := "i" + currNum
		a := "a" + currNum
		b := "b" + currNum
		cin := "c" + prevNum
		cout := "c" + currNum

		lg.CreateAlias(x, XOR, y, i)
		lg.CreateAlias(i, XOR, cin, o)
		lg.CreateAlias(x, AND, y, a)
		lg.CreateAlias(i, AND, cin, b)
		lg.CreateAlias(a, OR, b, cout)
		lg.PrintLogic(x, 1)
		lg.PrintLogic(y, 1)
		lg.PrintLogic(z, 1)
		lg.PrintLogic(i, 1)
		lg.PrintLogic(o, 1)
		lg.PrintLogic(a, 1)
		lg.PrintLogic(b, 1)
		lg.PrintLogic(cout, 1)

		fmt.Println()
	}

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
