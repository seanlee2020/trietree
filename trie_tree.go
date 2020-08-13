package trie

import "strings"

// QueryTrie is a trie of runes with string keys and interface{} values.
// Note that internal nodes have nil values so a stored nil value will not
// be distinguishable and will not be included in Walks.

type TrieTree struct {
	root *TrieNode
}

// NewTrieTree allocates and returns a new *TrieTree.
func NewTrieTree() *TrieTree {
	trieTree := new(TrieTree)
	trieTree.root = NewTrieNode()
	return trieTree
}

func (trieTree *TrieTree) Get(query string) *TrieNode {
	if empty(query) {
		return nil
	}
	tokens := strings.Fields(query)
	curNode := trieTree.root
	for idx, token := range tokens {
		if curNode.getChildren() != nil {
			if curNode.getChildren()[token] != nil {
				curNode = curNode.getChildren()[token]
				if idx == len(tokens)-1 {
					return curNode
				}
			}

		}

	}
	return nil
}

func (trieTree *TrieTree) GetChildren(query string) map[string]*TrieNode {
	trieNode := trieTree.Get(query)
	if trieNode != nil {
		return trieNode.getChildren()
	}
	return nil
}

func (trieTree *TrieTree) Insert(query string) {
	if empty(query) {
		return
	}
	tokens := strings.Fields(query)

	curNode := trieTree.root

	for idx, token := range tokens {
		if curNode.children == nil {
			curNode.children = make(map[string]*TrieNode)
		}
		if curNode.getChildren()[token] == nil {
			curNode.children[token] = NewTrieNode()
		}
		curNode = curNode.children[token]

		if idx == len(tokens)-1 {
			curNode.endQ = true
		}
	}
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
