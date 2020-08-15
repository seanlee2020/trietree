package main

import (
	"fmt"
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

	/*
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		})
		log.Println("listening on port 8088")
		log.Fatal(http.ListenAndServe(":8088", nil))
	*/

	r := mux.NewRouter()
	r.HandleFunc("/pills/{query}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		q := vars["query"]

		//fmt.Fprintf(w, "You've requested the book: %s on page %s\n", title, page)

		children := QueryTrie.GetChildren(q)

		fmt.Fprint(w, "size of children is", len(children))

		var nodeList = []*trie.TrieNode{}

		for _, node := range children {
			if node.EndQ {
				nodeList = append(nodeList, node)
			}
		}

		fmt.Fprint(w, "\nsize of nodeList is", len(nodeList))

		/*
			for idx, node := range nodeList {
				fmt.Fprint(w,"\nidx=", idx)
				fmt.Fprint(w,"\nnode.token=", node.Token)
			}*/

		sortNodes(nodeList)

		for idx, node := range nodeList {
			fmt.Fprint(w, "\nidx=", idx)
			fmt.Fprint(w, "\nnode.token=", node.Token)
		}

	})

	http.ListenAndServe(":80", r)
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
