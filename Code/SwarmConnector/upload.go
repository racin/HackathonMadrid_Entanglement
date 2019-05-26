package SwarmConnector

//import requiered libraries
import (
	//    "fmt"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	//bzzclient "https://github.com/ethereum/go-ethereum/tree/master/swarm/api/client/client.go"
	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

type FileStructure map[string]string

const UploadDirectory = "../files/"

func UploadAllChunks() {
	//define Swarm client
	client := bzzclient.NewClient("http://swarm.dappnode")
	//Create file retrieval log info

	newFile, err = os.Create("../retrives.txt")
	//Read directory

	files, err := ioutil.ReadDir(UploadDirectory)
	if err != nil {
		log.Fatal(err)
	}

	var fs map[string]string = make(map[string]string)
	for _, file := range files {
		if file.Name() == ".DS_Store" {
			continue
		}
		if file.Size() == 0 {
			continue
		}
		path := UploadDirectory + file.Name()
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

		fs[filepath.Base(file.Name())] = manifestHash
		byteArr, err := json.Marshal(fs)
		ioutil.WriteFile("../retrives.txt", byteArr, 0644)
	}

}
