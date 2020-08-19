package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/seanlee2020/trietree/trie"
	"net/http"
	"sort"
	"strings"
)

var QueryTrie *trie.TrieTree

func init() {
	QueryTrie = InitTrie()
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
	//fmt.Fprint(w, "\nsize of nodeList is", len(nodeList))

	newNodeList := nodeList
	if de {
		sortNodesByAlphabetic(newNodeList)
		newNodeList = duplicateRemoval(newNodeList)
		sortNodesByPopularity(newNodeList)
	} else {
		sortNodesByPopularity(newNodeList)
	}

	var pillList = []*Pill{}

	for _, node := range newNodeList {
		pill := new(Pill)
		pill.Token = node.Token
		pill.Query = q + " " + pill.Token

		if explain {
			pill.NumUsers = node.NumUsers
			pill.NumSessions = node.NumSessions
		}
		pillList = append(pillList, pill)

		/*fmt.Fprint(w, "\nidx=", idx)
		fmt.Fprint(w, "\nnode.token=", node.Token)
		fmt.Fprint(w, "\npill.query=", pill.Query)
		*/
	}

	js, err := json.Marshal(pillList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)

}

func InitTrie() *trie.TrieTree {

	tt := trie.NewTrieTree()

	tt.LoadData("/Users/seanl/data/search_browse/queries_nu_ns_nuser_2_nsession_2.csv")
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

/*
select one node from clusters as
[book, books]
[classic, classical, lassics]
*/
func duplicateRemoval(nodeList []*trie.TrieNode) []*trie.TrieNode {

	var newNodeList = []*trie.TrieNode{}
	preNode := nodeList[0]
	idx := 1
	for idx < len(nodeList) {
		node := nodeList[idx]

		//nextNode := nodeList[idx+1]
		//nextNode = nil

		if len(node.Token) <= (len(preNode.Token)+2) && len(node.Token) > len(preNode.Token) && node.Token[:len(preNode.Token)] == preNode.Token {
			winner := selectWinner(preNode, node)
			newNodeList = append(newNodeList, winner)

			if idx+1 >= len(nodeList) {
				break
			}
			preNode = nodeList[idx+1]
			idx += 2
		} else {
			newNodeList = append(newNodeList, preNode)
			preNode = node

			idx += 1
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
