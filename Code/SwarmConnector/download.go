package SwarmConnector

//import requiered libraries
import (
	"encoding/json"
	"fmt"
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	//bzzclient "https://github.com/ethereum/go-ethereum/tree/master/swarm/api/client/client.go"
	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

const file_to_retrieve = "D6"
const index = 6

var newFile *os.File
var err error

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

func DownloadAndReconstruct(filePath string, dataIndexes ...bool) (string, error) {

	client := bzzclient.NewClient("http://swarm.dappnode")
	config, _ := LoadFileStructure("../retrives.txt")
	var allChunks [][]byte
	var err error

	// Regular download .
	//lastData := 1
	var missingDataIndex []int
	for i := 1; i <= len(dataIndexes); i++ {
		if dataIndexes[i-1] == false {
			missingDataIndex = append(missingDataIndex, i)
			fmt.Print("Missing: " + strconv.Itoa(i) + "\n")
			continue
		}
		fmt.Print("NOT Missing: " + strconv.Itoa(i) + "\n")
		hash := config["d"+strconv.Itoa(i)]
		fmt.Println(hash)
		dataChunk, err := client.Download(hash, "")
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println(".OK??")
		content, err := ioutil.ReadAll(dataChunk)
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Println("!!!!!!")
		allChunks = append(allChunks, content)
		//lastData = i + 1
	}

	for i := 0; i < len(missingDataIndex); i++ {
		br, _, _ := Entangler.GetBackwardNeighbours(missingDataIndex[i])
		fr, _, _ := Entangler.GetForwardNeighbours(missingDataIndex[i])

		//Getting filenames to XOR
		values1 := []string{"p", strconv.Itoa(br), "_", strconv.Itoa(missingDataIndex[i])}
		file1 := strings.Join(values1, "")
		values2 := []string{"p", strconv.Itoa(fr), "_", strconv.Itoa(missingDataIndex[i])}
		file2 := strings.Join(values2, "")

		HashBck := config[file1]
		HashFwd := config[file2]

		fileA, _ := client.Download(HashBck, "")
		fileB, _ := client.Download(HashFwd, "")

		contentA, _ := ioutil.ReadAll(fileA)
		contentB, _ := ioutil.ReadAll(fileB)
		//fmt.Println(string(contentA)) // hello world
		//fmt.Println(err) // hello world

		//XOR PARITY CHUNKS
		Result, _ := Entangler.XORByteSlice(contentA, contentB)

		allChunks = append(allChunks, Result)

		//Create Result file
		_, err = os.Create(Entangler.DownloadDirectory + "d" + strconv.Itoa(missingDataIndex[i]))

		//Write XOR content to file
		ioutil.WriteFile(Entangler.DownloadDirectory+"d"+strconv.Itoa(missingDataIndex[i]), Result, 0644)
	}

	/*for i := 0; i < len(allChunks); i++ {
		dataChunk, err := Entangler.ReadChunk(Entangler.DownloadDirectory + "d" + strconv.Itoa(i))
		if err != nil {
			fmt.Println(err)
		}
		allChunks[i] = dataChunk
	}*/

	Entangler.RebuildFile(filePath, allChunks...)

	return filePath, err
}

func Download() {

	br, _, _ := Entangler.GetBackwardNeighbours(index)
	fr, _, _ := Entangler.GetForwardNeighbours(index)

	//Getting filenames to XOR
	values1 := []string{"p", strconv.Itoa(br), "_", strconv.Itoa(index)}
	file1 := strings.Join(values1, "")
	values2 := []string{"p", strconv.Itoa(fr), "_", strconv.Itoa(index)}
	file2 := strings.Join(values2, "")

	//Get hashes from file

	// assign values from config file to variables
	config, _ := LoadFileStructure("../retrives.txt")
	HashBck := config[file1]
	HashFwd := config[file2]
	/*
	   fmt.Println("F1:",HashBck)
	   fmt.Println("F2:",HashFwd)
	   fmt.Println("a", config["files/hello.txt"])
	*/
	//Retrive hashes
	client := bzzclient.NewClient("http://127.0.0.1:8500")

	fileA, err := client.Download(HashBck, "")
	fileB, err := client.Download(HashFwd, "")

	if err != nil {
		return
	}
	//file, err := client.Download(config["files/hello.txt"], "")
	contentA, err := ioutil.ReadAll(fileA)
	contentB, err := ioutil.ReadAll(fileB)
	//fmt.Println(string(contentA)) // hello world
	//fmt.Println(err) // hello world

	//XOR PARITY CHUNKS
	Result, _ := Entangler.XORByteSlice(contentA, contentB)
	//Create Result file
	//DataFile, err = os.Create("files/Result.txt")

	//Write XOR content to file
	ioutil.WriteFile("files/Results/Result.txt", Result, 0644)

}
