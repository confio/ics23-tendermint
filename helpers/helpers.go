package helpers

import (
	"github.com/tendermint/tendermint/crypto/merkle"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// SimpleResult contains a merkle.SimpleProof along with all data needed to build the confio/proof
type SimpleResult struct {
	Key      []byte
	Value    []byte
	Proof    *merkle.SimpleProof
	RootHash []byte
}

// GenerateRangeProof makes a tree of size and returns a range proof for one random element
//
// returns a range proof and the root hash of the tree
func GenerateRangeProof(size int) *SimpleResult {
	data := make(map[string][]byte)

	toValue := func(key string) []byte {
		return []byte("value_for_" + key)
	}

	// insert lots of info and store the bytes
	for i := 0; i < size; i++ {
		key := cmn.RandStr(20)
		data[key] = toValue(key)
	}

	root, proofs, allkeys := merkle.SimpleProofsFromMap(data)

	// grab a random key
	idx := cmn.RandInt() % size
	key := allkeys[idx]
	// along with it's proof
	proof := proofs[key]

	res := &SimpleResult{
		Key:      []byte(key),
		Value:    toValue(key),
		Proof:    proof,
		RootHash: root,
	}
	return res
}
