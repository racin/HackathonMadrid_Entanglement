package Entangler

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	//"time"
)

func (l *Lattice) RebuildFile(filePath string) error {
	/*if l.MissingDataBlocks != 0 {
		return errors.New("lattice is missing data blocks")
	}*/
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	w := bufio.NewWriter(f)

	for i := 0; i < l.NumDataBlocks; i++ {
		dat := l.Blocks[i].Data
		if dat == nil || len(dat) == 0 {
			return errors.New("data is nil. " + l.Blocks[i].String())
		}
		w.Write(dat)
	}
	w.Flush()
	return nil
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

func DebugPrint(format string, a ...interface{}) (int, error) {
	if len(a) == 0 {
		return 0, nil
	}
	z, ok := a[0].(string)

	if ok && z == "" {
		return 0, nil
	}
	if false {
		return fmt.Printf(format, a...)
	}
	return 0, nil
}

func (l *Lattice) HierarchicalRepair(block *Block, result chan *Block, path []*Block) *Block {
	if block == nil {
		DebugPrint("Block is nil ?? %v\n", block.String())
		return nil
	} else if block.HasData() {
		DebugPrint("Block has data. Returning. %v\n", block.String())
		if result != nil {
			result <- block
		}
		return block // No need to repair.
	} else {
		// Already attempted to repair this block in this tree.
		for i := 0; i < len(path); i++ {
			if block == path[i] {
				return block
			}
		}
	}
	var pathStr string = ""
	for i := 0; i < len(path); i++ {
		pathStr += "||| " + path[i].String()
	}
	DebugPrint("PATH: %v\n", pathStr)

	// Data repair
	if !block.IsParity {
		for i := 0; i < l.Alpha; i++ {
			DebugPrint("Repairing block. %v, Alpha: %v\n", block.String(), i)
			// Missing left or right. Repair one of them.
			if len(block.Right) > i {
				// Repair right
				l.HierarchicalRepair(block.Right[i], nil, append(path, block))
			} else {
				if result != nil {
					result <- nil
				}

				return block // Can not repair extreme (yet..)
			}

			if len(block.Left) > i {
				// Repair left
				l.HierarchicalRepair(block.Left[i], nil, append(path, block))
			} else if block.Right[i].HasData() {
				// First data is equal to first parity (yet...)
				block.Data = block.Right[i].Data
				if result != nil {
					result <- block
				}

				return block
			}

			if len(block.Left) > i && block.Left[i].HasData() && len(block.Right) > i && block.Right[i].HasData() {
				block, _ = l.XORBlocks(block.Left[i], block.Right[i])
				DebugPrint("Reconstructed block. %v\n", block.String())
			}

			if result != nil {
				if block.HasData() {
					DebugPrint("Sending block: %d, back on the channel\n", block.Position)
					result <- block
					return block
				} else if i == l.Alpha-1 {
					// Exhausted all possibilities.
					// Wait indefinately for file ?

					// Fatal error.. Could not repair
					DebugPrint("DID NOT FIND ANY WAY TO REPAIR BLOCK. ??? :-( %v\n", block.String())
					if result != nil {
						result <- block
					}
					//result <- nil
					//return block
				}
			}
		}

	} else {
		// Parity repair
		res := make(chan *Block)
		//defer close(res)
		l.DataRequest <- &DownloadRequest{Block: block, Result: res}
		for {
			select {
			case dl := <-res:
				if !dl.HasData() {
					// repair
					DebugPrint("Parity Block was missing. %v\n", dl.String())

					// Try to request parity
					//if block.Left[0].HasData() && block.Left[0].Left[block.Position].HasData() { // Closed lattice
					if len(block.Left) > 0 && block.Left[0].HasData() && len(block.Left[0].Left) > block.Position && block.Left[0].Left[block.Position].HasData() {
						DebugPrint("Parity repair 1. %v\n", dl.String())
						block, _ = l.XORBlocks(block.Left[0], block.Left[0].Left[block.Position])
						return block
						//} else if block.Right[0].HasData() && block.Right[0].Right[block.Position].HasData() { // Closed lattice
					} else if len(block.Right) > 0 && block.Right[0].HasData() && len(block.Right[0].Right) > block.Position && block.Right[0].Right[block.Position].HasData() {
						DebugPrint("Parity repair 2. %v\n", dl.String())
						block, _ = l.XORBlocks(block.Right[0].Right[block.Position], block.Right[0])
						return block
						//} else { // Closed lattice
					}

					if len(block.Left) > 0 && len(block.Left[0].Left) > block.Position {
						DebugPrint("Parity repair 3. %v\n", dl.String())
						leftData := l.HierarchicalRepair(block.Left[0], nil, append(path, block))
						DebugPrint("Parity repair 3. Got left data. %v\n", block.Left[0].String())
						leftParity := l.HierarchicalRepair(block.Left[0].Left[block.Position], nil, append(path, block))
						DebugPrint("Parity repair 3. Got left parity. %v\n", block.Left[0].Left[block.Position].String())
						if leftData.HasData() && leftParity.HasData() {
							block, _ = l.XORBlocks(leftData, leftParity)
							return block
						}
					} else {
						DebugPrint("Parity repair 5a. %v\n", block.String())
						l.DataRequest <- &DownloadRequest{Block: block, Result: res} // Because open lattice.
					}
					if len(block.Right) > 0 && len(block.Right[0].Right) > block.Position {
						DebugPrint("Parity repair 4. %v\n", dl.String())
						rightData := l.HierarchicalRepair(block.Right[0], nil, append(path, block))
						rightParity := l.HierarchicalRepair(block.Right[0].Right[block.Position], nil, append(path, block))
						if rightData.HasData() && rightParity.HasData() {
							block, _ = l.XORBlocks(rightData, rightParity)
							return block
						}
					} else {
						DebugPrint("Parity repair 5b. %v\n", block.String())
						l.DataRequest <- &DownloadRequest{Block: block, Result: res} // Because open lattice.
					}

					// Could not repair. Just return it.
					return block
				} else {
					DebugPrint("Parity download success. %v\n", dl.String())
					return dl
				}
			}
		}
	}

	return block
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

	if len(data.Right) > int(parity.Class) && data.Right[parity.Class] == parity { // Reconstruct left parity
		data.Left[parity.Class].Data, err = XORByteSlice(data.Data, parity.Data)
		return data.Left[parity.Class], err
	} else if len(data.Left) > int(parity.Class) && data.Left[parity.Class] == parity { // Reconstruct right parity
		data.Right[parity.Class].Data, err = XORByteSlice(data.Data, parity.Data)
		return data.Right[parity.Class], err
	} else {
		return nil, errors.New("blocks are not connected")
	}
}

func (l *Lattice) RoundrobinRepair(block *Block, result chan *Block, path []*Block) *Block {
	if block == nil {
		DebugPrint("Block is nil ?? %v\n", block.String())
		return nil
	} else if block.HasData() {
		DebugPrint("Block has data. Returning. %v\n", block.String())
		if result != nil {
			result <- block
		}
		return block // No need to repair.
	} else {
		// Already attempted to repair this block in this tree.
		for i := 0; i < len(path); i++ {
			if block == path[i] {
				return block
			}
		}
	}
	var pathStr string = ""
	for i := 0; i < len(path); i++ {
		pathStr += "||| " + path[i].String()
	}
	DebugPrint("PATH: %v\n", pathStr)

	// Data repair
	if !block.IsParity {
		for i := 0; i < l.Alpha; i++ {
			res := make(chan *Block, 2)
			if len(block.Left) > i {
				l.DataRequest <- &DownloadRequest{Block: block.Left[i], Result: res}
			}
			if len(block.Right) > i {
				l.DataRequest <- &DownloadRequest{Block: block.Right[i], Result: res}
			}

			j := 0
		roundrobinrepair:
			for {
				select {
				case dl := <-res:
					if !dl.HasData() {
						break roundrobinrepair
					}
					if j == 1 {
						block, _ = l.XORBlocks(block.Left[i], block.Right[i])
						if result != nil {
							if block.HasData() {
								DebugPrint("Sending block: %d, back on the channel\n", block.Position)
								result <- block
								return block
							}
						}
					} else {
						j++
					}
				}
			}
		}
		return l.HierarchicalRepair(block, result, path)
	}

	return block
}
