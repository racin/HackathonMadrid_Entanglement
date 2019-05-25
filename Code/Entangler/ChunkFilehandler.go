package Entangler

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

func WriteChunkToFile(data []byte, back, forward int) {
	filename := "p" + strconv.Itoa(back) + "_" + strconv.Itoa(forward)
	if _, err := os.Create(ChunkDirectory + filename); err == nil {

	} else {
		fmt.Println("Fatal error ... " + err.Error())
		os.Exit(1)
	}

	ioutil.WriteFile(ChunkDirectory+filename, data, os.ModeAppend)
}

func ReadChunk(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	fileinfo, _ := file.Stat()
	if err != nil {
		return nil, errors.New("Could not open chunk")
	}

	fileSize := fileinfo.Size()
	var output []byte = make([]byte, fileSize)

	file.Read(output)
	return output, nil
}
