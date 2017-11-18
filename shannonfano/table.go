package shannonfano

import (
	"github.com/braska/theory-of-information-and-coding-problems/tikshared"
	"math"
	"sort"
)

type symbolsOrderedMap struct {
	mapping map[string]int
	order   []string
	sum     int
	prefix  *tikshared.Code
}

type symbolCodePair struct {
	symbol string
	code   *tikshared.Code
}

func (som *symbolsOrderedMap) sortedInsert(f string, c int) {
	l := len(som.order)
	if l == 0 {
		som.order = []string{f}
	} else {
		i := sort.Search(l, func(i int) bool {
			return c > som.mapping[som.order[i]]
		})
		if i == l { // not found = new value is the smallest
			som.order = append(som.order, f)
		} else {
			tail := som.order[i:]
			tail = append([]string{f}, tail...)

			som.order = append(som.order[0:i], tail...)
		}
	}

	som.mapping[f] = c
	som.sum += c
}

func BuildTable(mode tikshared.EncodingMode, symbols map[string]int) (codes tikshared.Codes) {
	codes = make(tikshared.Codes)

	symbolsOrdered := &symbolsOrderedMap{mapping: make(map[string]int), prefix: new(tikshared.Code)}

	for key, value := range symbols {
		symbolsOrdered.sortedInsert(key, value)
	}

	c := make(chan *symbolCodePair)
	buildSubtree(symbolsOrdered, c)
	for i := 0; i < len(symbols); i++ {
		symbolCode := <-c
		codes[symbolCode.symbol] = symbolCode.code
	}

	return
}

func buildSubtree(symbols *symbolsOrderedMap, c chan *symbolCodePair) {
	if len(symbols.order) == 1 {
		c <- &symbolCodePair{symbol: symbols.order[0], code: symbols.prefix}
		return
	}

	part0 := symbolsOrderedMap{
		mapping: make(map[string]int),
		prefix: &tikshared.Code{
			CodeLen: symbols.prefix.CodeLen + 1,
			Code:    symbols.prefix.Code << 1,
		},
	}
	part0Ready := false

	part1 := symbolsOrderedMap{
		mapping: make(map[string]int),
		prefix: &tikshared.Code{
			CodeLen: symbols.prefix.CodeLen + 1,
			Code:    (symbols.prefix.Code << 1) + 1,
		},
	}

	middle := float64(symbols.sum) / 2

	for _, symbol := range symbols.order {
		symbolWeight := symbols.mapping[symbol]

		if !part0Ready {
			distanceFromMiddle := math.Abs(float64(part0.sum) - middle)
			distanceFromMiddleAfterAdditionSymbolWeight := math.Abs(float64(part0.sum+symbolWeight) - middle)

			if distanceFromMiddle > distanceFromMiddleAfterAdditionSymbolWeight {
				part0.sortedInsert(symbol, symbolWeight)
				continue
			} else {
				part0Ready = true
				go buildSubtree(&part0, c)
			}
		}
		part1.sortedInsert(symbol, symbolWeight)
	}
	go buildSubtree(&part1, c)
}
