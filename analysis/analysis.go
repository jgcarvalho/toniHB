package analysis

import (
	"math"
	"strings"

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

// TODO
func OdoHB(amd pdb.Amide, atm pdb.Atom, dist float64, angle float64) bool {
	return false
}

// TODO
func NdoHB(amd pdb.Amide, atm pdb.Atom, dist float64, angle float64) bool {
	return false
}

func Run(amd []pdb.Amide, atm []pdb.Atom, validAtoms string, maxDist float64, hbDist float64, hbAngle float64) {
	for i := 0; i < len(amd); i++ {
		for j := 0; j < len(atm); j++ {
			// skip if atomtype is not in the user input
			if strings.Contains(validAtoms, atm[j].Type) {
				// skip if is the same atom
				if amd[i].PDBNumber != atm[j].PDBNumber {
					// TODO should contact consider valid atom types?
					if inContact(amd[i], atm[j], maxDist) {
						amd[i].NumContacts++
					}
					if OdoHB(amd[i], atm[j], hbDist, hbAngle) {
						amd[i].NumOhb++
					}
					if NdoHB(amd[i], atm[j], hbDist, hbAngle) {
						amd[i].NumNhb++
					}
				}
			}
		}
	}
}
