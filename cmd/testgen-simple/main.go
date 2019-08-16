package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"

	tmproofs "github.com/confio/proofs-tendermint"
	"github.com/confio/proofs-tendermint/helpers"
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
	data := helpers.BuildMap(157)
	root := helpers.CalcRoot(data)

	keys := helpers.SortedKeys(data)
	key := []byte(helpers.GetKey(keys, helpers.Middle))

	converted, err := tmproofs.CreateMembershipProof(data, key)
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
		"root":      hex.EncodeToString(root),
		"existence": hex.EncodeToString(binary),
	}
	out, err := json.MarshalIndent(res, "", "  ")
	if err != nil {
		fmt.Printf("Error: json encoding: %+v\n", err)
		os.Exit(1)
	}

	fmt.Println(string(out))
}
