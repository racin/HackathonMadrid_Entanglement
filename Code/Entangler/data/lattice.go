package data

// s Horizontal strands. p Helical strands
type Lattice struct {
	DataNodes   []DataBlock
	ParityNodes []ParityBlock
	alpha       int
	s, p        int
}

func newLattice(esize, alpha, s, p int) *Lattice {
	return &Lattice{
		DataNodes:   make([]DataBlock, esize),
		ParityNodes: make([]ParityBlock, alpha*esize),
		alpha:       alpha,
		s:           s,
		p:           p,
	}
}

type DataBlock struct {
	left     []ParityBlock
	right    []ParityBlock
	data     []byte
	position int
}

type ParityBlock struct {
	left   DataBlock
	right  DataBlock
	strand int
	class  StrandClass
}

type StrandClass int

const (
	Horizontal StrandClass = iota
	Right
	Left
)
