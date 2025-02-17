package BPlusTree

type BTree struct {
	root uint64 //page num of root node

	get func(uint64) []byte // Function to read a page from disk
	new func([]byte) uint64 // Function to allocate and write a new page
	del func(uint64)        // Function to deallocate a page
}

// Insert a new key or update an existing key
func (tree *BTree) Insert(key []byte, val []byte) {

}

// Delete a key and returns whether the key was there
func (tree *BTree) Delete(key []byte) bool {

}
