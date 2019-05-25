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
const FileDirectory string = "../../files/"

func ChunkFile(file *os.File) ([]byte, error) {
	defer file.Close()
	fileinfo, _ := file.Stat()
	fileSize := float64(fileinfo.Size())

	totalChunks := uint64(math.Ceil(fileSize / MaxChunkSize))

	for i := uint64(0); i < totalChunks; i = i + 1 {
		currChunkSize := math.Min(MaxChunkSize, fileSize-(float64(i)*MaxChunkSize))
		currChunk := make([]byte, int(currChunkSize))

		file.Read(currChunk)
		filename := "d" + strconv.Itoa(int(i))
		if _, err := os.Create(FileDirectory + filename); err == nil {

		} else {
			fmt.Println("Fatal error ... " + err.Error())
			os.Exit(1)
		}

		ioutil.WriteFile(FileDirectory+filename, currChunk, os.ModeAppend)
	}

	return nil, nil
}
