package main

import (
	"fmt"
	"strconv"
)

type transition struct {
	source int
	destination int
	input byte
}

func main(){
	var transitionArray []transition
	var deltaArray []byte

	var src, dst = 0,0
	var inp byte = 0
	fmt.Println("Enter transitions in this form <source input destination>\nEx: 0 a 1\nEnter 0 0 0 to end input")
	for {
		_, _ = fmt.Scanf("%d %c %d",&src, &inp, &dst)

		if src==0 && inp == '0' && dst == 0 {
			break
		}
		if !itemExists(deltaArray, inp) {
			deltaArray = append(deltaArray, inp)
		}

		transitionArray = append(transitionArray, transition{src, dst, inp})
	}

	dfaArray := map[string][]string{}

	// GENERATE INITIAL PERMUTATIONS WITH SINGLE NODE
	var generateElement bool
	var generationArray []string
	for _, elem := range transitionArray {
		generateElement = false
		elemStr := strconv.Itoa(elem.source)
		if dfaArray[elemStr] == nil {
			dfaArray[elemStr] = make([]string, len(deltaArray))
		}
		if dfaArray[elemStr][elem.input - 97] != "" {
			generateElement = true
			dfaArray[elemStr][elem.input - 97] += ","
		}
		dfaArray[elemStr][elem.input - 97] += strconv.Itoa(elem.destination)

		if generateElement {
			generationArray = append(generationArray, dfaArray[elemStr][elem.input - 97])
		}
	}

	recursiveGenerate(dfaArray, generationArray)

	fmt.Printf("N\t|")
	for _,dlt := range deltaArray{
		fmt.Printf("\t%9c\t|", dlt)
	}
	fmt.Printf("\n")
	for key, elem := range dfaArray {
		fmt.Printf("q%s\t|", key)
		for _, entry := range elem{
			fmt.Printf("\t%9s\t|", setToString(entry))
		}
		fmt.Printf("\n")
	}
}

func itemExists(s []byte, str byte) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func recursiveGenerate(dfaMap map[string][]string, generationArray []string){
	if len(generationArray) == 0 {
		return
	}
	var newGenerationArray []string
	for _, elem := range generationArray {
		if dfaMap[elem] != nil{
			continue
		}
		for _, point := range elem {
			if dfaMap[elem] == nil {
				dfaMap[elem] = make([]string, len(dfaMap[string(point)]))
			}
			for key, delta := range dfaMap[string(byte(point))] {
				if len(dfaMap[elem][key]) == 0 {
					dfaMap[elem][key] += delta
				}else {
					dfaMap[elem][key] += "," + delta
				}
				newGenerationArray = append(newGenerationArray, dfaMap[elem][key])
			}
			continue
		}
	}
	recursiveGenerate(dfaMap, newGenerationArray)
}

func setToString(set string) string {
	var output string

	for idx, elem := range set {
		if idx % 2 == 0 {
			output += "q"
		}
		output += string(elem)
	}
	if len(output) == 0 {
		output = "-"
	}
	return output
}