//
// Created by foxfurry on 17/02/2021.
//

#include "finite_recursive.h"

std::pair<char, std::string> strToProdSet(std::string inp) {
    std::pair<char, std::string> output;
    output.first = inp[0];
    output.second = inp.substr(2, inp.size() - 2);
    return output;
}

bool checkInArray(char inp_c, std::vector<std::string> inp_v) {
    for (int idx = 0; idx < inp_v.size(); idx++) {
        if (inp_v[idx].find(inp_c) != std::string::npos)return true;
    }
    return false;
}

bool recursiveCheck(std::string source_string, char current_nonterm, std::map<char, std::vector<std::string>> prod_set) {
    if (source_string.size() == 0)return true;  // Check if string is abolutely trimmed

    for (int idx = 0; idx < prod_set[current_nonterm].size(); idx++) {  // Go thro every production of current_nonterm
        bool res = false;
        std::string prodSetVal = prod_set[current_nonterm][idx]; // Current production, for ex: abcS
        if (prodSetVal == source_string)return true;        // If current production is same as source string - exit
        if (isSubstr(prodSetVal, source_string)) {          // If current production is substring of source string -
            res = recursiveCheck(                           // trim source string and recursively check for next nonterm
                    source_string.substr(prodSetVal.size() - 1, source_string.size() - (prodSetVal.size() - 1)),
                    *(prodSetVal.end() - 1), prod_set);
        }
        if (res)return true;
    }
    return false;   // TRY DELETING THIS LINES IF IT FAILS, MAYBE IT WILL FIX IT
}

bool isSubstr(std::string src, std::string dst) {
    for (int idx = 0; idx < src.size(); idx++) {
        if (isupper(src[idx]) && idx != 0)return true;
        if (src[idx] != dst[idx])return false;
    }
    return true;
}