package Entangler

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

const (
	RightStrands      = 5
	LeftStrands       = 5
	HorizontalStrands = 5
	Alpha             = 3
)

var S = int(HorizontalStrands)
var P = int(RightStrands)

var ParityMemory [RightStrands + LeftStrands + HorizontalStrands][]byte

func GetTotalStrands() int {
	return RightStrands + LeftStrands + HorizontalStrands
}

func init() {
	fmt.Println("USed?")
	for i := 0; i < GetTotalStrands(); i++ {
		ParityMemory[i] = make([]byte, int(MaxChunkSize))
		//ParityMemory[i] = []byte{0}
	}
}

// TODO: Fix underflow naming errors on the nodes on the extreme of the lattice.
func GetBackwardNeighbours(index int) (r, h, l int) {
	// Check is it top, center or bottom in the lattice
	// 1 -> Top, 0 -> Bottom, else Center
	var nodePos = index % S

	if nodePos == 1 {
		r = index - (S * P) + int((math.Pow(float64(S), 2) - 1))
		h = index - S
		l = index - (S - 1)
	} else if nodePos == 0 {
		r = index - (S + 1)
		h = index - S
		l = index - (S * P) + int(math.Pow(float64(S-1), 2))
	} else {
		r = index - (S + 1)
		h = index - S
		l = index - (S - 1)
	}
	return
}

func GetMemoryPosition(index int) (r, h, l int) {
	// Get the position in the ParityMemory array where the parity is located
	// For now this will recursively call the GetBackwardNeighbours function

	h = ((index - 1) % S) + S
	r, l = index, index

	for ; r > S; r, _, _ = GetBackwardNeighbours(r) {
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

	for ; l > S; _, _, l = GetBackwardNeighbours(l) {
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
	var nodePos = index % S

	if nodePos == 1 {
		r = index + S + 1
		h = index + S
		l = index + (S * P) - int(math.Pow(float64(S-1), 2))
	} else if nodePos == 0 {
		r = index + (S * P) - int(math.Pow(float64(S), 2)-1)
		h = index + S
		l = index + S - 1
	} else {
		r = index + S + 1
		h = index + S
		l = index + (S - 1)
	}
	return
}

func entangle(datachunk []byte, index int) {
	r, h, l := GetMemoryPosition(index)
	rBack, hBack, lBack := GetBackwardNeighbours(index)
	rParity := ParityMemory[r]
	hParity := ParityMemory[h]
	lParity := ParityMemory[l]

	WriteChunkToFile(rParity, rBack, index)
	WriteChunkToFile(hParity, hBack, index)
	WriteChunkToFile(lParity, lBack, index)

	rNext, _ := XORByteSlice(datachunk, rParity)
	ParityMemory[r] = rNext

	hNext, _ := XORByteSlice(datachunk, hParity)
	ParityMemory[h] = hNext

	lNext, _ := XORByteSlice(datachunk, lParity)
	ParityMemory[l] = lNext
}

func EntangleFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		os.Exit(1)
	}
	numChunks, err := ChunkFile(file)

	for i := 1; i <= numChunks; i++ {
		dataChunk, err := ReadChunk(ChunkDirectory + "d" + strconv.Itoa(i))
		if err != nil {
			return err
		}
		entangle(dataChunk, i)
	}
	// File -> Data chunks
	// Datachunks ->
	return nil
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
