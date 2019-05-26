package main

//import requiered libraries
import (
	//    "fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	//bzzclient "https://github.com/ethereum/go-ethereum/tree/master/swarm/api/client/client.go"
	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

type FileStructure struct {
	nameToHash map[string]string `json:"nametoHash"`
}

func LoadFileStructure(path string) (*FileStructure, error) {

	var fs FileStructure
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(conf, &fs); err != nil {
		return nil, err
	}

	return &fs, nil
}

var (
	newFile *os.File
	err     error
)

func main() {

	var path string
	//define Swarm client
	client := bzzclient.NewClient("http://swarm.dappnode:8500")
	//Create file retrieval log info

	newFile, err = os.Create("files/retrives.txt")
	//Read directory

	files, err := ioutil.ReadDir("files/")
	if err != nil {
		log.Fatal(err)
	}

	var fs map[string]string = make(map[string]string)
	for _, file := range files {
		//fmt.Println(file.Name())
		path = "files/"
		path += file.Name()
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		fileSw, err := bzzclient.Open(path)
		if err != nil {
			log.Fatal(err)
		}
		//Upload file[n]
		manifestHash, err := client.Upload(fileSw, "", false)
		if err != nil {
			log.Fatal(err)
		}
		//Log retrieval information

		fs[file.Name()] = manifestHash
    byteArr, err := json.Marshal(fs)
    ioutil.WriteFile("files/retrives.txt", byteArr, 0644)
	}


}
