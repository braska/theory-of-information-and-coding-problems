package main

import (
	"flag"
	"fmt"
	"github.com/braska/theory-of-information-and-coding-problems/entropy"
	"github.com/braska/theory-of-information-and-coding-problems/huffman"
	"github.com/braska/theory-of-information-and-coding-problems/tikshared"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [command] [inputfile] [outputfile]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func encode(mode tikshared.EncodingMode, inputf *os.File, outputf *os.File) {
	table := tablewriter.NewWriter(os.Stdout)

	countSymbols, symbols, combinations, countCombinations := tikshared.ParseFile(inputf)

	table.Append([]string{"Runes", fmt.Sprint(countSymbols)})
	table.Append([]string{"Alphabet power", fmt.Sprint(len(symbols))})

	resultEntropy := entropy.Entropy(countSymbols, symbols)
	table.Append([]string{"Entropy", fmt.Sprint(resultEntropy)})

	conditionalEntropy := entropy.ConditionalEntropy(countSymbols, symbols, combinations, countCombinations)
	table.Append([]string{"Conditional entropy", fmt.Sprint(conditionalEntropy)})

	tableOfCodes := tablewriter.NewWriter(os.Stdout)

	var huffmanTreeVertices map[string]int

	if mode == tikshared.MODE_ONE {
		huffmanTreeVertices = symbols
	} else {
		huffmanTreeVertices = combinations
	}

	codes := huffman.BuildTable(mode, huffmanTreeVertices)

	for r, code := range codes {
		tableOfCodes.Append([]string{string(r), fmt.Sprint(code.AsString())})
	}

	bytesBefore, bytesAfter := tikshared.Encode(mode, inputf, outputf, codes)

	tableOfStats := tablewriter.NewWriter(os.Stdout)

	tableOfStats.Append([]string{"Input file size", humanize.Bytes(uint64(bytesBefore))})
	tableOfStats.Append([]string{"Output file size", humanize.Bytes(uint64(bytesAfter))})
	tableOfStats.Append([]string{"Bits per rune in input file", fmt.Sprint(float64(bytesBefore*8) / float64(countSymbols))})
	tableOfStats.Append([]string{"Bits per rune in output file", fmt.Sprint(float64(bytesAfter*8) / float64(countSymbols))})

	table.Render()
	tableOfCodes.Render()
	tableOfStats.Render()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Command is missing.")
		os.Exit(1)
	}

	if len(args) < 2 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}

	if len(args) < 3 {
		fmt.Println("Output file is missing.")
		os.Exit(1)
	}

	inputPath := os.Args[2]
	outputPath := os.Args[3]

	inputf, inputerr := os.Open(inputPath)
	tikshared.ErrorCheck(inputerr)

	outputf, outputerr := os.Create(outputPath)
	tikshared.ErrorCheck(outputerr)

	defer inputf.Close()
	defer outputf.Close()

	if os.Args[1] == "encode" {
		encode(tikshared.MODE_ONE, inputf, outputf)
	} else if os.Args[1] == "encode2" {
		encode(tikshared.MODE_TWO, inputf, outputf)
	} else if os.Args[1] == "decode" {
		tikshared.Decode(inputf, outputf)
	} else {
		fmt.Fprintf(os.Stderr, "%s: '%s' is not a %s command.\n", os.Args[0], os.Args[1], os.Args[0])
		os.Exit(1)
	}
}
