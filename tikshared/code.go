package tikshared

import "bytes"

type Code struct {
	Code    uint32
	CodeLen uint8
}

func (code Code) AsString() string {
	var buffer bytes.Buffer
	for i := code.CodeLen; i > 0; i-- {
		n := code.Code & (1 << (i - 1))
		var t string
		if n > 0 {
			t = "1"
		} else {
			t = "0"
		}
		buffer.WriteString(t)
	}

	return buffer.String()
}
