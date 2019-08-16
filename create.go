package tmproofs


import (
	"fmt"
	"sort"

	proofs "github.com/confio/proofs/go"
	"github.com/confio/proofs-tendermint/helpers"
	"github.com/tendermint/tendermint/crypto/merkle"
)

/*
CreateMembershipProof will produce a CommitmentProof that the given key (and queries value) exists in the iavl tree.
If the key doesn't exist in the tree, this will return an error.
*/
func CreateMembershipProof(data map[string][]byte, key []byte) (*proofs.CommitmentProof, error) {
	exist, err := createExistenceProof(data, key)
	if err != nil {
		return nil, err
	}
	proof := &proofs.CommitmentProof{
		Proof: &proofs.CommitmentProof_Exist{
			Exist: exist,
		},
	}
	return proof, nil
}

/*
CreateNonMembershipProof will produce a CommitmentProof that the given key doesn't exist in the iavl tree.
If the key exists in the tree, this will return an error.
*/
func CreateNonMembershipProof(data map[string][]byte, key []byte) (*proofs.CommitmentProof, error) {
	// ensure this key is not in the store
	if _, ok := data[string(key)]; ok {
		return nil, fmt.Errorf("Cannot create non-membership proof if key is in map")
	}

	keys := helpers.SortedKeys(data)
	rightidx := sort.SearchStrings(keys, string(key))

	var err error
	nonexist := &proofs.NonExistenceProof{
		Key: key,
	}

	// include left proof unless key is left of entire map
	if rightidx >= 1 {
		leftkey := keys[rightidx-1]
		nonexist.Left, err = createExistenceProof(data, []byte(leftkey))
		if err != nil {
			return nil, err
		}
	}

	// include right proof unless key is right of entire map
	if rightidx < len(keys) {
		rightkey := keys[rightidx]
		nonexist.Right, err = createExistenceProof(data, []byte(rightkey))
		if err != nil {
			return nil, err
		}

	}

	proof := &proofs.CommitmentProof{
		Proof: &proofs.CommitmentProof_Nonexist{
			Nonexist: nonexist,
		},
	}
	return proof, nil
}

func createExistenceProof(data map[string][]byte, key []byte) (*proofs.ExistenceProof, error) {
	value, ok := data[string(key)]
	if !ok {
		return nil, fmt.Errorf("cannot make existence proof if key is not in map")
	}
	
	_, proofs, _ := merkle.SimpleProofsFromMap(data)
	proof := proofs[string(key)]
	if proof == nil {
		return nil, fmt.Errorf("returned no proof for key")
	}

	return convertExistenceProof(proof, key, value)
}

