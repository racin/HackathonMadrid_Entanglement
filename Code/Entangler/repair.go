package Entangler

import (
	"errors"
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler/data"
)

type Lattice data.Lattice
type Data data.DataBlock
type Parity data.ParityBlock

func (l *Lattice) Download(block *data.Block) {

}
func (l *Lattice) Reconstruct() ([]byte, error) {
	out := make([]byte, l.NumBlocks)
	for i := 0; i < l.NumBlocks; i++ {
		b := l.Blocks[i]
		if b.IsParity {
			continue
		}
		if b.Data == nil {
			return nil, errors.New("missing data block")
		}
		out = append(out, b.Data[:]...)
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

func (l *Lattice) HierarchicalRepair(block *data.Block) *data.Block {
	if block == nil {
		return nil
	} else if block.Data != nil {
		return block // No need to repair.
	}

	// Data repair

	if !block.IsParity {
		strandMatch := make([]int, Alpha)
		for i := 0; i < Alpha; i++ {
			if block.Left[i] != nil {
				strandMatch[i] = strandMatch[i] + 1
			}
			if block.Right[i] != nil {
				strandMatch[i] = strandMatch[i] + 1
			}
		}

		mI, mC := getMaxStrandMatch(strandMatch)

		// Result, _ := Entangler.XORByteSlice(contentA, contentB)

		// If its 2 we already have the parities we need.
		if mC == 2 {
			// XOR Left & Right
			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])
		} else if mC == 1 {
			// Missing left or right. Repair one of them.
			if block.Left[mI].Data != nil {
				// Repair right
				l.HierarchicalRepair(block.Right[mI])
			} else {
				// Repair left
				l.HierarchicalRepair(block.Left[mI])
			}
			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])

			// XOR recovered with already existing
		} else {
			l.HierarchicalRepair(block.Right[mI])
			l.HierarchicalRepair(block.Left[mI])

			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])
		}
		return block

	} else {
		// Parity repair

		// Try to request parity
		if block.Left[0].Data != nil && block.Left[0].Left[block.Position].Data != nil {
			block, _ = l.XORBlocks(block.Left[0], block.Left[0].Left[block.Position])
		} else if block.Right[0].Data != nil && block.Right[0].Right[block.Position].Data != nil {
			block, _ = l.XORBlocks(block.Right[0].Right[block.Position], block.Right[0])
		}

		return block
	}
}

// XORBlocks Figures out which is the related block between a and b and attempts to repair it.
// At least one of them needs to be a parity block.
func (l *Lattice) XORBlocks(a *data.Block, b *data.Block) (*data.Block, error) {

	var err error
	// Case 1: Both is data (Invalid case)
	if !a.IsParity && !b.IsParity {
		return nil, errors.New("at least one block must be parity")
	}

	// Case 2: Both are Parity
	if a.IsParity && b.IsParity {
		if a.Right[0] == b.Left[0] {
			a.Right[0].Data, err = XORByteSlice(a.Data, b.Data)
			return a.Right[0], err
		} else if a.Left[0] == b.Right[0] {
			a.Left[0].Data, err = XORByteSlice(a.Data, b.Data)
			return a.Left[0], err
		} else {
			return nil, errors.New("blocks are not connected")
		}
	}

	// Case 3: One is Parity, one is Data
	var data, parity *data.Block
	if !a.IsParity {
		data, parity = a, b
	} else {
		data, parity = b, a
	}

	if data.Right[parity.Class] == parity { // Reconstruct left parity
		data.Left[parity.Class].Data, err = XORByteSlice(data.Data, parity.Data)
		return data.Left[parity.Class], err
	} else if data.Left[parity.Class] == parity { // Reconstruct right parity
		data.Right[parity.Class].Data, err = XORByteSlice(data.Data, parity.Data)
		return data.Right[parity.Class], err
	} else {
		return nil, errors.New("blocks are not connected")
	}
	/*
		if !a.IsParity {
			if a.Right[b.Class] == b { // Reconstruct left parity
				//a.Left[b.Class].Data, err = XORByteSlice(a.Right[b.Class].Data, b.Data)
				//r = a.Left[b.Class], x = a.Right[b.Class], y = b
				r, x, y = a.Left[b.Class], a.Right[b.Class], b
				//a.Left[b.Class].Data, err = XORByteSlice(a.Right[b.Class].Data, b.Data)
				//p1 = b
			} else if a.Left[b.Class] == b { // Reconstruct right parity
				r, x, y = a.Right[b.Class], a.Left[b.Class], b
				//a.Right[b.Class].Data, err = XORByteSlice(a.Left[b.Class].Data, b.Data)
				//p2 = b
			} else {
				return nil, errors.New("blocks are not connected")
			}
		} else if !b.IsParity {
			d = b
			if b.Right[a.Class] == a { // Reconstruct left parity
				b.Left[a.Class].Data, err = XORByteSlice(b.Right[a.Class].Data, a.Data)
				//p1 = a
			} else if b.Left[a.Class] == a { // Reconstruct right parity
				b.Right[a.Class].Data, err = XORByteSlice(b.Left[a.Class].Data, a.Data)
				//p2 = a
			} else {
				return nil, errors.New("blocks are not connected")
			}
		} else {

		}
		r.Data, err = XORByteSlice(x, y)
		return r
		if !a.IsParity {
			if !b.IsParity {
				return nil, errors.New("at least one block must be parity")
			}
			if a.Right[b.Class] != b {
				return nil, errors.New("blocks are not connected")
			}
			if a.Left[b.Class].Data != nil {
				return a.Left[b.Class], nil
			}
			// Have data block and its right parity. Can restore the left parity.
			a.Left[b.Class].Data, err = XORByteSlice(a.Data, a.Right[b.Class].Data)
		} else if !b.IsParity {
			if b.Right[b.Class] != b {
				return nil, errors.New("blocks are not connected")
			}
			if a.Left[b.Class].Data != nil {
				return nil, errors.New("(a) parity already contains data")
			}
			// Have data block and its left parity. Can restore the right parity.
			a.Left[b.Class].Data, err = XORByteSlice(a.Data, a.Right[b.Class].Data)
		}
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
		}*/
}
func (l *Lattice) RoundrobinRepair(block *data.LatticeBlock) {

}
