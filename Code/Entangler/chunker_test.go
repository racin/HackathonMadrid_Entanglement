package Entangler

import (
	"fmt"
	"os"
	"testing"
)

func TestChunkFile(t *testing.T) {
	filePath := "../../resources/images/ArraySatelite.jpg"

	t.Run("TestChunker", func(t *testing.T) {
		dir, err := os.Getwd()
		fmt.Println(dir)
		file, err := os.Open(filePath)
		if err != nil {
			t.Fatal("Error opening file " + err.Error())
		}
		ChunkFile(file)
	})
}
