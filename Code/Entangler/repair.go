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

func getMaxStrandMatch(arr []int) (index, max int) {
	for i := 0; i < len(arr); i++ {
		if arr[i] > max {
			max = arr[i]
			index = i
		}
	}
	return
}

func (l *Lattice) HierarchicalRepair(block *Block) *Block {
	if block == nil {
		return nil
	} else if block.Data != nil {
		return block // No need to repair.
	}

	// Data repair

	if d, ok := block.Base.(*Data); ok {
		if data == nil {
			return block
		}
		strandMatch := make([]int, Alpha)
		for i := 0; i < Alpha; i++ {
			if d.Left[i] != nil {
				strandMatch[i] = strandMatch[i] + 1
			}
			if d.Right[i] == nil {
				strandMatch[i] = strandMatch[i] + 1
			}
		}

		mI, mC := getMaxStrandMatch(strandMatch)

		// Result, _ := Entangler.XORByteSlice(contentA, contentB)

		// If its 2 we already have the parities we need.
		if mC == 2 {
			// XOR Left & Right
		} else if mC == 1 {
			// Missing left or right. Repair one of them.
			b := l.HierarchicalRepair(&Block{Base: data.Left, Data: nil})
			// XOR recovered with already existing
		} else {
			bLeft := l.HierarchicalRepair(&Block{Base: data.Left, Data: nil})
			bRight := l.HierarchicalRepair(&Block{Base: data.Right, Data: nil})
		}
		return &Block{Base: d, Data: nil}

	}

	// Parity repair

	if parity, ok := block.Base.(*Parity); ok {
		return &Block{Base: parity, Data: block.Data}
	}

	return block
}

func (l *Lattice) XORBlocks(left *Block, right *Block) (*Block, error) {
	// Case 1: Both are Parity
	// Case 2: One is Parity, one is Data
	var retBlock *Block
	var err error
	var ok bool
	var d *data.DataBlock
	var p1, p2 *Parity

	d, ok = left.Base.(*data.DataBlock)
	if !ok {
		p1, ok = left.Base.(*Parity)
		if !ok {
			return nil, errors.New("cannot type cast left block")
		}
	} else {
		d, ok = right.Base.(*data.DataBlock)
		if !ok {
			p2, ok = right.Base.(*Parity)
			if !ok {
				return nil, errors.New("cannot type cast right block")
			}
		}
	}

	if p1 != nil && d != nil {
		// Left parity and data block. Can reconstruct right parity
		if p1.Right != d {
			return nil, errors.New("left parity is not associated with data block")
		}

		d.Right[p1.Class].Data, err = XORByteSlice(p1.Data, d.Data)
		if err != nil {
			return nil, err
		}

		return &Block{Base: d.Right[p1.Class]}, nil

	} else if p2 != nil && d != nil {
		// Right parity and data block. Can reconstruct left parity
		if p2.Left != d {
			return nil, errors.New("right parity is not associated with data block")
		}

		d.Left[p2.Class].Data, err = XORByteSlice(d.Data, p2.Data)
		if err != nil {
			return nil, err
		}

		return &Block{Base: d.Left[p2.Class]}, nil
	} else {
		// Two parities
		if p1.Right != p2.Left {
			return nil, errors.New("parity blocks do not match")
		}

		p1.Right.Data, err = XORByteSlice(p1.Data, p2.Data)
		if err != nil {
			return nil, err
		}

		return &Block{Base: p1.Right}, nil
	}
}
func (l *Lattice) RoundrobinRepair(block *data.LatticeBlock) {

}
