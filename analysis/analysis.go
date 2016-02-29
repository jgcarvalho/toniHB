package analysis

import (
	"math"
	"strings"

	"github.com/jgcarvalho/toniHB/pdb"
)

// Acceptors are defined in analysis/acceptors.go

func isAcceptor(atm pdb.Atom) bool {
	for i := 0; i < len(Acceptors); i++ {
		if Acceptors[i].ResName == atm.ResName && Acceptors[i].AtmName == atm.Name {
			return true
		}
	}
	return false
}

func hbangle(Namd [3]float64, Hamd [3]float64, atm [3]float64) float64 {
	// v1 = vector H->N and v2 = vector H->acceptor
	v1 := [3]float64{Namd[0] - Hamd[0], Namd[1] - Hamd[1], Namd[2] - Hamd[2]}
	v2 := [3]float64{atm[0] - Hamd[0], atm[1] - Hamd[1], atm[2] - Hamd[2]}

	v1mag := math.Sqrt(v1[0]*v1[0] + v1[1]*v1[1] + v1[2]*v1[2])
	v1norm := [3]float64{v1[0] / v1mag, v1[1] / v1mag, v1[2] / v1mag}

	v2mag := math.Sqrt(v2[0]*v2[0] + v2[1]*v2[1] + v2[2]*v2[2])
	v2norm := [3]float64{v2[0] / v2mag, v2[1] / v2mag, v2[2] / v2mag}

	res := v1norm[0]*v2norm[0] + v1norm[1]*v2norm[1] + v1norm[2]*v2norm[2]
	angle := math.Pi - math.Acos(res)
	return angle
}

func distance(amd [3]float64, atm [3]float64) float64 {
	sqX := math.Pow(amd[0]-atm[0], 2)
	sqY := math.Pow(amd[1]-atm[1], 2)
	sqZ := math.Pow(amd[2]-atm[2], 2)
	return math.Sqrt(sqX + sqY + sqZ)
}

// Contacts are calculated with N atom (from NH) as sphere center
func inContact(amd pdb.Amide, atm pdb.Atom, dist float64) bool {
	return distance(amd.Nxyz, atm.XYZ) < dist
}

// Tests if there are a HB
func doHB(amd pdb.Amide, atm pdb.Atom, dist float64, angle float64) bool {
	// test if atom is an acceptor
	if isAcceptor(atm) {
		// test if distance H to acceptor is between 1.5 and dist (user defined)
		if distance(amd.Hxyz, atm.XYZ) > 1.5 && distance(amd.Hxyz, atm.XYZ) < dist {
			// fmt.Printf("%s %s %s %s %s %f %f ", amd.ResName, amd.ResNumber, atm.ResName, atm.ResNumber, atm.Name, distance(amd.Hxyz, atm.XYZ), hbangle(amd.Nxyz, amd.Hxyz, atm.XYZ))
			// test if angle is smaller than angle (user defined)
			if hbangle(amd.Nxyz, amd.Hxyz, atm.XYZ) < angle {
				// fmt.Printf("true\n")
				return true
				// } else {
				// 	fmt.Printf("false\n")
			}
		}
	}
	// if some tests above fails
	return false
}

func Run(amd []pdb.Amide, atm []pdb.Atom, validAtoms string, maxDist float64, hbDist float64, hbAngle float64) {
	for i := 0; i < len(amd); i++ {
		for j := 0; j < len(atm); j++ {

			// skip if is the same atom
			if amd[i].PDBNumber != atm[j].PDBNumber {
				// skip if atomtype is not valid (user input)
				if atm[j].Type != "" && strings.Contains(validAtoms, atm[j].Type) {
					if inContact(amd[i], atm[j], maxDist) {
						amd[i].NumContacts++
					}
				}
				if atm[j].Type == "O" {
					if doHB(amd[i], atm[j], hbDist, hbAngle) {
						amd[i].NumOhb++
					}
				}
				if atm[j].Type == "N" {
					if doHB(amd[i], atm[j], hbDist, hbAngle) {
						amd[i].NumNhb++
					}
				}
			}

		}
		if amd[i].NumOhb > 0 || amd[i].NumNhb > 0 {
			amd[i].DoHB = true
		}
	}
}
