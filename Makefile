.PHONY: build test testgen

GENDIR ?= ./testdata

# make sure we turn on go modules
export GO111MODULE := on

build:
	go build -mod=readonly ./cmd/testgen-simple

test:
	go test -mod=readonly .

testgen:
	# Usage: GENDIR=CONFIO/PROOFS/testdata/tendermint make testgen
	@mkdir -p "$(GENDIR)"
	go run -mod=readonly ./cmd/testgen-simple exist left 987 > "$(GENDIR)"/exist_left.json
	go run -mod=readonly ./cmd/testgen-simple exist middle 812 > "$(GENDIR)"/exist_middle.json
	go run -mod=readonly ./cmd/testgen-simple exist right 1261 > "$(GENDIR)"/exist_right.json
	go run -mod=readonly ./cmd/testgen-simple nonexist left 813 > "$(GENDIR)"/nonexist_left.json
	go run -mod=readonly ./cmd/testgen-simple nonexist middle 691 > "$(GENDIR)"/nonexist_middle.json
	go run -mod=readonly ./cmd/testgen-simple nonexist right 1535 > "$(GENDIR)"/nonexist_right.json