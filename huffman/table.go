package huffman

import "github.com/braska/theory-of-information-and-coding-problems/tikshared"

func BuildTable(mode tikshared.EncodingMode, huffmanTreeVertices map[string]int) (codes tikshared.Codes) {
	codes = make(tikshared.Codes)

	for len(huffmanTreeVertices) > 1 {
		key1 := tikshared.FindMinKeyInMap(huffmanTreeVertices)
		value1 := huffmanTreeVertices[key1]
		delete(huffmanTreeVertices, key1)
		key2 := tikshared.FindMinKeyInMap(huffmanTreeVertices)
		value2 := huffmanTreeVertices[key2]
		delete(huffmanTreeVertices, key2)

		// if this is just leaf in huffman tree (single rune)
		if (mode == tikshared.MODE_ONE && len([]rune(key1)) == 1) || (mode == tikshared.MODE_TWO && len([]rune(key1)) == 2) {
			// key1 is a string from left branch => code=0
			codes[key1] = &tikshared.Code{Code: 0, CodeLen: 1}
		} else {
			if mode == tikshared.MODE_ONE {
				// for all runes from key1 add ZERO as a prefix for code
				// adding zero as a prefix is the same as increment codeLen
				for _, r := range key1 {
					codes[string(r)].CodeLen += 1
				}
			} else {
				var prevSym string
				for _, r := range key1 {
					if prevSym != "" {
						combination := prevSym + string(r)
						codes[combination].CodeLen += 1
						prevSym = ""
					} else {
						prevSym = string(r)
					}
				}
			}
		}

		if (mode == tikshared.MODE_ONE && len([]rune(key2)) == 1) || (mode == tikshared.MODE_TWO && len([]rune(key2)) == 2) {
			codes[key2] = &tikshared.Code{Code: 1, CodeLen: 1}
		} else {
			if mode == tikshared.MODE_ONE {
				for _, r := range key2 {
					k := string(r)
					// set bit to 1 at specific position
					codes[k].Code |= (1 << codes[k].CodeLen)
					codes[k].CodeLen += 1
				}
			} else {
				var prevSym string
				for _, r := range key2 {
					if prevSym != "" {
						combination := prevSym + string(r)
						// set bit to 1 at specific position
						codes[combination].Code |= (1 << codes[combination].CodeLen)
						codes[combination].CodeLen += 1
						prevSym = ""
					} else {
						prevSym = string(r)
					}
				}
			}
		}

		huffmanTreeVertices[key1+key2] = value1 + value2
	}

	return
}
