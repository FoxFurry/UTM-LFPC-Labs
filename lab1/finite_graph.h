//
// Created by foxfurry on 17/02/2021.
//

#ifndef LAB1_FINITE_GRAPH_H
#define LAB1_FINITE_GRAPH_H

#include "automata_dependicies.h"


class Graph {
private:
    std::map<int, std::vector<int>> adjList;
    std::map<int, std::vector<std::string>> graph;
    std::vector<char> nonTermMap;

    // Map non-terminal characters to uniq integer
    int mapVnToInt(char vn) {
        auto it = std::find(nonTermMap.begin(), nonTermMap.end(), vn);

        if (it != nonTermMap.end()) {
            return it - nonTermMap.begin();
        } else {
            nonTermMap.push_back(vn);
            return nonTermMap.size() - 1;
        }
    }

    // Convert production sets into array of nodex (graph)
    void populateGraph(const std::map<char, std::vector<std::string>> &prod_set) {
        for (auto it: prod_set) {
            int varIndex = mapVnToInt(it.first);
            graph[varIndex] = it.second;
        }

        if (DEBUG) {
            std::cout << "Graph populated:\n";
            for (auto it_node: graph) {
                std::cout << it_node.first;
                for (auto it_str: it_node.second) {
                    std::cout << "\t" << it_str << "\n";
                }
            }
            std::cout << "--------------\n";
        }

    }

    // Conver production sets into list of directed edges (adjacent list)
    void generateAdjList(const std::map<char, std::vector<std::string>> &prod_set) {
        for (auto it_prod: prod_set) {
            int source = mapVnToInt(it_prod.first);     // First element of production set is source of edge

            for (auto it_str: it_prod.second) {
                char lastElem = *(it_str.end() - 1);    // Get last element of production set

                if (!isupper(lastElem) || !isalpha(lastElem)) {   // is last character of production is not a variable
                    continue;                                   // skip it
                }

                int destination = mapVnToInt(lastElem);     // Destination element of edge
                if (std::find(adjList[source].begin(), adjList[source].end(),
                              destination) != adjList[source].end()) {
                    continue;   // If such edge already exists - skip
                }
                adjList[source].push_back(destination);
            }
        }

        if (DEBUG) {
            std::cout << "Adjacent list generated:\n";
            for (auto it_src: adjList) {
                std::cout << it_src.first;
                for (auto it_dst: it_src.second) {
                    std::cout << "\t -> " << it_dst << "\n";
                }
            }
            std::cout << "--------------\n";
        }

        if (DEBUG) {
            std::cout << "Indexes of variables:\n";
            for (auto it: nonTermMap) {
                std::cout << it << "\t -> " << mapVnToInt(it) << '\n';
            }
            std::cout << "--------------\n";
        }
    }

    // Iterate thro productions of variable varIndex and check if word can be generated
    bool checkVariable(int varIndex, std::string word, int wordInxex) {
        bool result = false;
        if (wordInxex >= word.size()) {
            return false;   // If index is more than current word - probably something is wrong
        }

        for (auto it_prod: graph[varIndex]) {
            int tmpIndex = wordInxex;   // Each production has its own wordIndex

            int i;
            for (i = 0; i < it_prod.size(); i++) {    // Go thro production
                if (isupper(it_prod[i])) {            // If it is variable - go in recursion
                    result = checkVariable(mapVnToInt(it_prod[i]), word, tmpIndex);
                    if (result) {
                        return true;
                    }
                    break;  // Switch to next production
                } else if (it_prod[i] == word[tmpIndex]) {
                    tmpIndex++;     // If terminal elements match - increase match index
                } else {
                    break;  // Switch to next production
                }
            }

            if (tmpIndex == word.size() &&
                i == it_prod.size()) { // If whole word matches and production has no extra elements
                return true;
            }
        }
        return false;
    }

public:
    // Graph Constructor
    Graph(std::map<char, std::vector<std::string>> prod_set) {
        mapVnToInt('S');    // Populate Vn array with S, so S will be always mapped to 0 (zero)
        populateGraph(prod_set);
        generateAdjList(prod_set);
    }

    bool checkWord(std::string word) {
        return checkVariable(0, word, 0);
    }
};

#endif //LAB1_FINITE_GRAPH_H
