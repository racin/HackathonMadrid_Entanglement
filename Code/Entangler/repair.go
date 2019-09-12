package Entangler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	//"time"
)

func (l *Lattice) RebuildFile(filePath string) error {
	if l.MissingDataBlocks != 0 {
		return errors.New("lattice is missing data blocks")
	}
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	for i := 0; i < l.NumDataBlocks; i++ {
		dat := l.Blocks[i].Data
		if dat == nil || len(dat) == 0 {
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

func (l Lattice) HierarchicalRepair(block *Block, result chan *Block) *Block {
	if block == nil {
		return nil
	} else if block.Data != nil && len(block.Data) != 0 {
		return block // No need to repair.
	}

	// Data repair
	if !block.IsParity {
		strandMatch := make([]int, Alpha)
		for i := 0; i < Alpha; i++ {
			if len(block.Left) > i && block.Left[i] != nil && len(block.Left[i].Data) > 0 {
				strandMatch[i] = strandMatch[i] + 1
			}
			fmt.Printf("i is: %d, len is: %d, IsParity: %c\n", i, len(block.Right), block.IsParity)
			if len(block.Right) > i && block.Right[i] != nil && len(block.Right[i].Data) > 0 {
				strandMatch[i] = strandMatch[i] + 1
			}
		}
		fmt.Printf("Want to repairing block. Pos: %d\n", block.Position)
		mI, mC := getMaxStrandMatch(strandMatch)
		fmt.Printf("Want to repairing block. Pos: %d, mI: %d, mC: %d\n", block.Position, mI, mC)
		// Result, _ := Entangler.XORByteSlice(contentA, contentB)

		// If its 2 we already have the parities we need.
		if mC == 2 {
			// XOR Left & Right
			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])
		} else if mC == 1 {
			// Missing left or right. Repair one of them.
			if block.Left[mI].Data != nil {
				// Repair right
				l.HierarchicalRepair(block.Right[mI], nil)
			} else {
				// Repair left
				l.HierarchicalRepair(block.Left[mI], nil)
			}
			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])

			// XOR recovered with already existing
		} else {
			fmt.Printf("Repairing block. Pos: %d\n", block.Position)
			l.HierarchicalRepair(block.Right[mI], nil)
			l.HierarchicalRepair(block.Left[mI], nil)

			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])
		}

		if result != nil && block.Data != nil && len(block.Data) != 0 {
			fmt.Printf("Sending block: %d, back on the channel\n", block.Position)
			result <- block
		}
		return block

	} else {
		// Parity repair
		fmt.Println("test")
		res := make(chan *Block)
		defer close(res)
		l.DataRequest <- &DownloadRequest{Block: block, Result: res}
		select {
		case dl := <-res:
			if dl.Data == nil || len(dl.Data) == 0 {
				return nil
				// repair
				fmt.Printf("Block was missing. Position: %d\n", dl.Position)

				// Try to request parity
				if block.Left[0].Data != nil && block.Left[0].Left[block.Position].Data != nil {
					block, _ = l.XORBlocks(block.Left[0], block.Left[0].Left[block.Position])
					return block
				} else if block.Right[0].Data != nil && block.Right[0].Right[block.Position].Data != nil {
					block, _ = l.XORBlocks(block.Right[0].Right[block.Position], block.Right[0])
					return block
				} else {
					leftData := l.HierarchicalRepair(block.Left[0], nil)
					leftParity := l.HierarchicalRepair(block.Left[0].Left[block.Position], nil)
					block, _ = l.XORBlocks(leftData, leftParity)
					return block
				}
			} else {
				fmt.Printf("Parity download success. Position: %d_%d\n", dl.Left[0].Position, dl.Right[0].Position)
				return dl
			}
		}

		return block
	}
}

// func (l *Lattice) NewDownload(block *Block, f func(*Block, error)) {
// 	l.DataRequest <- &DownloadRequest{Block: block, Key: GetSwarmHash(block, &l.Config)}
// 	go func() {
// 		select {
// 		case <-time.After(30 * time.Second):
// 			f(nil, errors.New("download timeout expired"))
// 		case b := <-l.DataStream:
// 			if block.IsParity == b.IsParity && block.Position == b.Position && block.Class == b.Class {
// 				f(b, nil)
// 			}
// 		}
// 	}()
// }

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

func (l *Lattice) RoundrobinRepair(block *Block, result chan *Block) *Block {
	if block == nil {
		return nil
	} else if block.Data != nil && len(block.Data) != 0 {
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
				l.HierarchicalRepair(block.Right[mI], nil)
			} else {
				// Repair left
				l.HierarchicalRepair(block.Left[mI], nil)
			}
			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])

			// XOR recovered with already existing
		} else {
			l.HierarchicalRepair(block.Right[mI], nil)
			l.HierarchicalRepair(block.Left[mI], nil)

			block, _ = l.XORBlocks(block.Left[mI], block.Right[mI])
		}
		if result != nil {
			result <- block
		}
		return block

	} else {
		// Parity repair
		res := make(chan *Block)
		defer close(res)
		l.DataRequest <- &DownloadRequest{Block: block, Result: res}
		select {
		case dl := <-res:
			if dl.Data == nil || len(dl.Data) == 0 {
				// repair
				fmt.Printf("Block was missing. Position: %d\n", dl.Position)

				// Try to request parity
				if block.Left[0].Data != nil && block.Left[0].Left[block.Position].Data != nil {
					block, _ = l.XORBlocks(block.Left[0], block.Left[0].Left[block.Position])
					return block
				} else if block.Right[0].Data != nil && block.Right[0].Right[block.Position].Data != nil {
					block, _ = l.XORBlocks(block.Right[0].Right[block.Position], block.Right[0])
					return block
				} else {
					leftData := l.HierarchicalRepair(block.Left[0], nil)
					leftParity := l.HierarchicalRepair(block.Left[0].Left[block.Position+1], nil)
					block, _ = l.XORBlocks(leftData, leftParity)
					return block
				}
			} else {
				fmt.Printf("Parity download success. Position: %d_%d\n", dl.Left[0].Position, dl.Right[0].Position)
				return dl
			}
		}

		return block
	}
}
