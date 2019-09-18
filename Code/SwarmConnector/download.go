package SwarmConnector

//import requiered libraries
import (
	"encoding/json"
	"fmt"
	e "github.com/racin/HackathonMadrid_Entanglement/Code/Entangler"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	//bzzclient "https://github.com/ethereum/go-ethereum/tree/master/swarm/api/client/client.go"
	bzzclient "github.com/ethereum/go-ethereum/swarm/api/client"
)

type Lattice e.Lattice

const file_to_retrieve = "D6"
const index = 6

var newFile *os.File
var err error

type Downloader struct {
	Client *bzzclient.Client
}

type DownloadPool struct {
	lock     sync.Mutex       // Locking
	resource chan *Downloader // Channel to obtain resource from the pool
	lattice  *e.Lattice       // Shared map of retrieved blocks
	Capacity int              // Maximum capacity of the pool.
	count    int              // Current count of allocated resources.
	//Filepath     string           // Final output location
	endpoint     string
	Datarequests chan *e.DownloadRequest
	datastream   chan *e.DownloadResponse
}

func NewDownloadPool(capacity int, endpoint string) *DownloadPool {
	d := &DownloadPool{
		resource:     make(chan *Downloader, capacity),
		Datarequests: make(chan *e.DownloadRequest),
		//lattice:  data.NewLattice(Entangler.Alpha, Entangler.S, Entangler.P),
		Capacity: capacity,
		count:    0,
		//Filepath: filepath,
		endpoint: endpoint,
	}
	go func() {
		for {
			select {
			case request := <-d.Datarequests:
				// 	fmt.Printf("GOT DATA REQUEST. IsParity:%c, Pos: %d, Left: %d, Right: %d\n",
				// 		request.Block.IsParity, request.Block.Left[0].Position,
				// 		request.Block.Right[0].Position, request.Block.Position)
				go d.DownloadBlock(request.Block, request.Result)
			}
		}
	}()
	return d
}

var unAvailableData map[int]bool = map[int]bool{
	// 7:  true,
	// 11: true,
	// 15: true,
	// 16: true,
}
var unAvailableParity map[int][]int = map[int][]int{
	// 1:  []int{6, 10},
	// 6:  []int{11, 15},
	// 7:  []int{12, 11},
	// 11: []int{16, 20},
	// 16: []int{21, 25},
}

func (p *DownloadPool) DownloadBlock(block *e.Block, result chan *e.Block) {
	e.DebugPrint("GOT DATA REQUEST. %v\n", block.String())

	if unAvailableData[block.Position] ||
		(block.IsUnavailable && len(block.Left) > 0 && len(block.Left[0].Left) > 0 &&
			len(block.Right) > 0 && len(block.Right[0].Right) > 0) {
		e.DebugPrint("UNAVAILABLE BLOCK %v\n", block.String())
		e.DebugPrint("Len left: %v, Len right: %v\n", len(block.Left[0].Left), len(block.Right[0].Right))
		result <- block
		return
	}
	if block.IsParity && len(block.Left) > 0 && len(block.Right) > 0 {
		if _, ok := unAvailableParity[block.Left[0].Position]; ok {
			for i := 0; i < len(unAvailableParity[block.Left[0].Position]); i++ {
				if unAvailableParity[block.Left[0].Position][i] == block.Right[0].Position {
					e.DebugPrint("UNAVAILABLE PARITY BLOCK %v\n", block.String())
					result <- block
					return
				}
			}
		}
	}

	if block.HasData() {
		//block.DownloadStatus = 2
		e.DebugPrint("Block data already known. %v\n", block.String())
		result <- block
		return
	}

	if block.DownloadStatus != 0 {
		e.DebugPrint("Block download already queued. %v\n", block.String())
		result <- block
		return
	}
	start := time.Now().UnixNano()
	block.DownloadStatus = 1

	dl := p.reserve()

	content := make(chan []byte, 1) // Buffered chan is non-blocking
	//defer close(content)
	go func() {
		if file, err := dl.Client.Download(block.Identifier, ""); err == nil {
			if contentA, err := ioutil.ReadAll(file); err == nil {
				//block.DownloadStatus = 2
				e.DebugPrint("Completed download of block. %v\n", block.String())

				// Use Result if we get it.
				block.WasDownloaded = true
				block.Data = contentA
				content <- contentA
				p.lattice.DataStream <- block

			}
		} else {
			fmt.Println(err.Error())
		}

		p.release(dl)
		block.DownloadStatus = 0
		fmt.Printf("%t,%d,%d,%d,%t,%d,%d\n", block.IsParity, block.Position, block.LeftPos(0), block.RightPos(0), block.HasData(), start, time.Now().UnixNano())
		//result <- block

		// Dont use result
		//content <- contentA
	}()
	select {
	//case <-time.After(1 * time.Second):
	case <-time.After(5000 * time.Millisecond):
		//e.DebugPrint("TIMEOUT.%v\n", block.String())
		result <- block
	case <-content:
		if block.HasData() {
			result <- block
		}
	}
}
func (p *DownloadPool) DownloadFile(config, output string) error {
	// 1. Construct lattice
	lattice := e.NewLattice(e.Alpha, e.S, e.P, config, p.Datarequests)
	p.lattice = lattice

	// 2. Attempt to download Data Blocks
	for i := 0; i < lattice.NumDataBlocks; i++ {
		go p.DownloadBlock(lattice.Blocks[i], lattice.DataStream)
	}

	datablocks, parityblocks := 0, 0
	// 3. Issue repairs if neccesary
	go func() {
		time.Sleep(120 * time.Second)
		fmt.Println("TIMEOUT - NOT REPAIRABLE LATTICE -- IGNORE")
		os.Exit(0)
	}()
repairs:
	for {
		select {
		case dl := <-lattice.DataStream:
			if dl == nil {
				fmt.Println("FATAL ERROR. STOPPING DOWNLOAD.")
				// Try new strategy?
				os.Exit(0)
				break repairs
			} else if !dl.HasData() {
				// repair
				//e.DebugPrint("Block was missing. Position: %d\n", dl.Position)
				if dl.Position < 6 || (dl.Position > 242 && dl.Position != 245) { // Closed lattice ..
					go p.DownloadBlock(dl, lattice.DataStream)
				} else {
					//go p.DownloadBlock(dl, lattice.DataStream)
					//go lattice.HierarchicalRepair(dl, lattice.DataStream, make([]*e.Block, 0))
					go lattice.RoundrobinRepair(dl, lattice.DataStream, make([]*e.Block, 0))
				}
				//go p.DownloadBlock(dl, lattice.DataStream)
			} else {
				e.DebugPrint("Download success. %v\n", dl.String())
				if !dl.IsParity && dl.DownloadStatus != 3 {
					dl.DownloadStatus = 3
					lattice.MissingDataBlocks--
					e.DebugPrint("Data block download success. Position: %d. Missing: %d\n", dl.Position, lattice.MissingDataBlocks)

					// Due to concurrency bug. TODO: Fix with lock, counter ?
					complete := true
					for i := 0; i < lattice.NumDataBlocks; i++ {
						if !lattice.Blocks[i].HasData() {
							complete = false
							//e.DebugPrint("BREAKING OUT....%v\n", lattice.Blocks[i])
							//go lattice.HierarchicalRepair(lattice.Blocks[i], lattice.DataStream, make([]*e.Block, 0))
							break
						}
					}

					if complete { //lattice.MissingDataBlocks == 0 {
						e.DebugPrint("Received all data blocks. Position: %d\n", dl.Position)
						for i := 0; i < len(lattice.Blocks); i++ {
							if lattice.Blocks[i].HasData() && lattice.Blocks[i].WasDownloaded {
								if lattice.Blocks[i].IsParity {
									parityblocks++
								} else {
									datablocks++
								}
							}
						}
						break repairs
					}
				}
			}
		}
	}
	// 4. Rebuild the file

	fmt.Printf("Downloaded total. Datablocks: %d, Parityblocks: %d\n", datablocks, parityblocks)
	return lattice.RebuildFile(output)
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
		e.DebugPrint("Generating new resource")
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
		Client: bzzclient.NewClient(endpoint),
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
// func (d *Downloader) AsyncDownloadAndReconstruct(content map[string][]byte, dataIndex string) error {
// 	if d.CanReconstruct(dataIndex) {

// 	}

// 	file, err := d.client.Download(dp.con)
// 	if err != nil {

// 	}
// }

func StringConvertBlockIndex(index ...int) string {
	lenI := len(index)

	switch lenI {
	case 1:
		return strconv.Itoa(index[0])
	case 2:
		return "p" + strconv.Itoa(index[0]) + strconv.Itoa(index[1])
	default:
		return ""
	}
}

// func (d *Downloader) CanReconstruct(content *map[string][]byte, dataIndex string) {
// 	br, bh, bl := e.GetBackwardNeighbours(dataIndex) // Right, Horizontal, Left
// 	fr, fh, fl := e.GetForwardNeighbours(dataIndex)

// 	return (*content)[dataIndex] != nil ||
// 		((*content)[br] != nil && (*content)[fr] != nil) ||
// 		((*content)[bh] != nil && (*content)[fh] != nil) ||
// 		((*content)[bl] != nil && (*content)[fl] != nil)
// }

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
			br, _, _ := e.GetBackwardNeighbours(i, e.S, e.P)
			fr, _, _ := e.GetForwardNeighbours(i, e.S, e.P)
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
				Result, _ := e.XORByteSlice(contentA, contentB)

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
	e.RebuildFile(filePath, allChunks...)

	return filePath, err
}

// func Download() {

// 	br, _, _ := Entangler.GetBackwardNeighbours(index)
// 	fr, _, _ := Entangler.GetForwardNeighbours(index)

// 	//Getting filenames to XOR
// 	values1 := []string{"p", strconv.Itoa(br), "_", strconv.Itoa(index)}
// 	file1 := strings.Join(values1, "")
// 	values2 := []string{"p", strconv.Itoa(fr), "_", strconv.Itoa(index)}
// 	file2 := strings.Join(values2, "")

// 	//Get hashes from file

// 	// assign values from config file to variables
// 	config, _ := LoadFileStructure("../retrives.txt")
// 	HashBck := config[file1]
// 	HashFwd := config[file2]
// 	/*
// 	   fmt.Println("F1:",HashBck)
// 	   fmt.Println("F2:",HashFwd)
// 	   fmt.Println("a", config["files/hello.txt"])
// 	*/
// 	//Retrive hashes
// 	client := bzzclient.NewClient("http://127.0.0.1:8500")

// 	fileA, err := client.Download(HashBck, "")
// 	fileB, err := client.Download(HashFwd, "")

// 	if err != nil {
// 		return
// 	}
// 	//file, err := client.Download(config["files/hello.txt"], "")
// 	contentA, err := ioutil.ReadAll(fileA)
// 	contentB, err := ioutil.ReadAll(fileB)
// 	//fmt.Println(string(contentA)) // hello world
// 	//fmt.Println(err) // hello world

// 	//XOR PARITY CHUNKS
// 	Result, _ := Entangler.XORByteSlice(contentA, contentB)
// 	//Create Result file
// 	//DataFile, err = os.Create("files/Result.txt")

// 	//Write XOR content to file
// 	ioutil.WriteFile("files/Results/Result.txt", Result, 0644)

// }
