package BPlusTree

type BTree struct {
	root uint64 //page num of root node

	get func(uint64) []byte // Function to read a page from disk
	new func([]byte) uint64 // Function to allocate and write a new page
	del func(uint64)        // Function to deallocate a page
}

// Insert a new key or update an existing key
func (tree *BTree) Insert(key []byte, val []byte) {
	if tree.root == 0 {
		newRoot := Node(make([]byte, MaxPageSize))
		newRoot.setHeader(BTreeLeaf, 2)
		nodeInsertKV(newRoot, 0, 0, nil, nil)
		nodeInsertKV(newRoot, 1, 0, key, val)
		tree.root = tree.new(newRoot)
		return
	}

	node := treeInsert(tree, tree.get(tree.root), key, val)
	nsplit, split := splitNode3(node)
	tree.del(tree.root)

	if nsplit > 1 {
		// the root was split, add a new level.
		root := Node(make([]byte, MaxPageSize))
		root.setHeader(BTreeInternal, nsplit)
		for i, knode := range split[:nsplit] {
			ptr, key := tree.new(knode), knode.getKey(0)
			nodeInsertKV(root, uint16(i), ptr, key, nil)
		}
		tree.root = tree.new(root)
	} else {
		tree.root = tree.new(split[0])
	}
}

// Delete a key and returns whether the key was there
func (tree *BTree) Delete(key []byte) bool {

}
