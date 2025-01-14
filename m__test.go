// // package main__test

// import (
// 	"math/rand"
// 	"testing"
// )

// // Define a node structure for the pointer-based binary tree.
// type Vector struct {
// 	X, Y, Z float32
// }

// type Triangle struct {
// 	Normal             Vector
// 	V1, V2, V3         Vector
// 	R, G, B, A         uint8
// 	Specular           float32
// 	DirectToScattering float32
// }

// type BoundingBox struct {
// 	Min, Max Vector
// }

// type TreeNode struct {
// 	Triangle    *Triangle
// 	BoundingBox BoundingBox
// 	Left, Right *TreeNode
// }

// type bBoxTreeNode struct {
// 	BoundingBox BoundingBox
// 	Triangle    *Triangle
// }

// type bBoxTreeNodeTriangle struct {
// 	BoundingBox BoundingBox
// 	Triangle    Triangle
// }

// type bBoxTreeNodeIndex struct {
// 	BoundingBox BoundingBox
// 	Triangle    int32
// }

// func BenchmarkBVH(b *testing.B) {
// 	tests := []struct {
// 		name string
// 		fn   func(b *testing.B)
// 	}{	
// 		{"BVH", testBvh},
// 		{"BVH_Standard", test_standard},
// 		{"BVH_Pointer", test_pointer},
// 		{"BVH_Struct", test_struct},
// 		{"BVH_Index", test_index},
// 	}

// 	for _, tt := range tests {
// 		b.Run(tt.name, tt.fn)
// 	}
// }

// func (b *TreeNode) GenerateBVH(depth int) {
// 	if depth == 0 {
// 		return
// 	}
// 	b.Left.GenerateBVH(depth - 1)
// 	b.Right.GenerateBVH(depth - 1)
// }

// func GenerateBVH_Pointer(depth int) (tree []bBoxTreeNode) {
// 	numNodes := 1 << depth
// 	tree = make([]bBoxTreeNode, numNodes)
// 	triangle := make([]Triangle, numNodes)
// 	for i := 0; i < numNodes; i++ {
// 		tree[i].Triangle = &triangle[i]
// 	}
// 	return tree
// }

// func GenerateBVH_Struct(depth int) (tree []bBoxTreeNodeTriangle) {
// 	numNodes := 1 << depth
// 	tree = make([]bBoxTreeNodeTriangle, numNodes)
// 	for i := 0; i < numNodes; i++ {
// 		tree[i].Triangle = Triangle{}
// 	}
// 	return tree
// }

// func GenerateBVH_Index(depth int) (tree []bBoxTreeNodeIndex, triangles []Triangle) {
// 	numNodes := 1 << depth
// 	tree = make([]bBoxTreeNodeIndex, numNodes)
// 	triangles = make([]Triangle, numNodes)
// 	for i := 0; i < numNodes; i++ {
// 		tree[i].Triangle = int32(i)
// 	}
// 	return tree, triangles
// }

// const depth = 10

// func (n *TreeNode) getNode(leftRight bool) *TreeNode {
// 	if leftRight {
// 		return n.Left
// 	}
// 	return n.Right
// }

// func testBvh(b *testing.B) {
// 	tree := &TreeNode{}
// 	tree.GenerateBVH(depth)
// 	finalNode := &TreeNode{}

// 	// start the timer
// 	b.ResetTimer()
// 	for i := 0; i < depth; i++ {
// 		// generate random choice of left and right
// 		if rand.Intn(2) == 0 {
// 			finalNode = tree.getNode(true)
// 			_ = *finalNode.Triangle
// 		} else {
// 			finalNode = tree.getNode(false)
// 			_ = *finalNode.Triangle
// 		}
// 	}

// }

// func test_standard(b *testing.B) {
// 	tree := GenerateBVH_Pointer(depth)

// 	// start the timer
// 	b.ResetTimer()

// 	for i := 0; i < depth; i++ {
// 		left := depth * 2
// 		right := depth*2 + 1
// 		if rand.Intn(2) == 0 {
// 			t := tree[left]
// 			_ = t.Triangle
// 		} else {
// 			t := tree[right]
// 			_ = t.Triangle
// 		}
// 	}

// }

// func test_pointer(b *testing.B) {
// 	tree := GenerateBVH_Pointer(depth)

// 	// start the timer
// 	b.ResetTimer()

// 	for i := 0; i < depth; i++ {
// 		left := depth * 2
// 		right := depth*2 + 1
// 		// generate random choice of left and right
// 		if rand.Intn(2) == 0 {
// 			t := tree[left]
// 			_ = *t.Triangle
// 		} else {
// 			t := tree[right]
// 			_ = t.Triangle
// 		}
// 	}

// }

// func test_struct(b *testing.B) {
// 	tree := GenerateBVH_Struct(depth)

// 	// start the timer
// 	b.ResetTimer()

// 	for i := 0; i < depth; i++ {
// 		left := depth * 2
// 		right := depth*2 + 1
// 		// generate random choice of left and right
// 		if rand.Intn(2) == 0 {
// 			t := tree[left]
// 			_ = t.Triangle
// 		} else {
// 			t := tree[right]
// 			_ = t.Triangle
// 		}
// 	}

// }

// func test_index(b *testing.B) {
// 	tree, triangles := GenerateBVH_Index(depth)

// 	// start the timer
// 	b.ResetTimer()

// 	for i := 0; i < depth; i++ {
// 		left := depth * 2
// 		right := depth*2 + 1
// 		// generate random choice of left and right
// 		if rand.Intn(2) == 0 {
// 			t := tree[left]
// 			_ = triangles[t.Triangle]
// 		} else {
// 			t := tree[right]
// 			_ = triangles[t.Triangle]
// 		}
// 	}
// }