package trie

type TrieNode struct {
	Token       string
	children    map[string]*TrieNode
	isLeaf      bool
	EndQ        bool
	NumUsers    int
	NumSessions int
	NumHits     int
	Traffic     int
	Reverse     bool
}

// NewTrieNode allocates and returns a new *TrieNode.
func NewTrieNode() *TrieNode {
	return new(TrieNode)
}

func (trieNode *TrieNode) getChildren() map[string]*TrieNode {
	return trieNode.children
}
