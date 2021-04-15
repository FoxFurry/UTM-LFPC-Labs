package main

import (
	"bufio"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	fileData, fileError := os.Open("input.txt")
	if fileError != nil { panic(fileError) }

	dfaArray := map[string][]string{}

	print("Reading productions\n")

	readFile(fileData, &dfaArray)
	printDfa(dfaArray)

	print("\nEliminating epsilon productions\n\n")

	eliminateEps(&dfaArray)
	printDfa(dfaArray)

	print("\nEliminating renaming productions\n\n")

	eliminateRenaming(&dfaArray)
	printDfa(dfaArray)

	print("\nEliminating nonproductive symbols\n\n")

	eliminateNonProductive(&dfaArray)
	printDfa(dfaArray)

	print("\nEliminating inaccessible symbols\n\n")

	eliminateInaccessible(&dfaArray)
	printDfa(dfaArray)

}

/*
Reads production set from file
If file format is bad it will panic

Void
 */
func readFile(pInputFile *os.File, pDfaMap *map[string][]string) {

	productionRegex := regexp.MustCompile("^(\\s*)[Pp](\\s*)=(\\s*){(\\s*)$")	// 'P={' or 'p  =    {'
	setRegex := regexp.MustCompile("^[a-zA-Z](\\s*)->(\\s*)[a-zA-Z0]+(,?)$")	// 'A->ABC' or 'A->ABC,'
	setSplitRegex := regexp.MustCompile("-\\>|,")								// Splits 'A->ABC,' to ['A', 'ABC']

	lineRead := bufio.NewScanner(pInputFile)

	productionRead := false
	lineCounter := 0

	for lineRead.Scan() {
		lineCounter++

		strippedLine := strings.TrimSpace(lineRead.Text())		// Trim whitespaces

		if strippedLine == "" { continue }

		if productionRead { // If production sets found -> read them
			switch {
				case strippedLine[0] == '}':					// End of reading file
					break

				case setRegex.MatchString(strippedLine):		// Found a production set
					setElements := setSplitRegex.Split(strippedLine, -1)
					lhs:= setElements[0]						// Left hand side e.g 'S'
					rhs:= setElements[1]						// Right hand side e.g 'ABC'

						if (*pDfaMap)[lhs] == nil { 					// If we never had such key
							(*pDfaMap)[lhs] = []string{rhs} 			// Create it
						} else {										// Else
						(*pDfaMap)[lhs] = append((*pDfaMap)[lhs], rhs)	// Append a new set for this key
					}
					//println(lhs + "\t" + rhs)

				default:										// Every other case - error
					panic("Unexpected character at line " + strconv.Itoa(lineCounter))

			} //switch
		} else { // Search for production sets
			switch {
				case strippedLine[0] == '#':					// Skip comment
					continue

				case productionRegex.MatchString(strippedLine):	// Beginning of production sets
					println("Production sets found:")
					productionRead = true

				default:										// Every other case - error
					panic("Unexpected character at line " + strconv.Itoa(lineCounter))

			} // switch
		} //else
	} //for
} //readFile


/*
Deletes epsilon productions from map and generates inverse map for this set
Return value shows if epsilon productions were found in map

False -> No epsilon production were found and map is unchanged
True -> Epsilon production were eliminated and map was changed
 */
func eliminateEps(pDfaMap *map[string][]string){
	for key, value := range *pDfaMap {
		for idx, set := range value {
			if set == "0" {		// If it is an epsilon set - eliminate it and generate inverse
				(*pDfaMap)[key] = fastRemove((*pDfaMap)[key], idx)
				//idx--
				addInverseProductions(pDfaMap, rune(key[0]))
			}
		} //forMapValue
	} //forMap
} //eliminateEps

/*
Я конечно постарался сделать это адекватно и читабельно, но получилось как всегда. Time complexity в данном случае
очень странное, что-то между O(NlogN) и O(N^2). Мне кажется, что это чуть быстрее, чем рекурсивно, потому что здесь
генерируются комбинации сразу для всех productions, а в рекурсивном ты вызываешь рекурсию для каждого.
Memory complexity 100% больше, чем у рекурсивной, но всего на чуть-чуть. Я вообще в шоке, что оно работает
 */
func addInverseProductions(pDfaMap *map[string][]string, searchElem rune){
	println("Epsilon key is: " + string(searchElem))
	for key, value := range *pDfaMap {
		for _, set := range value {					// Traversing each production set
			runeCounter := 0						// Combinations counter
			var keyIndexes []int

			for idx, element := range set{
				if element == searchElem {
					runeCounter++ 					// Count number of combinations and add every index to array
					keyIndexes = append(keyIndexes, idx)
				}
			} //forMapValueElement

			if runeCounter == 0 || keyIndexes == nil { continue }

			var combinationCounter = int(math.Pow(2, float64(runeCounter)) - 1)
			println("Generating " + strconv.Itoa(combinationCounter) + " combination(s) for: " + set)

			for bcd := 0; bcd < combinationCounter; bcd++ {			// Generate combinations from i to 2^N-1
				var setCopy = set
				for bitNumber := 0; bitNumber < runeCounter; bitNumber++ {	// Generate a string for this combination
					if (bcd & (1 << bitNumber)) == 0 {						// Check if such bit is set to 0
						var index = keyIndexes[bitNumber]                   // Find index of this key in string
						setCopy = setCopy[:index] + "\t" + setCopy[index+1:]// Basically I use whitespace character as
																			// a marker for further removal
					}
				} //forBCDBit
				setCopy = strings.ReplaceAll(setCopy, "\t", "")	// Remove all marked characters
				(*pDfaMap)[key] = append((*pDfaMap)[key], setCopy)			// Finally, push new production to our map
			} //forBCD
		} //forMapValue
	} //forMap
} //addInverseProductions


/*
Deletes renamings
Return value shows if renaming
 */
func eliminateRenaming(pDfaMap *map[string][]string){
	for key, value := range *pDfaMap {
		for idx, set := range value {
			if len(set) == 1 && set[0] >= 'A' && set[0] <= 'Z' {	// If we have a renaming
				(*pDfaMap)[key] = fastRemove((*pDfaMap)[key], idx)
				generateRenamingProductions(pDfaMap, key, set)
			}
		} //forMapValue
	} //forMap
} //eliminateRenaming


func generateRenamingProductions(pDfaMap *map[string][]string, targetKey string, renameKey string){
	println("Eliminating renaming " + targetKey + "->" + renameKey)

	for _, set := range (*pDfaMap)[renameKey] {
		_, isPresent := find((*pDfaMap)[targetKey], set)
		if !isPresent {
			(*pDfaMap)[targetKey] = append((*pDfaMap)[targetKey], set)
		}
	}
} //generateRenamingProductions

func eliminateNonProductive(pDfaMap *map[string][]string){
	isVnProductiveMap := map[string]bool{}
	var hasChanged = false
	for {
		hasChanged = false
		for key, value := range *pDfaMap {
			if isVnProductiveMap[key] { continue }
			for _, set := range value {
				var setStatus = true
				for _, element := range set {
					if element >= 'A' && element <= 'Z' {
						if isVnProductiveMap[string(element)] == false {
							setStatus = false
							break
						}
					}
				} //forMapValueElement
				if setStatus {
					isVnProductiveMap[key] = true
					hasChanged = true
					break
				}
			} //forMapValue
		} //forMap
		if !hasChanged { break }
	} //for

	print("Productive sets = { ")
	for key := range isVnProductiveMap {
		print(key + " ")
	}
	println("}\n")

	for key, value := range *pDfaMap {
		valueRange := len(value)
		for idx := 0; idx < valueRange; idx++ {
			set := value[idx]
			for _, element := range set {
				if element >= 'A' && element <= 'Z' {
					if isVnProductiveMap[string(element)] == false {
						println("Deleting set: " + key + "->" + set)
						(*pDfaMap)[key] = fastRemove((*pDfaMap)[key], idx)
						if len((*pDfaMap)[key]) == 0 {
							delete(*pDfaMap, key)
						}
						idx--
						valueRange--
						break
					}
				}
			} //forMapValueElement
		} //forMapValue
	} //forMap

} //eliminateNonProductive

func eliminateInaccessible(pDfaMap *map[string][]string){
	isVnAccessibleMap := map[string]bool{}

	for _, value := range *pDfaMap {
		for _, set := range value {
			for _, element := range set {
				if element >= 'A' && element <= 'Z' {
					isVnAccessibleMap[string(element)] = true
				}
			} //forMapElement
		} //forMapValue
	} //forMap

	for key, _ := range *pDfaMap {
		if !isVnAccessibleMap[key] {
			println(key + " is inaccessible. Deleting it")
			delete(*pDfaMap, key)
		}
	}

} //eliminateInaccessible

// -----------------------------------
// UTILS

func printDfa(_dfaMap map[string][]string){
	println()
	for key, value := range _dfaMap {
		println(key + " -> ")
		for _, set := range value {
			println("     " + set)
		}
	}
} //printDfa

func fastRemove(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
} //fastRemove

func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
} //find

func isLowerString(val string) bool {
	return val == strings.ToLower(val)
} //isLowerString

func isLowerChar(val rune) bool {
	return val >= 'a' && val <= 'z'
} //isLowerChar