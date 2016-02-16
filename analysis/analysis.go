package analysis

import (
	"math"

	"github.com/jgcarvalho/toniHB/pdb"
)

func distance(amd [3]float64, atm [3]float64) float64 {
	sqX := math.Pow(amd[0]-atm[0], 2)
	sqY := math.Pow(amd[1]-atm[1], 2)
	sqZ := math.Pow(amd[2]-atm[2], 2)
	return math.Sqrt(sqX + sqY + sqZ)
}

// TODO should contact consider valid atom types?
func inContact(amd pdb.Amide, atm pdb.Atom, dist float64) bool {
	return distance(amd.XYZ, atm.XYZ) < dist
}

func Run(amd []pdb.Amide, atm []pdb.Atom, maxDist float64, hbDist float64, hbAngle float64) {
	for i := 0; i < len(amd); i++ {
		for j := 0; j < len(atm); j++ {
			// skip if is the same atom
			if amd[i].PDBNumber != atm[j].PDBNumber {
				// TODO should contact consider valid atom types?
				if inContact(amd[i], atm[j], maxDist) {
					amd[i].NumContacts++
				}
			}
		}
	}
}
