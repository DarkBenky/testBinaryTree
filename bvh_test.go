package main__test

import (
	"math/rand"
	"testing"
	"unsafe"
)

// Define a node structure for the pointer-based binary tree.
type Vector struct {
	X, Y, Z float32
}

type Triangle struct {
	Normal             Vector
	V1, V2, V3         Vector
	R, G, B, A         uint8
	Specular           float32
	DirectToScattering float32
}

type BoundingBox struct {
	Min, Max Vector
}

type TreeNode struct {
	Value       int
	Triangle    *Triangle
	BoundingBox BoundingBox
	Left, Right *TreeNode
}

type TreeNodeNoPtr struct {
	Value       int
	Triangle    Triangle
	BoundingBox BoundingBox
	Left, Right *TreeNodeNoPtr
}

type TreeNodePtr struct {
	Value       int
	Triangle    *Triangle
	BoundingBox *BoundingBox
	Left, Right *TreeNodePtr
}

type TreeArrayNodePtr struct {
	Triangle    *Triangle
	BoundingBox *BoundingBox
}

type TreeArrayNodeNoPtr struct {
	Triangle        Triangle
	boundingBoxFlag bool
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

var t = Triangle{
	Normal:             Vector{X: 1, Y: 2, Z: 3},
	V1:                 Vector{X: 4, Y: 5, Z: 6},
	V2:                 Vector{X: 7, Y: 8, Z: 9},
	V3:                 Vector{X: 10, Y: 11, Z: 12},
	R:                  13,
	G:                  14,
	B:                  15,
	A:                  16,
	Specular:           17,
	DirectToScattering: 18,
}

func buildPointerTreeNoPtr(depth int) *TreeNodeNoPtr {
	if depth == 0 {
		return nil
	}
	return &TreeNodeNoPtr{
		Value:    depth,
		Triangle: t, // Direct assignment since it's not a pointer
		BoundingBox: BoundingBox{
			Min: Vector{X: 1, Y: 2, Z: 3},
			Max: Vector{X: 4, Y: 5, Z: 6},
		},
		Left:  buildPointerTreeNoPtr(depth - 1),
		Right: buildPointerTreeNoPtr(depth - 1),
	}
}

func buildPointerTreePtr(depth int) *TreeNodePtr {
	if depth == 0 {
		return nil
	}
	return &TreeNodePtr{
		Value:    depth,
		Triangle: &t, // Assign the pointer to the triangle
		BoundingBox: &BoundingBox{
			Min: Vector{X: 1, Y: 2, Z: 3},
			Max: Vector{X: 4, Y: 5, Z: 6},
		},
		Left:  buildPointerTreePtr(depth - 1),
		Right: buildPointerTreePtr(depth - 1),
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

func ConvertToFixedArray(root *TreeNode, index int) [65536]int {
	var array [65536]int
	if root == nil {
		return array
	}
	array[index] = root.Value
	ConvertToFixedArray(root.Left, index*2)
	ConvertToFixedArray(root.Right, index*2+1)
	return array
}

func ConvertToFixedArrayNoPtr(root *TreeNodeNoPtr, index int) [65536]TreeArrayNodeNoPtr {
	var array [65536]TreeArrayNodeNoPtr
	if root == nil {
		return array
	}
	array[index] = TreeArrayNodeNoPtr{
		Triangle:        root.Triangle,
		boundingBoxFlag: true,
	}
	ConvertToFixedArrayNoPtr(root.Left, index*2)
	ConvertToFixedArrayNoPtr(root.Right, index*2+1)
	return array
}

func ConvertToFixedArrayPtr(root *TreeNodePtr, index int) [65536]TreeArrayNodePtr {
	var array [65536]TreeArrayNodePtr
	if root == nil {
		return array
	}
	array[index] = TreeArrayNodePtr{
		Triangle:    root.Triangle,
		BoundingBox: root.BoundingBox,
	}
	ConvertToFixedArrayPtr(root.Left, index*2)
	ConvertToFixedArrayPtr(root.Right, index*2+1)
	return array
}

// Recursive traversal for the pointer-based binary tree.
func traversePointerTree(node *TreeNode) (sum int) {
	if node == nil {
		return 0
	}

	if node.Triangle != nil && node.Triangle.A == 16 {
		sum += 1
	}
	if node.BoundingBox.Max.Z == 6 {
		sum += 1
	}
	sum += node.Value

	return node.Value + traversePointerTree(node.Left) + traversePointerTree(node.Right)
}

func traversePointerTreeNoPtr(node *TreeNodeNoPtr) (sum int) {
	if node == nil {
		return 0
	}

	if node.Triangle.A == 16 {
		sum += 1
	}
	if node.BoundingBox.Max.Z == 6 {
		sum += 1
	}
	sum += node.Value

	return node.Value + traversePointerTreeNoPtr(node.Left) + traversePointerTreeNoPtr(node.Right)
}

func traversePointerTreePtr(node *TreeNodePtr) (sum int) {
	if node == nil {
		return 0
	}
	if node.Triangle != nil && node.Triangle.A == 16 {
		sum += 1
	}
	if node.BoundingBox != nil && node.BoundingBox.Max.Z == 6 {
		sum += 1
	}

	sum += node.Value

	return node.Value + traversePointerTreePtr(node.Left) + traversePointerTreePtr(node.Right)
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

func traverseArrayTreeFixed(tree *[65536]int, index int) int {
	if index >= len(tree) {
		return 0
	}
	return tree[index] + traverseArrayTreeFixed(tree, 2*index) + traverseArrayTreeFixed(tree, 2*index+1)
}

func traverseArrayTreeFixedNoPtr(tree *[65536]TreeArrayNodeNoPtr, index int) (sum int) {
	if index >= len(tree) {
		return 0
	}
	if tree[index].boundingBoxFlag && tree[index].Triangle.A == 16 {
		sum += 1
	}
	if tree[index].Triangle.A == 16 {
		sum += 1
	}

	sum += traverseArrayTreeFixedNoPtr(tree, 2*index) + traverseArrayTreeFixedNoPtr(tree, 2*index+1)

	return sum
}

func traverseArrayTreeFixedPtr(tree *[65536]TreeArrayNodePtr, index int) (sum int) {
	if index >= len(tree) {
		return 0
	}
	if tree[index].Triangle != nil && tree[index].Triangle.A == 16 {
		sum += 1
	}
	if tree[index].BoundingBox != nil && tree[index].BoundingBox.Max.Z == 6 {
		sum += 1
	}
	sum += traverseArrayTreeFixedPtr(tree, 2*index) + traverseArrayTreeFixedPtr(tree, 2*index+1)

	return sum
}

func traverseArrayTreeFixedIterative(tree *[65536]int) int {
	sum := 0
	stack := []int{1} // Start with index 1 to match recursive version

	for len(stack) > 0 {
		index := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if index >= len(tree) {
			continue
		}

		sum += tree[index]

		// Match the recursive version's index calculations
		leftChild := 2 * index
		rightChild := 2*index + 1

		// Push children if they're within bounds
		if rightChild < len(tree) {
			stack = append(stack, rightChild)
		}
		if leftChild < len(tree) {
			stack = append(stack, leftChild)
		}
	}

	return sum
}

func traverseArrayTreeUnsafe(tree *ArrayTree, index int, nodeSize uintptr) int {
	if index >= len(tree.Nodes) {
		return 0
	}
	// nodeSize := unsafe.Sizeof(tree.Nodes[0])  // Get the size of an int (node) in bytes
	treePtr := unsafe.Pointer(&tree.Nodes[0]) // Get the pointer to the first node
	nodeSize = unsafe.Sizeof(tree.Nodes[0])
	return GetNode(index, nodeSize, treePtr) + traverseArrayTreeUnsafe(tree, 2*index, nodeSize) + traverseArrayTreeUnsafe(tree, 2*index+1, nodeSize)
}

func GetNode(index int, nodeSize uintptr, treePointer unsafe.Pointer) int {
	ptr := unsafe.Pointer(uintptr(treePointer) + uintptr(index)*nodeSize)
	return *(*int)(ptr)
}

// func BenchmarkFixedArrayTreeTraversalIterative(b *testing.B) {
// 	tree := buildPointerTree(16) // Adjust depth as needed
// 	treeArray := ConvertToFixedArray(tree, 1)
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traverseArrayTreeFixedIterative(&treeArray)
// 	}
// 	// fmt.Println(res, "res")
// }

// func BenchmarkFixedArrayTreeTraversal(b *testing.B) {
// 	tree := buildPointerTree(16) // Adjust depth as needed
// 	treeArray := ConvertToFixedArray(tree, 1)
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traverseArrayTreeFixed(&treeArray, 1)
// 	}
// 	// fmt.Println(res, "res")
// }

// // Benchmark for pointer-based tree traversal.
// func BenchmarkPointerTreeTraversal(b *testing.B) {
// 	tree := buildPointerTree(16) // Adjust depth as needed
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traversePointerTree(tree)
// 	}
// 	// fmt.Println(res, "res")
// }

// // Benchmark for array-based tree traversal.
// func BenchmarkArrayTreeTraversal(b *testing.B) {
// 	// tree := buildArrayTree(8) // Adjust depth as needed
// 	tree := buildPointerTree(16) // Adjust depth as needed
// 	treeArray := ArrayTree{Nodes: ConvertTreeToArray(tree, 1)}
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traverseArrayTree(&treeArray, 1)
// 	}
// 	// fmt.Println(res, "res")
// }

// // Benchmark for unsafe array-based tree traversal.
// func BenchmarkArrayTreeTraversalUnsafe(b *testing.B) {
// 	tree := buildPointerTree(16) // Adjust depth as needed
// 	treeArray := ArrayTree{Nodes: ConvertTreeToArray(tree, 1)}
// 	b.ResetTimer()
// 	// res := 0
// 	nodeSize := unsafe.Sizeof(treeArray.Nodes[0])
// 	for i := 0; i < b.N; i++ {
// 		_ = traverseArrayTreeUnsafe(&treeArray, 1, nodeSize)
// 	}
// 	// fmt.Println(res, "res")
// }

// func BenchmarkPointerTreeTraversalNoPtr(b *testing.B) {
// 	tree := buildPointerTreeNoPtr(16) // Adjust depth as needed
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traversePointerTreeNoPtr(tree)
// 	}
// 	// fmt.Println(res, "res")
// }

// func BenchmarkPointerTreeTraversalPtr(b *testing.B) {
// 	tree := buildPointerTreePtr(16) // Adjust depth as needed
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traversePointerTreePtr(tree)
// 	}
// 	// fmt.Println(res, "res")
// }

// func BenchmarkFixedArrayTreeTraversalPtr(b *testing.B) {
// 	tree := buildPointerTreePtr(16) // Adjust depth as needed
// 	treeArray := ConvertToFixedArrayPtr(tree, 1)
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traverseArrayTreeFixedPtr(&treeArray, 1)
// 	}
// 	// fmt.Println(res, "res")
// }

// func BenchmarkFixedArrayTreeTraversalNoPtr(b *testing.B) {
// 	tree := buildPointerTreeNoPtr(16) // Adjust depth as needed
// 	treeArray := ConvertToFixedArrayNoPtr(tree, 1)
// 	b.ResetTimer()
// 	// res := 0
// 	for i := 0; i < b.N; i++ {
// 		_ = traverseArrayTreeFixedNoPtr(&treeArray, 1)
// 	}
// 	// fmt.Println(res, "res")
// }

type bBoxTreeNode struct {
	BoundingBox BoundingBox
	Triangle    *Triangle
}

type bBoxTreeNodeTriangle struct {
	BoundingBox BoundingBox
	Triangle    Triangle
}

type bBoxTreeNodeIndex struct {
	BoundingBox BoundingBox
	Triangle    int32
}

func (b *TreeNode) GenerateBVH(depth int) {
	if depth == 0 {
		return
	}

	t := Triangle{
		Normal:             Vector{X: 1, Y: 2, Z: 3},
		V1:                 Vector{X: 4, Y: 5, Z: 6},
		V2:                 Vector{X: 7, Y: 8, Z: 9},
		V3:                 Vector{X: 10, Y: 11, Z: 12},
		R:                  13,
		G:                  14,
		B:                  15,
		A:                  16,
		Specular:           17,
		DirectToScattering: 18,
	}

	// Initialize Left and Right nodes if they don't exist
	if b.Left == nil {
		b.Left = &TreeNode{
			Triangle: &t,
			BoundingBox: BoundingBox{
				Min: Vector{X: 0, Y: 0, Z: 0},
				Max: Vector{X: 1, Y: 1, Z: 1},
			},
		}
	}

	if b.Right == nil {
		b.Right = &TreeNode{
			Triangle: &t,
			BoundingBox: BoundingBox{
				Min: Vector{X: 0, Y: 0, Z: 0},
				Max: Vector{X: 1, Y: 1, Z: 1},
			},
		}
	}

	b.Left.GenerateBVH(depth - 1)
	b.Right.GenerateBVH(depth - 1)
}

func GenerateBVH_Pointer(depth int) (tree []bBoxTreeNode) {
	numNodes := 1 << depth
	tree = make([]bBoxTreeNode, numNodes)
	triangle := make([]Triangle, numNodes)
	for i := 0; i < numNodes; i++ {
		tree[i].Triangle = &triangle[i]
	}
	return tree
}

func GenerateBVH_Struct(depth int) (tree []bBoxTreeNodeTriangle) {
	numNodes := 1 << depth
	tree = make([]bBoxTreeNodeTriangle, numNodes)
	for i := 0; i < numNodes; i++ {
		tree[i].Triangle = Triangle{}
	}
	return tree
}

func GenerateBVH_Index(depth int) (tree []bBoxTreeNodeIndex, triangles []Triangle) {
	numNodes := 1 << depth
	tree = make([]bBoxTreeNodeIndex, numNodes)
	triangles = make([]Triangle, numNodes)
	for i := 0; i < numNodes; i++ {
		tree[i].Triangle = int32(i)
	}
	return tree, triangles
}

const d = 20

func (n *TreeNode) getNode(leftRight bool) *TreeNode {
	if leftRight {
		return n.Left
	}
	return n.Right
}

func BenchmarkBvh(b *testing.B) {
	// Create root node
	root := &TreeNode{
		Triangle: &Triangle{},
		BoundingBox: BoundingBox{
			Min: Vector{X: 0, Y: 0, Z: 0},
			Max: Vector{X: 1, Y: 1, Z: 1},
		},
	}

	root.GenerateBVH(d)

	b.ResetTimer()
	for i := 0; i < d; i++ {
		if rand.Intn(2) == 0 {
			root = root.getNode(true)
			_ = root.Triangle
		} else {
			root = root.getNode(false)
			_ = root.Triangle
		}
	}
}

func BenchmarkStandard(b *testing.B) {
	tree := GenerateBVH_Pointer(d)

	// start the timer
	b.ResetTimer()

	for i := 0; i < d; i++ {
		left := i * 2
		right := i*2 + 1
		if rand.Intn(2) == 0 {
			t := tree[left]
			_ = *t.Triangle
		} else {
			t := tree[right]
			_ = *t.Triangle
		}
	}
	// fmt.Println("final result : ", tt)
}

func BenchmarkPointer(b *testing.B) {
	tree := GenerateBVH_Pointer(d)

	// start the timer
	b.ResetTimer()

	for i := 0; i < d; i++ {
		left := i * 2
		right := i*2 + 1
		// generate random choice of left and right
		if rand.Intn(2) == 0 {
			t := tree[left]
			_ = *t.Triangle
		} else {
			t := tree[right]
			_ = *t.Triangle
		}
	}
	// fmt.Println("final result : ", tt)
}

func BenchmarkStruct(b *testing.B) {
	tree := GenerateBVH_Struct(d)
	// start the timer
	b.ResetTimer()

	for i := 0; i < d; i++ {
		left := i * 2
		right := i*2 + 1
		// generate random choice of left and right
		if rand.Intn(2) == 0 {
			t := tree[left]
			_ = t.Triangle
		} else {
			t := tree[right]
			_ = t.Triangle
		}
	}
	// fmt.Println("final result : ", tt)
}

func BenchmarkIndex(b *testing.B) {
	tree, triangles := GenerateBVH_Index(d)
	// start the timer
	b.ResetTimer()

	for i := 0; i < d; i++ {
		left := i * 2
		right := i*2 + 1
		// generate random choice of left and right

		if rand.Intn(2) == 0 {
			t := tree[left]
			_ = triangles[t.Triangle]
		} else {
			t := tree[right]
			_ = triangles[t.Triangle]
		}
	}

	// fmt.Println("final result : ", tt)
}
