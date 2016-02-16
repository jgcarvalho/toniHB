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
	NumO        int
	NumN        int
	DoHB        bool
	// used internally
	PDBNumber int
	XYZ       [3]float64
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

			if atm.Name == "N" {
				countAmides++
				var amd Amide
				amd.Number = countAmides
				amd.PDBNumber = atm.PDBNumber
				amd.ResName = atm.ResName
				amd.ResNumber = atm.ResNumber
				amd.XYZ = atm.XYZ
				amides = append(amides, amd)
			}
		}
	}
	return amides, atoms, nil
}
