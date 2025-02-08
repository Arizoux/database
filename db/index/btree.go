package index

import "personalDB/db/storage"

//Implements B-Tree indexing

type BTree struct {
	root uint64 //page num of root node

	get func(uint64) []byte //func to read a page from disk
	new func([]byte) uint64 // Function to allocate and write a new page
	del func(uint64)        // Function to deallocate a page
}

func leafInsert(n storage.Node, old storage.Node, idx uint16, key []byte, value []byte) {}
