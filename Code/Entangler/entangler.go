package Entangler

import (
	"fmt"
	"math"
	"os"
)

const (
	rightStrands      = 5
	leftStrands       = 5
	horizontalStrands = 5
)

var s = int(horizontalStrands)
var p = int(rightStrands)

var ParityMemory [rightStrands + leftStrands + horizontalStrands][]byte

func GetTotalStrands() int {
	return rightStrands + leftStrands + horizontalStrands
}

func init() {
	for i := 0; i < GetTotalStrands(); i++ {
		ParityMemory[i] = make([]byte, int(MaxChunkSize))
	}
}

func GetBackwardNeighbours(index int) (r, h, l int) {
	// Check is it top, center or bottom in the lattice
	// 1 -> Top, 0 -> Bottom, else Center
	var nodePos = index % s

	if nodePos == 1 {
		r = index - (s * p) + int((math.Pow(float64(s), 2) - 1))
		h = index - s
		l = index - (s - 1)
	} else if nodePos == 0 {
		r = index - (s + 1)
		h = index - s
		l = index - (s * p) + int(math.Pow(float64(s-1), 2))
	} else {
		r = index - (s + 1)
		h = index - s
		l = index - (s - 1)
	}
	return
}

func GetMemoryPosition(index int) (r, h, l int) {
	// Get the position in the ParityMemory array where the parity is located
	// For now this will recursively call the GetBackwardNeighbours function

	h = ((index - 1) % s) + s
	r, l = index, index

	for ; r > s; r, _, _ = GetBackwardNeighbours(r) {
	}

	switch r {
	case 1:
		r = 0
		break
	case 2:
		r = 4
		break
	case 3:
		r = 3
		break
	case 4:
		r = 2
		break
	case 5:
		r = 1
		break
	}

	for ; l > s; _, _, l = GetBackwardNeighbours(l) {
	}

	switch l {
	case 1:
		l = 11
		break
	case 2:
		l = 12
		break
	case 3:
		l = 13
		break
	case 4:
		l = 14
		break
	case 5:
		l = 10
		break
	}

	return
}

func GetForwardNeighbours(index int) (r, h, l int) {
	// Check is it top, center or bottom in the lattice
	// 1 -> Top, 0 -> Bottom, else Center
	var nodePos = index % s

	if nodePos == 1 {
		r = index + s + 1
		h = index + s
		l = index + (s * p) - int(math.Pow(float64(s-1), 2))
	} else if nodePos == 0 {
		r = index + (s * p) - int(math.Pow(float64(s), 2)-1)
		h = index + s
		l = index + s - 1
	} else {
		r = index + s + 1
		h = index + s
		l = index + (s - 1)
	}
	return
}

func entangle(datachunk []byte, index int) {
	r, h, l := GetMemoryPosition(index)
	rParity := ParityMemory[r]
	hParity := ParityMemory[h]
	lParity := ParityMemory[l]

	rNext, _ := XORByteSlice(datachunk, rParity)
	ParityMemory[r] = rNext

	hNext, _ := XORByteSlice(datachunk, hParity)
	ParityMemory[h] = hNext

	lNext, _ := XORByteSlice(datachunk, lParity)
	ParityMemory[l] = lNext

	//WriteFile(rNext, index)
}

func EntangleFile(filename string) {
	// Input file
	filePath := "../../resources/images/ArraySatelite.jpg"
	file, err := os.Open(filePath)
	if err != nil {
		os.Exit(1)
	}
	ChunkFile(file)

	// File -> Data chunks
	// Datachunks ->
}
func XORByteSlice(a []byte, b []byte) ([]byte, error) {
	if len(a) != len(b) {
		return nil, fmt.Errorf("length of byte slices is not equivalent: %d != %d", len(a), len(b))
	}

	buf := make([]byte, len(a))

	for i, _ := range a {
		buf[i] = a[i] ^ b[i]
	}

	return buf, nil
}
