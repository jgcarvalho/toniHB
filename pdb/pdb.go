package pdb

import (
	"bufio"
	"errors"
	"fmt"
	"os"
)

type Amide struct {
	Number      int
	ResName     string
	ResNumber   string // residue number may contain insertion code ex. 101A
	NumContacts int
	NumOhb      int
	NumNhb      int
	DoHB        bool
	// used internally
	PDBNumber int
	Nxyz      [3]float64
	Hxyz      [3]float64
}

type Atom struct {
	PDBNumber int
	Name      string
	ResName   string
	ResNumber string // residue number may contain insertion code ex. 101A
	Chain     string
	XYZ       [3]float64
	Type      string
}

func LoadFile(pdbfile string) ([]Amide, []Atom, error) {
	var amidesTmp []Amide
	var amides []Amide
	var atoms []Atom

	if pdbfile == "" {
		return nil, nil, errors.New("You have to enter a valid PDB file name")
	}

	f, err := os.Open(pdbfile)
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	defer f.Close()

	data := bufio.NewReader(f)
	var (
		line    []byte
		strLine string
	)

	countAmides := 0
	for {
		line, _, err = data.ReadLine()
		if err != nil {
			break
		}

		// TODO need check if lipids are HETATM
		if string(line[:6]) == "ATOM  " || string(line[:6]) == "HETATM" {
			strLine = string(line)
			var atm Atom

			fmt.Sscanf(strLine[6:11], "%d", &atm.PDBNumber)
			fmt.Sscanf(strLine[12:16], "%s", &atm.Name)
			fmt.Sscanf(strLine[17:20], "%s", &atm.ResName)
			fmt.Sscanf(strLine[21:22], "%s", &atm.Chain)
			fmt.Sscanf(strLine[22:27], "%s", &atm.ResNumber)

			fmt.Sscanf(strLine[30:38], "%f", &atm.XYZ[0])
			fmt.Sscanf(strLine[38:46], "%f", &atm.XYZ[1])
			fmt.Sscanf(strLine[46:54], "%f", &atm.XYZ[2])
			fmt.Sscanf(strLine[76:78], "%s", &atm.Type)

			atoms = append(atoms, atm)

			// get N atom from protein backbone
			if string(line[:6]) == "ATOM  " && atm.Name == "N" {
				countAmides++
				var amd Amide
				amd.Number = countAmides
				amd.PDBNumber = atm.PDBNumber
				amd.ResName = atm.ResName
				amd.ResNumber = atm.ResNumber
				amd.Nxyz = atm.XYZ
				amidesTmp = append(amidesTmp, amd)
			}

			// get H atom from amides protein backbone
			if string(line[:6]) == "ATOM  " && atm.Name == "H" {
				// check if H atom is from the same residue that last amide add
				// WARNING: this can be invalid if H atom is before N atom in the PDB file
				// but gromacs apparently put the N before H
				if amidesTmp[len(amidesTmp)-1].ResNumber == atm.ResNumber {
					amidesTmp[len(amidesTmp)-1].Hxyz = atm.XYZ
				}
			}
		}
	}

	// check for NH3+ from N-terminal
	// NH3+ won't have H (they have H1, H2 and H3)
	// so N-terminal residues are EXCLUDED from analysis (Do they need be included?)
	for _, x := range amidesTmp {
		if x.Hxyz != [3]float64{0.0, 0.0, 0.0} {
			amides = append(amides, x)
		}
	}
	return amides, atoms, nil
}
