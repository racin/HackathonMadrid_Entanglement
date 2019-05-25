package Entangler

import (
	"os"
	"strconv"
	"testing"
)

func TestRebuildFile(t *testing.T) {
	filePath := "../../resources/images/ArraySatelite.jpg"

	file, err := os.Open(filePath)
	if err != nil {
		os.Exit(1)
	}
	numChunks, err := ChunkFile(file)
	allChunks := make([][]byte, numChunks)

	for i := 0; i < numChunks; i++ {
		dataChunk, err := ReadChunk(ChunkDirectory + "d" + strconv.Itoa(i))
		if err != nil {
			t.Fatal("Fail")
		}
		allChunks[i] = dataChunk
	}

	RebuildFile("RebuildedFile.jpg", allChunks...)
}

func TestRebuildFileFromParities(t *testing.T) {
	numChunks := 16
	allChunks := make([][]byte, numChunks)
	var dataChunk []byte
	var err error
	for i := 1; i <= numChunks; i++ {
		if i == 8 {
			p1, err := ReadChunk(TestingDirectory + "p4_8")
			if err != nil {
				t.Fatal("Fail")
			}
			p2, err := ReadChunk(TestingDirectory + "p8_12")
			if err != nil {
				t.Fatal("Fail")
			}
			dataChunk, _ = XORByteSlice(p1, p2)
		} else {
			dataChunk, err = ReadChunk(TestingDirectory + "d" + strconv.Itoa(i))
			if err != nil {
				t.Fatal("Fail")
			}
		}

		allChunks[i-1] = dataChunk
	}

	RebuildFile("RebuildFromChunk.jpg", allChunks...)
}
