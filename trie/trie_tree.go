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

func (trieTree *TrieTree) LoadData(dataFile string, reverse bool) {
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
			query := strings.ToLower(fields[0])
			nu := 0
			ns := 0
			nt := 0
			if len(fields) == 3 {
				tmpnu, _ := strconv.Atoi(fields[1])
				tmpns, _ := strconv.Atoi(fields[2])
				nu = tmpnu
				ns = tmpns
			} else if len(fields) == 2 {
				tmpnt, _ := strconv.Atoi(fields[1])
				nt = tmpnt
			}
			nh := 0
			if !reverse {
				trieTree.Insert(query, nu, ns, nh, nt)
			} else {
				revertedQuery := revertQuery(query)
				trieTree.Insert(revertedQuery, nu, ns, nh, nt)
			}
		} else {
			head = false
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func revertQuery(query string) string {
	tokens := strings.Fields(query)
	idx := len(tokens) - 1
	ret := tokens[idx]
	idx--
	for idx >= 0 {
		ret += " " + tokens[idx]
		idx--
	}
	return ret
}

func (trieTree *TrieTree) Insert(query string, nu int, ns int, nh int, nt int) {
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
			curNode.NumHits = nh
			curNode.Traffic = nt
		}
	}
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
