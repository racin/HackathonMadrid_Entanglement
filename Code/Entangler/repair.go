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
		l.DataRequest <- &data.DownloadRequest{}

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
}

func (l *Lattice) RoundrobinRepair(block *data.LatticeBlock) {

}
