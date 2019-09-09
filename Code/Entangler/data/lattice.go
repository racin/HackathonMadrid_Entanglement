package data

// s Horizontal strands. p Helical strands
type Lattice struct {
	// DataNodes   []*DataBlock
	// ParityNodes []*ParityBlock
	Blocks      []*Block
	NumBlocks   int
	Alpha       int
	S, P        int
	DataRequest chan *DownloadRequest
	DataStream  chan *DownloadResponse
	Config      map[string]string
	confpath    string
}

func NewLattice(esize, alpha, s, p int, confpath string) *Lattice {
	numBlocks := (1 + alpha) * esize
	conf, _ := LoadFileStructure(confpath)
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
