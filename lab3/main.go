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

	dfaAfterRead := map[string][]string{}
	for k,v := range dfaArray { dfaAfterRead[k] = v}

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

	print("\nObtaining CNF\n\n")

	obtainCNF(&dfaArray)
	printDfa(dfaArray)

	print("\n---------------------------\nBefore CNF\n\n")
	printDfa(dfaAfterRead)

	print("\n---------------------------\nAfter CNF\n\n")
	printDfa(dfaArray)


}

/*
Reads production set from file
If file format is bad it will panic

Parameters:
	pInputFile - pointer to file from which production sets will be read
	pDfaMap - pointer of map string->[]string which will be populated with production sets
Return:
	void
 */
func readFile(pInputFile *os.File, pDfaMap *map[string][]string) {
	productionRegex := regexp.MustCompile("^(\\s*)[Pp](\\s*)=(\\s*){(\\s*)$")	// 'P={' or 'p  =    {'
	setRegex := regexp.MustCompile("^[a-zA-Z](\\s*)->(\\s*)[a-zA-Z0]+(,?)$")	// 'A->ABC' or 'A->ABC,'
	setSplitRegex := regexp.MustCompile("->|,")								// Splits 'A->ABC,' to ['A', 'ABC']

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

Parameters:
	pDfaMap - pointer of map string->[]string from which epsilon productions will be eliminated
Return:
	void
 */
func eliminateEps(pDfaMap *map[string][]string){
	for key, value := range *pDfaMap {
		for idx, set := range value {
			if set == "0" {		// If it is an epsilon set - eliminate it and generate inverse
				(*pDfaMap)[key] = fastRemove((*pDfaMap)[key], idx)
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

Parameters:
	pDfaMap - pointer of map string->[]string from which renaming productions will be eliminated
Return:
	void
 */
func eliminateRenaming(pDfaMap *map[string][]string){
	for key, value := range *pDfaMap {
		for idx, set := range value {
			if len(set) == 1 && set[0] >= 'A' && set[0] <= 'Z' {	// If we have a renaming
				(*pDfaMap)[key] = fastRemove((*pDfaMap)[key], idx)	// Delete the renaming set
				generateRenamingProductions(pDfaMap, key, set)		// Generate all necessary productions after renaming
			}
		} //forMapValue
	} //forMap
} //eliminateRenaming

/*
Generates productions after deleting renaming

If "S->B" is our production, 'S' is targetKey and 'B' is renameKey
They are actually characters, not strings, but my map keys are strings so I left them as strings

Parameters:
	pDfaMap - pointer of map string->[]string for which will be generated productions
	targetKey - key for which will be generated
	renameKey - key from which will be generated
Return:
	void
 */
func generateRenamingProductions(pDfaMap *map[string][]string, targetKey string, renameKey string){
	println("Eliminating renaming " + targetKey + "->" + renameKey)

	for _, set := range (*pDfaMap)[renameKey] {				// Go through every renameKey set
		_, isPresent := find((*pDfaMap)[targetKey], set)	// Check is such set already exists in targetKey
		if !isPresent {										// If not -> add it
			(*pDfaMap)[targetKey] = append((*pDfaMap)[targetKey], set)
		}
	} //forSet
} //generateRenamingProductions


/*
Eliminates non productive sets from map

Basically it searched for productive keys and after search - eliminates all other keys

Parameters:
	pDfaMap - pointer of map string->[]string from which non productive sets will be eliminated
Return:
	void
 */
func eliminateNonProductive(pDfaMap *map[string][]string){
	isVnProductiveMap := map[string]bool{}					// This map will contain productive keys
															// I know it can be changed to just array, but maps have
															// O(nlogn) search complexity and I search a lot
	var hasChanged = false

	for {													// Here we start searching productive keys
		hasChanged = false									// We are iterating infinitely and this flag will stop us
		for key, value := range *pDfaMap {					// Take every key and set for this key
			if isVnProductiveMap[key] { continue }			// Is this key is already productive - skip
			for _, set := range value {						// Go through every set
				var setStatus = true						// This flag is set to false if this set is non productive
				for _, element := range set {				// Go through set elements
					if element >= 'A' && element <= 'Z' {	// If we see non terminal symbol
						if isVnProductiveMap[string(element)] == false {	// Check is non terminal is productive
							setStatus = false				// If not - all set is non productive
							break
						}									// Basically I check if every symbol in set is productive
					}
				} //forMapValueElement
				if setStatus {								// If it is productive
					isVnProductiveMap[key] = true			// Set key as true
					hasChanged = true						// We have changed out productive map, so we should iterate
					break									// one more time
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
		for idx := 0; idx < valueRange; idx++ {		// I have to use normal loop because of "going back" in loop
			set := value[idx]
			for _, element := range set {
				if element >= 'A' && element <= 'Z' {	// If we found a non terminal symbol
					if isVnProductiveMap[string(element)] == false {	// Check if it is not productive
						println("Deleting set: " + key + "->" + set)	// If it is not - delete it
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

/*
Eliminates inaccessible sets from map

Parameters:
	pDfaMap - pointer of map string->[]string from which inaccessible sets will be eliminated
Return:
	void
 */
func eliminateInaccessible(pDfaMap *map[string][]string){
	isVnAccessibleMap := map[string]bool{}

	for _, value := range *pDfaMap {
		for _, set := range value {
			for _, element := range set {
				if element >= 'A' && element <= 'Z' {
					isVnAccessibleMap[string(element)] = true	// Add to map every non terminal symbol we see in RHS
				}
			} //forMapValueElement
		} //forMapValue
	} //forMap

	for key := range *pDfaMap {
		if !isVnAccessibleMap[key] {	// If such non terminal is not in map -> it is inaccessible
			println(key + " is inaccessible. Deleting it")
			delete(*pDfaMap, key)
		}
	} //forMap

} //eliminateInaccessible

/*
Finally - transforms our grammar to CNF

Parameters:
	pDfaMap - pointer of map string->[]string which will obtain CNF after this function
Return:
	void
*/
func obtainCNF(pDfaMap *map[string][]string) {
	chosmkyMap := map[string]string{}
	chomskyRegex := regexp.MustCompile("^([A-Z]{2})$|^([a-z])$") 	// This regex checks if production is CNF

	lastX := 1
	lastY := 1

	for key, value := range *pDfaMap{
		for idx, set := range value{
			if !chomskyRegex.MatchString(set) {												// If current set is not CNF
				(*pDfaMap)[key][idx] = recursiveOptimize(&chosmkyMap, set, &lastX, &lastY)	// Optimize it
			}
		} //forMapValue
	} //forMap

	for key, value := range chosmkyMap {	// After optimizing each set, add sets from chomskyMap to pDfaMap
		(*pDfaMap)[key] = []string{value}
	}
}

/*
This function optimizes set recursively till it is CNF, it will also create CNF for every subset of this set

Parameters:
	pChomskyMap - pointer on map with CNF elements
	set - a set string which will be CNF-ed
	addType - either 'X' or 'Y'
	addTypeCounter - counter for addType, used to control number of addType elements
Return:
	string - a new set, basically a CNF of original set
*/
func recursiveOptimize(pChomskyMap *map[string]string, set string, pLastX *int, pLastY *int) string{
	var returnSet = ""
	// By the way, '== 2' should also work in this condition, cuz sets of len 1 are passed by regex and cannot be here
	if len(set) <= 2 {						// There are 2 cases, if len of set if <=2 - we can finally get CNF
		for idx, element := range set{
			if isLowerChar(element) {  		// Even with <=2 symbols we can still have terminal symbols
				returnSet = set[:idx] + set[idx+1:]	// Delete terminal character
				returnSet += addOrGetValue(pChomskyMap, string(element),'X', pLastX) // Add a key with value...
																						//...of this terminal character
				return returnSet
			}
		}
	}else {									// If len is >2 - we need to optimize it further
		for idx, element := range set {		// Go through every element
			if element >= 'A' && element <= 'Z' {	// If it is non terminal
				returnSet = string(element)			// We add to return state
				set = set[:idx] + set[idx+1:]		// Take the rest of string (without this exact non terminal)
				break								// ... and break
			}
		}
		// Boy this function is very spooky and dangerous, but it basically does all the job
		// It optimizes our set to "maximum" possible level and after this we create a value for it (or get if exist)
		// You can see I am suing 'Y' type, because Y will always contain 'AB' form (2 non term) as value
		return returnSet + addOrGetValue(pChomskyMap, recursiveOptimize(pChomskyMap, set, pLastX, pLastY), 'Y', pLastY)
	}
	return "If you see this -> something terrible happened on run-time"
}


/*
This function has 2 use cases:
If 'set' does not exist in map - add it to the map, increase counters and return new key fo this value
If 'set' exists - return key for it

Return key is always of form 'addType + [0 .. addTypeCounter]' e.g 'X1', 'Y5', 'X666'

Parameters:
	pChomskyMap - pointer on map with CNF elements
	set - a set string which we will search inside of map
	addType - either 'X' or 'Y'
	addTypeCounter - counter for addType, used to control number of addType elements
Return:
	key
 */
func addOrGetValue(pChomskyMap *map[string]string, set string,addType rune, addTypeCounter *int) string{
	isPresent := containsValue(*pChomskyMap, set)
	if isPresent != "" {
		return isPresent
	}else{
		var chomskyKey = string(addType) + strconv.Itoa(*addTypeCounter)
		*addTypeCounter++
		(*pChomskyMap)[chomskyKey] = set
		return chomskyKey
	}
}

/*
Checks if map contains some specific value in any key.

Parameters:
	m - map string->string in which we search the value
Return:
	key - if value is found in map
	"" - otherwise
 */
func containsValue(m map[string]string, v string) string {
	for key, x := range m {
		if x == v {
			return key
		}
	}
	return ""
}

// -----------------------------------
// UTILS

// Do I need to explain this?
func printDfa(_dfaMap map[string][]string){
	println()
	for key, value := range _dfaMap {
		for _, set := range value {
			println(key + " -> " + set)
		}
	}
} //printDfa

/*
Removes element from a string slice and returns new string
 */
func fastRemove(s []string, i int) []string {
	return append(s[:i], s[i+1:]...)
} //fastRemove

/*
Finds if 'val' is present in slice
 */
func find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
} //find

func isLowerChar(val rune) bool {
	return val >= 'a' && val <= 'z'
} //isLowerChar


