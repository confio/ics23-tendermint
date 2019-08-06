package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	tmproofs "github.com/confio/proofs-tendermint"
)

/**
testgen-simple will generate a json struct on stdout (meant to be saved to file for testdata).
this will be an auto-generated existence proof in the form:

{
	"root": "<hex encoded root hash of tree>",
	"existence": "<hex encoded protobuf marshaling of an existence proof>"
}
**/

func main() {
	proof := tmproofs.GenerateRangeProof(157)

	converted, err := tmproofs.ConvertExistenceProof(proof.Proof, proof.Key, proof.Value)
	if err != nil {
		fmt.Printf("Error: convert proof: %+v\n", err)
		os.Exit(1)
	}

	binary, err := converted.Marshal()
	if err != nil {
		fmt.Printf("Error: protobuf marshal: %+v\n", err)
		os.Exit(1)
	}

	res := map[string]interface{}{
		"root":      hex.EncodeToString(proof.RootHash),
		"existence": hex.EncodeToString(binary),
	}
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("Error: json encoding: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}
