package BTree

type BTree struct {
	root uint64 //page num of root node

	get func(uint64) []byte // Function to read a page from disk
	new func([]byte) uint64 // Function to allocate and write a new page
	del func(uint64)        // Function to deallocate a page
}
