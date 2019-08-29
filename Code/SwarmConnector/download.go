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
	"sync"

	//bzzclient "https://github.com/ethereum/go-ethereum/tree/master/swarm/api/client/client.go"
	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

const file_to_retrieve = "D6"
const index = 6

var newFile *os.File
var err error

type Downloader struct {
	client *bzzclient.Client
}

type DownloadPool struct {
	lock     sync.Mutex        // Locking
	resource chan *Downloader  // Channel to obtain resource from the pool
	content  map[string][]byte // Shared map of retrieved blocks
	Capacity int               // Maximum capacity of the pool.
	count    int               // Current count of allocated resources.
	Filepath string            // Final output location
	endpoint string
}

func NewDownloadPool(capacity int, filepath string, endpoint string) *DownloadPool {
	return &DownloadPool{
		resource: make(chan *Downloader, capacity),
		content:  make(map[string][]byte),
		Capacity: capacity,
		count:    0,
		Filepath: filepath,
		endpoint: endpoint,
	}
}

// Drain drains the pool until it has no more than n resources
func (p *DownloadPool) Drain(n int) {
	p.lock.Lock()
	defer p.lock.Unlock()
	for len(p.resource) > n {
		<-p.resource
		p.count--
	}
}

// Reserve is blocking until it returns an available tree
// it reuses free trees or creates a new one if size is not reached
// TODO: should use a context here
func (p *DownloadPool) reserve() *Downloader {
	p.lock.Lock()
	defer p.lock.Unlock()
	var d *Downloader
	if p.count == p.Capacity {
		return <-p.resource
	}
	select {
	case d = <-p.resource:
	default:
		d = newDownloader(p.endpoint)
		p.count++
	}
	return d
}

// release gives back a Downloader to the pool
func (p *DownloadPool) release(d *Downloader) {
	p.resource <- d // can never fail ...
}

// Initizalites a new downloader
func newDownloader(endpoint string) *Downloader {
	return &Downloader{
		client: bzzclient.NewClient(endpoint),
	}
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

/// Strategy 1: Hierarchical
/// Strategy 2: Round-robin (Tail latency)
func (d *Downloader) AsyncDownloadAndReconstruct(dp DownloadPool, dataIndex string) error {
	if d.CanReconstruct(dataIndex) {

	}

	file, err := d.client.Download(dp.con)
	if err != nil {

	}
}

func (d *Downloader) CanReconstruct(dp DownloadPool, dataIndex string) {
	br, bh, bl := Entangler.GetBackwardNeighbours(dataIndex) // Right, Horizontal, Left
	fr, fh, fl := Entangler.GetForwardNeighbours(dataIndex)

	return dp.content[dataIndex] != nil ||
		(dp.content[br] != nil && dp.content[fr] != nil) ||
		(dp.content[bh] != nil && dp.content[fh] != nil) ||
		(dp.content[bl] != nil && dp.content[fl] != nil)
}

func DownloadAndReconstruct(filePath string, dataIndexes ...bool) (string, error) {
	client := bzzclient.NewClient("https://swarm-gateways.net")
	config, _ := LoadFileStructure("../retrives.txt")
	var allChunks [][]byte
	var err error

	// Regular download .
	//lastData := 1
	//var missingDataIndex []int
	for i := 1; i <= len(dataIndexes); i++ {
		if dataIndexes[i-1] == false {
			//missingDataIndex = append(missingDataIndex, i)
			fmt.Print("Missing: " + strconv.Itoa(i) + "\n")
			br, _, _ := Entangler.GetBackwardNeighbours(i)
			fr, _, _ := Entangler.GetForwardNeighbours(i)
			//Getting filenames to XOR
			values1 := []string{"p", strconv.Itoa(br), "_", strconv.Itoa(i)}
			file1 := strings.Join(values1, "")
			values2 := []string{"p", strconv.Itoa(i), "_", strconv.Itoa(fr)}
			file2 := strings.Join(values2, "")
			fmt.Println(file1)
			fmt.Println(file2)
			HashBck := config[file1]
			HashFwd := config[file2]

			if HashBck == "" || HashFwd == "" {
				hash := config["d"+strconv.Itoa(i)]
				fmt.Println(hash)
				dataChunk, err := client.Download(hash, "")
				if err != nil {
					fmt.Println(err.Error())
				}
				content, err := ioutil.ReadAll(dataChunk)
				if err != nil {
					fmt.Println(err.Error())
				}
				allChunks = append(allChunks, content)
			} else {
				fileA, _ := client.Download(HashBck, "")
				fileB, _ := client.Download(HashFwd, "")

				contentA, _ := ioutil.ReadAll(fileA)
				contentB, _ := ioutil.ReadAll(fileB)
				//fmt.Println(string(contentA)) // hello world
				//fmt.Println(err) // hello world

				//XOR PARITY CHUNKS
				Result, _ := Entangler.XORByteSlice(contentA, contentB)

				allChunks = append(allChunks, Result)
			}

			//Create Result file
			//_, err = os.Create(Entangler.DownloadDirectory + "d" + strconv.Itoa(i))

			//Write XOR content to file
			//ioutil.WriteFile(Entangler.DownloadDirectory+"d"+strconv.Itoa(i), Result, 0644)
			continue
		}
		fmt.Print("NOT Missing: " + strconv.Itoa(i) + "\n")
		hash := config["d"+strconv.Itoa(i)]
		fmt.Println(hash)
		dataChunk, err := client.Download(hash, "")
		if err != nil {
			fmt.Println(err.Error())
		}
		content, err := ioutil.ReadAll(dataChunk)
		if err != nil {
			fmt.Println(err.Error())
		}
		allChunks = append(allChunks, content)
		//lastData = i + 1
	}
	fmt.Println("Length of dataIndexes: " + string(strconv.Itoa(len(dataIndexes))))
	fmt.Println("Length of allChunks: " + string(strconv.Itoa(len(allChunks))))
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
