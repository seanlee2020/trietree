package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/seanlee2020/trietree/trie"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
)

var QueryTrie *trie.TrieTree
var ReverseQueryTrie *trie.TrieTree

var BlockList map[string]bool

func init() {
	BlockList = InitBockList()
	QueryTrie = InitTrie()
	ReverseQueryTrie = InitReverseTrie()
}

func main() {

	r := mux.NewRouter()

	r.HandleFunc("/searchpills", processPills)

	http.ListenAndServe(":8088", r)
}

//mem var need to be epxorted ( fist char upper cased)
type Pill struct {
	Token       string
	Query       string
	NumUsers    int `json:",omitempty"`
	NumSessions int `json:",omitempty"`
}

func processPills(w http.ResponseWriter, r *http.Request) {
	reqParams := r.URL.Query()
	//fmt.Fprintln(w, "reqParams", reqParams)

	q := reqParams["query"][0]

	explain := false

	if reqParams["explain"] != nil && len(reqParams["explain"]) == 1 && strings.EqualFold(reqParams["explain"][0], "true") {
		explain = true
	}

	de := false

	if reqParams["de"] != nil && len(reqParams["de"]) == 1 && strings.EqualFold(reqParams["de"][0], "true") {
		de = true
	}

	children := QueryTrie.GetChildren(q)

	var nodeList = []*trie.TrieNode{}

	for _, node := range children {
		if node.EndQ {
			nodeList = append(nodeList, node)
		}
	}

	revertedQ := revertQuery(q)
	reverseChildren := ReverseQueryTrie.GetChildren(revertedQ)
	var reverseNodeList = []*trie.TrieNode{}
	for _, node := range reverseChildren {
		if node.EndQ {
			node.Reverse = true
			reverseNodeList = append(reverseNodeList, node)
		}
	}

	//fmt.Fprint(w, "\nsize of nodeList is", len(nodeList))

	newNodeList := removeBlockedQ(nodeList, q)
	reverseNewNodeList := removeBlockedQ(reverseNodeList, q)
	newNodeList = append(newNodeList, reverseNewNodeList...)
	if de {
		sortNodesByAlphabetic(newNodeList)
		newNodeList = duplicateRemoval(newNodeList)
		newNodeList = duplicateRemovalOrdering(newNodeList)
		sortNodesByPopularity(newNodeList)
	} else {
		sortNodesByPopularity(newNodeList)
	}

	var pillList = []*Pill{}

	for _, node := range newNodeList {
		pill := new(Pill)
		pill.Token = node.Token
		if !node.Reverse {
			pill.Query = q + " " + pill.Token
		} else {
			pill.Query = pill.Token + " " + q
		}
		if explain {
			pill.NumUsers = node.NumUsers
			pill.NumSessions = node.NumSessions
		}
		pillList = append(pillList, pill)
	}

	js, err := json.Marshal(pillList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func InitBockList() map[string]bool {

	blockList := make(map[string]bool)

	dataFile := "/Users/seanl/data/search_browse/blocklist.csv"

	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	head := true
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if head {
			head = false
			continue
		} else {
			if count%10000 == 0 {
				fmt.Println(line)
			}
			fields := strings.Split(line, ",")
			query := fields[0]
			blockList[query] = true
			count++
		}
	}

	return blockList
}

func InitTrie() *trie.TrieTree {

	tt := trie.NewTrieTree()

	tt.LoadData("/Users/seanl/data/search_browse/queries_nu_ns_nuser_2_nsession_2.csv", false)
	return tt

}

func InitReverseTrie() *trie.TrieTree {

	tt := trie.NewTrieTree()
	tt.LoadData("/Users/seanl/data/search_browse/queries_nu_ns_nuser_2_nsession_2.csv", true)
	return tt

}

func sortNodesByPopularity(nl []*trie.TrieNode) {
	sort.SliceStable(nl, func(i, j int) bool {
		nodeI, nodeJ := nl[i], nl[j]
		scoreI := getScore(nodeI)
		scoreJ := getScore(nodeJ)
		return scoreI > scoreJ
	})
}

func sortNodesByAlphabetic(nl []*trie.TrieNode) {
	sort.SliceStable(nl, func(i, j int) bool {
		nodeI, nodeJ := nl[i], nl[j]
		return nodeI.Token < nodeJ.Token
	})
}

func removeBlockedQ(nl []*trie.TrieNode, query string) []*trie.TrieNode {
	var newNodeList = []*trie.TrieNode{}

	for _, node := range nl {

		q := query + " " + node.Token
		if !BlockList[q] && !BlockList[node.Token] && !BlockList[query] {
			newNodeList = append(newNodeList, node)
		}
	}
	return newNodeList

}

/*
select one node from clusters as
[book, books]
[classic, classical, lassics]
*/
func duplicateRemoval(nodeList []*trie.TrieNode) []*trie.TrieNode {

	var newNodeList = []*trie.TrieNode{}
	if len(nodeList) == 0 {
		return nodeList
	}
	preNode := nodeList[0]
	idx := 1
	for idx < len(nodeList) {
		node := nodeList[idx]

		//nextNode := nodeList[idx+1]
		//nextNode = nil

		//if len(node.Token) <= (len(preNode.Token)+2) && len(node.Token) > len(preNode.Token) && node.Token[:len(preNode.Token)] == preNode.Token {
		if containsAndSimilar(node.Token, preNode.Token) {
			winner := selectWinner(preNode, node)

			j := idx + 1
			for j < len(nodeList) {
				node := nodeList[j]
				//if len(node.Token) <= (len(preNode.Token)+2) && len(node.Token) > len(preNode.Token) && node.Token[:len(preNode.Token)] == preNode.Token {
				if containsAndSimilar(node.Token, preNode.Token) {
					winner = selectWinner(winner, node)
				} else {
					break
				}

				j++
			}
			newNodeList = append(newNodeList, winner)

			if j >= len(nodeList) {
				break
			}

			preNode = nodeList[j]
			idx = j + 1
		} else {
			newNodeList = append(newNodeList, preNode)
			preNode = node

			idx += 1
		}
	}
	return newNodeList
}

func duplicateRemovalOrdering(nodeList []*trie.TrieNode) []*trie.TrieNode {
	var newNodeList = []*trie.TrieNode{}
	if len(nodeList) == 0 {
		return nodeList
	}
	token2Nodes := make(map[string][]*trie.TrieNode)
	for _, node := range nodeList {
		if token2Nodes[node.Token] == nil {
			list := []*trie.TrieNode{}
			token2Nodes[node.Token] = list
		}
		token2Nodes[node.Token] = append(token2Nodes[node.Token], node)
	}

	for _, val := range token2Nodes {

		if len(val) == 1 {
			newNodeList = append(newNodeList, val[0])
		} else {
			score0 := getScore(val[0])
			score1 := getScore(val[1])
			if score0 > score1 {
				newNodeList = append(newNodeList, val[0])
			} else {
				newNodeList = append(newNodeList, val[1])
			}
		}
	}
	return newNodeList

}

func selectWinner(nodeA *trie.TrieNode, nodeB *trie.TrieNode) *trie.TrieNode {
	scoreA := getScore(nodeA)
	scoreB := getScore(nodeB)

	if scoreA >= scoreB {
		return nodeA
	} else {
		return nodeB
	}
}

func getScore(node *trie.TrieNode) int {
	score := node.NumUsers*10 + node.NumSessions*3
	return score
}

func containsAndSimilar(tokenA string, tokenB string) bool {

	if len(tokenA) <= (len(tokenB)+2) && len(tokenA) > len(tokenB) && tokenA[:len(tokenB)] == tokenB {
		return true
	}
	return false

}

func revertQuery(query string) string {
	tokens := strings.Fields(query)
	idx := len(tokens) - 1
	ret := tokens[idx]
	idx--
	for idx >= 0 {
		ret += " " + tokens[idx]
	}
	return ret
}
