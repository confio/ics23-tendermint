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

// buildPath returns a list of steps from leaf to root
// in each step, true means index is left side, false index is right side
func buildPath(idx int, total int) []bool {
	if total < 2 {
		return nil
	}
	numLeft := getSplitPoint(total)
	goLeft := idx < numLeft

	// we put goLeft at the end of the array, as we recurse from top to bottom,
	// and want the leaf to be first in array, root last
	if goLeft {
		return append(buildPath(idx, numLeft), goLeft)
	}
	return append(buildPath(idx-numLeft, total-numLeft), goLeft)
}

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
