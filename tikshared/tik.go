package tikshared

import (
	"bufio"
	"os"
)

func FindMinKeyInMap(m map[string]int) (k string) {
	var minValue int
	isFirstIter := true
	for key, value := range m {
		if isFirstIter {
			minValue = value
			k = key
			isFirstIter = false
		}

		if value < minValue {
			minValue = value
			k = key
		}
	}
	return
}

func ErrorCheck(e error) {
	if e != nil {
		panic(e)
	}
}

func ParseFile(inputf *os.File) (countSymbols int, symbols map[string]int, combinations map[string]int, countCombinations int) {
	countSymbols = 0
	symbols = make(map[string]int)
	combinations = make(map[string]int)
	var prevSymbol string

	scanner := bufio.NewScanner(inputf)
	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		s := scanner.Text()
		symbols[s] += 1
		countSymbols += 1
		if prevSymbol != "" {
			combination := prevSymbol + s
			combinations[combination] += 1
		}
		prevSymbol = s
	}
	countCombinations = countSymbols - 1
	return
}
