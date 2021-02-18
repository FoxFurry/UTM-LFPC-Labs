//
// Created by foxfurry on 17/02/2021.
//

#ifndef LAB1_FINITE_RECURSIVE_H

#include "automata_dependicies.h"

// Converts string "S=abcS" into pair S/abcS
std::pair<char, std::string> strToProdSet(std::string inp);

// Checks if this char is inside of any string in inp_v vector
bool checkInArray(char inp_c, std::vector<std::string> inp_v);

// check is src is sub string of dst starting from first character
bool isSubstr(std::string src, std::string dst);

/* Recursively check is a production sets of current_nonterm can trim source string
 * To check if input string can be created by defined alphabet I am trying to trim it down to 0 using productions set
 * example:
 * S=bbP
 * P=bd
 * input=bbbd
 * 1) bbbd - S = bd with current_nonterm P
 * 2) bd = P[0] -> true
 */
bool recursiveCheck(std::string source_string, char current_nonterm, std::map<char, std::vector<std::string>> prod_set);

#endif //LAB1_FINITE_RECURSIVE_H
