.PHONY: test

# make sure we turn on go modules
export GO111MODULE := on

build:
	go build -mod=readonly ./cmd/testgen-simple

test:
	go test -mod=readonly .

testgen:
	# Usage: make testgen > CONFIO/PROOFS/testdata/tendermint/existence1.json
	@go run -mod=readonly ./cmd/testgen-simple