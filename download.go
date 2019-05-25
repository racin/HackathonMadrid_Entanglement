package main

//import requiered libraries
import (
   "fmt"
    "io/ioutil"
    "strings"
    "strconv"
    "github.com/racin/HackathonMadrid_Entanglement/Code/Entangler"
    "encoding/json"

    //bzzclient "https://github.com/ethereum/go-ethereum/tree/master/swarm/api/client/client.go"
    bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

const file_to_retrieve = "D6"
const index = 6



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

func main() {

  br, _, _ := Entangler.GetBackwardNeighbours (index)
  fr, _, _ := Entangler.GetForwardNeighbours (index)

  //Getting filenames to XOR
  values1 := []string{"p",strconv.Itoa(br),"_",strconv.Itoa(index)}
  file1:= strings.Join(values1, "")
  values2 := []string{"p",strconv.Itoa(fr),"_",strconv.Itoa(index)}
  file2:= strings.Join(values2, "")

//Get hashes from file


// assign values from config file to variables
config,_ := LoadFileStructure("files/retrives.txt")
HashBck := config[file1]
HashFwd := config[file2]
 fmt.Println("F1:",HashBck)
 fmt.Println("F2:",HashFwd)
 fmt.Println("a", config["files/hello.txt"])

 //Retrive hashes
 client := bzzclient.NewClient("http://127.0.0.1:8500")

      fileB, err := client.Download(HashBck, "")
      fileF, err := client.Download(HashFwd, "")
//file, err := client.Download(config["files/hello.txt"], "")
     content, err := ioutil.ReadAll(file)
fmt.Println(string(content)) // hello world
fmt.Println(err) // hello world
//XOR PARITY CHUNKS
XORByteSlice(fileB,fileF)

}
