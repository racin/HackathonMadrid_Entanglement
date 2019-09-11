package Entangler

import (
	"bufio"
	"errors"
	"os"
	"time"
)

func (l *Lattice) RebuildFile(filePath string) error {
	if l.MissingDataBlocks != 0 {
		return errors.New("lattice is missing data blocks")
	}
	f, err := os.Create(filePath)
	if err != nil {
		os.Exit(1)
	}
	w := bufio.NewWriter(f)

	for i := 0; i < l.NumDataBlocks; i++ {
		dat := l.Blocks[i].Data
		if dat == nil {
			return errors.New("data is nil")
		}
		w.Write(dat)
	}
	w.Flush()
	return nil
}

func (l *Lattice) Download(block *Block) {

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

func (l Lattice) HierarchicalRepair(block *Block) *Block {
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
		l.NewDownload(block, func(b *Block, err error) {
			if err != nil {
				return
			}
			// Try to request parity
			if block.Left[0].Data != nil && block.Left[0].Left[block.Position].Data != nil {
				block, _ = l.XORBlocks(block.Left[0], block.Left[0].Left[block.Position])
			} else if block.Right[0].Data != nil && block.Right[0].Right[block.Position].Data != nil {
				block, _ = l.XORBlocks(block.Right[0].Right[block.Position], block.Right[0])
			}
		})

		return block
	}
}

func (l *Lattice) NewDownload(block *Block, f func(*Block, error)) {
	l.DataRequest <- &DownloadRequest{Block: block, Key: GetSwarmHash(block, &l.Config)}
	go func() {
		select {
		case <-time.After(30 * time.Second):
			f(nil, errors.New("download timeout expired"))
		case b := <-l.DataStream:
			if block.IsParity == b.Block.IsParity && block.Position == b.Block.Position && block.Class == b.Block.Class {
				f(b.Block, nil)
			}
		}
	}()
}

// XORBlocks Figures out which is the related block between a and b and attempts to repair it.
// At least one of them needs to be a parity block.
func (l *Lattice) XORBlocks(a *Block, b *Block) (*Block, error) {

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
	var data, parity *Block
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

func (l *Lattice) RoundrobinRepair(block LatticeBlock) {

}
