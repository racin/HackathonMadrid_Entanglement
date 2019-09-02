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

	if data, ok := block.Base.(*Data); ok {
		if data == nil {
			return block
		}
		strandMatch := make([]int, Alpha)
		for i := 0; i < Alpha; i++ {
			if data.Left[i] != nil {
				strandMatch[i] = strandMatch[i] + 1
			}
			if data.Right[i] == nil {
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
		return &Block{Base: data, Data: nil}

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
	var data *Data
	var p1, p2 *Parity

	data, ok = left.Base.(*Data)
	if !ok {
		p1, ok = left.Base.(*Parity)
		if !ok {
			return nil, errors.New("cannot type cast left block")
		}
	} else {
		data, ok = right.Base.(*Data)
		if !ok {
			p2, ok = right.Base.(*Parity)
			if !ok {
				return nil, errors.New("cannot type cast right block")
			}
		}
	}

	if p1 != nil && data != nil {
		// Two parities
		if p1.Right != p2.Left {
			return nil, errors.New("parity blocks do not match")
		}

		p1.Right.Data, err = XORByteSlice(p1.Data, p2.Data)
		if err != nil {
			return nil, err
		}

	} else if p2 != nil && data != nil {
		// Two parities
		if p1.Right != p2.Left {
			return nil, errors.New("parity blocks do not match")
		}

		p1.Right.Data, err = XORByteSlice(p1.Data, p2.Data)
		if err != nil {
			return nil, err
		}

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

	return
}
func (l *Lattice) RoundrobinRepair(block *data.LatticeBlock) {

}
