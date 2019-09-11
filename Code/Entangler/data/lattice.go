package data

import (
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

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
	DataStream        chan *DownloadResponse
	Config            map[string]string
	confpath          string
	MissingDataBlocks int
}

// TODO: Exact calculations for strandLen
func sortConfigKeys(keys []reflect.Value, alpha, s, p int) ([]reflect.Value, []reflect.Value, []reflect.Value, []reflect.Value) {
	// sortedKeys := make([]reflect.Value,len(keys))
	strandLen := len(keys) / (alpha + 1)
	dataKeys := make([]reflect.Value, strandLen)
	hpKeys := make([]reflect.Value, strandLen)
	rpKeys := make([]reflect.Value, strandLen)
	lpKeys := make([]reflect.Value, strandLen)

	for _, key := range keys {
		keyStr := key.String()
		isParity := keyStr[:1] == "p"
		position := keyStr[1:]
		if isParity {
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
	class StrandClass) {
	for _, key := range keys {
		keyStr := key.String()
		fmt.Println("Key: " + keyStr)
		position := keyStr[1:]
		leftright := strings.Split(position, "-")
		left, _ := strconv.Atoi(leftright[0])
		right, _ := strconv.Atoi(leftright[1])

		b := &Block{IsParity: true, Class: class, Identifier: conf[keyStr]}

		if dataLeft := blocks[left-1]; dataLeft != nil {
			b.Left = []*Block{dataLeft}
			dataLeft.Right = append(dataLeft.Right, b)
		}

		if dataRight := blocks[right-1]; dataRight != nil {
			b.Right = []*Block{dataRight}
			dataRight.Left = append(dataRight.Left, b)
		}

		blocks = append(blocks, b)
	}
}

func createDataBlocks(conf map[string]string, keys []reflect.Value,
	blocks []*Block, alpha int) {
	for _, key := range keys {
		keyStr := key.String()
		position := keyStr[1:]
		pos, _ := strconv.Atoi(position)
		b := &Block{Position: pos, IsParity: false,
			Left:       make([]*Block, alpha),
			Right:      make([]*Block, alpha),
			Identifier: conf[keyStr]}
		blocks = append(blocks, b)
	}
}
func NewLattice(alpha, s, p int, confpath string) *Lattice {
	//numBlocks := (1 + alpha) * esize
	conf, _ := LoadFileStructure(confpath)
	dataKeys, hpKeys, rpKeys, lpKeys := sortConfigKeys(reflect.ValueOf(conf).MapKeys(), alpha, s, p)
	blocks := make([]*Block, len(conf))
	//datablocks := make(map[string]*Block, len(dataKeys))
	//datablocks := make([]*Block, len(dataKeys))

	createDataBlocks(conf, dataKeys, blocks, alpha)
	//copy(datablocks, blocks) // Blocks should be sorted already.

	createParities(conf, hpKeys, blocks, Horizontal)
	createParities(conf, rpKeys, blocks, Right)
	createParities(conf, lpKeys, blocks, Left)

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
	Left       []*Block
	Right      []*Block
	Position   int
	Data       []byte
	IsParity   bool
	Class      StrandClass
	Identifier string
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
