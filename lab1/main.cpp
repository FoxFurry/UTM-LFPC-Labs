/*
 * What's bad about this implementation?
 * - No error handling, almost
 * - Big input restrictions (absolutely covers task inputs, but in global term they are big)
 * IMPORTANT! IF IT WILL EVENTUALLY FAIL - PLEASE DELETE LINE NUMBER 172 <return false;> I KNOW THIS SOUNDS STUPID SORRY
 */


#include <iostream>
#include <vector>
#include <map>
#include <algorithm>

using namespace std;

//some extra prints
#define DEBUG 1

// Converts string "S=abcS" into pair S/abcS
pair<char, string> strToProdSet(string inp);

// Checks if this char is inside of any string in inp_v vector
bool checkInArray(char inp_c, vector<string> inp_v);

/* Recursively check is a production sets of current_nonterm can trim source string
 * To check if input string can be created by defined alphabet I am trying to trim it down to 0 using productions set
 * example:
 * S=bbP
 * P=bd
 * input=bbbd
 * 1) bbbd - S = bd with current_nonterm P
 * 2) bd = P[0] -> true
 */
bool recursiveCheck(string source_string, char current_nonterm, map<char, vector<string>> prod_set);

// check is src is sub string of dst starting from first character
bool isSubstr(string src, string dst);

int main() {

    // Temporary values
    char tmpNonTerm = 0;
    string tmpTerm = "";
    string tmpProd = "";
    string input = "";
    pair<char, string> tmpProdPair;

    // Alphabet values
    vector<char> nonTerminal;
    vector<string> terminal;
    map<char, vector<string>> productionSet;

    // Scanning non-terminating (Vn) symbols. 0 for end-scan
    printf("Enter non-terminating symbols (only uppercase chars are accepted)\nEnter 0 (zero) to terminate\n");
    while (true) {
        printf(" >> ");
        cin >> tmpNonTerm;
        if (tmpNonTerm == '0') {
            if (nonTerminal.size() == 0) {
                printf("No nonterminal characters found. Terminating...");
                return 1;
            }
            break;
        } else if (!isalpha(tmpNonTerm)) {
            printf("Input is outside of range, ignoring\n");
            continue;
        } else if (islower(tmpNonTerm)) {
            tmpNonTerm = toupper(tmpNonTerm);
        }
        nonTerminal.push_back(tmpNonTerm);
    }
    if (DEBUG) {
        printf("Vn = {");
        for (int idx = 0; idx < nonTerminal.size() - 1; idx++) {
            printf("%c,", nonTerminal[idx]);
        }
        printf("%c}\n\n", *(nonTerminal.end() - 1));
    }

    // Scanning terminating symbols (Vt). 0 for end-scan
    printf("Enter terminating symbols (only lowercase strings are accepted)\nEnter 0 (zero) to terminate\n");
    while (true) {
        printf(" >> ");
        cin >> tmpTerm;

        if (tmpTerm == "0") {
            if (terminal.size() == 0) {
                printf("No terminal strings found. Terminating...");
                return 1;
            }
            break;
        }
        transform(tmpTerm.begin(), tmpTerm.end(), tmpTerm.begin(), ::tolower);  // Input string is always converted
                                                                                // to lower
        terminal.push_back(tmpTerm);
    }
    if (DEBUG) {
        cout << "Vе = {";
        for (int idx = 0; idx < nonTerminal.size() - 1; idx++) {
            cout << terminal[idx] << ",";
        }
        cout << terminal[terminal.size()-1] << "}\n\n";
    }


    // Scanning production sets (P={})
    printf("Enter set of productions in this format: <uppercase char> = <string>\n");
    while (true) {
        printf(" >> ");
        cin >> tmpProd;

        if (tmpProd == "0") {
            break;
        }

        tmpProdPair = strToProdSet(tmpProd);

        productionSet[tmpProdPair.first].push_back(tmpProdPair.second);
    }

    // Input word
    cin >> input;

    // Checking if such word can even be possibly generated
    for (auto idx: input) {
        if (isupper(idx)) {
            continue;
        }
        if (!checkInArray(idx, terminal)) {
            printf("Cannot create such word from alphabet. Terminating...");
            return 0;
        }
    }

    // Checking if such word can be actually generated
    bool result = recursiveCheck(input, 'S', productionSet);    // We start with nonterm S


    if (result == true) {
        cout << "Such word can be generated" << '\n';
    } else {
        cout << "Such word cannot be generated" << '\n';
    }
    return 0;
}

pair<char, string> strToProdSet(string inp) {
    pair<char, string> output;
    output.first = inp[0];
    output.second = inp.substr(2, inp.size() - 2);
    return output;
}

bool checkInArray(char inp_c, vector<string> inp_v) {
    for (int idx = 0; idx < inp_v.size(); idx++) {
        if (inp_v[idx].find(inp_c) != string::npos)return true;
    }
    return false;
}

bool recursiveCheck(string source_string, char current_nonterm, map<char, vector<string>> prod_set) {
    if (source_string.size() == 0)return true;  // Check if string is abolutely trimmed

    for (int idx = 0; idx < prod_set[current_nonterm].size(); idx++) {  // Go thro every production of current_nonterm
        bool res = false;
        string prodSetVal = prod_set[current_nonterm][idx]; // Current production, for ex: abcS
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

bool isSubstr(string src, string dst) {
    for (int idx = 0; idx < src.size(); idx++) {
        if (isupper(src[idx]) && idx != 0)return true;
        if (src[idx] != dst[idx])return false;
    }
    return true;
}