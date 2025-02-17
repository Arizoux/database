package BPlusTree

import "encoding/binary"

/*
this function inserts a new kv-pair into a leaf node. This is not done in-place (copy-on-write) in order to prevent corruption,
but instead a new node is created from the old node (so that, should an error happen during writing, data is not lost
*/

func leafInsert(new Node, old Node, idx uint16, key []byte, value []byte) {
	new.setHeader(BTreeLeaf, old.numKeys()+1)
	nodeCopyRange(new, old, 0, 0, idx)                     //copy all the kv-pairs from old node to new node up to the index where we want to insert the new kv-pair
	nodeInsertKV(new, idx, 0, key, value)                  //the pointer is set to 0 because the node is a leaf node (but we still need a value in order for our node functions to work)
	nodeCopyRange(new, old, idx+1, idx, old.numKeys()-idx) //copy the rest of the kv-pairs. Add one to the index of the new node because we inserted the new key at its place.
}

func leafUpdate(new Node, old Node, idx uint16, key []byte, val []byte) {
	new.setHeader(BTreeLeaf, old.numKeys())
	nodeCopyRange(new, old, 0, 0, idx)
	nodeInsertKV(new, idx, 0, key, val)
	nodeCopyRange(new, old, idx+1, idx+1, old.numKeys()-(idx+1))
}

// insert a kv-pair at a certain index in a node
func nodeInsertKV(n Node, idx uint16, ptr uint64, key []byte, value []byte) {
	n.setPointer(idx, ptr)

	kvPos := n.kvPos(idx)
	binary.LittleEndian.PutUint16(n[kvPos:], uint16(len(key)))     //set the key length at the kv-pos + 0. Cast the len func to uint16 because it normally returns an int64
	binary.LittleEndian.PutUint16(n[kvPos+2:], uint16(len(value))) //set the value at kv-pos + 2

	copy(n[kvPos+4:], key) // copy the key and value to the correct positions
	copy(n[kvPos+4+uint16(len(key)):], value)

	n.setOffset(idx+1, n.getOffset(idx)+4+uint16(len(key)+len(value)))
}

// copy kv-pairs from an old node to a new one
func nodeCopyRange(new Node, old Node, idxNew uint16, idxOld uint16, n uint16) {
	for i := uint16(0); i < n; i++ {
		nodeInsertKV(new, idxNew+i, old.getPointer(idxOld+i), old.getKey(idxOld+i), old.getVal(idxOld+i))
	}
}

// replace a link with one or multiple links
func nodeUpdateKidsInternal(tree *BTree, new Node, old Node, idx uint16, kids ...Node) {
	inc := uint16(len(kids))
	new.setHeader(BTreeInternal, old.numKeys()+inc-1)

	nodeCopyRange(new, old, 0, 0, idx)
	for i, node := range kids {
		nodeInsertKV(new, idx+uint16(i), tree.new(node), node.getKey(0), nil)
	}
	nodeCopyRange(new, old, idx+inc, idx+1, old.numKeys()-(idx+1))
}

// split a node into 2 new ones
func splitNode2(old Node, left Node, right Node) {

}

// split a node if it's too big. the results are 1~3 nodes
func splitNode3(old Node) (uint16, [3]Node) {
	if old.nbytes() <= MaxPageSize {
		// guarantees that the node is not bigger than MaxPageSize
		old = old[:MaxPageSize]
		return 1, [3]Node{old}
	}

	left := Node(make([]byte, 2*MaxPageSize)) //this node might get split again
	right := Node(make([]byte, MaxPageSize))

	splitNode2(old, left, right)

	if left.nbytes() <= MaxPageSize {
		left = left[:MaxPageSize]
		return 2, [3]Node{left, right}
	}

	leftleft := Node(make([]byte, MaxPageSize))
	middle := Node(make([]byte, MaxPageSize))

	splitNode2(leftleft, middle, left)

	if leftleft.nbytes() > MaxPageSize {
		panic("leftleft node is still too big after splitting")
	}

	return 3, [3]Node{leftleft, middle, right}
}
