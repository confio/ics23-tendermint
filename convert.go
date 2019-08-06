package iavlproofs

import (
	"fmt"
	"math/bits"

	proofs "github.com/confio/proofs/go"
	"github.com/tendermint/tendermint/crypto/merkle"
)

// ConvertExistenceProof will convert the given proof into a valid
// existence proof, if that's what it is.
//
// This is the simplest case of the range proof and we will focus on
// demoing compatibility here
func ConvertExistenceProof(p *merkle.SimpleProof, key, value []byte) (*proofs.ExistenceProof, error) {
	proof := &proofs.ExistenceProof{
		Key:   key,
		Value: value,
		Leaf:  convertLeafOp(),
		Path:  convertInnerOps(p),
	}
	return proof, nil
}

// this is adapted from merkle/hash.go:leafHash()
// and merkle/simple_map.go:KVPair.Bytes()
func convertLeafOp() *proofs.LeafOp {
	prefix := []byte{0}

	return &proofs.LeafOp{
		Hash:         proofs.HashOp_SHA256,
		PrehashKey:   proofs.HashOp_NO_HASH,
		PrehashValue: proofs.HashOp_SHA256,
		Length:       proofs.LengthOp_VAR_PROTO,
		Prefix:       prefix,
	}
}

func convertInnerOps(p *merkle.SimpleProof) []*proofs.InnerOp {
	fmt.Printf("%d of %d\n", p.Index, p.Total)

	idx := p.Index
	// total := p.Total

	var inners []*proofs.InnerOp
	for _, aunt := range p.Aunts {
		// TODO: determine if left or right
		// code adapted from merkle/simple_proof.go:computeHashFromAunts
		auntLeft := idx%2 == 1
		idx = idx / 2

		// combine with: 0x01 || lefthash || righthash
		inner := &proofs.InnerOp{Hash: proofs.HashOp_SHA256}
		if auntLeft {
			inner.Prefix = append([]byte{1}, aunt...)
		} else {
			inner.Prefix = []byte{1}
			inner.Suffix = aunt
		}
		inners = append(inners, inner)
	}
	return inners
}

// // we cannot get the proofInnerNode type, so we need to do the whole path in one function
// func convertInnerOps(path iavl.PathToLeaf) []*proofs.InnerOp {
// 	steps := make([]*proofs.InnerOp, 0, len(path))

// 	// lengthByte is the length prefix prepended to each of the sha256 sub-hashes
// 	var lengthByte byte = 0x20

// 	// we need to go in reverse order, iavl starts from root to leaf,
// 	// we want to go up from the leaf to the root
// 	for i := len(path) - 1; i >= 0; i-- {
// 		// this is adapted from iavl/proof.go:proofInnerNode.Hash()
// 		prefix := aminoVarInt(int64(path[i].Height))
// 		prefix = append(prefix, aminoVarInt(path[i].Size)...)
// 		prefix = append(prefix, aminoVarInt(path[i].Version)...)

// 		var suffix []byte
// 		if len(path[i].Left) > 0 {
// 			// length prefixed left side
// 			prefix = append(prefix, lengthByte)
// 			prefix = append(prefix, path[i].Left...)
// 			// prepend the length prefix for child
// 			prefix = append(prefix, lengthByte)
// 		} else {
// 			// prepend the length prefix for child
// 			prefix = append(prefix, lengthByte)
// 			// length-prefixed right side
// 			suffix = []byte{lengthByte}
// 			suffix = append(suffix, path[i].Right...)
// 		}

// 		op := &proofs.InnerOp{
// 			Hash:   proofs.HashOp_SHA256,
// 			Prefix: prefix,
// 			Suffix: suffix,
// 		}
// 		steps = append(steps, op)
// 	}
// 	return steps
// }

func getSplitPoint(length int) int {
	if length < 1 {
		panic("Trying to split a tree with size < 1")
	}
	uLength := uint(length)
	bitlen := bits.Len(uLength)
	k := 1 << uint(bitlen-1)
	if k == length {
		k >>= 1
	}
	return k
}

// func aminoVarInt(orig int64) []byte {
// 	// amino-specific byte swizzling
// 	i := uint64(orig) << 1
// 	if orig < 0 {
// 		i = ^i
// 	}

// 	// avoid multiple allocs for normal case
// 	res := make([]byte, 0, 8)

// 	// standard protobuf encoding
// 	for i >= 1<<7 {
// 		res = append(res, uint8(i&0x7f|0x80))
// 		i >>= 7
// 	}
// 	res = append(res, uint8(i))
// 	return res
// }
