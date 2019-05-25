package Entangler

import (
	"fmt"
	"io/ioutil"
	"math"
	"os"
	//"path/filepath"
	"strconv"
)

const MaxChunkSize float64 = 3900
const ChunkDirectory string = "../../files/"

// Chunks the file. Sends back how many chunks it made
func ChunkFile(file *os.File) (int, error) {
	defer file.Close()
	fileinfo, _ := file.Stat()
	fileSize := float64(fileinfo.Size())

	totalChunks := int(math.Ceil(fileSize / MaxChunkSize))

	for i := 0; i < totalChunks; i = i + 1 {
		currChunkSize := math.Min(MaxChunkSize, fileSize-(float64(i)*MaxChunkSize))
		currChunk := make([]byte, int(currChunkSize))

		file.Read(currChunk)
		filename := "d" + strconv.Itoa(int(i))
		if _, err := os.Create(ChunkDirectory + filename); err == nil {

		} else {
			fmt.Println("Fatal error ... " + err.Error())
			os.Exit(1)
		}

		ioutil.WriteFile(ChunkDirectory+filename, currChunk, os.ModeAppend)
	}

	return totalChunks, nil
}
