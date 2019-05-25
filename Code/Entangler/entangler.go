package Entangler

import (
	"fmt"
)

const (
	rightStrands      = 5
	leftStrands       = 5
	horizontalStrands = 5
	strands           = 5
)

var ParityMemory = make([]int, 15)

func GivemeInput(index int) []byte {
	// Check is it top, center or bottom?
	// 1 -> Top, 0 -> Bottom, else Center
	//var strandPos int = index % strands

	return nil

}
func Entangle(chunk []byte, index int) {

}
func XORByteSlice(a, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}

	buf := make([]byte, len(a))

	for i, _ := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}
