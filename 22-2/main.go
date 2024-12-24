package main

import (
	"fmt"
	"log"
	"time"
)

const NumSeqs = 19 * 19 * 19 * 19
const NumSamples = 2000

type Seq uint32

func (in Seq) Unpack() [4]PriceDiff {
	return [4]PriceDiff{
		PriceDiff((in&0xFF00_0000)>>24) - 9,
		PriceDiff((in&0x00FF_0000)>>16) - 9,
		PriceDiff((in&0x0000_FF00)>>8) - 9,
		PriceDiff((in&0x0000_00FF)>>0) - 9,
	}
}

func (in Seq) AddAndShift(diff PriceDiff) Seq {
	in = in << 8
	in = in & 0xFFFFFF00
	diff = diff + 9
	in = in | Seq(diff)
	return in
}

func (in Seq) String() string {
	unpacked := in.Unpack()
	return fmt.Sprintf("%#08x: % d, % d, % d, % d", uint32(in),
		unpacked[0], unpacked[1], unpacked[2], unpacked[3])
}

func Generate(in uint32) uint32 {
	var ret uint64 = uint64(in)
	temp := ret << 6        // ret * 64
	ret = ret ^ temp        // ret XOR temp (mix)
	ret = ret & (1<<24 - 1) // ret MOD 2^24 (prune)

	temp = ret >> 5         // ret / 32
	ret = ret ^ temp        // ret XOR temp (mix)
	ret = ret & (1<<24 - 1) // ret MOD 2^24 (prune)

	temp = ret << 11        // ret * 2048
	ret = ret ^ temp        // ret XOR temp (mix)
	ret = ret & (1<<24 - 1) // ret MOD 2^24 (prune)

	return uint32(ret)
}

func GenerateN(i uint32, n int) uint32 {
	for j := 0; j < n; j++ {
		i = Generate(i)
	}
	return i
}

type Price int8
type PriceDiff int8
type Prices []Price
type SeqResults struct {
	SeqPrices map[Seq]Prices
	NumInputs int
}

func (in *SeqResults) Init(numInputs int) {
	in.SeqPrices = make(map[Seq]Prices)
	in.NumInputs = numInputs
}

func (in Price) Diff(other Price) PriceDiff {
	return PriceDiff(in - other)
}

func (in *SeqResults) PopulateSeqResults(inputIndex int, secret uint32) {
	var diffs Seq = 0xFFFFFFFF               // Packed diffs between value and prev. +9 to avoid negative values.
	var prevPrice Price = Price(secret % 10) // Previous value

	for j := 0; j < NumSamples; j++ {
		// Compute the next in the sequence
		secret = Generate(secret)

		price := Price(secret % 10)

		diff := price.Diff(prevPrice)
		diffs = diffs.AddAndShift(diff)

		DebugLogf("%8d: %d (% d)", secret, price, diff)

		if j >= 3 {
			in.Add(inputIndex, diffs, price)
			DebugLogf(" Added (%v)", diffs)
		}

		DebugLogf("\n")

		prevPrice = price
	}
}

// Add a restult to the SeqResults.
//
//	inputIndex is the index of the input that generated the result.
//	seq is the sequence that generated the result.
//	price is the price of the result.
//
// Returns true if the result was added, false if it was already present.
func (sr *SeqResults) Add(inputIndex int, seq Seq, price Price) bool {

	// Get the prices for the sequence
	prices, ok := sr.SeqPrices[seq]
	if !ok {
		prices = make(Prices, sr.NumInputs)
		for i := range prices {
			prices[i] = -1
		}
		sr.SeqPrices[seq] = prices
	}

	// If we don't have a price already then set it
	currPrice := prices[inputIndex]
	if currPrice == -1 {
		prices[inputIndex] = price
		return true
	}

	return false
}

func (in *SeqResults) FindBestSeq() (Seq, int) {
	bestSeq := Seq(0xFFFF_FFFF)
	bestTotal := 0

	for seq, prices := range in.SeqPrices {
		total := 0
		for _, price := range prices {
			if price != -1 {
				total += int(price)
			}
		}

		if total > bestTotal {
			bestSeq = seq
			bestTotal = total
		}
	}

	return bestSeq, bestTotal
}

//-------------------------------------------------------------------------

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	timeStart := time.Now()

	Debug = false

	lines := ReadFileLines("input.txt")
	var inputs []uint32
	for _, line := range lines {
		inputs = append(inputs, uint32(MustAtoi(line)))
	}

	// Initialize SeqResults
	var results SeqResults
	results.Init(len(inputs))

	for inputIndex, secret := range inputs {
		results.PopulateSeqResults(inputIndex, secret)
	}

	fmt.Println("Unique sequences:", len(results.SeqPrices))

	bestSeq, bestTotal := results.FindBestSeq()
	fmt.Println("Best sequence:", bestSeq, "Total:", bestTotal)

	fmt.Println("Elapsed time:", time.Since(timeStart))
}
