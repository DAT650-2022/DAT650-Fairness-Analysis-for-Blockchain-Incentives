package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
)

// MerkleTree represents a Merkle tree
type MerkleTree struct {
	RootNode *Node
	Leafs    []*Node
}

// Node represents a Merkle tree node
type Node struct {
	Parent *Node
	Left   *Node
	Right  *Node
	Hash   []byte
}

const (
	leftNode = iota
	rightNode
)

// MerkleProof represents way to prove element inclusion on the merkle tree
type MerkleProof struct {
	proof [][]byte
	index []int64
}

// NewMerkleTree creates a new Merkle tree from a sequence of data
func NewMerkleTree(data [][]byte) *MerkleTree {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	tree := MerkleTree{}
	// create data nodes
	if len(data) == 0 {
		panic("No merkle tree nodes")
	} else if len(data) == 1 {
		// return single data node as root
		tree.RootNode = NewMerkleNode(nil, nil, data[0])
		return &tree
	}
	// add all the data to the bottom of the tree
	for _, j := range data {
		tree.Leafs = append(tree.Leafs, NewMerkleNode(nil, nil, j))
	}
	// add duplicate nodes in case of not power of 2
	// this is a lazy way of creating duplicate inner nodes to fill the tree
	flag := CheckNumberPowerOfTwo(len(tree.Leafs))
	for flag != 0 {
		tree.Leafs = append(tree.Leafs, NewMerkleNode(nil, nil, data[len(data)-1]))
		flag = CheckNumberPowerOfTwo(len(tree.Leafs))
	}
	//if len(tree.Leafs) % 2 == 1{tree.Leafs = append(tree.Leafs, NewMerkleNode(nil, nil, data[len(data)-1]))}
	// builds the inside and the root of the tree
	tree = *buildInner(&tree, tree.Leafs)
	return &tree
}

// checks if power of 2
// from https://www.tutorialspoint.com/golang-program-to-check-whether-given-positive-number-is-power-of-2-or-not-without-using-any-branching-or-loop
// cause lazy
func CheckNumberPowerOfTwo(n int) int {
	return n & (n - 1)
}

func buildInner(tree *MerkleTree, nodes []*Node) *MerkleTree {
	// the set of nodes that were created in the current operation
	var intNodes []*Node
	// iterate by 2 because of left+right
	for i := 0; i < len(nodes); i += 2 {
		var left, right = nodes[i], nodes[i+1]
		// the node with the hash being left and right
		node := NewMerkleNode(left, right, nil)
		// adding the new node to the tree
		tree.Leafs = append(tree.Leafs, node)
		// adding the node to the intermediate nodes for recursive function
		intNodes = append(intNodes, node)
		// wildly shit way of doing this but w/e
		for i, leaf := range tree.Leafs {
			if bytes.Equal(leaf.Hash, left.Hash) {
				leaf.Parent = node
				tree.Leafs[i+1].Parent = node
			}
		}
		if len(nodes) == 2 {
			// if we`re on the last step we can return the tree
			tree.RootNode = node
			return tree
		}
	}
	return buildInner(tree, intNodes)
}

// NewMerkleNode creates a new Merkle tree node
func NewMerkleNode(left, right *Node, data []byte) *Node {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	var hash [32]byte
	// if its a leaf it stores a hash of the data
	if left == nil && right == nil {
		hash = sha256.Sum256(data)
	} else {
		// if its an inner node it stores a hash of its children
		var parents []byte
		parents = append(parents, left.Hash...)
		parents = append(parents, right.Hash...)
		hash = sha256.Sum256(parents)
	}
	return &Node{
		Left:  left,
		Right: right,
		Hash:  hash[:],
	}
}

// MerkleRootHash return the hash of the merkle root node
func (mt *MerkleTree) MerkleRootHash() []byte {
	return mt.RootNode.Hash
}

// MakeMerkleProof returns a list of hashes and indexes required to
// reconstruct the merkle path of a given hash
//
// @param hash represents the hashed data (e.g. transaction ID) stored on
// the leaf node
// @return the merkle proof (list of intermediate hashes), a list of indexes
// indicating the node location in relation with its parent (using the
// constants: leftNode or rightNode), and a possible error.
func (mt *MerkleTree) MakeMerkleProof(hash []byte) ([][]byte, []int64, error) {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	for _, leaf := range mt.Leafs {
		if bytes.Equal(leaf.Hash, hash) {
			currentParent := leaf.Parent
			var merklePath [][]byte
			var indexes []int64
			for currentParent != nil {
				if bytes.Equal(currentParent.Left.Hash, leaf.Hash) {
					merklePath = append(merklePath, currentParent.Right.Hash)
					indexes = append(indexes, 1)
				} else {
					merklePath = append(merklePath, currentParent.Left.Hash)
					indexes = append(indexes, 0)
				}
				leaf = currentParent
				currentParent = currentParent.Parent
			}
			return merklePath, indexes, nil
		}
	}
	return [][]byte{}, []int64{}, fmt.Errorf("Node %x not found", hash)
}

// VerifyProof verifies that the correct root hash can be retrieved by
// recreating the merkle path for the given hash and merkle proof.
//
// @param rootHash is the hash of the current root of the merkle tree
// @param hash represents the hash of the data (e.g. transaction ID)
// to be verified
// @param mProof is the merkle proof that contains the list of intermediate
// hashes and their location on the tree required to reconstruct
// the merkle path.
func VerifyProof(rootHash []byte, hash []byte, mProof MerkleProof) bool {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	var hashed [32]byte
	for i, j := range mProof.proof {
		var temphash []byte
		if mProof.index[i] == 1 {
			temphash = append(temphash, hash...)
			temphash = append(temphash, j...)
		} else {
			temphash = append(temphash, j...)
			temphash = append(temphash, hash...)
		}
		hashed = sha256.Sum256(temphash)
		hash = hashed[:]
	}
	return bytes.Equal(rootHash, hash)
}

// this one wasn't here before but I added it anyway just in case

// VerifyHash verifies whether a hash exists in the merkle tree
// by iterating over all leaves.
//
// @param hash represents the hash of the data (e.g. transaction ID)
// to be verified
func (mt *MerkleTree) VerifyHash(hash []byte) bool {
	// TODO(student) -- YOU DON'T NEED TO CHANGE YOUR PREVIOUS METHOD
	// wdym don't need to change previous, we hadn't made this before?
	for _, j := range mt.Leafs {
		if bytes.Equal(j.Hash, hash) {
			return true
		}
	}
	return false
}
