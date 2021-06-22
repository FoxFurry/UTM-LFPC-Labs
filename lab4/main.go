package main

import (
	"fmt"
)

type ll1 struct {
	dfaMap map[string][]string
	vn     []byte
	vt     []byte
}

func main() {
	var27 := ll1{
		map[string][]string{
			"S": {"EdF"},
			"E": {"ebD"},
			"F": {"FaL"},
			"D": {"FbE"},
			"L": {"aL", "a"},
		},
		[]byte{'S', 'F', 'L', 'E'},
		[]byte{'a', 'b', 'c', 'd', 'e'},
	}

	var27.generateLL1(true)
}

func (ll *ll1) generateLL1(isDebug bool) {
	//var matrixSize int = len(ll.vn) + len(ll.vt)
	//
	//result := make([][]int, matrixSize)
	//for i := 0; i < matrixSize; i++ {
	//	result[i] = make([]int, matrixSize)
	//}

	if isDebug {
		ll.print()
	}
	ll.eliminateLeftRecursion()
	if isDebug {
		ll.print()
	}
	ll.generateFirst()
	if isDebug {
		ll.print()
	}
	ll.generateFollow()
	if isDebug {
		ll.print()
	}
	ll.constructParsingTable()
	if isDebug {
		ll.print()
	}

	//return result
}

func (ll *ll1) eliminateLeftRecursion() {
	for _, nonterm := range ll.vn {
		if ll.dfaMap[string(nonterm)] == nil {

		}
	}
}

func (ll *ll1) generateFirst() {

}

func (ll *ll1) generateFollow(){

}

func (ll *ll1) constructParsingTable(){

}

func (ll *ll1) print(){
	fmt.Printf("DFA\n{\n")
	for key, val := range ll.dfaMap {
		fmt.Printf("\t%v -> ", key)
		for _, set := range val {
			fmt.Printf("%v\t", set)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("}\nVN\n{ ")
	for _, val := range ll.vn {
		fmt.Printf("%c ", val)
	}
	fmt.Printf("}\nVT\n{ ")
	for _, val := range ll.vt {
		fmt.Printf("%c ", val)
	}
	fmt.Printf("}")
}
