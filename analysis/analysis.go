package analysis

import (
	"math"
	"strings"

	"github.com/jgcarvalho/toniHB/pdb"
)

// THIS FUNCTION IS IMPORTANT
// ACCEPTORS MUST BE DEFINED HERE

func isAcceptor(atm pdb.Atom) bool {
	for i := 0; i < len(Acceptors); i++ {
		if Acceptors[i].ResName == atm.ResName && Acceptors[i].AtmName == atm.Name {
			return true
		}
	}
	return false
	// switch atm.Name {
	// // CHARMM27 and OPLS/AA protein HB acceptors
	// case "O", "OD1", "OD2", "OE1", "OE2", "OG1", "OH":
	// 	// ASPP and GLUP are protonated so they don't have acceptors
	// 	if atm.ResName == "ASPP" || atm.ResName == "GLUP" {
	// 		return false
	// 	}
	// 	return true
	// case "ND1", "NE2":
	// 	// HSD, HSE and HSP are used to indicate which atom is protonated
	// 	if atm.ResName == "HSD" && atm.Name == "NE2" {
	// 		return true
	// 	} else if atm.ResName == "HSE" && atm.Name == "ND1" {
	// 		return true
	// 	} else if atm.ResName == "HSP" {
	// 		return false
	// 		// HIS don't indicate which atom is protonated
	// 	} else if atm.ResName == "HIS" {
	// 		return true
	// 	} else {
	// 		return false
	// 	}
	// 	// nucleic acid acceptors
	// case "O6", "N3", "N7", "O1P", "O2P", "O2'", "O3'", "O4'", "O5'", "N1", "O2", "O4":
	// 	if atm.ResName == "DA" || atm.ResName == "DC" || atm.ResName == "DG" || atm.ResName == "DT" || atm.ResName == "A" || atm.ResName == "C" || atm.ResName == "G" || atm.ResName == "T" {
	// 		return true
	// 		// Palmitate has a O2 atom acceptor
	// 	} else if atm.ResName == "PALM" && atm.Name == "O2" {
	// 		return true
	// 	} else {
	// 		return false
	// 	}
	// 	// lipids
	// case "OS1", "OS2", "OS3", "OS4", "OH2":
	// 	return true
	// case "O1":
	// 	if atm.ResName == "PALM" {
	// 		return true
	// 	} else {
	// 		return false
	// 	}
	// default:
	// 	return false
	// }
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
			// test if angle is smaller than angle (user defined)
			if hbangle(amd.Nxyz, amd.Hxyz, atm.XYZ) < angle {
				return true
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
