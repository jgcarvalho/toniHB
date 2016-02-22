package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jgcarvalho/toniHB/analysis"

	"github.com/jgcarvalho/toniHB/pdb"
)

// Analyse all pdb files
func Analyse(dir, validAtoms string, radius, dist, angle float64, output string) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}

	pdbFiles, err := filepath.Glob(absDir + "/*")
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	if len(pdbFiles) == 0 {
		log.Println("ERROR: PDB files not found in path", absDir)
	}
	fmt.Printf("%d PDB files founded at %s\n", len(pdbFiles), absDir)

	allPDBamides := make([][]pdb.Amide, len(pdbFiles))
	allPDBatoms := make([][]pdb.Atom, len(pdbFiles))
	for i := 0; i < len(pdbFiles); i++ {
		allPDBamides[i], allPDBatoms[i], err = pdb.LoadFile(pdbFiles[i])
		analysis.Run(allPDBamides[i], allPDBatoms[i], validAtoms, radius, dist, angle)
	}
	// write result to file
	Report(allPDBamides, output)
	fmt.Println("DONE")
}

// Write result to output file
func Report(allPDBamides [][]pdb.Amide, output string) {
	f, err := os.Create(output)
	if err != nil {
		log.Println("ERROR:", err)
		return
	}
	defer f.Close()

	nPDB := len(allPDBamides)
	nAmides := len(allPDBamides[0])
	// print header
	fmt.Fprint(f, "Number, Residue")
	for i := 0; i < nPDB; i++ {
		fmt.Fprintf(f, ", Contacts_%d, O_%d, N_%d, H-bond_%d", i+1, i+1, i+1, i+1)
	}
	fmt.Fprint(f, "\n")
	// print data
	for i := 0; i < nAmides; i++ {
		for j := 0; j < nPDB; j++ {
			if j == 0 {
				fmt.Fprintf(f, "%d, \"%s %s\"", allPDBamides[j][i].Number, allPDBamides[j][i].ResName, allPDBamides[j][i].ResNumber)
			}
			fmt.Fprintf(f, ", %d, %d, %d, %t", allPDBamides[j][i].NumContacts, allPDBamides[j][i].NumOhb, allPDBamides[j][i].NumNhb, allPDBamides[j][i].DoHB)
			if j == (nPDB - 1) {
				fmt.Fprintf(f, "\n")
			}
		}
	}
}

func main() {
	// parse parameters

	validAtoms := flag.String("type", "", "Valid atom types. Ex. -type CNOPH or -type \"C N O P H\"")
	radius := flag.Float64("radius", 0.0, "Radius of the sphere to count atoms around N atom")
	dist := flag.Float64("dist", 0.0, "Maximum distance between H and acceptor")
	angle := flag.Float64("angle", -1.0, "Maximum angle between vectors H->N and H->acceptor (radians)")
	dir := flag.String("dir", "", "PDB files directory (the files should be named like somename_0001.pdb ... somename_9999.pdb)")
	output := flag.String("o", "", "Output file")
	flag.Parse()

	if *dir == "" {
		fmt.Println("Please, set the pdb files path using the flag -dir")
		return
	}
	if *validAtoms == "" {
		fmt.Println("Please, set the valid atom types using the flag -type")
		return
	}
	if *radius == 0.0 {
		fmt.Println("Please, set the radius of the sphere using the flag -radius")
		return
	}
	if *dist == 0.0 {
		fmt.Println("Please, set the maximum HB distance using the flag -dist")
		return
	}
	if *angle == -1.0 {
		fmt.Println("Please, set the maximum HB angle using the flag -angle")
		return
	}
	if *output == "" {
		fmt.Println("Please, set the output file using the flag -o")
		return
	}

	fmt.Println("PDB path:", *dir)
	fmt.Println("Valid atoms:", *validAtoms)
	fmt.Println("Radius of the sphere to count atoms around N:", *radius)
	fmt.Println("Maximum distance between H and accepto:", *dist)
	fmt.Println("Maximum angle between vectors H->N and H->acceptor (radians):", *angle)

	fmt.Println("\nRunning...")

	Analyse(*dir, *validAtoms, *radius, *dist, *angle, *output)
}
