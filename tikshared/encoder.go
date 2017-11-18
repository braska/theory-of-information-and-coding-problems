package tikshared

import (
	"bufio"
	"encoding/gob"
	"os"
)

func Encode(mode EncodingMode, inputf *os.File, outputf *os.File, codes Codes) (bytesBefore int64, bytesAfter int64) {
	enc := gob.NewEncoder(outputf)
	err := enc.Encode(codes)
	ErrorCheck(err)

	inputf.Seek(0, 0)
	scanner := bufio.NewScanner(inputf)
	scanner.Split(bufio.ScanRunes)
	var tmpByte byte
	var bitsInTmpByte uint8 = 0
	var prevSym string
	for scanner.Scan() {
		s := scanner.Text()

		if mode == MODE_TWO {
			if prevSym == "" {
				prevSym = s
				continue
			}
			s = prevSym + s
			prevSym = ""
		}

		hcode := *codes[s]
		for hcode.CodeLen > 0 {
			shift := int(hcode.CodeLen) - (8 - int(bitsInTmpByte))
			var bitsMustBeWritenToTmpByte uint8 = 0
			if shift >= 0 {
				tmpByte += byte((hcode.Code >> uint(shift)) & 0xff)
				bitsMustBeWritenToTmpByte = 8 - bitsInTmpByte
				bitsInTmpByte = 8
			} else if shift < 0 {
				tmpByte += byte((hcode.Code << uint(shift*-1)) & 0xff)
				bitsMustBeWritenToTmpByte = hcode.CodeLen
				bitsInTmpByte += hcode.CodeLen
			}

			if shift <= 0 {
				hcode.CodeLen = 0
			} else {
				hcode.CodeLen -= bitsMustBeWritenToTmpByte
			}

			if bitsInTmpByte == 8 {
				_, err = outputf.Write([]byte{tmpByte})
				ErrorCheck(err)
				bitsInTmpByte = 0
				tmpByte = 0
			}
		}
	}

	if bitsInTmpByte > 0 {
		_, err = outputf.Write([]byte{tmpByte})
		ErrorCheck(err)
		bitsInTmpByte = 0
		tmpByte = 0
	}

	outputf.Sync()

	inputfStat, istatErr := inputf.Stat()
	ErrorCheck(istatErr)
	outputfStat, ostatErr := outputf.Stat()
	ErrorCheck(ostatErr)

	bytesBefore = inputfStat.Size()
	bytesAfter = outputfStat.Size()

	return
}
