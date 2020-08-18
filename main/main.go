package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/seanlee2020/trietree/trie"
	"net/http"
	"sort"
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

	children := QueryTrie.GetChildren(q)

	var nodeList = []*trie.TrieNode{}

	for _, node := range children {
		if node.EndQ {
			nodeList = append(nodeList, node)
		}
	}
	//fmt.Fprint(w, "\nsize of nodeList is", len(nodeList))
	sortNodes(nodeList)

	var pillList = []*Pill{}

	for _, node := range nodeList {
		pill := new(Pill)
		pill.Token = node.Token
		pill.Query = q + " " + pill.Token
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

func sortNodes(nodeList []*trie.TrieNode) {
	sort.SliceStable(nodeList, func(i, j int) bool {
		nodeI, nodeJ := nodeList[i], nodeList[j]
		scoreI := nodeI.NumUsers*10 + nodeI.NumSessions*3
		scoreJ := nodeJ.NumUsers*10 + nodeJ.NumSessions*3
		return scoreI > scoreJ
	})

}
