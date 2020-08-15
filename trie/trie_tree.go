package trie

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

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

func (trieTree *TrieTree) LoadData(dataFile string) {
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
		if !head {
			if count%10000 == 0 {
				fmt.Println(line)
			}
			count++
			fields := strings.Split(line, ",")
			query := fields[0]
			nu, _ := strconv.Atoi(fields[1])
			ns, _ := strconv.Atoi(fields[2])
			nh := 100
			trieTree.Insert(query, nu, ns, nh)
		} else {
			head = false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func (trieTree *TrieTree) Insert(query string, nu int, ns int, nh int) {
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
			newNode := NewTrieNode()
			newNode.Token = token
			curNode.children[token] = newNode
		}
		curNode = curNode.children[token]
		if idx == len(tokens)-1 {
			curNode.EndQ = true
			curNode.NumUsers = nu
			curNode.NumSessions = ns
			curNode.numHits = nh
		}
	}
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
