package tikshared

import (
	"os"
	"fmt"
	"encoding/gob"
	"errors"
)

func Decode(inputf *os.File, outputf *os.File) {
	dec := gob.NewDecoder(inputf)
	var codes = new(Codes)
	err := dec.Decode(codes)
	ErrorCheck(err)
	fmt.Println(codes)
	panic(errors.New("Not implemented yet!"))
}
