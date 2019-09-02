package Entangler

import (
	"errors"
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler/data"
)

type Lattice data.Lattice
type Block data.LatticeBlock
type Data data.DataBlock
type Parity data.ParityBlock

func (l *Lattice) Reconstruct() ([]byte, error) {
	out := make([]byte, len(l.DataNodes))
	for i := 0; i < len(l.DataNodes); i++ {
		if l.DataNodes[i] == nil {
			return nil, errors.New("missing data block")
		}
		out = append(out, l.DataNodes[i].Data[:]...)
	}

	return out, nil
}

func (l *Lattice) HierarchicalRepair(block *Block) *Block {
	if block == nil {
		return nil
	}

	// Data repair

	if data, ok := block.Base.(*Data); ok {
		return &Block{Base: data, Data: nil}

	}

	// Parity repair

	if parity, ok := block.Base.(*Parity); ok {
		return &Block{Base: parity, Data: nil}
	}

	return nil
}

func (l *Lattice) RoundrobinRepair(block *data.LatticeBlock) {

}
