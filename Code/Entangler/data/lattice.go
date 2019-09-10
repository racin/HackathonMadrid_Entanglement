package data

import (
	"strings"
	"strconv"
	"math"
	"reflect"
	"fmt"
)
// s Horizontal strands. p Helical strands
type Lattice struct {
	// DataNodes   []*DataBlock
	// ParityNodes []*ParityBlock
	Blocks            []*Block
	NumBlocks         int
	Alpha             int
	S, P              int
	DataRequest       chan *DownloadRequest
	DataStream        chan *DownloadResponse
	Config            map[string]string
	confpath          string
	MissingDataBlocks int
}

// TODO: Exact calculations for strandLen
func sortConfigKeys(keys []reflect.Value, alpha, s, p int) []reflect.Value {
	sortedKeys := make([]reflect.Value,len(keys))
	strandLen := len(keys)/(alpha +1)
	dataKeys := make([]reflect.Value,strandLen)
	hpKeys := make([]reflect.Value,strandLen)
	rpKeys := make([]reflect.Value,strandLen)
	lpKeys := make([]reflect.Value,strandLen)

	for i, key := range keys {
		keyStr := key.String()
		isParity := keyStr[:1] == "p"
		if isParity {
			position := keyStr[1:]
			leftright := strings.Split(position, "-")
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
		} else {
			dataKeys = append(dataKeys, key)
		}
	}
	copy(sortedKeys, dataKeys)
	copy(sortedKeys[strandLen:], hpKeys)
	copy(sortedKeys[strandLen*2:], rpKeys)
	copy(sortedKeys[strandLen*3:], lpKeys)
	//return append(dataKeys[:], append(hpKeys[:], append(rpKeys[:], lpKeys[:]...)...)...)
	return sortedKeys
}

func NewLattice(esize, alpha, s, p int, confpath string) *Lattice {
	numBlocks := (1 + alpha) * esize
	conf, _ := LoadFileStructure(confpath)
	sortedKeys := sortConfigKeys(reflect.ValueOf(conf).MapKeys(), alpha, s, p)
	blocks := make([]*Block, numBlocks)

	for _, key := range sortedKeys {
		keyStr := key.String()
		fmt.Println("Key: " + keyStr)
		// Construct blocks
		isParity := keyStr[:1] == "p"
		position := keyStr[1:]
		if isParity {
			leftright := strings.Split(position, "-")
			left, _ := strconv.Atoi(leftright[0])
			right, _ := strconv.Atoi(leftright[1])
			var class StrandClass

			fr, fh, fl := GetForwardNeighbours(left, s, p)
			switch right {
			case fr:
				class = Right
			case fh:
				class = Horizontal
			case fl:
				class = Left
			}
			
			blocks = append(blocks, &Block {})
		} else {
			pos, _ := strconv.Atoi(position)
			
			blocks = append(blocks, &Block {})
		}
	}
	
	for key, _ := range conf {
		// Construct blocks
		isParity := key[:1] == "p"
		position := key[1:]
		if isParity {
			leftright := strings.Split(position, "-")
			left, _ := strconv.Atoi(leftright[0])
			right, _ := strconv.Atoi(leftright[1])
			var class StrandClass

			fr, fh, fl := GetForwardNeighbours(left, s, p)
			switch right {
			case fr:
				class = Right
			case fh:
				class = Horizontal
			case fl:
				class = Left
			}
			
			blocks = append(blocks, &Block {})
		} else {
			pos, _ := strconv.Atoi(position)
			
			blocks = append(blocks, &Block {})
		}

	}
	return &Lattice{
		// DataNodes:   make([]*DataBlock, esize),
		// ParityNodes: make([]*ParityBlock, alpha*esize),
		Blocks:   make([]*Block, numBlocks),
		Alpha:    alpha,
		S:        s,
		P:        p,
		confpath: confpath,
		Config:   conf,
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
	Left     []*Block
	Right    []*Block
	Position int
	Data     []byte
	IsParity bool
	Class    StrandClass
}

type StrandClass int

const (
	Horizontal StrandClass = iota
	Right
	Left
)

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
func GetBackwardNeighbours(index int) (r, h, l int) {
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

func GetMemoryPosition(index int) (r, h, l int) {
	// Get the position in the ParityMemory array where the parity is located
	// For now this will recursively call the GetBackwardNeighbours function

	h = ((index - 1) % S) + S
	r, l = index, index

	for ; r > S; r, _, _ = GetBackwardNeighbours(r) {
	}

	switch r {
	case 1:
		r = 0
		break
	case 2:
		r = 4
		break
	case 3:
		r = 3
		break
	case 4:
		r = 2
		break
	case 5:
		r = 1
		break
	}

	for ; l > S; _, _, l = GetBackwardNeighbours(l) {
	}

	switch l {
	case 1:
		l = 11
		break
	case 2:
		l = 12
		break
	case 3:
		l = 13
		break
	case 4:
		l = 14
		break
	case 5:
		l = 10
		break
	}

	return
}