/*
 * What's bad about this implementation?
 * - No error handling, almost
 * - Big input restrictions (absolutely covers task inputs, but in global term they are big)
*/


#include <iostream>
#include <vector>
#include <map>
#include <algorithm>

#include "finite_recursive.h"
#include "finite_graph.h"

using namespace std;


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
    cout << "\nEnter a string to check it by finite automaton\n >> ";
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

    bool result = false;
    if(METHOD == 1){
        Graph grph(productionSet);
        result = grph.checkWord(input);
    }
    else if(METHOD == 2){
        result = recursiveCheck(input, 'S', productionSet);
    }
    else{
        cout << "wrong method specified\n";
    }


    if (result == true) {
        cout << "Such word can be generated" << '\n';
    } else {
        cout << "Such word cannot be generated" << '\n';
    }
    return 0;
}

