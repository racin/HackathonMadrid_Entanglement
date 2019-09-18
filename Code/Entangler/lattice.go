package Entangler

import (
	"fmt"
	"math"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const MaxSizeChunk int = 3900

// s Horizontal strands. p Helical strands
type Lattice struct {
	// DataNodes   []*DataBlock
	// ParityNodes []*ParityBlock
	Blocks            []*Block
	DataBlocks        []*Block
	NumDataBlocks     int
	Alpha             int
	S, P              int
	DataRequest       chan *DownloadRequest
	DataStream        chan *Block
	Config            map[string]string
	confpath          string
	MissingDataBlocks int
	MaxChunkSize      int
}

// TODO: Exact calculations for strandLen
func sortConfigKeys(keys []reflect.Value, alpha, s, p int) ([]reflect.Value, []reflect.Value, []reflect.Value, []reflect.Value) {
	// sortedKeys := make([]reflect.Value,len(keys))
	strandLen := len(keys) / (alpha + 1)
	dataKeys := make([]reflect.Value, strandLen)
	hpKeys := make([]reflect.Value, 0, strandLen)
	rpKeys := make([]reflect.Value, 0, strandLen)
	lpKeys := make([]reflect.Value, 0, strandLen)

	for _, key := range keys {
		keyStr := key.String()
		isParity := keyStr[:1] == "p"
		position := keyStr[1:]
		if isParity {
			leftright := strings.Split(position, "_")
			left, _ := strconv.Atoi(leftright[0])
			right, _ := strconv.Atoi(leftright[1])
			fr, fh, fl := GetForwardNeighbours(left, s, p)
			switch right {
			case fr:
				rpKeys = append(rpKeys, key)
			case fh:
				hpKeys = append(hpKeys, key)
			case fl:
				lpKeys = append(lpKeys, key)
			}
		} else if keyStr[:1] == "d" {
			// Sort the data keys
			pos, _ := strconv.Atoi(position)
			dataKeys[pos-1] = key
			//dataKeys = append(dataKeys, key)
		}
	}

	return dataKeys, hpKeys, rpKeys, lpKeys
}

func createParities(conf map[string]string,
	keys []reflect.Value, blocks []*Block,
	class StrandClass, unAvailablePct int) []*Block {
	for _, key := range keys {
		keyStr := key.String()
		position := keyStr[1:]
		leftright := strings.Split(position, "_")
		left, _ := strconv.Atoi(leftright[0])
		right, _ := strconv.Atoi(leftright[1])
		rnd := rand.Intn(100)
		var unavail bool
		if rnd < unAvailablePct {
			unavail = true
		}

		b := &Block{IsParity: true,
			Class:         class,
			Identifier:    conf[keyStr],
			Data:          make([]byte, 0, MaxSizeChunk),
			Left:          make([]*Block, 0, 1),
			Right:         make([]*Block, 0, 1),
			IsUnavailable: unavail}

		if left > 0 {
			if dataLeft := blocks[left-1]; dataLeft != nil {
				b.Left = []*Block{dataLeft}
				dataLeft.Right = append(dataLeft.Right, b)
			}
		}

		if dataRight := blocks[right-1]; dataRight != nil {
			b.Right = []*Block{dataRight}
			dataRight.Left = append(dataRight.Left, b)
		}

		blocks = append(blocks, b)
	}
	return blocks
}

func createDataBlocks(conf map[string]string, keys []reflect.Value,
	blocks []*Block, alpha int, unAvailablePct int) []*Block {
	for _, key := range keys {
		keyStr := key.String()
		position := keyStr[1:]
		pos, _ := strconv.Atoi(position)
		rnd := rand.Intn(100)
		var unavail bool
		if rnd < unAvailablePct {
			unavail = true
		}

		b := &Block{Position: pos, IsParity: false,
			Left:          make([]*Block, 0, alpha),
			Right:         make([]*Block, 0, alpha),
			Identifier:    conf[keyStr],
			Data:          make([]byte, 0, MaxSizeChunk),
			IsUnavailable: unavail}
		blocks = append(blocks, b)
	}
	return blocks
}
func NewLattice(alpha, s, p int, confpath string, datarequest chan *DownloadRequest) *Lattice {
	//numBlocks := (1 + alpha) * esize
	conf, _ := LoadFileStructure(confpath)
	dataKeys, hpKeys, rpKeys, lpKeys := sortConfigKeys(reflect.ValueOf(conf).MapKeys(), alpha, s, p)
	blocks := make([]*Block, 0, len(conf))
	rand.Seed(time.Now().UnixNano())
	//datablocks := make(map[string]*Block, len(dataKeys))
	//datablocks := make([]*Block, len(dataKeys))

	blocks = createDataBlocks(conf, dataKeys, blocks, alpha, 15)
	//copy(datablocks, blocks) // Blocks should be sorted already.

	blocks = createParities(conf, hpKeys, blocks, Horizontal, 15)
	blocks = createParities(conf, rpKeys, blocks, Right, 15)
	blocks = createParities(conf, lpKeys, blocks, Left, 15)

	return &Lattice{
		// DataNodes:   make([]*DataBlock, esize),
		// ParityNodes: make([]*ParityBlock, alpha*esize),
		NumDataBlocks:     len(dataKeys),
		MissingDataBlocks: len(dataKeys),
		Blocks:            blocks,
		Alpha:             alpha,
		S:                 s,
		P:                 p,
		confpath:          confpath,
		DataStream:        make(chan *Block, len(conf)*5),
		MaxChunkSize:      3900,
		DataRequest:       datarequest,
		//Config:   conf,
	}
}

type LatticeBlock struct {
	Data []byte
	Base interface{}
}

type DataBlock struct {
	LatticeBlock
	Left     []*ParityBlock
	Right    []*ParityBlock
	Position int
}

type ParityBlock struct {
	LatticeBlock
	Left   *DataBlock
	Right  *DataBlock
	Strand int
	Class  StrandClass
}

type Block struct {
	Left           []*Block
	Right          []*Block
	Position       int
	Data           []byte
	IsParity       bool
	Class          StrandClass
	Identifier     string
	DownloadStatus int
	WasDownloaded  bool
	IsUnavailable  bool
}

type StrandClass int

const (
	Horizontal StrandClass = iota
	Right
	Left
)

func (b *Block) LeftPos(class int) int {
	if len(b.Left) > class {
		return b.Left[class].Position
	}
	return 0
}

func (b *Block) RightPos(class int) int {
	if len(b.Right) > class {
		return b.Right[class].Position
	}
	return 0
}

func (b *Block) String() string {
	/*if b.Position != 11 && b.LeftPos(0) != 11 && b.RightPos(0) != 11 {
		return ""
	}*/
	return fmt.Sprintf("IsParity:%t, Pos: %d, Left: %d, Right: %d, HasData: %t, DownloadStatus: %d",
		b.IsParity, b.Position, b.LeftPos(0),
		b.RightPos(0), b.HasData(),
		b.DownloadStatus)
}

func (b *Block) HasData() bool {
	return b != nil && b.Data != nil && len(b.Data) != 0
}

// Alpha = 3
func GetForwardNeighbours(index, S, P int) (r, h, l int) {
	// Check is it top, center or bottom in the lattice
	// 1 -> Top, 0 -> Bottom, else Center
	var nodePos = index % S

	if nodePos == 1 {
		r = index + S + 1
		h = index + S
		l = index + (S * P) - int(math.Pow(float64(S-1), 2))
	} else if nodePos == 0 {
		r = index + (S * P) - int(math.Pow(float64(S), 2)-1)
		h = index + S
		l = index + S - 1
	} else {
		r = index + S + 1
		h = index + S
		l = index + (S - 1)
	}
	return
}

// TODO: Fix underflow naming errors on the nodes on the extreme of the lattice.
func GetBackwardNeighbours(index, S, P int) (r, h, l int) {
	// Check is it top, center or bottom in the lattice
	// 1 -> Top, 0 -> Bottom, else Center
	var nodePos = index % S

	if nodePos == 1 {
		r = index - (S * P) + int((math.Pow(float64(S), 2) - 1))
		h = index - S
		l = index - (S - 1)
	} else if nodePos == 0 {
		r = index - (S + 1)
		h = index - S
		l = index - (S * P) + int(math.Pow(float64(S-1), 2))
	} else {
		r = index - (S + 1)
		h = index - S
		l = index - (S - 1)
	}
	return
}

func GetMemoryPosition(index, S, P int) (r, h, l int) {
	// Get the position in the ParityMemory array where the parity is located
	// For now this will recursively call the GetBackwardNeighbours function

	h = ((index - 1) % S) + S
	r, l = index, index

	for ; r > S; r, _, _ = GetBackwardNeighbours(r, S, P) {
	}

	switch r {
	case 1:
		r = 0
	case 2:
		r = 4
	case 3:
		r = 3
	case 4:
		r = 2
	case 5:
		r = 1
	}

	for ; l > S; _, _, l = GetBackwardNeighbours(l, S, P) {
	}

	switch l {
	case 1:
		l = 11
	case 2:
		l = 12
	case 3:
		l = 13
	case 4:
		l = 14
	case 5:
		l = 10
	}

	return
}
