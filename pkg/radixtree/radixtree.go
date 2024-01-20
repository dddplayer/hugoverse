package radixtree

import (
	"fmt"
	"sort"
	"strings"
)

type edge struct {
	name  string
	start *node
	end   *node
}

type edges []*edge

type node struct {
	val      any
	prefix   *edge
	suffixes edges
}

type Tree struct {
	root *node
}

// New returns an empty Tree
func New() *Tree {
	return &Tree{root: &node{
		val:      "TreeRoot",
		prefix:   nil,
		suffixes: edges{},
	}}
}

func (t *Tree) Get(k string) (interface{}, bool) {
	n := t.root
	path := k
	for {
		if len(path) == 0 {
			return n.val, true
		}

		e := n.getEdge(path)
		if e == nil {
			break
		} else if e.name == "" {
			break
		}

		if strings.HasPrefix(path, e.name) {
			path = path[len(e.name):]
			n = e.end
		} else {
			break
		}
	}
	return nil, false
}

// Insert is used to add/update a value to radix tree
func (t *Tree) Insert(k string, v interface{}) bool {
	var parent *node
	n := t.root
	path := k

	if len(path) == 0 {
		panic("empty key not supported yet.")
	}

	for {
		if len(path) == 0 {
			n.val = v
			return true
		}

		// Look for the edge
		parent = n
		e := n.getEdge(path)

		// No same prefix edge found, create one
		if e == nil {
			e := newNodeEdge()
			e.name = path
			e.start = parent
			e.end.val = v

			parent.addEdge(e)
			return true
		}

		// Determine the longest prefix of the search key on match
		commonPrefix := longestPrefix(path, e.name)
		// Edge found with fully overlap
		// Look into the right depth
		if commonPrefix == len(e.name) {
			path = path[commonPrefix:]
			n = e.end
			continue
		}

		// Right depth found
		// Split the current node
		// Create common edge
		commonEdge := newNodeEdge()
		commonEdge.name = path[:commonPrefix]
		commonEdge.start = e.start

		e.start.addEdge(commonEdge)
		e.start.delEdge(e)
		commonEdge.end.addEdge(e)

		// Update edge name with uncommon part
		e.name = e.name[commonPrefix:]

		// Create the new joined one
		freshEdgeName := path[commonPrefix:]
		if len(freshEdgeName) == 0 {
			commonEdge.end.val = v
		} else {
			freshEdge := newNodeEdge()
			freshEdge.name = freshEdgeName
			freshEdge.start = commonEdge.end
			freshEdge.end.val = v

			commonEdge.end.addEdge(freshEdge)
		}
		return true
	}

	return false
}

func newNodeEdge() *edge {
	e := &edge{
		name:  "",
		start: nil,
		end: &node{
			val:      nil,
			prefix:   nil,
			suffixes: edges{},
		},
	}
	e.end.prefix = e
	return e
}

func (n *node) getEdge(path string) *edge {
	return findEdgeWithSamePrefix(
		getFirstByte(path), n.suffixes)
}

func (n *node) addEdge(e *edge) {
	e.start = n
	n.suffixes = append(n.suffixes, e)
	sortAscending(n.suffixes)
}

func (n *node) delEdge(e *edge) {
	if e.start == n {
		var pos int
		for i, item := range n.suffixes {
			if item == e {
				pos = i
				break
			}
		}
		n.suffixes = append(n.suffixes[:pos],
			n.suffixes[pos+1:]...)
	}
}

func sortAscending(es edges) {
	sort.Slice(es, func(i, j int) bool {
		return getFirstByte(
			es[i].name) < getFirstByte(es[j].name)
	})
}

func getFirstByte(v string) byte {
	return v[0]
}

func findEdgeWithSamePrefix(
	firstByte byte, es edges) *edge {
	for _, e := range es {
		if getFirstByte(e.name) == firstByte {
			return e
		}
	}
	return nil
}

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int {
	longest := len(k1)
	if l := len(k2); l < longest {
		longest = l
	}
	var i int
	for i = 0; i < longest; i++ {
		if k1[i] != k2[i] {
			break
		}
	}
	return i
}

func (t *Tree) Print() {
	recursivePrint(t.root)
}

func recursivePrint(n *node) {
	printNode(n)
	for _, e := range n.suffixes {
		recursivePrint(e.end)
	}
}

func printNode(n *node) {
	prefix := ""
	if n.prefix != nil {
		prefix = n.prefix.name
	}
	fmt.Printf("node: prefix %s, value %s\n",
		prefix, n.val)
}

// WalkFn is used when walking the tree. Takes a
// key and value, returning if iteration should
// be terminated.
type WalkFn func(s string, v interface{}) bool

func (t *Tree) Walk(fn WalkFn) {
	recursiveWalk(t.root.prefix.name, t.root, fn)
}

// recursiveWalk is used to do a pre-order walk of a node
// recursively. Returns true if the walk should be terminated
func recursiveWalk(prefix string, n *node, fn WalkFn) bool {
	pre := prefix + n.prefix.name
	if fn(pre, n.val) {
		return true
	}
	for _, e := range n.suffixes {
		if recursiveWalk(pre, e.end, fn) {
			return true
		}
	}
	return false
}
