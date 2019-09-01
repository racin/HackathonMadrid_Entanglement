package Entangler

import (
	"errors"
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler/data"
)

type Lattice data.Lattice

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
