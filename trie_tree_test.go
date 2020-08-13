package trie

import (
	"fmt"
	"testing"
)

func TestTrieTree(t *testing.T) {
	trie := NewTrieTree()
	testTrieTree(t, trie)
}

func testTrieTree(t *testing.T, trie *TrieTree) {
	q1 := "harry"
	q2 := "harry potter"
	q3 := "harry potter book"
	q4 := "harry potter movie"

	trie.Insert(q1, 100, 30, 100)
	trie.Insert(q2, 9, 20, 200)
	trie.Insert(q3, 40, 50, 100)
	trie.Insert(q4, 60, 100, 300)

	node2 := trie.Get(q2)

	if node2 != nil {
		fmt.Print("node2 is not null\n")
	}
	children := trie.GetChildren(q2)
	fmt.Print("size of children is", len(children))

	var nodeList = []*TrieNode{}

	for _, node := range children {
		nodeList = append(nodeList, node)
	}

	fmt.Print("\nsize of nodeList is", len(nodeList))

	for idx, node := range nodeList {

		fmt.Print("\nidx=", idx)
		fmt.Print("\nnode.token=", node.token)

	}

}
