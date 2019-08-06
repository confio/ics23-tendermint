package iavlproofs

import (
	"bytes"
	"fmt"
	"testing"
)

func TestLeafOp(t *testing.T) {
	proof := GenerateRangeProof(20)

	converted, err := ConvertExistenceProof(proof.Proof, proof.Key, proof.Value)
	if err != nil {
		t.Fatal(err)
	}

	leaf := converted.GetLeaf()
	if leaf == nil {
		t.Fatalf("Missing leaf node")
	}

	hash, err := leaf.Apply(converted.Key, converted.Value)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(hash, proof.Proof.LeafHash) {
		t.Errorf("Calculated: %X\nExpected:   %X", hash, proof.Proof.LeafHash)
	}
}

func XTestConvertProof(t *testing.T) {
	for i := 0; i < 2; i++ {
		t.Run(fmt.Sprintf("Run %d", i), func(t *testing.T) {
			proof := GenerateRangeProof(200)

			converted, err := ConvertExistenceProof(proof.Proof, proof.Key, proof.Value)
			if err != nil {
				t.Fatal(err)
			}

			calc, err := converted.Calculate()
			if err != nil {
				t.Fatal(err)
			}

			if !bytes.Equal(calc, proof.RootHash) {
				t.Errorf("Calculated: %X\nExpected:   %X", calc, proof.RootHash)
			}
		})
	}
}
