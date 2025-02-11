package BTree

import "encoding/binary"

//Implements B-Tree indexing

type BTree struct {
	root uint64 //page num of root node

	get func(uint64) []byte //func to read a page from disk
	new func([]byte) uint64 // Function to allocate and write a new page
	del func(uint64)        // Function to deallocate a page
}

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
