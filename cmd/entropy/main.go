package main

import (
	"flag"
	"fmt"
	"github.com/braska/theory-of-information-and-coding-problems/entropy"
	"github.com/braska/theory-of-information-and-coding-problems/tikshared"
	"github.com/olekukonko/tablewriter"
	"os"
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: %s [inputfile]\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Input file is missing.")
		os.Exit(1)
	}

	f, err := os.Open(os.Args[1])
	tikshared.ErrorCheck(err)
	table := tablewriter.NewWriter(os.Stdout)

	countSymbols, symbols, combinations, countCombinations := tikshared.ParseFile(f)

	table.Append([]string{"Runes", fmt.Sprint(countSymbols)})
	table.Append([]string{"Alphabet power", fmt.Sprint(len(symbols))})

	resultEntropy := entropy.Entropy(countSymbols, symbols)
	table.Append([]string{"Entropy", fmt.Sprint(resultEntropy)})

	conditionalEntropy := entropy.ConditionalEntropy(countSymbols, symbols, combinations, countCombinations)
	table.Append([]string{"Conditional entropy", fmt.Sprint(conditionalEntropy)})

	table.Render()

	f.Close()
}
