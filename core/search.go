package core

import (
	"github.com/cbergoon/merkletree"
)

func insertContent(slice []*entry, tree []merkletree.Content, newId *entry) ([]*entry, []merkletree.Content) {
	length := len(slice)
	currIdx := 0
	maxIdx := length - 1

	if maxIdx == -1 {
		slice = append(slice, newId)
		tree = append(tree, newId)
		return slice, tree
	}

	for {
		if currIdx >= maxIdx {
			slice = append(slice, nil)
			copy(slice[currIdx+1:], slice[currIdx:])

			if slice[currIdx].cmp(newId) == -1 {
				currIdx += 1
			}
			slice[currIdx] = newId

			tree = append(tree, nil)
			copy(tree[currIdx+1:], tree[currIdx:])
			tree[currIdx] = newId

			return slice, tree
		}
		mid := (currIdx + maxIdx) / 2

		if slice[mid].cmp(newId) == -1 {
			currIdx = mid + 1
		} else {
			maxIdx = mid - 1
		}
	}
}
