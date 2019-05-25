package Entangler

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"strconv"
)

const maxChunkSize float64 = 3900

func ChunkFile(file *os.File) ([]byte, error) {
	defer file.Close()
	fileinfo, _ := file.Stat()
	fileSize := float64(fileinfo.Size())

	totalChunks := uint64(math.Ceil(fileSize / maxChunkSize))

	for i := uint64(0); i < totalChunks; i = i + 1 {
		currChunkSize := math.Min(maxChunkSize, fileSize-(float64(i)*maxChunkSize))
		currChunk := make([]byte, int(currChunkSize))

		file.Read(currChunk)
		filename := strconv.Itoa(int(i)) + "_" + filepath.Base(file.Name())
		if _, err := os.Create(filename); err == nil {

		} else {
			fmt.Println("Fatal error ... " + err.Error())
			os.Exit(1)
		}

		ioutil.WriteFile(filename, currChunk, os.ModeAppend)
	}

	return nil, nil
}
