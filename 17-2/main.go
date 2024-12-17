package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"
)

// ---------------------------------------------------------------------
type Machine struct {
	A, B, C int
	InstPtr int
	Prog    []int
	Output  []int
}

type Instruction int

const (
	ADV Instruction = iota
	BXL
	BST
	JNZ
	BXC
	OUT
	BDV
	CDV
)

func (i Instruction) String() string {
	return [...]string{"ADV", "BXL", "BST", "JNZ", "BXC", "OUT", "BDV", "CDV"}[i]
}

func (m *Machine) Clone() *Machine {
	m2 := *m
	m2.Prog = make([]int, len(m.Prog))
	copy(m2.Prog, m.Prog)
	m2.Output = make([]int, len(m.Output))
	copy(m2.Output, m.Output)
	return &m2
}

func (m *Machine) Run() {
	for {
		if m.InstPtr >= len(m.Prog) {
			break
		}

		// Read the instruction
		inst := Instruction(m.Prog[m.InstPtr])
		m.InstPtr++

		// Read the operand
		operandLiteral := m.Prog[m.InstPtr]
		m.InstPtr++

		// Compute the combo operand
		operandCombo := operandLiteral
		switch operandCombo {
		case 0, 1, 2, 3:
			// Do nothing
		case 4:
			operandCombo = m.A
		case 5:
			operandCombo = m.B
		case 6:
			operandCombo = m.C
		}

		switch inst {
		// The adv instruction (opcode 0) performs division. The numerator is the
		// value in the A register. The denominator is found by raising 2 to the
		// power of the instruction's combo operand. (So, an operand of 2 would
		// divide A by 4 (2^2); an operand of 5 would divide A by 2^B.) The result
		// of the division operation is truncated to an integer and then written to
		// the A register.
		case ADV:
			m.A = m.A >> operandCombo

		// The bxl instruction (opcode 1) calculates the bitwise XOR of register B
		// and the instruction's literal operand, then stores the result in register
		// B.
		case BXL:
			m.B ^= operandLiteral

		// The bst instruction (opcode 2) calculates the value of its combo operand
		// modulo 8 (thereby keeping only its lowest 3 bits), then writes that value
		// to the B register.
		case BST:
			m.B = operandCombo & 0x07

		// The jnz instruction (opcode 3) does nothing if the A register is 0.
		// However, if the A register is not zero, it jumps by setting the
		// instruction pointer to the value of its literal operand; if this
		// instruction jumps, the instruction pointer is not increased by 2 after
		// this instruction.
		case JNZ:
			if m.A != 0 {
				m.InstPtr = operandLiteral
			}

		// The bxc instruction (opcode 4) calculates the bitwise XOR of register B
		// and register C, then stores the result in register B. (For legacy
		// reasons, this instruction reads an operand but ignores it.)
		case BXC:
			m.B = m.B ^ m.C

		// The out instruction (opcode 5) calculates the value of its combo operand
		// modulo 8, then outputs that value. (If a program outputs multiple values,
		// they are separated by commas.)
		case OUT:
			m.Output = append(m.Output, operandCombo&0x07)

		// The bdv instruction (opcode 6) works exactly like the adv instruction
		// except that the result is stored in the B register. (The numerator is
		// still read from the A register.)
		case BDV:
			m.B = m.A >> operandCombo

		// The cdv instruction (opcode 7) works exactly like the adv instruction
		// except that the result is stored in the C register. (The numerator is
		// still read from the A register.)
		case CDV:
			m.C = m.A >> operandCombo
		}

		// fmt.Println("After:")
		// fmt.Println("Inst:", inst)
		// fmt.Println("Literal operand: ", operandLiteral)
		// fmt.Println("Combo operand: ", operandCombo)
		// fmt.Println("Machine state:", m)
		// fmt.Println()
	}
}

func LoadMachine(infile string) *Machine {
	f, err := os.Open(infile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)

	m := Machine{}

	// Register A: 729
	scan.Scan()
	line := scan.Text()
	re := regexp.MustCompile(`Register A: (\d+)`)
	matches := re.FindStringSubmatch(line)
	if len(matches) != 2 {
		log.Fatalf("Couldn't parse line: %q", line)
	}
	m.A = MustAtoi(matches[1])

	// Register B: 0
	scan.Scan()
	line = scan.Text()
	re = regexp.MustCompile(`Register B: (\d+)`)
	matches = re.FindStringSubmatch(line)
	if len(matches) != 2 {
		log.Fatalf("Couldn't parse line: %q", line)
	}
	m.B = MustAtoi(matches[1])

	// Register C: 0
	scan.Scan()
	line = scan.Text()
	re = regexp.MustCompile(`Register C: (\d+)`)
	matches = re.FindStringSubmatch(line)
	if len(matches) != 2 {
		log.Fatalf("Couldn't parse line: %q", line)
	}
	m.C = MustAtoi(matches[1])

	// Blank line
	scan.Scan()
	line = scan.Text()
	if line != "" {
		log.Fatalf("Expected blank line, got: %q", line)
	}

	// Program: 0,1,5,4,3,0
	scan.Scan()
	line = scan.Text()
	re = regexp.MustCompile(`Program: (.+)`)
	matches = re.FindStringSubmatch(line)
	if len(matches) != 2 {
		log.Fatalf("Couldn't parse line: %q", line)
	}
	for _, s := range strings.Split(matches[1], ",") {
		m.Prog = append(m.Prog, MustAtoi(s))
	}

	return &m
}

func (m *Machine) String() string {
	return fmt.Sprintf("A=%d B=%d C=%d Inst=%d Prog=%v Output=%v", m.A, m.B, m.C, m.InstPtr, m.Prog, m.Output)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	m := LoadMachine("input.txt")

	// Work backward from the last instruction and search 3 bits at a time
	candidates := make([]int, 0)
	candidates = append(candidates, 0)

	for i := len(m.Prog) - 1; i >= 0; i-- {
		fmt.Printf("Program element %d is %d\n", i, m.Prog[i])

		nextCandidates := make([]int, 0)
		for _, c := range candidates {
			c = c << 3
			for j := 0; j < 8; j++ {
				aInit := c | j
				m2 := m.Clone()
				m2.A = aInit
				m2.Run()

				fmt.Printf("Candidate: %d, Output: %v\n", aInit, m2.Output)

				if slices.Compare(m2.Output, m.Prog[i:]) == 0 {
					nextCandidates = append(nextCandidates, aInit)
				}
			}
		}

		candidates = nextCandidates

		fmt.Println("Candidates:", candidates)
	}

	fmt.Println(slices.Min(candidates))
}
