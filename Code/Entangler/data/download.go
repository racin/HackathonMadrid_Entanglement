package data

import (
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

type Config map[string]string

func LoadFileStructure(path string) (map[string]string, error) {

	var fs map[string]string = make(map[string]string)
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(conf, &fs); err != nil {
		return nil, err
	}

	return fs, nil
}

func (l *Lattice) NewDownload(block *Block) *DownloadRequest {
	return &DownloadRequest{Block: block, Key: getSwarmHash(block)}
}
