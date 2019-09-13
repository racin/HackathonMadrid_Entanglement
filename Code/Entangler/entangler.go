package Entangler

import (
	"fmt"
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
	for i := 0; i < GetTotalStrands(); i++ {
		ParityMemory[i] = make([]byte, int(MaxChunkSize))
		//ParityMemory[i] = []byte{0}
	}
}

func entangle(datachunk []byte, index int) {
	r, h, l := GetMemoryPosition(index, S, P)
	rBack, hBack, lBack := GetBackwardNeighbours(index, S, P)
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
