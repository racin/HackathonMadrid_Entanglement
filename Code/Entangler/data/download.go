package data

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
)

type DownloadRequest struct {
	Key   string
	Block *Block
}

func getSwarmHash(block *Block) string {
	config, _ := LoadFileStructure("../retrives.txt")
	if block.IsParity {
		return config["p"+strconv.Itoa(block.Left[0].Position)+"_"+strconv.Itoa(block.Right[0].Position)]
	}
	return config["d"+strconv.Itoa(block.Position)]
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

func (l *Lattice) NewDownload(block *Block) *DownloadRequest {
	return &DownloadRequest{Block: block, Key: getSwarmHash(block)}
}
