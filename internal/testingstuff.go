package internal

import (
	"bytes"
	"strconv"
)

// SomeRandomStuff this is to check
func SomeRandomStuff() int {
	x := getFib()
	return x
}

func getFib() int {
	var x bytes.Buffer
	_, e := x.WriteString("9")
	if e != nil {
		return -1
	}
	stringInNumbe := x.String()
	v, er := strconv.ParseInt(stringInNumbe, 10, 64)
	if er != nil {
		return -1

	}
	return int(v)
}
