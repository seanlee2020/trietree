package trie

type TrieNode struct {
	token       string
	children    map[string]*TrieNode
	isLeaf      bool
	endQ        bool
	numUsers    int
	numSessions int
	numHits     int
	reverse     bool
}

// NewTrieNode allocates and returns a new *TrieNode.
func NewTrieNode() *TrieNode {
	return new(TrieNode)
}

func (trieNode *TrieNode) getChildren() map[string]*TrieNode {
	return trieNode.children
}
