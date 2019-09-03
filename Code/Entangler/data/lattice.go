package data

// s Horizontal strands. p Helical strands
type Lattice struct {
	DataNodes   []*DataBlock
	ParityNodes []*ParityBlock
	Blocks      []*LatticeBlock
	Alpha       int
	S, P        int
}

func NewLattice(esize, alpha, s, p int) *Lattice {
	return &Lattice{
		DataNodes:   make([]*DataBlock, esize),
		ParityNodes: make([]*ParityBlock, alpha*esize),
		Alpha:       alpha,
		S:           s,
		P:           p,
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

type StrandClass int

const (
	Horizontal StrandClass = iota
	Right
	Left
)
