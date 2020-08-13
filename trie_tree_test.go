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

	trie.Insert(q1)
	trie.Insert(q2)
	trie.Insert(q3)
	trie.Insert(q4)

	node2 := trie.Get(q2)

	if node2 != nil {
		fmt.Print("node2 is not null\n")
	}
	children := trie.GetChildren(q2)
	fmt.Print("size of children is", len(children))
}
