package main__test

import (
	"testing"
	"unsafe"
)

// Define a node structure for the pointer-based binary tree.
type TreeNode struct {
	Value       int
	Left, Right *TreeNode
}

// Recursive function to create a pointer-based binary tree.
func buildPointerTree(depth int) *TreeNode {
	if depth == 0 {
		return nil
	}
	return &TreeNode{
		Value: depth,
		Left:  buildPointerTree(depth - 1),
		Right: buildPointerTree(depth - 1),
	}
}

var array = make([]int, 2)

func ConvertTreeToArray(root *TreeNode, index int) []int {
	if root == nil {
		return array
	}
	if index >= len(array) {
		// make a new array with double the size
		NewArray := make([]int, len(array)*2)
		copy(NewArray, array)
		array = NewArray
	}

	array[index] = root.Value
	ConvertTreeToArray(root.Left, index*2)
	ConvertTreeToArray(root.Right, index*2+1)
	return array
}

// Recursive traversal for the pointer-based binary tree.
func traversePointerTree(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return node.Value + traversePointerTree(node.Left) + traversePointerTree(node.Right)
}

// Array-based binary tree representation.
type ArrayTree struct {
	Nodes []int
}

// Builds an array-based binary tree.
func buildArrayTree(depth int) ArrayTree {
	size := (1 << depth) - 1 // Calculate the total number of nodes for a full binary tree
	nodes := make([]int, size)
	for i := 0; i < size; i++ {
		nodes[i] = i + 1
	}
	return ArrayTree{Nodes: nodes}
}

// Traverses the array-based binary tree in a regular way.
func traverseArrayTree(tree *ArrayTree, index int) int {
	if index >= len(tree.Nodes) {
		return 0
	}
	return tree.Nodes[index] + traverseArrayTree(tree, 2*index) + traverseArrayTree(tree, 2*index+1)
}

func traverseArrayTreeUnsafe(tree *ArrayTree, index int) int {
	if index >= len(tree.Nodes) {
		return 0
	}
	nodeSize := unsafe.Sizeof(tree.Nodes[0])  // Get the size of an int (node) in bytes
	treePtr := unsafe.Pointer(&tree.Nodes[0]) // Get the pointer to the first node
	return GetNode(index, nodeSize, treePtr) + traverseArrayTreeUnsafe(tree, 2*index) + traverseArrayTreeUnsafe(tree, 2*index+1)
}

func GetNode(index int, nodeSize uintptr, treePointer unsafe.Pointer) int {
	ptr := unsafe.Pointer(uintptr(treePointer) + uintptr(index)*nodeSize)
	return *(*int)(ptr)
}

// Benchmark for pointer-based tree traversal.
func BenchmarkPointerTreeTraversal(b *testing.B) {
	tree := buildPointerTree(24) // Adjust depth as needed
	b.ResetTimer()
	// res := 0
	for i := 0; i < b.N; i++ {
		_ = traversePointerTree(tree)
	}
	// fmt.Println(res, "res")
}

// Benchmark for array-based tree traversal.
func BenchmarkArrayTreeTraversal(b *testing.B) {
	// tree := buildArrayTree(8) // Adjust depth as needed
	tree := buildPointerTree(24) // Adjust depth as needed
	treeArray := ArrayTree{Nodes: ConvertTreeToArray(tree, 1)}
	b.ResetTimer()
	// res := 0
	for i := 0; i < b.N; i++ {
		_ = traverseArrayTree(&treeArray, 1)
	}
	// fmt.Println(res, "res")
}

// Benchmark for unsafe array-based tree traversal.
func BenchmarkArrayTreeTraversalUnsafe(b *testing.B) {
	tree := buildPointerTree(24) // Adjust depth as needed
	treeArray := ArrayTree{Nodes: ConvertTreeToArray(tree, 1)}
	b.ResetTimer()
	// res := 0
	for i := 0; i < b.N; i++ {
		_ = traverseArrayTreeUnsafe(&treeArray, 1)
	}
	// fmt.Println(res, "res")
}
