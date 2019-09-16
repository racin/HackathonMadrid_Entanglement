package SwarmConnector

//import requiered libraries
import (
	//    "fmt"
	"encoding/json"
	"fmt"
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
	//Create file retrieval log info

	newFile, err = os.Create("../retrives.txt")
	//Read directory

	files, err := ioutil.ReadDir(UploadDirectory)
	if err != nil {
		log.Fatal(err)
	}

	var fs map[string]string = make(map[string]string)
	for _, file := range files {

		//Upload file[n]
		manifestHash := SwarmUpload(file)
		if manifestHash == "" {
			continue
		}
		//Log retrieval information

		fs[filepath.Base(file.Name())] = manifestHash
		byteArr, _ := json.Marshal(fs)
		ioutil.WriteFile("../retrives.txt", byteArr, 0644)
	}

}

func UploadLargeFile() {
	//define Swarm client
	//Create file retrieval log info

	newFile, err = os.Create("../retrives.txt")
	//Read directory

	files, err := ioutil.ReadDir(UploadDirectory)
	if err != nil {
		log.Fatal(err)
	}

	var fs map[string]string = make(map[string]string)
	for _, file := range files {

		//Upload file[n]
		manifestHash := SwarmUpload(file)
		if manifestHash == "" {
			continue
		}
		//Log retrieval information

		fs[filepath.Base(file.Name())] = manifestHash
		byteArr, _ := json.Marshal(fs)
		ioutil.WriteFile("../retrives.txt", byteArr, 0644)
	}

}

func SwarmUpload(fileInfo os.FileInfo) string {
	if fileInfo.Name() == ".DS_Store" {
		return ""
	}
	if fileInfo.Size() == 0 {
		return ""
	}
	path := UploadDirectory + fileInfo.Name()
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	fileSw, err := bzzclient.Open(path)
	if err != nil {
		fmt.Println("ERR 1" + err.Error())
		log.Fatal(err)
	}

	client := bzzclient.NewClient("https://swarm-gateways.net")
	manifestHash, err := client.Upload(fileSw, "", false)
	if err != nil {
		fmt.Println("ERR 2" + err.Error())
		return SwarmUpload(fileInfo)
	}
	return manifestHash
}
