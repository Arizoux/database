package BPlusTree

import "errors"

// remove a key from a leaf node by skipping over it (this is possible due to copy-on-write)
func leafDelete(new Node, old Node, idx uint16) {
	new.setHeader(BTreeLeaf, old.numKeys()-1)
	nodeCopyRange(new, old, 0, 0, idx)
	nodeCopyRange(new, old, idx, idx+1, old.numKeys()-idx)
}

// merge 2 nodes into 1
func nodeMerge(new Node, left Node, right Node) error {
	if left.getNodeType() != right.getNodeType() {
		return errors.New("node types are not equal (nodeMerge)")
	}

	new.setHeader(left.getNodeType(), right.numKeys()+left.numKeys())
	nodeCopyRange(new, left, 0, 0, left.numKeys())
	nodeCopyRange(new, right, left.numKeys()+1, 0, right.numKeys())
	return nil
}

// replace 2 adjacent links with 1
func nodeReplace2Links(new Node, old Node, idx uint16, ptr uint64, key []byte) {
	new.setHeader(BTreeInternal, old.numKeys()-1)
	nodeCopyRange(new, old, 0, 0, idx-1)
	nodeInsertKV(new, idx, ptr, key, nil)
	nodeCopyRange(new, old, idx+1, idx+2, old.numKeys()-idx-2)
}

/*
returns which sibling (left or right) to merge with
parent: the parent node of the updated node
idx: the index of the updated node in the parent node
updated: the node that was updated
*/
func shouldMerge(tree *BTree, parent Node, idx uint16, updated Node) (int, Node) {
	if updated.nbytes() > MaxPageSize/4 {
		return 0, Node{}
	}
	// checks if the node is not the left most node. If it isn't then it is merged with the left sibling
	if idx > 0 {
		sibling := Node(tree.get(parent.getPointer(idx - 1)))
		mergedSize := sibling.nbytes() + updated.nbytes() - Header
		if mergedSize <= MaxPageSize {
			return -1, sibling
		}
	}
	// checks if the node is not the right most node. If it isn't then it is merged with the right sibling
	if idx+1 < parent.numKeys() {
		sibling := Node(tree.get(parent.getPointer(idx + 1)))
		merged := sibling.nbytes() + updated.nbytes() - Header
		if merged <= MaxPageSize {
			return 1, sibling // right
		}
	}
	//if no merge can happen because the other nodes are too big, or it has no siblings, return 0
	return 0, Node{}
}
