package Entangler

import (
	"github.com/racin/HackathonMadrid_Entanglement/Code/Entangler/data"
)

type Lattice data.Lattice

func (l *Lattice) Reconstruct() ([]byte, error) {
	out := make([]byte, len(l.DataNodes))
	for i := 0; i < len(l.DataNodes); i++ {
		out = append(out, l.DataNodes[i].Data[:]...)
	}

	return out, nil
}
