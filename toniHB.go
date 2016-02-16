package main

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

// Test if two atoms are doing a HB
func isHB() {

}

func main() {
	// parse parameters

	// calculate
}
