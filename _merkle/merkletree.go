package merkle

import "io"

type Tree struct {
	root *node

	depth    uint64
	numLeafs uint64
}

type node struct {
	hash                []byte
	parent, left, right *node

	isLeaf bool
	data   []byte
}

/*
type leaf struct {
}
*/

func newTree() *Tree {
	return &Tree{
		root: &node{},
	}
}

func createTree(data io.Reader) *Tree {
	return nil
}

func (t *Tree) Add(data []byte) error {
	if t.root.left == nil && t.root.right == nil {
		t.root.isLeaf = true
		t.root.data = data
	} else {

	}
}

func (t *Tree) Del(data []byte) error {

}

func (t *Tree) Exist(data []byte) error {

}

func (t *Tree) CompareRoot(hash []byte) (bool, error) {

}

/*
func (t *Tree) Read(p []byte) (int, error) {

}

func (t *Tree) Write(data []byte) (int, error) {

}
*/
