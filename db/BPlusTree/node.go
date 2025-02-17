package BPlusTree

import (
	"bytes"
	"encoding/binary"
)

//Manages fixed-size nodes in b+ tree

/*
Fixed node size:
	type: 2bytes					node type (internal or leaf)
	nkeys: 2bytes					number of keys stored
	pointers: nkeys * 8bytes		child pointer (only for internal nodes)
	offsets: nkeys * 2bytes			offset list for key-value pairs
	key-values:	variable			actual keys and values
	unused: remaining				extra space for future growth

Each key-value pair is stored like this:
	klen: 2b						length of key
	vlen: 2b						length of value
	key: klen
	value: vlen
*/

const (
	MaxPageSize = 4096
	Header      = 4
)

const (
	BTreeInternal = 1
	BTreeLeaf     = 2
)

type Node []byte

// returns the page type that is stored in the first 2 bytes
func (n Node) getNodeType() uint16 {
	return binary.LittleEndian.Uint16(n[0:2])
}

// returns the number of keys in the page
func (n Node) numKeys() uint16 {
	return binary.LittleEndian.Uint16(n[2:4])
}

// sets the data in the header of the page (type of page (leaf or internal) in the first two bytes, num of keys in the 3-4 bytes)
func (n Node) setHeader(nodeType uint16, numKeys uint16) {
	binary.LittleEndian.PutUint16(n[0:2], nodeType)
	binary.LittleEndian.PutUint16(n[2:4], numKeys)
}

// gets the 8 pointer bytes from the page (which located at the header + 8 times (8 bytes is the size of a pointer) the index
func (n Node) getPointer(index uint16) uint64 {
	if index > n.numKeys() {
		panic("getPointer: index out of range")
	}
	pntPos := Header + 8*index
	return binary.LittleEndian.Uint64(n[pntPos : pntPos+8])
}

// sets an 8 byte pointer to the given index + header
func (n Node) setPointer(index uint16, pointer uint64) {
	if index > n.numKeys() {
		panic("setPointer: index out of range")
	}
	pntPos := Header + 8*index
	binary.LittleEndian.PutUint64(n[pntPos:pntPos+8], pointer)
}

// calculates the offset position of the index
func offsetPos(n Node, index uint16) uint16 {
	if index > n.numKeys() || index < 1 {
		panic("offsetPos: index out of range")
	}
	return Header + 8*n.numKeys() + 2*(index-1)
}

// gets the offset of an index
func (n Node) getOffset(index uint16) uint16 {
	if index == 0 {
		return 0
	}
	offset := offsetPos(n, index)
	return binary.LittleEndian.Uint16(n[offset : offset+2])
}

// sets the offset to an index
func (n Node) setOffset(index uint16, offset uint16) {
	binary.LittleEndian.PutUint16(n[offsetPos(n, index):], offset)
}

// gets the POSITION of a kv-pair (not the actual value)
func (n Node) kvPos(index uint16) uint16 {
	if index > n.numKeys() {
		panic("kvPos: index out of range")
	}
	return Header + 8*n.numKeys() + 2*n.numKeys() + n.getOffset(index)
}

/* gets the key for an index by first getting the kv position, with that then decoding the key len
   (which is located directly at the kv position) and then returns the key.
   The key is stored after the first 4 bytes (klen and vlen). The key is klen bytes long.
*/

func (n Node) getKey(index uint16) []byte {
	if index >= n.numKeys() {
		panic("getKey: index out of range")
	}
	kvPos := n.kvPos(index)
	klen := binary.LittleEndian.Uint16(n[kvPos:])
	return n[kvPos+4 : kvPos+4+klen]
}

/* gets the val for an index by first getting the kv position, with that then decoding the value len
   (which is located after the first 2 bytes that contain the key len) and then returns the value.
   The value is stored after the first 4 bytes (klen and vlen) + klen. The value is vlen bytes long.
*/

func (n Node) getVal(index uint16) []byte {
	if index >= n.numKeys() {
		panic("getVal: index out of range")
	}
	kvPos := n.kvPos(index)
	klen := binary.LittleEndian.Uint16(n[kvPos:])
	vlen := binary.LittleEndian.Uint16(n[kvPos+2:])
	return n[kvPos+4+klen : kvPos+4+klen+vlen]
}

// returns the size of the used data in a node (used space)
func (n Node) nbytes() uint16 {
	return n.kvPos(n.numKeys())
}

// performs binary search in the node, returns the index of the key if found, else it returns an error
func nodeLookupBS(n Node, key []byte) uint16 {
	minKey := uint16(0)
	maxKey := n.numKeys()
	midKey := uint16(0)

	for minKey <= maxKey {
		midKey = minKey + (maxKey-minKey)/2

		if bytes.Equal(n.getKey(midKey), key) {
			return midKey
		}
		cmp := bytes.Compare(n.getKey(midKey), key)

		if cmp < 0 {
			minKey = midKey + 1
		}

		if cmp > 0 {
			maxKey = midKey - 1
		}
	}
	return midKey
}
