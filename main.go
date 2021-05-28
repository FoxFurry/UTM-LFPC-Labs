package main

import "path/filepath"



func main(){
	dfaMap := map[string][]string{		// Variant 27
		"S": {"EdF"},
		"E": {"ebD"},
		"F": {"FaL"},
		"D": {"FbE"},
		"L": {"aL", "a"},
	}
	Vn := []byte{'S','F','L','E'}
	Vt := []byte{'a','b','c','d','e'}

	presedenceMatrix := generatePredenceProto(Vn, Vt)
}

func generatePredenceProto(Vn []byte, Vt []byte) [][]int{
	var matrixSize int = len(Vn) + len(Vt)
	result := make([][]int, matrixSize)
	for i := 0; i < matrixSize; i++ {
		result[i] = make([]int, matrixSize)
	}

	for _, nonTerminal := range Vn {
		symbolToIndex(nonTerminal)
	}

	for _, terminal := range Vt {
		symbolToIndex(terminal)
	}

	return result
}

func symbolToIndex(byte){

}
