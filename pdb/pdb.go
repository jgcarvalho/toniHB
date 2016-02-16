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
	ResNumber   int
	NumContacts int
	NumO        int
	NumN        int
	DoHB        bool
	// used internally
	PDBNumber int
	Coord     [3]float64
}

type Atom struct {
	PDBNumber int
	AtomName  string
	ResName   string
	ResNumber int
	Coord     [3]float64
	AtomType  string
}

func LoadFile(pdbfile string) ([]Amide, []Atom, error) {
	// var p Protein
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

		chain     string
		curNumber int
		resNumber int
		icode     string
		curIcode  string
		nres      int
	)

	nres = 1
	for {
		line, _, err = data.ReadLine()
		if err != nil {
			break
		}

		if string(line[:4]) == "ATOM" {
			strLine = string(line)

			fmt.Sscanf(strLine[22:26], "%d", &resNumber)
			fmt.Sscanf(strLine[26:27], "%s", &icode)
			if resNumber != curNumber || chain != strLine[21:22] || icode != curIcode {

				curNumber = resNumber
				curIcode = icode

				if chain != strLine[21:22] {
					chain = strLine[21:22]
					p.Chains = append(p.Chains, chain)
				}

				var r Residue
				r.N = nres
				r.Npdb = resNumber
				r.ICode = icode
				r.Chain = chain

				fmt.Sscanf(strLine[17:20], "%s", &r.Code)

				var a Atom
				fmt.Sscanf(strLine[6:11], "%d", &a.N)
				fmt.Sscanf(strLine[12:16], "%s", &a.Name)
				fmt.Sscanf(strLine[76:78], "%s", &a.Type)
				fmt.Sscanf(strLine[30:38], "%f", &a.XYZ[0])
				fmt.Sscanf(strLine[38:46], "%f", &a.XYZ[1])
				fmt.Sscanf(strLine[46:54], "%f", &a.XYZ[2])
				fmt.Sscanf(strLine[54:60], "%f", &a.Occ)
				fmt.Sscanf(strLine[60:66], "%f", &a.Bfactor)
				fmt.Sscanf(strLine[16:17], "%s", &a.AltLoc)
				if len(strLine) > 79 {
					fmt.Sscanf(strLine[78:80], "%f", &a.Charge)
				}

				// r.Atoms = append(r.Atoms, a)
				// p.Residues = append(p.Residues, r)
				nres += 1
			} else {
				var a Atom
				fmt.Sscanf(strLine[6:11], "%d", &a.N)
				fmt.Sscanf(strLine[12:16], "%s", &a.Name)
				fmt.Sscanf(strLine[76:78], "%s", &a.Type)
				fmt.Sscanf(strLine[30:38], "%f", &a.XYZ[0])
				fmt.Sscanf(strLine[38:46], "%f", &a.XYZ[1])
				fmt.Sscanf(strLine[46:54], "%f", &a.XYZ[2])
				fmt.Sscanf(strLine[54:60], "%f", &a.Occ)
				fmt.Sscanf(strLine[60:66], "%f", &a.Bfactor)
				fmt.Sscanf(strLine[16:17], "%s", &a.AltLoc)
				if len(strLine) > 79 {
					fmt.Sscanf(strLine[78:80], "%f", &a.Charge)
				}

				// p.Residues[len(p.Residues)-1].Atoms = append(p.Residues[len(p.Residues)-1].Atoms, a)
			}
		}
	}
	// return &p, nil
	return amides, atoms, nil
}
