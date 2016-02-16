package main

import (
	"fmt"

	"github.com/jgcarvalho/toniHB/analysis"

	"github.com/jgcarvalho/toniHB/pdb"
)

func main() {
	// parse parameters
	amd, atm, err := pdb.LoadFile("./test/lysozime/em.pdb")
	// calculate
	analysis.Run(amd, atm, "CONPH", 3.6, 3.6, 0.5)

	fmt.Println("ATOMS")
	fmt.Println(atm)
	fmt.Println("AMIDES")
	fmt.Println(amd)

	fmt.Println(err)
}
