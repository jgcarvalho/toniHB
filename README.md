# toniHB

## Installation

After install [Go](https://golang.org/) and set a [Workspace (GOPATH environment variable)](https://golang.org/doc/code.html#Workspaces). Do `go get github.com/jgcarvalho/toniHB` to install or `go get -u github.com/jgcarvalho/toniHB` to update.

## Usage

Parameters:
- -dir: PDB files directory (the files should be named like somename_0001.pdb ... somename_9999.pdb)
- -type: Valid atom types
- -radius: Radius of the sphere to count atoms around N atom
- -dist: Maximum distance between H and acceptor
- -angle Maximum angle between vectors H->N and H->acceptor (radians)
- -o: Output file

Examples:

`toniHB -dir ./pdbfiles/ -radius 3.6 -dist 2.8 -angle 1.0 -type COHPN -o output.csv`

or

`toniHB -dir ./pdbfiles/ -radius 3.6 -dist 2.8 -angle 1.0 -type "C O H P N" -o output.csv`

## Specifications

Contacts are atoms inside a sphere centered at N backbone atom.

HB distance is calculated between H and acceptor.

HB angle is calculated as Pi minus the angle between vector H->N and vector H->acceptor. So, when the atoms are aligned N-H---O the angle is ZERO.    

## Notes

- PDB files in the PBD path **must** have the extension **.pdb**
- The PDB files are sorted alphabetically so you should use
```
 protein_001.pdb
 protein_002.pdb
 ...
 protein_999.pdb
 ```
 and **NOT**
 ```
  protein_1.pdb
  protein_2.pdb
  ...
  protein_9.pdb
  protein_10.pdb
  ...
  protein_999.pdb
  ```
