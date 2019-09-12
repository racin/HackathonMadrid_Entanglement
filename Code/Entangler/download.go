package Entangler

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type DownloadRequest struct {
	Result chan *Block
	Block  *Block
}

func GetSwarmHash(block *Block, fs *map[string]string) string {
	if block.IsParity {
		return (*fs)["p"+strconv.Itoa(block.Left[0].Position)+"_"+strconv.Itoa(block.Right[0].Position)]
	}
	return (*fs)["d"+strconv.Itoa(block.Position)]
}

type DownloadResponse struct {
	*DownloadRequest
	Value []byte
}

func LoadFileStructure(path string) (map[string]string, error) {
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var fs map[string]string = make(map[string]string)
	if err = json.Unmarshal(conf, &fs); err != nil {
		return nil, err
	}

	return fs, nil
}
